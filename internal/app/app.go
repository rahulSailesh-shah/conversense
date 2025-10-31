package app

import (
	"context"
	"fmt"

	"github.com/rahulSailesh-shah/converSense/internal/service"
	"github.com/rahulSailesh-shah/converSense/pkg/config"
	"github.com/rahulSailesh-shah/converSense/pkg/database"
)


type App struct {
	Config *config.AppConfig
	DB database.DB
	Service *service.Service
}


func NewApp(ctx context.Context, cfg *config.AppConfig) (*App, error) {
	db := database.NewPostgresDB(ctx, &cfg.DB)
	if err := db.Connect(); err != nil {
		return  nil, err
	}

	dbInstance := db.GetDB()
	if dbInstance == nil {
		return nil,  fmt.Errorf("database not initialize")
	}

	services := service.NewService(dbInstance)

	return &App{
		Config: cfg,
		DB:     db,
		Service: services,
	}, nil

}
