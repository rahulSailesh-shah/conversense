package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rahulSailesh-shah/converSense/internal/db/repo"
	"github.com/rahulSailesh-shah/converSense/internal/dto"
	"github.com/rahulSailesh-shah/converSense/pkg/config"
	"github.com/rahulSailesh-shah/converSense/pkg/livekit"
)

type MeetingService interface {
	CreateMeeting(ctx context.Context, request dto.CreateMeetingRequest) (*dto.MeetingResponse, error)
	UpdateMeeting(ctx context.Context, request dto.UpdateMeetingRequest) (*dto.MeetingResponse, error)
	GetMeetings(ctx context.Context, request dto.GetMeetingsRequest) (*dto.PaginatedMeetingsResponse, error)
	GetMeeting(ctx context.Context, request dto.GetMeetingRequest) (*dto.MeetingResponse, error)
	DeleteMeeting(ctx context.Context, request dto.DeleteMeetingRequest) error
	StartMeeting(ctx context.Context, request dto.StartMeetingRequest) (string, error)
}

type meetingService struct {
	queries *repo.Queries
	db      *pgxpool.Pool

	// LiveKit configuration
	lkConfig     *config.LiveKitConfig
	geminiConfig *config.GeminiConfig
	awsConfig    *config.AWSConfig
}

func NewMeetingService(
	db *pgxpool.Pool,
	queries *repo.Queries,
	lkConfig *config.LiveKitConfig,
	geminiConfig *config.GeminiConfig,
	awsConfig *config.AWSConfig,
) MeetingService {
	return &meetingService{
		db:           db,
		queries:      queries,
		lkConfig:     lkConfig,
		geminiConfig: geminiConfig,
		awsConfig:    awsConfig,
	}
}

func (s *meetingService) CreateMeeting(ctx context.Context, request dto.CreateMeetingRequest) (*dto.MeetingResponse, error) {
	newMeeting, err := s.queries.CreateMeeting(ctx, repo.CreateMeetingParams{
		Name:    request.Name,
		UserID:  request.UserID,
		AgentID: request.AgentID,
	})
	if err != nil {
		return nil, err
	}
	return toMeetingResponse(newMeeting), nil
}

func (s *meetingService) UpdateMeeting(ctx context.Context, request dto.UpdateMeetingRequest) (*dto.MeetingResponse, error) {
	currentMeeting, err := s.queries.GetMeeting(ctx, repo.GetMeetingParams{
		ID:     request.ID,
		UserID: request.UserID,
	})
	if err != nil {
		return nil, fmt.Errorf("-- failed to get meeting --: %w", err)
	}

	if request.Name != "" {
		currentMeeting.Name = request.Name
	}

	if request.AgentID != uuid.Nil {
		currentMeeting.AgentID = request.AgentID
	}

	if request.Status != "" {
		currentMeeting.Status = request.Status
	}

	if request.StartTime != nil {
		currentMeeting.StartTime = request.StartTime
	}

	if request.EndTime != nil {
		currentMeeting.EndTime = request.EndTime
	}

	if request.TranscriptURL != nil {
		currentMeeting.TranscriptUrl = request.TranscriptURL
	}

	if request.RecordingURL != nil {
		currentMeeting.RecordingUrl = request.RecordingURL
	}

	if request.Summary != nil {
		currentMeeting.Summary = request.Summary
	}

	data, err := json.MarshalIndent(currentMeeting, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("-- failed to marshal meeting --: %w", err)
	}
	fmt.Println(string(data))

	updatedMeeting, err := s.queries.UpdateMeeting(ctx, repo.UpdateMeetingParams{
		ID:            currentMeeting.ID,
		UserID:        currentMeeting.UserID,
		Name:          currentMeeting.Name,
		AgentID:       currentMeeting.AgentID,
		Status:        currentMeeting.Status,
		StartTime:     currentMeeting.StartTime,
		EndTime:       currentMeeting.EndTime,
		TranscriptUrl: currentMeeting.TranscriptUrl,
		RecordingUrl:  currentMeeting.RecordingUrl,
		Summary:       currentMeeting.Summary,
	})
	if err != nil {
		return nil, fmt.Errorf("-- failed to update meeting --: %w", err)
	}

	return toMeetingResponse(updatedMeeting), nil
}

func (s *meetingService) GetMeetings(ctx context.Context, request dto.GetMeetingsRequest) (*dto.PaginatedMeetingsResponse, error) {
	rows, err := s.queries.GetMeetings(ctx, repo.GetMeetingsParams{
		UserID:  request.UserID,
		Column2: request.Search,
		Limit:   request.Limit,
		Offset:  request.Offset,
	})
	if err != nil {
		return nil, err
	}

	var totalCount int32
	if len(rows) > 0 {
		totalCount = int32(rows[0].TotalCount)
	}
	meetings := make([]dto.MeetingResponse, 0, len(rows))
	for _, row := range rows {
		meetings = append(meetings, dto.MeetingResponse{
			ID:        row.ID,
			UserID:    row.UserID,
			Name:      row.Name,
			AgentID:   row.AgentID,
			Status:    row.Status,
			CreatedAt: row.CreatedAt,
			UpdatedAt: row.UpdatedAt,
			AgentDetails: &dto.AgentDetails{
				Name:         row.AgentName,
				Instructions: row.AgentInstructions,
			},
		})
	}

	currentPage := (request.Offset / request.Limit) + 1
	totalPages := (totalCount + request.Limit - 1) / request.Limit

	return &dto.PaginatedMeetingsResponse{
		Meetings:        meetings,
		HasNextPage:     currentPage < totalPages,
		HasPreviousPage: currentPage > 1,
		TotalCount:      totalCount,
		CurrentPage:     currentPage,
		TotalPages:      totalPages,
	}, nil
}

