package livekit

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/livekit/media-sdk"
	"github.com/livekit/protocol/auth"
	"github.com/livekit/protocol/livekit"
	"github.com/livekit/protocol/logger"
	lksdk "github.com/livekit/server-sdk-go/v2"
	lkmedia "github.com/livekit/server-sdk-go/v2/pkg/media"
	"github.com/pion/webrtc/v4"
	"github.com/rahulSailesh-shah/converSense/internal/db/repo"
	"github.com/rahulSailesh-shah/converSense/pkg/config"
	sentimentanalyzer "github.com/rahulSailesh-shah/converSense/pkg/sentiment-analyzer"
)

type SessionCallbacks struct {
	OnMeetingEnd func(meetingID string, recordingURL string, transcriptURL string, err error)
}

type StreamTextData struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

type LiveKitSession struct {
	meetingDetails  *repo.GetMeetingRow
	userDetails     *repo.User
	room            *lksdk.Room
	handler         *GeminiRealtimeAPIHandler
	egressInfo      *livekit.EgressInfo
	lkConfig        *config.LiveKitConfig
	geminiConfig    *config.GeminiConfig
	awsConfig       *config.AWSConfig
	ctx             context.Context
	cancel          context.CancelFunc
	callbacks       SessionCallbacks
	recordingURL    string
	transcriptURL   string
	stopOnce        sync.Once
	textStreamQueue chan StreamTextData
}

func NewLiveKitSession(
	meetingDetails *repo.GetMeetingRow,
	userDetails *repo.User,
	lkConfig *config.LiveKitConfig,
	geminiConfig *config.GeminiConfig,
	awsConfig *config.AWSConfig,
	callbacks SessionCallbacks,
) *LiveKitSession {
	ctx, cancel := context.WithCancel(context.Background())

	return &LiveKitSession{
		meetingDetails:  meetingDetails,
		userDetails:     userDetails,
		lkConfig:        lkConfig,
		geminiConfig:    geminiConfig,
		awsConfig:       awsConfig,
		ctx:             ctx,
		cancel:          cancel,
		callbacks:       callbacks,
		stopOnce:        sync.Once{},
		textStreamQueue: make(chan StreamTextData, 100),
	}
}

func (s *LiveKitSession) Start() error {
	if err := s.connectBot(); err != nil {
		return fmt.Errorf("failed to connect bot: %w", err)
	}
	return nil
}

func (s *LiveKitSession) Stop() error {
	var stopErr error
	meetingId := s.meetingDetails.ID.String()
	s.stopOnce.Do(func() {
		s.cancel()
		if s.egressInfo != nil {
			// if err := s.stopRecording(s.egressInfo.EgressId); err != nil {
			// 	stopErr = fmt.Errorf("failed to stop recording: %w", err)
			// }
		}
		if s.textStreamQueue != nil {
			close(s.textStreamQueue)
		}
		if s.room != nil {
			s.room.Disconnect()
		}
		if s.handler != nil {
			s.handler.Close()
		}
		if s.callbacks.OnMeetingEnd != nil {
			s.callbacks.OnMeetingEnd(meetingId, s.recordingURL, s.transcriptURL, stopErr)
		}
	})
	return stopErr
}

