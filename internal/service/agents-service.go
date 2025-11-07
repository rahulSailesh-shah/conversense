package service

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rahulSailesh-shah/converSense/internal/db/repo"
	"github.com/rahulSailesh-shah/converSense/internal/dto"
)

type AgentService interface {
	CreateAgent(ctx context.Context, request dto.CreateAgentRequest) (*dto.AgentResponse, error)
	UpdateAgent(ctx context.Context, request dto.UpdateAgentRequest) (*dto.AgentResponse, error)
	GetAgents(ctx context.Context, request dto.GetAgentsRequest) ([]*dto.AgentResponse, error)
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

func (s *agentService) GetAgents(ctx context.Context, request dto.GetAgentsRequest) ([]*dto.AgentResponse, error) {
	agents, err := s.queries.GetAgents(ctx, request.UserID)
	if err != nil {
		fmt.Println("Error fetching agents:", err)
		return nil, err
	}

	var responses []*dto.AgentResponse
	for _, agent := range agents {
		responses = append(responses, toAgentResponse(agent))
	}
	return responses, nil
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
