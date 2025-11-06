package http

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/lestrrat-go/jwx/v3/jwk"
	"github.com/rahulSailesh-shah/converSense/internal/app"
	"github.com/rahulSailesh-shah/converSense/internal/transport/handler"
	"github.com/rahulSailesh-shah/converSense/internal/transport/http/middleware"
)

func RegisterRoutes(r *gin.Engine, authKeys jwk.Set, app *app.App) {
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://127.0.0.1:5173", "http://localhost:9001", "http://127.0.0.1:9001"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
	}))

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	// Middlewares
	protected := r.Group("")
	protected.Use(middleware.AuthMiddleware(authKeys))

	// Agent routes
	agentRoutes := protected.Group("/agents")
	agentHandler := handler.NewAgentHandler(app.Service.Agent)
	{
		agentRoutes.POST("", agentHandler.CreateAgent)
		agentRoutes.PUT("/:id", agentHandler.UpdateAgent)
		agentRoutes.GET("", agentHandler.GetAgents)
		agentRoutes.DELETE("/:id", agentHandler.DeleteAgent)
	}
}
