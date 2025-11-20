package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rahulSailesh-shah/converSense/internal/dto"
	"github.com/rahulSailesh-shah/converSense/internal/service"
)

type MeetingHandler struct {
	meetingService service.MeetingService
}

func NewMeetingHandler(meetingService service.MeetingService) *MeetingHandler {
	return &MeetingHandler{
		meetingService: meetingService,
	}
}

func (h *MeetingHandler) CreateMeeting(c *gin.Context) {
	var req dto.CreateMeetingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	req.UserID = c.MustGet("userId").(string)
	meeting, err := h.meetingService.CreateMeeting(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "Failed to create meeting",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "Meeting created successfully",
		Data:    meeting,
	})
}

func (h *MeetingHandler) UpdateMeeting(c *gin.Context) {
	var req dto.UpdateMeetingRequest
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
	meeting, err := h.meetingService.UpdateMeeting(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "Failed to update meeting",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "Meeting updated successfully",
		Data:    meeting,
	})
}

func (h *MeetingHandler) GetMeetings(c *gin.Context) {
	search := c.DefaultQuery("search", "")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))

	meetings, err := h.meetingService.GetMeetings(c.Request.Context(), dto.GetMeetingsRequest{
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
		Message: "Meetings retrieved successfully",
		Data:    meetings,
	})
}

func (h *MeetingHandler) GetMeeting(c *gin.Context) {
	agentId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Invalid agent ID",
			Error:   err.Error(),
		})
		return
	}

	meeting, err := h.meetingService.GetMeeting(c.Request.Context(), dto.GetMeetingRequest{
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
		Message: "Meeting retrieved successfully",
		Data:    meeting,
	})
}

func (h *MeetingHandler) DeleteMeeting(c *gin.Context) {
	agentId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Invalid agent ID",
			Error:   err.Error(),
		})
		return
	}

	err = h.meetingService.DeleteMeeting(c.Request.Context(), dto.DeleteMeetingRequest{
		ID: agentId,
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

func (h *MeetingHandler) StartMeeting(c *gin.Context) {
	meetingId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Invalid meeting ID",
			Error:   err.Error(),
		})
		return
	}

	token, err := h.meetingService.StartMeeting(c.Request.Context(), dto.StartMeetingRequest{
		ID:     meetingId,
		UserID: c.MustGet("userId").(string),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "Failed to start meeting",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "Meeting started successfully",
		Data: map[string]string{
			"token": token,
		},
	})
}
