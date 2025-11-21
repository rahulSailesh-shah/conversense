package service

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rahulSailesh-shah/converSense/internal/db/repo"
	"github.com/rahulSailesh-shah/converSense/pkg/config"
	"github.com/rahulSailesh-shah/converSense/pkg/inngest"
)

type Service struct {
	Agent   AgentService
	Meeting MeetingService
	Gemini  GeminiService
}

func NewService(db *pgxpool.Pool, queries *repo.Queries, inngest *inngest.Inngest, cfg *config.AppConfig) *Service {
	// Initialize Services
	agentService := NewAgentService(db, queries)
	meetingService := NewMeetingService(db, queries, inngest, &cfg.LiveKit, &cfg.Gemini, &cfg.AWS)
	geminiService := NewGeminiService(queries, &cfg.Gemini, &cfg.AWS)

	return &Service{
		Agent:   agentService,
		Meeting: meetingService,
		Gemini:  geminiService,
	}
}
