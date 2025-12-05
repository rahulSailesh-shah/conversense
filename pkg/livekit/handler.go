package livekit

import (
	"context"
	"encoding/binary"
	"fmt"
	"strings"
	"time"

	"github.com/livekit/media-sdk"
	"github.com/rahulSailesh-shah/converSense/internal/db/repo"
	"github.com/rahulSailesh-shah/converSense/pkg/config"
	sentimentanalyzer "github.com/rahulSailesh-shah/converSense/pkg/sentiment-analyzer"
	"google.golang.org/genai"
)

type GeminiRealtimeAPIHandler struct {
	client             *genai.Client
	session            *genai.Session
	ctx                context.Context
	cancel             context.CancelFunc
	cb                 *GeminiRealtimeAPIHandlerCallbacks
	sentimentAnalyzer  sentimentanalyzer.SentimentAnalyzer
	transcript         *SessionTranscript
	userDetails        *repo.User
	meetingDetails     *repo.GetMeetingRow
	currentUserContent string // Accumulate user chunks
	currentBotContent  string // Accumulate bot chunks
	currentTurnStart   time.Time
}

type TranscriptDataStream struct {
	Role      string    `json:"role"`      // "user" or "ai"
	Name      string    `json:"name"`      // Speaker's name
	Content   string    `json:"content"`   // Transcript text
	Timestamp time.Time `json:"timestamp"` // When the segment was captured
}

type GeminiRealtimeAPIHandlerCallbacks struct {
	OnAudioReceived  func(audio media.PCM16Sample)
	OnUserSentiment  func(result *sentimentanalyzer.SentimentResult)
	OnUserTranscript func(result *TranscriptDataStream)
}

type SessionTranscript struct {
	Segments []SessionTranscriptSegment `json:"segments"`
}

type SessionTranscriptSegment struct {
	Role      string    `json:"role"`      // "user" or "ai"
	Name      string    `json:"name"`      // Speaker's name
	Content   string    `json:"content"`   // Transcript text
	Timestamp time.Time `json:"timestamp"` // When the segment was captured
}

func NewGeminiRealtimeAPIHandler(parentCtx context.Context,
	config *config.GeminiConfig,
	userDetails *repo.User,
	meetingDetails *repo.GetMeetingRow,
	cb *GeminiRealtimeAPIHandlerCallbacks,
	sentimentAnalyzer sentimentanalyzer.SentimentAnalyzer,
) (*GeminiRealtimeAPIHandler, error) {
	ctx, cancel := context.WithCancel(parentCtx)

	apiKey := config.APIKey
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: apiKey,
	})
	if err != nil {
		cancel()
		return nil, fmt.Errorf("failed to create Gemini client: %w", err)
	}

	model := config.RealtimeModel
	if model == "" {
		model = "gemini-2.5-flash-native-audio-preview-09-2025"
	}
	thinkingBudget := int32(0)
	systemInstructions := meetingDetails.AgentInstructions
	session, err := client.Live.Connect(ctx, model, &genai.LiveConnectConfig{
		SystemInstruction:        genai.NewContentFromText(systemInstructions, genai.RoleUser),
		ResponseModalities:       []genai.Modality{genai.ModalityAudio},
		InputAudioTranscription:  &genai.AudioTranscriptionConfig{},
		OutputAudioTranscription: &genai.AudioTranscriptionConfig{},
		ThinkingConfig: &genai.ThinkingConfig{
			ThinkingBudget: &thinkingBudget,
		},
		Tools: []*genai.Tool{{
			GoogleSearchRetrieval: &genai.GoogleSearchRetrieval{},
		}},
	})
	if err != nil {
		cancel()
		return nil, fmt.Errorf("failed to connect session: %w", err)
	}

	h := &GeminiRealtimeAPIHandler{
		client:            client,
		session:           session,
		ctx:               ctx,
		cancel:            cancel,
		cb:                cb,
		sentimentAnalyzer: sentimentAnalyzer,
		userDetails:       userDetails,
		meetingDetails:    meetingDetails,
		transcript: &SessionTranscript{
			Segments: make([]SessionTranscriptSegment, 0),
		},
		currentTurnStart: time.Now(),
	}

	go h.readMessages()
	return h, nil
}

