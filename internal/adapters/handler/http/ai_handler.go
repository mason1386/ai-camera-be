package http

import (
	"net/http"
	"strconv"
	"time"

	"app/internal/core/domain"
	"app/internal/core/ports"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AIHandler struct {
	service ports.AIService
}

func NewAIHandler(service ports.AIService) *AIHandler {
	return &AIHandler{service: service}
}

// GetConfig godoc
// @Summary Get AI configuration for a camera
// @Tags ai
// @Accept json
// @Produce json
// @Param cameraId path string true "Camera ID"
// @Success 200 {object} domain.AIConfig
// @Router /ai-configs/camera/{cameraId} [get]
func (h *AIHandler) GetConfig(c *gin.Context) {
	cameraID, err := uuid.Parse(c.Param("cameraId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid Camera ID"})
		return
	}

	config, err := h.service.GetConfig(c.Request.Context(), cameraID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	if config == nil {
		// Return defaults or 404
		c.JSON(http.StatusOK, domain.AIConfig{CameraID: cameraID, AIEnabled: false})
		return
	}

	c.JSON(http.StatusOK, config)
}

// UpdateConfig godoc
// @Summary Update AI configuration
// @Tags ai
// @Accept json
// @Produce json
// @Param request body domain.AIConfig true "AI Config Info"
// @Success 200 {object} domain.AIConfig
// @Router /ai-configs [post]
func (h *AIHandler) UpdateConfig(c *gin.Context) {
	var req domain.AIConfig
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	if err := h.service.UpdateConfig(c.Request.Context(), &req); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, req)
}

// ListEvents godoc
// @Summary List AI events
// @Tags ai
// @Accept json
// @Produce json
// @Param camera_id query string false "Camera ID"
// @Param event_type query string false "Event Type"
// @Param status query string false "Status"
// @Param from_date query string false "From Date"
// @Param to_date query string false "To Date"
// @Success 200 {object} PaginatedResponse
// @Router /events [get]
func (h *AIHandler) ListEvents(c *gin.Context) {
	filter := &ports.EventFilter{
		Limit:  10,
		Offset: 0,
	}

	if cid := c.Query("camera_id"); cid != "" {
		if uid, err := uuid.Parse(cid); err == nil {
			filter.CameraID = &uid
		}
	}
	if et := c.Query("event_type"); et != "" {
		t := domain.EventType(et)
		filter.EventType = &t
	}
	if st := c.Query("status"); st != "" {
		s := domain.EventStatus(st)
		filter.Status = &s
	}
	if from := c.Query("from_date"); from != "" {
		if t, err := time.Parse(time.RFC3339, from); err == nil {
			filter.FromDate = &t
		}
	}
	if to := c.Query("to_date"); to != "" {
		if t, err := time.Parse(time.RFC3339, to); err == nil {
			filter.ToDate = &t
		}
	}

	if l := c.Query("limit"); l != "" {
		if val, err := strconv.Atoi(l); err == nil {
			filter.Limit = int32(val)
		}
	}
	if o := c.Query("offset"); o != "" {
		if val, err := strconv.Atoi(o); err == nil {
			filter.Offset = int32(val)
		}
	}

	events, err := h.service.ListEvents(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, events)
}

// UpdateEventStatus godoc
// @Summary Update AI event status
// @Tags ai
// @Accept json
// @Produce json
// @Param id path string true "Event ID"
// @Param request body object true "Status object"
// @Success 200 {object} domain.AIEvent
// @Router /events/{id} [patch]
func (h *AIHandler) UpdateEventStatus(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid ID"})
		return
	}

	var req struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	var resolvedBy uuid.UUID
	if val, exists := c.Get("userID"); exists {
		resolvedBy, _ = uuid.Parse(val.(string))
	}

	event, err := h.service.UpdateEventStatus(c.Request.Context(), id, domain.EventStatus(req.Status), resolvedBy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, event)
}

// GetDashboardStats godoc
// @Summary Get dashboard statistics
// @Tags ai
// @Accept json
// @Produce json
// @Success 200 {object} domain.DashboardStats
// @Router /stats/dashboard [get]
func (h *AIHandler) GetDashboardStats(c *gin.Context) {
	stats, err := h.service.GetDashboardStats(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, stats)
}
