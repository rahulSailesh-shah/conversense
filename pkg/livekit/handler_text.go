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

// GeminiRealtimeTextHandler is a text-output variant of the realtime handler.
// It streams audio input to Gemini Live and consumes text responses instead of audio.
// A small JSON context payload is sent upfront (hardcoded for now; replace with DB-derived state later).
type GeminiRealtimeTextHandler struct {
	client             *genai.Client
	session            *genai.Session
	ctx                context.Context
	cancel             context.CancelFunc
	cb                 *GeminiRealtimeAPIHandlerCallbacks
	sentimentAnalyzer  sentimentanalyzer.SentimentAnalyzer
	transcript         *SessionTranscript
	userDetails        *repo.User
	meetingDetails     *repo.GetMeetingRow
	currentUserContent string
	currentBotContent  string
	currentTurnStart   time.Time
	contextJSON        string
}

// NewGeminiRealtimeTextHandler connects a Live session configured for text output.
// It also sends an initial text message containing the current JSON context.
func NewGeminiRealtimeTextHandler(
	parentCtx context.Context,
	cfg *config.GeminiConfig,
	userDetails *repo.User,
	meetingDetails *repo.GetMeetingRow,
	cb *GeminiRealtimeAPIHandlerCallbacks,
	sentimentAnalyzer sentimentanalyzer.SentimentAnalyzer,
) (*GeminiRealtimeTextHandler, error) {
	ctx, cancel := context.WithCancel(parentCtx)

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: cfg.APIKey,
	})
	if err != nil {
		cancel()
		return nil, fmt.Errorf("failed to create Gemini client: %w", err)
	}

	model := cfg.RealtimeModel
	if model == "" {
		model = "gemini-2.5-flash-native-audio-preview-09-2025"
	}
	thinkingBudget := int32(0)

	// Hardcoded JSON context for now; replace with DB-sourced state when ready.
	contextJSON := fmt.Sprintf(`{"meeting_id":"%s","user":"%s","agent":"%s"}`,
		meetingDetails.ID.String(), userDetails.Name, meetingDetails.AgentName)

	session, err := client.Live.Connect(ctx, model, &genai.LiveConnectConfig{
		SystemInstruction:       genai.NewContentFromText(meetingDetails.AgentInstructions, genai.RoleUser),
		ResponseModalities:      []genai.Modality{genai.ModalityText}, // request text output
		InputAudioTranscription: &genai.AudioTranscriptionConfig{
			// empty config enables transcription
		},
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

	h := &GeminiRealtimeTextHandler{
		client:            client,
		session:           session,
		ctx:               ctx,
		cancel:            cancel,
		cb:                cb,
		sentimentAnalyzer: sentimentAnalyzer,
		userDetails:       userDetails,
		meetingDetails:    meetingDetails,
		contextJSON:       contextJSON,
		transcript: &SessionTranscript{
			Segments: make([]SessionTranscriptSegment, 0),
		},
		currentTurnStart: time.Now(),
	}

	// Send initial JSON context as text input.
	if err := h.SendText(contextJSON); err != nil {
		fmt.Printf("[-] Failed to send initial context JSON: %v\n", err)
	}

	go h.readMessages()
	return h, nil
}

// SendAudioChunk streams PCM16 audio to the model.
func (h *GeminiRealtimeTextHandler) SendAudioChunk(sample media.PCM16Sample) error {
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

// SendText sends a text chunk (used here for JSON context or any extra user text).
func (h *GeminiRealtimeTextHandler) SendText(text string) error {
	return h.session.SendRealtimeInput(genai.LiveRealtimeInput{
		Text: text,
	})
}

func (h *GeminiRealtimeTextHandler) readMessages() {
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

func (h *GeminiRealtimeTextHandler) handleMessage(response *genai.LiveServerMessage) {
	if response.ServerContent == nil {
		return
	}

	// Capture model text output from ModelTurn parts.
	if response.ServerContent.ModelTurn != nil {
		for _, part := range response.ServerContent.ModelTurn.Parts {
			if part == nil {
				continue
			}
			if part.Text != "" {
				h.currentBotContent += " " + part.Text
			}
		}
		if h.currentBotContent != "" {
			h.cb.OnUserTranscript(&TranscriptDataStream{
				Role:      "ai",
				Name:      h.meetingDetails.AgentName,
				Content:   strings.TrimSpace(h.currentBotContent),
				Timestamp: h.currentTurnStart,
			})
		}
	}

	// Accumulate input transcription chunks from the user.
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

	// On turn completion, persist accumulated content and run sentiment.
	if response.ServerContent.TurnComplete {
		fmt.Println("✅ Turn complete - ready for next input (text output mode)")
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
			} else {
				h.cb.OnUserSentiment(res)
			}
		}
	}
}

func (h *GeminiRealtimeTextHandler) GetTranscript() *SessionTranscript {
	return h.transcript
}

func (h *GeminiRealtimeTextHandler) Close() error {
	h.cancel()
	return h.session.Close()
}
