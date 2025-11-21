package dto

import (
	"time"

	"github.com/google/uuid"
)

// Requests
type CreateMeetingRequest struct {
	Name    string    `json:"name" binding:"required"`
	UserID  string    `json:"-"`
	AgentID uuid.UUID `json:"agentId" binding:"required"`
}

type StartMeetingRequest struct {
	ID     uuid.UUID `json:"-"`
	UserID string    `json:"-"`
}

type UpdateMeetingRequest struct {
	ID            uuid.UUID  `json:"-"`
	UserID        string     `json:"-"`
	Name          string     `json:"name,omitempty"`
	AgentID       uuid.UUID  `json:"agentId,omitempty"`
	Status        string     `json:"status,omitempty"`
	StartTime     *time.Time `json:"startTime,omitempty"`
	EndTime       *time.Time `json:"endTime,omitempty"`
	TranscriptURL *string    `json:"transcriptUrl,omitempty"`
	RecordingURL  *string    `json:"recordingUrl,omitempty"`
	Summary       *string    `json:"summary,omitempty"`
}

type GetMeetingsRequest struct {
	UserID string `form:"userId"`
	Search string `form:"search"`
	Limit  int32  `form:"limit"`
	Offset int32  `form:"offset"`
}

type GetMeetingRequest struct {
	ID     uuid.UUID `json:"-"`
	UserID string    `json:"-"`
}

type DeleteMeetingRequest struct {
	ID     uuid.UUID `json:"-"`
	UserID string    `json:"-"`
}

type GetPreSignedRecordingURLRequest struct {
	MeetingID uuid.UUID `json:"-" `
	UserID    string    `json:"-" `
	FileType  string    `json:"fileType" binding:"required,oneof=recording transcript"` // "recording" or "transcript"
}

// Responses

type MeetingResponse struct {
	ID            uuid.UUID     `db:"id" json:"id"`
	Name          string        `db:"name" json:"name"`
	UserID        string        `db:"user_id" json:"userId"`
	AgentID       uuid.UUID     `db:"agent_id" json:"agentId"`
	StartTime     *time.Time    `db:"start_time" json:"startTime"`
	EndTime       *time.Time    `db:"end_time" json:"endTime"`
	Status        string        `db:"status" json:"status"`
	TranscriptUrl *string       `db:"transcript_url" json:"transcriptUrl"`
	RecordingUrl  *string       `db:"recording_url" json:"recordingUrl"`
	Summary       *string       `db:"summary" json:"summary"`
	CreatedAt     time.Time     `db:"created_at" json:"createdAt"`
	UpdatedAt     time.Time     `db:"updated_at" json:"updatedAt"`
	AgentDetails  *AgentDetails `json:"agentDetails,omitempty"`
}

type PaginatedMeetingsResponse struct {
	Meetings        []MeetingResponse `json:"meetings"`
	HasNextPage     bool              `json:"hasNextPage"`
	HasPreviousPage bool              `json:"hasPreviousPage"`
	TotalCount      int32             `json:"totalCount"`
	CurrentPage     int32             `json:"currentPage"`
	TotalPages      int32             `json:"totalPages"`
}