func (h *GeminiRealtimeAPIHandler) SendAudioChunk(sample media.PCM16Sample) error {
	bytes := make([]byte, len(sample)*2)
	for i, s := range sample {
		binary.LittleEndian.PutUint16(bytes[i*2:], uint16(s))
	}

	err := h.session.SendRealtimeInput(genai.LiveRealtimeInput{
		Audio: &genai.Blob{
			Data:     bytes,
			MIMEType: "audio/pcm;rate=16000",
		},
	})
	if err != nil {
		return fmt.Errorf("error sending audio: %w", err)
	}

	return nil
}

func (h *GeminiRealtimeAPIHandler) readMessages() {
	for {
		response, err := h.session.Receive()
		if err != nil {
			select {
			case <-h.ctx.Done():
				fmt.Println("[-] Session closed")
				return
			default:
				fmt.Println("[-] Error receiving message:", err)
			}
			return
		}
		h.handleMessage(response)
	}
}

func (h *GeminiRealtimeAPIHandler) handleMessage(response *genai.LiveServerMessage) {
	if response.ServerContent == nil {
		return
	}

	// Accumulate output transcription chunks from the bot
	if response.ServerContent.OutputTranscription != nil {
		text := response.ServerContent.OutputTranscription.Text
		if text != "" {
			h.currentBotContent += " " + text
		}
		h.cb.OnUserTranscript(&TranscriptDataStream{
			Role:      "ai",
			Name:      h.meetingDetails.AgentName,
			Content:   strings.TrimSpace(h.currentBotContent),
			Timestamp: h.currentTurnStart,
		})
	}

	// Accumulate input transcription chunks from the user
	if response.ServerContent.InputTranscription != nil {
		text := response.ServerContent.InputTranscription.Text
		if text != "" {
			h.currentUserContent += " " + text
		}
		h.cb.OnUserTranscript(&TranscriptDataStream{
			Role:      "user",
			Name:      h.userDetails.Name,
			Content:   strings.TrimSpace(h.currentUserContent),
			Timestamp: h.currentTurnStart,
		})
	}

	// Handle audio from the bot
	if response.ServerContent.ModelTurn != nil {
		for _, part := range response.ServerContent.ModelTurn.Parts {
			if part.InlineData != nil && part.InlineData.Data != nil {
				audioBytes := part.InlineData.Data
				audioPCM16 := make(media.PCM16Sample, len(audioBytes)/2)
				for i := 0; i < len(audioBytes); i += 2 {
					audioPCM16[i/2] = int16(binary.LittleEndian.Uint16(audioBytes[i : i+2]))
				}
				h.cb.OnAudioReceived(audioPCM16)
			}
		}
	}

	// On turn completion, create segments from accumulated content
	if response.ServerContent.TurnComplete {
		fmt.Println("✅ Turn complete - ready for next input")
		if h.currentUserContent != "" {
			userSegment := SessionTranscriptSegment{
				Role:      "user",
				Name:      h.userDetails.Name,
				Content:   strings.TrimSpace(h.currentUserContent),
				Timestamp: h.currentTurnStart,
			}
			h.transcript.Segments = append(h.transcript.Segments, userSegment)
		}

		if h.currentBotContent != "" {
			botSegment := SessionTranscriptSegment{
				Role:      "ai",
				Name:      h.meetingDetails.AgentName,
				Content:   strings.TrimSpace(h.currentBotContent),
				Timestamp: time.Now(),
			}
			h.transcript.Segments = append(h.transcript.Segments, botSegment)
		}

		userMessage := h.currentUserContent
		h.currentUserContent = ""
		h.currentBotContent = ""
		h.currentTurnStart = time.Now()
		if userMessage != "" {
			res, err := h.sentimentAnalyzer.Analyze(h.ctx, userMessage, h.userDetails.Name)
			if err != nil {
				fmt.Println("Error analyzing sentiment:", err)
			}
			h.cb.OnUserSentiment(res)

		}

	}
}

func (h *GeminiRealtimeAPIHandler) GetTranscript() *SessionTranscript {
	return h.transcript
}

func (h *GeminiRealtimeAPIHandler) Close() error {
	h.cancel()
	return h.session.Close()
}
