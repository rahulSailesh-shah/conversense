package service

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rahulSailesh-shah/converSense/internal/db/repo"
	"github.com/rahulSailesh-shah/converSense/pkg/config"
)

type Service struct {
	Agent   AgentService
	Meeting MeetingService
}

func NewService(db *pgxpool.Pool, queries *repo.Queries, cfg *config.AppConfig) *Service {
	return &Service{
		Agent:   NewAgentService(db, queries),
		Meeting: NewMeetingService(db, queries, &cfg.LiveKit, &cfg.Gemini, &cfg.AWS),
	}
}
