package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateAgentRequest struct {
	Name         string `json:"name" binding:"required"`
	UserID       string `json:"-"`
	Instructions string `json:"instructions" binding:"required"`
}

type UpdateAgentRequest struct {
	ID           uuid.UUID `json:"-"`
	UserID       string    `json:"-"`
	Name         string    `json:"name,omitempty"`
	Instructions string    `json:"instructions,omitempty"`
}

type GetAgentsRequest struct {
	UserID string `form:"userId"`
	Search string `form:"search"`
	Limit  int32  `form:"limit"`
	Offset int32  `form:"offset"`
}

type GetAgentRequest struct {
	ID     uuid.UUID `json:"-"`
	UserID string    `json:"-"`
}

type DeleteAgentRequest struct {
	ID     uuid.UUID `json:"-"`
	UserID string    `json:"-"`
}

type AgentDetails struct {
	Name         string `db:"agent_name" json:"name"`
	Instructions string `db:"agent_instructions" json:"instructions"`
}

type AgentResponse struct {
	ID           uuid.UUID `db:"id" json:"id"`
	Name         string    `db:"name" json:"name"`
	UserID       string    `db:"user_id" json:"userId"`
	Instructions string    `db:"instructions" json:"instructions"`
	CreatedAt    time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt    time.Time `db:"updated_at" json:"updatedAt"`
}

type PaginatedAgentsResponse struct {
	Agents          []AgentResponse `json:"agents"`
	HasNextPage     bool            `json:"hasNextPage"`
	HasPreviousPage bool            `json:"hasPreviousPage"`
	TotalCount      int32           `json:"totalCount"`
	CurrentPage     int32           `json:"currentPage"`
	TotalPages      int32           `json:"totalPages"`
}
