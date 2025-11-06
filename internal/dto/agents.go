package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateAgentRequest struct {
	Name         string `json:"name" binding:"required"`
	UserID       string `json:"userId" binding:"required"`
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
}

type DeleteAgentRequest struct {
	ID     uuid.UUID `json:"-"`
	UserID string    `json:"-"`
}

type AgentResponse struct {
	ID           uuid.UUID `db:"id" json:"id"`
	Name         string    `db:"name" json:"name"`
	UserID       string    `db:"user_id" json:"userId"`
	Instructions string    `db:"instructions" json:"instructions"`
	CreatedAt    time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt    time.Time `db:"updated_at" json:"updatedAt"`
}
