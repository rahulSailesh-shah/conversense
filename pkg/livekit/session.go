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
	"github.com/rahulSailesh-shah/converSense/pkg/config"
)

type SessionCallbacks struct {
	OnMeetingEnd func(meetingID string, recordingURL string, transcriptURL string)
}

type SessionTranscript struct {
	Segments []SessionTranscriptSegment
}

type SessionTranscriptSegment struct {
	User string
	Bot  string
}

type LiveKitSession struct {
	meetingID      string
	agentName      string
	userID         string
	room           *lksdk.Room
	handler        *GeminiRealtimeAPIHandler
	egressInfo     *livekit.EgressInfo
	lkConfig       *config.LiveKitConfig
	geminiConfig   *config.GeminiConfig
	awsConfig      *config.AWSConfig
	ctx            context.Context
	cancel         context.CancelFunc
	done           chan struct{}
	callbacks      SessionCallbacks
	recordingURL   string
	transcriptURL  string
	transcriptData *SessionTranscript
	stopOnce       sync.Once
}

func NewLiveKitSession(
	meetingID string,
	userID string,
	agentName string,
	lkConfig *config.LiveKitConfig,
	geminiConfig *config.GeminiConfig,
	awsConfig *config.AWSConfig,
	callbacks SessionCallbacks,
) *LiveKitSession {
	ctx, cancel := context.WithCancel(context.Background())

	return &LiveKitSession{
		meetingID:    meetingID,
		userID:       userID,
		agentName:    agentName,
		lkConfig:     lkConfig,
		geminiConfig: geminiConfig,
		awsConfig:    awsConfig,
		ctx:          ctx,
		cancel:       cancel,
		done:         make(chan struct{}),
		callbacks:    callbacks,
		stopOnce:     sync.Once{},
		transcriptData: &SessionTranscript{
			Segments: make([]SessionTranscriptSegment, 0),
		},
	}
}

func (s *LiveKitSession) Start() error {
	if err := s.connectBot(); err != nil {
		return fmt.Errorf("failed to connect bot: %w", err)
	}
	fmt.Println("[-] LiveKit session started successfully", "meetingID", s.meetingID)
	return nil
}

func (s *LiveKitSession) Stop() error {
	var stopErr error
	s.stopOnce.Do(func() {
		if s.egressInfo != nil {
			if err := s.stopRecording(s.egressInfo.EgressId); err != nil {
				stopErr = fmt.Errorf("failed to stop recording: %w", err)
			}
		}

		if s.room != nil {
			s.room.Disconnect()
		}
		if s.handler != nil {
			s.handler.Close()
		}
		close(s.done)
		if s.callbacks.OnMeetingEnd != nil {
			fmt.Println("[-] Meeting ended, starting post-processing", "meetingID", s.meetingID)
			s.callbacks.OnMeetingEnd(s.meetingID, s.recordingURL, s.transcriptURL)
		}
	})
	return stopErr
}

func (s *LiveKitSession) GenerateUserToken() (string, error) {
	at := auth.NewAccessToken(s.lkConfig.APIKey, s.lkConfig.APISecret)
	grant := &auth.VideoGrant{
		RoomJoin: true,
		Room:     s.meetingID,
	}
	at.SetVideoGrant(grant).
		SetIdentity(s.userID).
		SetValidFor(time.Hour)
	token, err := at.ToJWT()
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s *LiveKitSession) connectBot() error {
	audioWriterChan := make(chan media.PCM16Sample, 100)
	handler, err := NewGeminiRealtimeAPIHandler(&GeminiRealtimeAPIHandlerCallbacks{
		OnAudioReceived: func(audio media.PCM16Sample) {
			select {
			case audioWriterChan <- audio:
			case <-s.ctx.Done():
				return
			}
		},
	}, s.geminiConfig, s.transcriptData)
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

	egressInfo, err := s.startRecording()
	if err != nil {
		logger.Errorw("Failed to start recording", err, "meetingID", s.meetingID)
	} else {
		s.egressInfo = egressInfo
	}

	fmt.Println("[-] Bot connected successfully", "meetingID", s.meetingID)
	return nil
}

func (s *LiveKitSession) connectToRoom() error {
	room, err := lksdk.ConnectToRoom(s.lkConfig.Host, lksdk.ConnectInfo{
		APIKey:              s.lkConfig.APIKey,
		APISecret:           s.lkConfig.APISecret,
		RoomName:            s.meetingID,
		ParticipantIdentity: s.agentName,
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
			fmt.Println("[-] Participant disconnected", "meetingID", s.meetingID, "participantID", participant.Identity())
			s.Stop()
		},
		OnDisconnected: func() {
			fmt.Println("[-] Bot disconnected from room", "meetingID", s.meetingID)
			if pcmRemoteTrack != nil {
				pcmRemoteTrack.Close()
				pcmRemoteTrack = nil
			}
		},
		OnDisconnectedWithReason: func(reason lksdk.DisconnectionReason) {
			fmt.Println("[-] Bot disconnected with reason", "meetingID", s.meetingID, "reason", reason)
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
		logger.Errorw("Failed to create publish track", err, "meetingID", s.meetingID)
		return
	}
	defer func() {
		publishTrack.ClearQueue()
		publishTrack.Close()
		close(audioWriterChan)
	}()

	if _, err = s.room.LocalParticipant.PublishTrack(publishTrack, &lksdk.TrackPublicationOptions{
		Name: s.agentName,
	}); err != nil {
		logger.Errorw("Failed to publish track", err, "meetingID", s.meetingID)
		return
	}

	for {
		select {
		case sample, ok := <-audioWriterChan:
			if !ok {
				return
			}
			if err := publishTrack.WriteSample(sample); err != nil {
				logger.Errorw("Failed to write sample", err, "meetingID", s.meetingID)
			}
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
		logger.Errorw("Failed to create remote track", err, "meetingID", s.meetingID)
		return nil, err
	}

	return trackWriter, nil
}

func (s *LiveKitSession) startRecording() (*livekit.EgressInfo, error) {
	fmt.Println("[-] Starting recording", "meetingID", s.meetingID)
	req := &livekit.RoomCompositeEgressRequest{
		RoomName:  s.meetingID,
		Layout:    "speaker",
		AudioOnly: false,
	}
	outputPath := fmt.Sprintf("%s/%s/recording.mp4", s.userID, s.meetingID)
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

	fmt.Println("[-] Recording started", "meetingID", s.meetingID, "egressID", res.EgressId)
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

	fmt.Println("[-] Recording stopped", "meetingID", s.meetingID, "egressID", egressID)

	// TODO: Get the recording URL from egress info and set s.recordingURL
	// This would require querying the egress status to get the final S3 URL
	s.recordingURL = fmt.Sprintf("s3://%s/recordings/%s.mp4", s.awsConfig.Bucket, s.meetingID)

	if s.handler != nil && s.transcriptData != nil {
		jsonBytes, err := json.MarshalIndent(s.handler.transcript, "", "  ")
		if err != nil {
			fmt.Println("Error marshaling transcript:", err)
		} else {
			s3Key := fmt.Sprintf("%s/%s/transcript.json", s.userID, s.meetingID)

			// Build AWS config dynamically using your access keys
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
			if err != nil {
				fmt.Println("Failed to upload transcript:", err)
			} else {
				s.transcriptURL = fmt.Sprintf("s3://%s/%s", s.awsConfig.Bucket, s3Key)
				fmt.Println("[-] Transcript uploaded to S3", "url", s.transcriptURL)
			}
		}
	}

	return nil
}
