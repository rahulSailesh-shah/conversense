package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/rahulSailesh-shah/converSense/internal/db/repo"
	"github.com/rahulSailesh-shah/converSense/pkg/config"
	"github.com/rahulSailesh-shah/converSense/pkg/livekit"
)

type ChatService interface {
	Chat(ctx context.Context, meetingID uuid.UUID, userID string, message string) (<-chan string, error)
	GetChatHistory(ctx context.Context, meetingID uuid.UUID, userID string) ([]repo.MeetingChatMessages, error)
}

type chatService struct {
	queries         *repo.Queries
	config          *config.OpenAIConfig
	awsConfig       *config.AWSConfig
	transcriptCache sync.Map // Cache transcripts by meetingID
}

func NewChatService(queries *repo.Queries, config *config.OpenAIConfig, awsConfig *config.AWSConfig) ChatService {
	return &chatService{
		queries:         queries,
		config:          config,
		awsConfig:       awsConfig,
		transcriptCache: sync.Map{},
	}
}

func (s *chatService) Chat(ctx context.Context, meetingID uuid.UUID, userID string, message string) (<-chan string, error) {
	meeting, err := s.queries.GetMeeting(ctx, repo.GetMeetingParams{
		ID:     meetingID,
		UserID: userID,
	})
	if err != nil {
		return nil, fmt.Errorf("unauthorized: user does not have access to this meeting or meeting does not exist")
	}

	_, err = s.queries.CreateChatMessage(ctx, repo.CreateChatMessageParams{
		MeetingID: meetingID,
		UserID:    userID,
		Role:      "user",
		Content:   message,
	})
	if err != nil {
		fmt.Printf("Failed to save user message: %v\n", err)
	}

	const maxHistoryMessages = 20
	historyRows, err := s.queries.GetRecentChatMessages(ctx, repo.GetRecentChatMessagesParams{
		MeetingID: meetingID,
		Limit:     maxHistoryMessages,
	})
	if err != nil {
		fmt.Printf("Failed to fetch chat history: %v\n", err)
	}

	history := make([]repo.MeetingChatMessages, len(historyRows))
	for i, msg := range historyRows {
		history[len(historyRows)-1-i] = msg
	}

	var transcriptContext string
	if meeting.TranscriptUrl != nil && meeting.Status == "completed" {
		if cached, ok := s.transcriptCache.Load(meetingID.String()); ok {
			transcriptContext = cached.(string)
		} else {
			transcript, err := s.fetchTranscript(ctx, *meeting.TranscriptUrl)
			if err != nil {
				fmt.Printf("Failed to fetch transcript: %v\n", err)
			} else {
				transcriptContext = formatTranscript(transcript)
				s.transcriptCache.Store(meetingID.String(), transcriptContext)
			}
		}
	}

	client := openai.NewClient(
		option.WithAPIKey(s.config.APIKey),
		option.WithBaseURL(s.config.BaseURL),
	)

	systemPrompt := fmt.Sprintf(`
      You are an AI assistant helping the user revisit a recently completed meeting.
      Below is a summary of the meeting, generated from the transcript:

      %s

      The following are your original instructions from the live meeting assistant. Please continue to follow these behavioral guidelines as you assist the user:

      %s

      The user may ask questions about the meeting, request clarifications, or ask for follow-up actions.
      Always base your responses on the meeting summary above.

      You also have access to the recent conversation history between you and the user. Use the context of previous messages to provide relevant, coherent, and helpful responses. If the user's question refers to something discussed earlier, make sure to take that into account and maintain continuity in the conversation.

      If the summary does not contain enough information to answer a question, politely let the user know.

      Be concise, helpful, and focus on providing accurate information from the meeting and the ongoing conversation.
      `, transcriptContext, meeting.AgentInstructions)

	messages := []openai.ChatCompletionMessageParamUnion{
		openai.SystemMessage(systemPrompt),
	}

	for _, msg := range history {
		if msg.Role == "ai" {
			messages = append(messages, openai.AssistantMessage(msg.Content))
		} else {
			messages = append(messages, openai.UserMessage(msg.Content))
		}
	}

	messages = append(messages, openai.UserMessage(message))

	stream := make(chan string)

	go func() {
		defer close(stream)

		openaiStream := client.Chat.Completions.NewStreaming(ctx, openai.ChatCompletionNewParams{
			Model:    "z-ai/glm4.7",
			Messages: messages,
		})

		var fullResponse strings.Builder

		for openaiStream.Next() {
			chunk := openaiStream.Current()
			if len(chunk.Choices) > 0 {
				text := chunk.Choices[0].Delta.Content
				if text == "" {
					continue
				}
				fullResponse.WriteString(text)
				select {
				case <-ctx.Done():
					return
				case stream <- text:
				}
			}
		}

		if err := openaiStream.Err(); err != nil {
			fmt.Printf("Error in openai stream: %v\n", err)
			return
		}

		_, err = s.queries.CreateChatMessage(ctx, repo.CreateChatMessageParams{
			MeetingID: meetingID,
			UserID:    "ai",
			Role:      "ai",
			Content:   fullResponse.String(),
		})
		if err != nil {
			fmt.Printf("Failed to save AI message: %v\n", err)
		}
	}()

	return stream, nil
}

func (s *chatService) GetChatHistory(ctx context.Context, meetingID uuid.UUID, userID string) ([]repo.MeetingChatMessages, error) {
	// Authorization check
	_, err := s.queries.GetMeeting(ctx, repo.GetMeetingParams{
		ID:     meetingID,
		UserID: userID,
	})
	if err != nil {
		return nil, fmt.Errorf("unauthorized")
	}

	return s.queries.GetChatMessages(ctx, meetingID)
}

func (s *chatService) fetchTranscript(ctx context.Context, s3URL string) (*livekit.SessionTranscript, error) {
	bucket, key, err := parseS3URL(s3URL)
	if err != nil {
		return nil, err
	}

	cfg, err := awsconfig.LoadDefaultConfig(ctx,
		awsconfig.WithRegion(s.awsConfig.Region),
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(s.awsConfig.AccessKey, s.awsConfig.SecretKey, "")),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load aws config: %w", err)
	}

	client := s3.NewFromConfig(cfg)

	result, err := client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get object from s3: %w", err)
	}
	defer result.Body.Close()

	var transcript livekit.SessionTranscript
	if err := json.NewDecoder(result.Body).Decode(&transcript); err != nil {
		return nil, fmt.Errorf("failed to decode transcript json: %w", err)
	}

	return &transcript, nil
}

func formatTranscript(transcript *livekit.SessionTranscript) string {
	var sb strings.Builder
	for _, segment := range transcript.Segments {
		sb.WriteString(fmt.Sprintf("[%s] %s: %s\n", segment.Timestamp.Format("15:04:05"), segment.Name, segment.Content))
	}
	return sb.String()
}
