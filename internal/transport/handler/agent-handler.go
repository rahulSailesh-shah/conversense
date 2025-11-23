package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rahulSailesh-shah/converSense/internal/dto"
	"github.com/rahulSailesh-shah/converSense/internal/service"
)

type AgentHandler struct {
	agentService service.AgentService
}

func NewAgentHandler(agentService service.AgentService) *AgentHandler {
	return &AgentHandler{
		agentService: agentService,
	}
}

func (h *AgentHandler) CreateAgent(c *gin.Context) {
	var req dto.CreateAgentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	req.UserID = c.MustGet("userId").(string)
	agent, err := h.agentService.CreateAgent(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "Failed to create agent",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "Agent created successfully",
		Data:    agent,
	})
}

func (h *AgentHandler) UpdateAgent(c *gin.Context) {
	var req dto.UpdateAgentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	var err error
	req.ID, err = uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Invalid agent ID",
			Error:   err.Error(),
		})
		return
	}
	req.UserID = c.MustGet("userId").(string)
	agent, err := h.agentService.UpdateAgent(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "Failed to update agent",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "Agent updated successfully",
		Data:    agent,
	})
}

func (h *AgentHandler) GetAgents(c *gin.Context) {
	search := c.DefaultQuery("search", "")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))

	agents, err := h.agentService.GetAgents(c.Request.Context(), dto.GetAgentsRequest{
		Search: search,
		Limit:  int32(limit),
		Offset: int32((page - 1) * limit),
		UserID: c.MustGet("userId").(string),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "Failed to get agents",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "Agents retrieved successfully",
		Data:    agents,
	})
}

func (h *AgentHandler) GetAgent(c *gin.Context) {
	agentId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Invalid agent ID",
			Error:   err.Error(),
		})
		return
	}

	agent, err := h.agentService.GetAgent(c.Request.Context(), dto.GetAgentRequest{
		ID:     agentId,
		UserID: c.MustGet("userId").(string),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "Failed to get agent",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "Agent retrieved successfully",
		Data:    agent,
	})
}

func (h *AgentHandler) DeleteAgent(c *gin.Context) {
	agentId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Invalid agent ID",
			Error:   err.Error(),
		})
		return
	}

	err = h.agentService.DeleteAgent(c.Request.Context(), dto.DeleteAgentRequest{
		ID:     agentId,
		UserID: c.MustGet("userId").(string),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "Failed to delete agent",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "Agent deleted successfully",
	})
}
