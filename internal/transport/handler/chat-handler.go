package handler

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rahulSailesh-shah/converSense/internal/dto"
	"github.com/rahulSailesh-shah/converSense/internal/service"
)

type ChatHandler struct {
	chatService service.ChatService
}

func NewChatHandler(chatService service.ChatService) *ChatHandler {
	return &ChatHandler{
		chatService: chatService,
	}
}

type ChatRequest struct {
	Message string `json:"message" binding:"required"`
}

func (h *ChatHandler) Chat(c *gin.Context) {
	meetingID, err := uuid.Parse(c.Param("meetingId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Invalid meeting ID",
			Error:   err.Error(),
		})
		return
	}

	var req ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Invalid request body",
			Error:   err.Error(),
		})
		return
	}

	userID := c.MustGet("userId").(string)

	stream, err := h.chatService.Chat(c.Request.Context(), meetingID, userID, req.Message)
	if err != nil {
		// Distinguish between auth error and other errors if possible,
		// but for now 500 or 403 based on error string is a simple heuristic
		status := http.StatusInternalServerError
		if err.Error() == "unauthorized: user does not have access to this meeting or meeting does not exist" {
			status = http.StatusForbidden
		}

		c.JSON(status, dto.ErrorResponse{
			Message: "Failed to process chat request",
			Error:   err.Error(),
		})
		return
	}

	// Set headers for SSE
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Transfer-Encoding", "chunked")

	c.Stream(func(w io.Writer) bool {
		if chunk, ok := <-stream; ok {
			fmt.Println("Chunk--->:", chunk)
			// Send data chunk
			c.SSEvent("message", chunk)
			return true
		}
		// Stream finished
		return false
	})
}

func (h *ChatHandler) GetHistory(c *gin.Context) {
	// 1. Parse meeting ID
	meetingIDStr := c.Param("meetingId")
	meetingID, err := uuid.Parse(meetingIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid meeting ID"})
		return
	}

	// 2. Get user ID from context
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// 3. Fetch chat history
	history, err := h.chatService.GetChatHistory(c.Request.Context(), meetingID, userID.(string))
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"messages": history})
}