func (s *LiveKitSession) GenerateUserToken() (string, error) {
	at := auth.NewAccessToken(s.lkConfig.APIKey, s.lkConfig.APISecret)
	grant := &auth.VideoGrant{
		RoomJoin: true,
		Room:     s.meetingDetails.ID.String(),
	}
	at.SetVideoGrant(grant).
		SetIdentity(s.userDetails.Name).
		SetValidFor(time.Hour)
	token, err := at.ToJWT()
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s *LiveKitSession) connectBot() error {
	sentimentAnalyzer, err := sentimentanalyzer.NewSentimentAnalyzer(sentimentanalyzer.AnalyzerTypeOllama)
	if err != nil {
		logger.Errorw("Failed to create sentiment analyzer", err, "meetingID", s.meetingDetails.ID.String())
		return fmt.Errorf("failed to create sentiment analyzer: %w", err)
	}
	audioWriterChan := make(chan media.PCM16Sample, 500)
	handler, err := NewGeminiRealtimeAPIHandler(s.ctx,
		s.geminiConfig,
		s.userDetails,
		s.meetingDetails,
		&GeminiRealtimeAPIHandlerCallbacks{
			OnAudioReceived: func(audio media.PCM16Sample) {
				select {
				case audioWriterChan <- audio:
				case <-s.ctx.Done():
					return
				default:
					logger.Warnw("audio writer channel is full", fmt.Errorf("audio writer channel is full"))
				}
			},
			OnUserSentiment: func(result *sentimentanalyzer.SentimentResult) {
				streamTextData := StreamTextData{
					Type: "sentiment",
					Data: result,
				}
				select {
				case s.textStreamQueue <- streamTextData:
				default:
					logger.Warnw("Text stream queue full, dropping sentiment message", nil)
				}
			},
			OnUserTranscript: func(result *TranscriptDataStream) {
				streamTextData := StreamTextData{
					Type: "transcript",
					Data: result,
				}
				select {
				case s.textStreamQueue <- streamTextData:
				default:
					logger.Warnw("Text stream queue full, dropping transcript message", nil)
				}
			},
		}, sentimentAnalyzer)
	if err != nil {
		close(audioWriterChan)
		return fmt.Errorf("failed to create Gemini handler: %w", err)
	}
	s.handler = handler

	if err := s.connectToRoom(); err != nil {
		s.handler.Close()
		close(audioWriterChan)
		return fmt.Errorf("failed to connect to room: %w", err)
	}

	go s.handlePublish(audioWriterChan)
	go s.handleTextStreamQueue()

	// egressInfo, err := s.startRecording()
	// if err != nil {
	// 	logger.Errorw("Failed to start recording", err, "meetingID", s.meetingDetails.ID.String())
	// } else {
	// 	s.egressInfo = egressInfo
	// }
	return nil
}

func (s *LiveKitSession) connectToRoom() error {
	room, err := lksdk.ConnectToRoom(s.lkConfig.Host, lksdk.ConnectInfo{
		APIKey:              s.lkConfig.APIKey,
		APISecret:           s.lkConfig.APISecret,
		RoomName:            s.meetingDetails.ID.String(),
		ParticipantIdentity: s.meetingDetails.AgentName,
	}, s.callbacksForRoom())
	if err != nil {
		return err
	}
	s.room = room
	return nil
}

func (s *LiveKitSession) callbacksForRoom() *lksdk.RoomCallback {
	var pcmRemoteTrack *lkmedia.PCMRemoteTrack

	return &lksdk.RoomCallback{
		ParticipantCallback: lksdk.ParticipantCallback{
			OnTrackSubscribed: func(track *webrtc.TrackRemote, publication *lksdk.RemoteTrackPublication, rp *lksdk.RemoteParticipant) {
				if pcmRemoteTrack != nil {
					return
				}
				pcmRemoteTrack, _ = s.handleSubscribe(track)
			},
		},
		OnParticipantDisconnected: func(participant *lksdk.RemoteParticipant) {
			s.Stop()
		},
		OnDisconnected: func() {
			if pcmRemoteTrack != nil {
				pcmRemoteTrack.Close()
				pcmRemoteTrack = nil
			}
		},
		OnDisconnectedWithReason: func(reason lksdk.DisconnectionReason) {
			if pcmRemoteTrack != nil {
				pcmRemoteTrack.Close()
				pcmRemoteTrack = nil
			}
		},
	}
}

func (s *LiveKitSession) handlePublish(audioWriterChan chan media.PCM16Sample) {
	publishTrack, err := lkmedia.NewPCMLocalTrack(24000, 1, logger.GetLogger())
	if err != nil {
		logger.Errorw("Failed to create publish track", err, "meetingID", s.meetingDetails.ID.String())
		return
	}
	defer func() {
		publishTrack.ClearQueue()
		publishTrack.Close()
		close(audioWriterChan)
	}()

	if _, err = s.room.LocalParticipant.PublishTrack(publishTrack, &lksdk.TrackPublicationOptions{
		Name: s.meetingDetails.AgentName,
	}); err != nil {
		logger.Errorw("Failed to publish track", err, "meetingID", s.meetingDetails.ID.String())
		return
	}

	for {
		select {
		case sample, ok := <-audioWriterChan:
			if !ok {
				return
			}
			if err := publishTrack.WriteSample(sample); err != nil {
				logger.Errorw("Failed to write sample", err, "meetingID", s.meetingDetails.ID.String())
			}
		case <-s.ctx.Done():
			return
		}
	}
}

