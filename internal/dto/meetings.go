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

type UpdateMeetingRequest struct {
	ID      uuid.UUID `json:"-"`
	UserID  string    `json:"-"`
	Name    string    `json:"name,omitempty"`
	AgentID uuid.UUID `json:"agentId,omitempty"`
	Status  string    `json:"status,omitempty"`
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

// Responses
type MeetingResponse struct {
	ID        uuid.UUID `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	UserID    string    `db:"user_id" json:"userId"`
	AgentID   uuid.UUID `db:"agent_id" json:"agentId"`
	Status    string    `db:"status" json:"status"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`
}

type PaginatedMeetingsResponse struct {
	Meetings        []MeetingResponse `json:"meetings"`
	HasNextPage     bool              `json:"hasNextPage"`
	HasPreviousPage bool              `json:"hasPreviousPage"`
	TotalCount      int32             `json:"totalCount"`
	CurrentPage     int32             `json:"currentPage"`
	TotalPages      int32             `json:"totalPages"`
}
