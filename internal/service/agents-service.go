package service

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rahulSailesh-shah/converSense/internal/db/repo"
	"github.com/rahulSailesh-shah/converSense/internal/dto"
)

type AgentService interface {
	CreateAgent(ctx context.Context, request dto.CreateAgentRequest) (*dto.AgentResponse, error)
	UpdateAgent(ctx context.Context, request dto.UpdateAgentRequest) (*dto.AgentResponse, error)
	GetAgents(ctx context.Context, request dto.GetAgentsRequest) (*dto.PaginatedAgentsResponse, error)
	GetAgent(ctx context.Context, request dto.GetAgentRequest) (*dto.AgentResponse, error)
	DeleteAgent(ctx context.Context, request dto.DeleteAgentRequest) error
}

type agentService struct {
	queries *repo.Queries
	db      *pgxpool.Pool
}

func NewAgentService(db *pgxpool.Pool, queries *repo.Queries) AgentService {
	return &agentService{
		db:      db,
		queries: queries,
	}
}

func (s *agentService) CreateAgent(ctx context.Context, request dto.CreateAgentRequest) (*dto.AgentResponse, error) {
	newAgent, err := s.queries.CreateAgent(ctx, repo.CreateAgentParams{
		Name:         request.Name,
		UserID:       request.UserID,
		Instructions: request.Instructions,
	})
	if err != nil {
		return nil, err
	}
	return toAgentResponse(newAgent), nil
}

func (s *agentService) UpdateAgent(ctx context.Context, request dto.UpdateAgentRequest) (*dto.AgentResponse, error) {
	currentAgent, err := s.queries.GetAgent(ctx, repo.GetAgentParams{
		ID:     request.ID,
		UserID: request.UserID,
	})
	if err != nil {
		return nil, err
	}

	if request.Name != "" {
		currentAgent.Name = request.Name
	}
	if request.Instructions != "" {
		currentAgent.Instructions = request.Instructions
	}

	updatedAgent, err := s.queries.UpdateAgent(ctx, repo.UpdateAgentParams{
		ID:           currentAgent.ID,
		Name:         currentAgent.Name,
		Instructions: currentAgent.Instructions,
	})
	if err != nil {
		return nil, err
	}

	return toAgentResponse(updatedAgent), nil
}

func (s *agentService) GetAgents(ctx context.Context, request dto.GetAgentsRequest) (*dto.PaginatedAgentsResponse, error) {
	rows, err := s.queries.GetAgents(ctx, repo.GetAgentsParams{
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
	agents := make([]dto.AgentResponse, 0, len(rows))
	for _, row := range rows {
		agents = append(agents, dto.AgentResponse{
			ID:           row.ID,
			UserID:       row.UserID,
			Name:         row.Name,
			Instructions: row.Instructions,
			CreatedAt:    row.CreatedAt,
			UpdatedAt:    row.UpdatedAt,
		})
	}

	currentPage := (request.Offset / request.Limit) + 1
	totalPages := (totalCount + request.Limit - 1) / request.Limit

	return &dto.PaginatedAgentsResponse{
		Agents:          agents,
		HasNextPage:     currentPage < totalPages,
		HasPreviousPage: currentPage > 1,
		TotalCount:      totalCount,
		CurrentPage:     currentPage,
		TotalPages:      totalPages,
	}, nil
}

func (s *agentService) DeleteAgent(ctx context.Context, request dto.DeleteAgentRequest) error {
	err := s.queries.DeleteAgent(ctx, request.ID)
	if err != nil {
		return err
	}
	return nil
}

func (s *agentService) GetAgent(ctx context.Context, request dto.GetAgentRequest) (*dto.AgentResponse, error) {
	agent, err := s.queries.GetAgent(ctx, repo.GetAgentParams{
		ID:     request.ID,
		UserID: request.UserID,
	})
	if err != nil {
		return nil, err
	}
	return toAgentResponse(agent), nil
}

func toAgentResponse(agent repo.Agent) *dto.AgentResponse {
	return &dto.AgentResponse{
		ID:           agent.ID,
		Name:         agent.Name,
		UserID:       agent.UserID,
		Instructions: agent.Instructions,
		CreatedAt:    agent.CreatedAt,
		UpdatedAt:    agent.UpdatedAt,
	}
}
