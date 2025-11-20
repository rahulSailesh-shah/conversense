package livekit

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"

	"github.com/livekit/media-sdk"
	"github.com/rahulSailesh-shah/converSense/pkg/config"
	"google.golang.org/genai"
)

type GeminiRealtimeAPIHandler struct {
	client  *genai.Client
	session *genai.Session
	ctx     context.Context
	cancel  context.CancelFunc
	cb      *GeminiRealtimeAPIHandlerCallbacks
}

type GeminiRealtimeAPIHandlerCallbacks struct {
	OnAudioReceived func(audio media.PCM16Sample)
}

func NewGeminiRealtimeAPIHandler(cb *GeminiRealtimeAPIHandlerCallbacks, config *config.GeminiConfig) (*GeminiRealtimeAPIHandler, error) {
	ctx, cancel := context.WithCancel(context.Background())

	apiKey := config.APIKey
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: apiKey,
	})
	if err != nil {
		cancel()
		return nil, fmt.Errorf("failed to create Gemini client: %w", err)
	}

	model := config.Model
	if model == "" {
		model = "gemini-2.5-flash-native-audio-preview-09-2025"
	}

	thinkingBudget := int32(0)
	session, err := client.Live.Connect(ctx, model, &genai.LiveConnectConfig{
		SystemInstruction:        genai.NewContentFromText("You are a helpful assistant and answer in a friendly tone.", genai.RoleUser),
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
		client:  client,
		session: session,
		ctx:     ctx,
		cancel:  cancel,
		cb:      cb,
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
		select {
		case <-h.ctx.Done():
			return
		default:
		}

		response, err := h.session.Receive()
		if err != nil {
			fmt.Println("Error receiving message:", err)
			return
		}

		h.handleMessage(response)
	}
}

func (h *GeminiRealtimeAPIHandler) handleMessage(response *genai.LiveServerMessage) {
	if response.ServerContent == nil {
		return
	}

	_, err := json.MarshalIndent(response.ServerContent, "", "  ")
	if err != nil {
		fmt.Println("Error marshalling model response:", err)
		return
	}

	// fmt.Println("Model Response:", string(modelResponse))

	if response.ServerContent.ModelTurn != nil {
		for _, part := range response.ServerContent.ModelTurn.Parts {
			if part.InlineData != nil && part.InlineData.Data != nil {
				audioBytes := part.InlineData.Data
				// Convert bytes to PCM16
				audioPCM16 := make(media.PCM16Sample, len(audioBytes)/2)
				for i := 0; i < len(audioBytes); i += 2 {
					audioPCM16[i/2] = int16(binary.LittleEndian.Uint16(audioBytes[i : i+2]))
				}

				h.cb.OnAudioReceived(audioPCM16)
			}
		}
	}

	if response.ServerContent.TurnComplete {
		fmt.Println("✅ Turn complete - ready for next input")
	}
}

func (h *GeminiRealtimeAPIHandler) Close() error {
	h.cancel()
	return h.session.Close()
}
