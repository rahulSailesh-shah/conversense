package app

import (
	"context"
	"fmt"

	"github.com/rahulSailesh-shah/converSense/internal/db/repo"
	"github.com/rahulSailesh-shah/converSense/internal/service"
	"github.com/rahulSailesh-shah/converSense/pkg/config"
	"github.com/rahulSailesh-shah/converSense/pkg/database"
	"github.com/rahulSailesh-shah/converSense/pkg/inngest"
)

type App struct {
	Config  *config.AppConfig
	DB      database.DB
	Service *service.Service
	Inngest *inngest.Inngest
}

func NewApp(ctx context.Context, cfg *config.AppConfig) (*App, error) {
	db := database.NewPostgresDB(ctx, &cfg.DB)
	if err := db.Connect(); err != nil {
		fmt.Println("Error connecting to database:", err)
		return nil, err
	}

	dbInstance := db.GetDB()
	if dbInstance == nil {
		fmt.Println("Database instance is nil")
		return nil, fmt.Errorf("database not initialize")
	}

	queries := repo.New(dbInstance)
	inngest, err := inngest.NewInngest(&cfg.AWS, &cfg.Gemini, &cfg.OpenAI, queries)
	if err != nil {
		return nil, err
	}
	services := service.NewService(dbInstance, queries, inngest, cfg)

	return &App{
		Config:  cfg,
		DB:      db,
		Service: services,
		Inngest: inngest,
	}, nil
}
