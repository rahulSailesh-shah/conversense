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
	"github.com/rahulSailesh-shah/converSense/internal/db/repo"
	"github.com/rahulSailesh-shah/converSense/pkg/config"
	"github.com/rahulSailesh-shah/converSense/pkg/livekit"
	"google.golang.org/genai"
)

type GeminiService interface {
	Chat(ctx context.Context, meetingID uuid.UUID, userID string, message string) (<-chan string, error)
	GetChatHistory(ctx context.Context, meetingID uuid.UUID, userID string) ([]repo.MeetingChatMessages, error)
}

type geminiService struct {
	queries         *repo.Queries
	config          *config.GeminiConfig
	awsConfig       *config.AWSConfig
	transcriptCache sync.Map // Cache transcripts by meetingID
}

func NewGeminiService(queries *repo.Queries, config *config.GeminiConfig, awsConfig *config.AWSConfig) GeminiService {
	return &geminiService{
		queries:         queries,
		config:          config,
		awsConfig:       awsConfig,
		transcriptCache: sync.Map{},
	}
}

func (s *geminiService) Chat(ctx context.Context, meetingID uuid.UUID, userID string, message string) (<-chan string, error) {
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

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: s.config.APIKey,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create gemini client: %w", err)
	}

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

	systemPrompt += "\n\nChat History:"
	for _, msg := range history {
		role := "User"
		if msg.Role == "ai" {
			role = "AI"
		}
		systemPrompt += fmt.Sprintf("\n%s: %s", role, msg.Content)
	}

	stream := make(chan string)

	go func() {
		defer close(stream)

		geminiStream := client.Models.GenerateContentStream(
			ctx,
			"gemini-2.0-flash",
			genai.Text(systemPrompt+"\n\nUser: "+message),
			nil,
		)

		var fullResponse strings.Builder

		for chunk, err := range geminiStream {
			if err != nil {
				fmt.Printf("Error in gemini stream: %v\n", err)
				return
			}
			if len(chunk.Candidates) > 0 && len(chunk.Candidates[0].Content.Parts) > 0 {
				part := chunk.Candidates[0].Content.Parts[0]
				text := part.Text
				fullResponse.WriteString(text)
				select {
				case <-ctx.Done():
					return
				case stream <- text:
				}
			}
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

func (s *geminiService) GetChatHistory(ctx context.Context, meetingID uuid.UUID, userID string) ([]repo.MeetingChatMessages, error) {
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

func (s *geminiService) fetchTranscript(ctx context.Context, s3URL string) (*livekit.SessionTranscript, error) {
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
