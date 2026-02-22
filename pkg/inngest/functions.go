package inngest

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	"github.com/inngest/inngestgo"
	"github.com/inngest/inngestgo/step"
	"github.com/rahulSailesh-shah/converSense/internal/db/repo"
	"google.golang.org/genai"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

func (i *Inngest) RegisterFunctions() error {
	err := i.postProcessMeeting()
	return err
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

func (i *Inngest) PostProcessMeeting(ctx context.Context, meetingId string, userId string) error {
	fmt.Println("[--] Meeting post-processing event sent", "meetingID", meetingId)
	_, err := i.client.Send(ctx, inngestgo.Event{
		Name: "conversense/post-process-meeting",
		Data: map[string]any{
			"meetingId": meetingId,
			"userId":    userId,
		},
	})
	return err
}

func (i *Inngest) postProcessMeeting() error {
	type PostProcessEventData struct {
		MeetingId string `json:"meetingId"`
		UserID    string `json:"userId"`
	}

	_, err := inngestgo.CreateFunction(
		i.client,
		inngestgo.FunctionOpts{
			ID:   "post-process-meeting",
			Name: "Post Process Meeting",
		},
		inngestgo.EventTrigger("conversense/post-process-meeting", nil),
		func(ctx context.Context, input inngestgo.Input[PostProcessEventData]) (any, error) {
			fmt.Println("[---] Meeting post-processing started", "meetingID", input.Event.Data.MeetingId)
			// Fetch meeting details
			meetingId, err := uuid.Parse(input.Event.Data.MeetingId)
			if err != nil {
				return nil, err
			}
			userId := input.Event.Data.UserID
			meetingDetails, err := step.Run(ctx, "fetch-data", func(ctx context.Context) (*repo.GetMeetingRow, error) {
				meetingDetails, err := i.queries.GetMeeting(ctx, repo.GetMeetingParams{
					ID:     meetingId,
					UserID: userId,
				})
				return &meetingDetails, err
			})
			if err != nil {
				return nil, err
			}

			// Fetch transcript
			transcriptURL := meetingDetails.TranscriptUrl
			if transcriptURL == nil {
				return "", fmt.Errorf("no transcript URL found for meeting")
			}
			transcriptData, err := step.Run(ctx, "fetch-transcript",
				func(ctx context.Context) (*SessionTranscript, error) {
					transcript, err := i.fetchTranscriptFromS3(ctx, *transcriptURL)
					return transcript, err
				})
			if err != nil {
				return nil, err
			}
			fmt.Println("[---] Transcript fetched successfully", "meetingID", meetingId)

			// Generate summary
			summary, err := step.Run(ctx, "generate-summary", func(ctx context.Context) (string, error) {
				summary, err := i.processTranscriptWithOpenAI(ctx, transcriptData)
				return summary, err
			})
			if err != nil {
				return nil, err
			}
			fmt.Println("[---] Summary generated successfully", "meetingID", meetingId)
			fmt.Println(summary)

			_, err = step.Run(ctx, "save-summary", func(ctx context.Context) (any, error) {
				_, err := i.queries.UpdateMeeting(ctx, repo.UpdateMeetingParams{
					ID:            meetingDetails.ID,
					UserID:        meetingDetails.UserID,
					Name:          meetingDetails.Name,
					AgentID:       meetingDetails.AgentID,
					Status:        "completed",
					Summary:       &summary,
					StartTime:     meetingDetails.StartTime,
					EndTime:       meetingDetails.EndTime,
					TranscriptUrl: meetingDetails.TranscriptUrl,
					RecordingUrl:  meetingDetails.RecordingUrl,
				})
				return nil, err
			})
			if err != nil {
				return nil, err
			}
			fmt.Println("[---] Summary saved to database successfully", "meetingID", meetingId)

			return summary, nil
		},
	)
	return err
}

func parseS3URL(s3URL string) (bucket, key string, err error) {
	if !strings.HasPrefix(s3URL, "s3://") {
		return "", "", fmt.Errorf("invalid S3 URL format")
	}
	path := strings.TrimPrefix(s3URL, "s3://")
	parts := strings.SplitN(path, "/", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid S3 URL format")
	}
	return parts[0], parts[1], nil
}

func (i *Inngest) fetchTranscriptFromS3(ctx context.Context, transcriptURL string) (*SessionTranscript, error) {
	awsCfg := aws.Config{
		Region:      i.awsConfig.Region,
		Credentials: credentials.NewStaticCredentialsProvider(i.awsConfig.AccessKey, i.awsConfig.SecretKey, ""),
	}

	s3Client := s3.NewFromConfig(awsCfg)
	bucket, key, err := parseS3URL(transcriptURL)
	if err != nil {
		return nil, err
	}

	result, err := s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get object from S3: %w", err)
	}
	defer result.Body.Close()

	body, err := io.ReadAll(result.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read S3 object body: %w", err)
	}

	var transcript SessionTranscript
	if err := json.Unmarshal(body, &transcript); err != nil {
		return nil, fmt.Errorf("failed to unmarshal transcript: %w", err)
	}

	return &transcript, nil
}

