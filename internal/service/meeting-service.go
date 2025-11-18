package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rahulSailesh-shah/converSense/internal/db/repo"
	"github.com/rahulSailesh-shah/converSense/internal/dto"
)

type MeetingService interface {
	CreateMeeting(ctx context.Context, request dto.CreateMeetingRequest) (*dto.MeetingResponse, error)
	UpdateMeeting(ctx context.Context, request dto.UpdateMeetingRequest) (*dto.MeetingResponse, error)
	GetMeetings(ctx context.Context, request dto.GetMeetingsRequest) (*dto.PaginatedMeetingsResponse, error)
	GetMeeting(ctx context.Context, request dto.GetMeetingRequest) (*dto.MeetingResponse, error)
	DeleteMeeting(ctx context.Context, request dto.DeleteMeetingRequest) error
}

type meetingService struct {
	queries *repo.Queries
	db      *pgxpool.Pool
}

func NewMeetingService(db *pgxpool.Pool, queries *repo.Queries) MeetingService {
	return &meetingService{
		db:      db,
		queries: queries,
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
		return nil, err
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

	updatedMeeting, err := s.queries.UpdateMeeting(ctx, repo.UpdateMeetingParams{
		ID:      currentMeeting.ID,
		Name:    currentMeeting.Name,
		AgentID: currentMeeting.AgentID,
		Status:  currentMeeting.Status,
	})
	if err != nil {
		return nil, err
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
