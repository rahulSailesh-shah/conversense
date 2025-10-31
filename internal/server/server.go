package server

import (
	"context"
	"fmt"
	"log"
	httpSrv "net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rahulSailesh-shah/converSense/internal/app"
	"github.com/rahulSailesh-shah/converSense/internal/transport/http"
	"github.com/rahulSailesh-shah/converSense/pkg/auth"
)

type Server struct {
	App *app.App
	ctx context.Context
	Engine *gin.Engine
	httpServer *httpSrv.Server
}


func NewServer(ctx context.Context, app *app.App) (*Server, error) {
	engine := gin.Default()

	authKeys, err := auth.LoadKeys(app.Config.Auth.JwksURL)
	if err != nil {
		return nil, fmt.Errorf("failed to load auth keys: %w", err)
	}

	// register routes here
	http.RegisterRoutes(engine, authKeys, app )

	srv := &httpSrv.Server{
		Addr:    fmt.Sprintf(":%d", app.Config.Server.Port),
		Handler: engine,
	}

	return &Server{
		App: app,
		ctx: ctx,
		Engine: engine,
		httpServer: srv,
	}, nil
}


func(s *Server) Run() error {
	done := make(chan bool, 1)
	go s.gracefulShutdown(done)

	log.Printf("Starting server on port %d", s.App.Config.Server.Port)
	if err := s.httpServer.ListenAndServe(); err != nil && err != httpSrv.ErrServerClosed {
		return fmt.Errorf("could not start server: %w", err)
	}

	<-done
	return nil
}


func (s *Server) gracefulShutdown(done chan bool) {
	ctx, stop := signal.NotifyContext(s.ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()

	log.Println("shutting down gracefully, press Ctrl+C again to force")
	defer s.App.DB.Close()
	stop()

	timeout := time.Duration(s.App.Config.Server.GracefulShutdownSec) * time.Second
	ctx, cancel := context.WithTimeout(s.ctx, timeout)
	defer cancel()
	if err := s.httpServer.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown with error: %v", err)
	}

	log.Println("Server exiting")
	done <- true
}