func (i *Inngest) processTranscriptWithGemini(ctx context.Context, transcript *SessionTranscript,
) (string, error) {
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  i.geminiConfig.APIKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return "", fmt.Errorf("failed to create Gemini client: %w", err)
	}

	var fullText strings.Builder
	for _, segment := range transcript.Segments {
		fullText.WriteString(fmt.Sprintf("[%s]: %s\n", segment.Name, segment.Content))
	}

	// Create prompt
	prompt := fmt.Sprintf(`
        You are an expert summarizer. You write readable, concise, simple content. You are given a transcript of a meeting and you need to summarize it.

        Use the following markdown structure for every output:

        ### Overview
        Provide a detailed, engaging summary of the session's content. Focus on major features, user workflows, and any key takeaways. Write in a narrative style, using full sentences. Highlight unique or powerful aspects of the product, platform, or discussion.

        ### Notes
        Break down key content into thematic sections with timestamp ranges. Each section should summarize key points, actions, or demos in bullet format.

        Example:
        #### Section Name
        - Main point or demo shown here
        - Another key insight or interaction
        - Follow-up tool or explanation provided

        #### Next Section
        - Feature X automatically does Y
        - Mention of integration with Z

        Transcript:\n
        %s`, fullText.String())

	model := "gemini-2.0-flash-lite"
	response, err := client.Models.GenerateContent(
		ctx,
		model,
		genai.Text(prompt),
		nil,
	)
	if err != nil {
		return "", fmt.Errorf("failed to generate content: %w", err)
	}

	return response.Text(), nil
}

func (i *Inngest) processTranscriptWithOpenAI(ctx context.Context, transcript *SessionTranscript) (string, error) {
	client := openai.NewClient(
		option.WithAPIKey(i.openaiConfig.APIKey),
		option.WithBaseURL(i.openaiConfig.BaseURL),
	)

	var fullText strings.Builder
	for _, segment := range transcript.Segments {
		fullText.WriteString(fmt.Sprintf("[%s]: %s\n", segment.Name, segment.Content))
	}

	prompt := fmt.Sprintf(`
        You are an expert summarizer. You write readable, concise, simple content. You are given a transcript of a meeting and you need to summarize it.

        Use the following markdown structure for every output:

        ### Overview
        Provide a detailed, engaging summary of the session's content. Focus on major features, user workflows, and any key takeaways. Write in a narrative style, using full sentences. Highlight unique or powerful aspects of the product, platform, or discussion.

        ### Notes
        Break down key content into thematic sections with timestamp ranges. Each section should summarize key points, actions, or demos in bullet format.

        Example:
        #### Section Name
        - Main point or demo shown here
        - Another key insight or interaction
        - Follow-up tool or explanation provided

        #### Next Section
        - Feature X automatically does Y
        - Mention of integration with Z

        Transcript:\n
        %s`, fullText.String())

	response, err := client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model: "z-ai/glm4.7",
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(prompt),
		},
	})
	if err != nil {
		return "", fmt.Errorf("failed to generate content: %w", err)
	}

	return response.Choices[0].Message.Content, nil
}