func (s *LiveKitSession) handleTextStreamQueue() {
	for {
		select {
		case data, ok := <-s.textStreamQueue:
			if !ok {
				// Channel closed, exit worker
				return
			}
			marshalData, err := json.Marshal(data)
			if err != nil {
				logger.Errorw("Failed to marshal text stream data", err, "meetingID", s.meetingDetails.ID.String())
				continue
			}
			s.room.LocalParticipant.SendText(string(marshalData), lksdk.StreamTextOptions{
				Topic: "room",
			})
		case <-s.ctx.Done():
			return
		}
	}
}

func (s *LiveKitSession) handleSubscribe(track *webrtc.TrackRemote) (*lkmedia.PCMRemoteTrack, error) {
	if track.Codec().MimeType != webrtc.MimeTypeOpus {
		logger.Warnw("Received non-opus track", nil, "track", track.Codec().MimeType)
	}

	writer := NewRemoteTrackWriter(s.handler)
	trackWriter, err := lkmedia.NewPCMRemoteTrack(track, writer, lkmedia.WithTargetSampleRate(16000))
	if err != nil {
		logger.Errorw("Failed to create remote track", err, "meetingID", s.meetingDetails.ID.String())
		return nil, err
	}

	return trackWriter, nil
}

func (s *LiveKitSession) startRecording() (*livekit.EgressInfo, error) {
	req := &livekit.RoomCompositeEgressRequest{
		RoomName:  s.meetingDetails.ID.String(),
		Layout:    "grid",
		AudioOnly: false,
	}
	outputPath := fmt.Sprintf("%s/%s/recording.mp4", s.userDetails.ID, s.meetingDetails.ID.String())
	req.FileOutputs = []*livekit.EncodedFileOutput{
		{
			Filepath: outputPath,
			Output: &livekit.EncodedFileOutput_S3{
				S3: &livekit.S3Upload{
					AccessKey:      s.awsConfig.AccessKey,
					Secret:         s.awsConfig.SecretKey,
					Region:         s.awsConfig.Region,
					Bucket:         s.awsConfig.Bucket,
					ForcePathStyle: false,
				},
			},
		},
	}

	egressClient := lksdk.NewEgressClient(
		s.lkConfig.Host,
		s.lkConfig.APIKey,
		s.lkConfig.APISecret,
	)
	res, err := egressClient.StartRoomCompositeEgress(context.Background(), req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *LiveKitSession) stopRecording(egressID string) error {
	egressClient := lksdk.NewEgressClient(
		s.lkConfig.Host,
		s.lkConfig.APIKey,
		s.lkConfig.APISecret,
	)

	_, err := egressClient.StopEgress(context.Background(), &livekit.StopEgressRequest{
		EgressId: egressID,
	})
	if err != nil {
		return err
	}

	s.recordingURL = fmt.Sprintf("s3://%s/%s/%s/recording.mp4", s.awsConfig.Bucket, s.userDetails.ID, s.meetingDetails.ID.String())

	if s.handler != nil {
		transcriptData := s.handler.GetTranscript()
		jsonBytes, err := json.MarshalIndent(transcriptData, "", "  ")
		if err == nil {
			s3Key := fmt.Sprintf("%s/%s/transcript.json", s.userDetails.ID, s.meetingDetails.ID.String())
			awsCfg := aws.Config{
				Region:      s.awsConfig.Region,
				Credentials: credentials.NewStaticCredentialsProvider(s.awsConfig.AccessKey, s.awsConfig.SecretKey, ""),
			}
			s3Client := s3.NewFromConfig(awsCfg)

			_, err := s3Client.PutObject(context.Background(), &s3.PutObjectInput{
				Bucket: &s.awsConfig.Bucket,
				Key:    &s3Key,
				Body:   bytes.NewReader(jsonBytes),
			})
			if err == nil {
				s.transcriptURL = fmt.Sprintf("s3://%s/%s", s.awsConfig.Bucket, s3Key)
			}
		}
	}

	return nil
}