func (s *meetingService) DeleteMeeting(ctx context.Context, request dto.DeleteMeetingRequest) error {
	err := s.queries.DeleteMeeting(ctx, request.ID)
	if err != nil {
		return err
	}
	return nil
}

func (s *meetingService) GetMeeting(ctx context.Context, request dto.GetMeetingRequest) (*dto.MeetingResponse, error) {
	meeting, err := s.queries.GetMeeting(ctx, repo.GetMeetingParams{
		ID:     request.ID,
		UserID: request.UserID,
	})
	if err != nil {
		return nil, err
	}
	return toMeetingAgentResponse(meeting), nil
}

func (s *meetingService) StartMeeting(ctx context.Context, request dto.StartMeetingRequest) (string, error) {
	// Get meeting details
	meeting, err := s.queries.GetMeeting(ctx, repo.GetMeetingParams{
		ID:     request.ID,
		UserID: request.UserID,
	})
	if err != nil {
		return "", fmt.Errorf("failed to get meeting: %w", err)
	}

	if meeting.Status != "upcoming" {
		return "", fmt.Errorf("meeting is not in upcoming state")
	}

	meetingIDStr := request.ID.String()
	session := livekit.NewLiveKitSession(
		meetingIDStr,
		request.UserID,
		meeting.AgentName,
		s.lkConfig,
		s.geminiConfig,
		s.awsConfig,
		livekit.SessionCallbacks{
			OnMeetingEnd: func(meetingID string, recordingURL string, transcriptURL string) {
				s.onMeetingEnd(meetingID, recordingURL, transcriptURL)
			},
		},
	)

	if err := session.Start(); err != nil {
		return "", fmt.Errorf("failed to start session: %w", err)
	}
	startTime := time.Now()
	_, err = s.UpdateMeeting(ctx, dto.UpdateMeetingRequest{
		ID:        request.ID,
		UserID:    request.UserID,
		Status:    "active",
		StartTime: &startTime,
	})
	if err != nil {
		session.Stop()
		return "", fmt.Errorf("failed to update meeting: %w", err)
	}
	token, err := session.GenerateUserToken(request.UserID)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return token, nil
}

func toMeetingAgentResponse(meeting repo.GetMeetingRow) *dto.MeetingResponse {
	return &dto.MeetingResponse{
		ID:        meeting.ID,
		Name:      meeting.Name,
		UserID:    meeting.UserID,
		AgentID:   meeting.AgentID,
		Status:    meeting.Status,
		CreatedAt: meeting.CreatedAt,
		UpdatedAt: meeting.UpdatedAt,
		AgentDetails: &dto.AgentDetails{
			Name:         meeting.AgentName,
			Instructions: meeting.AgentInstructions,
		},
	}
}

func toMeetingResponse(meeting repo.Meeting) *dto.MeetingResponse {
	return &dto.MeetingResponse{
		ID:            meeting.ID,
		Name:          meeting.Name,
		UserID:        meeting.UserID,
		AgentID:       meeting.AgentID,
		Status:        meeting.Status,
		TranscriptUrl: meeting.TranscriptUrl,
		RecordingUrl:  meeting.RecordingUrl,
		Summary:       meeting.Summary,
		CreatedAt:     meeting.CreatedAt,
		UpdatedAt:     meeting.UpdatedAt,
	}
}

// TODO: triggers Inngest background post-processing
func (s *meetingService) onMeetingEnd(meetingID string, recordingURL string, transcriptURL string) {
	fmt.Println("Meeting ended, starting post-processing", "meetingID", meetingID)

	// For now, update meeting directly
	go func() {
		ctx := context.Background()
		meetingUUID, err := uuid.Parse(meetingID)
		if err != nil {
			fmt.Println("Failed to parse meeting ID", err, "meetingID", meetingID)
			return
		}

		endTime := time.Now()
		updateReq := dto.UpdateMeetingRequest{
			ID:      meetingUUID,
			UserID:  "", // System update, no user ID needed
			Status:  "completed",
			EndTime: &endTime,
		}

		if recordingURL != "" {
			updateReq.RecordingURL = &recordingURL
		}
		if transcriptURL != "" {
			updateReq.TranscriptURL = &transcriptURL
		}

		// Get the meeting first to get the user ID
		meeting, err := s.queries.GetMeeting(ctx, repo.GetMeetingParams{
			ID:     meetingUUID,
			UserID: "", // We need to query without user filter
		})
		if err != nil {
			fmt.Println("Failed to get meeting for post-processing", err, "meetingID", meetingID)
			return
		}

		updateReq.UserID = meeting.UserID
		_, err = s.UpdateMeeting(ctx, updateReq)
		if err != nil {
			fmt.Println("Failed to update meeting after end", err, "meetingID", meetingID)
			return
		}

		fmt.Println("Meeting post-processing completed", "meetingID", meetingID)
	}()
}
