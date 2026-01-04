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

type AnalyticsHandler struct {
	service ports.AnalyticsService
}

func NewAnalyticsHandler(service ports.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{service: service}
}

// ListRecognitionLogs godoc
// @Summary List recognition logs
// @Tags analytics
// @Accept json
// @Produce json
// @Param identity_id query string false "Identity ID"
// @Param camera_id query string false "Camera ID"
// @Param from query string false "From Date (RFC3339)"
// @Param to query string false "To Date (RFC3339)"
// @Param limit query int false "Limit"
// @Success 200 {object} RecognitionLogResponse
// @Router /recognition/logs [get]
func (h *AnalyticsHandler) ListRecognitionLogs(c *gin.Context) {
	filter := &ports.RecognitionFilter{
		Limit:  10,
		Offset: 0,
	}

	if id := c.Query("identity_id"); id != "" {
		if uid, err := uuid.Parse(id); err == nil {
			filter.IdentityID = &uid
		}
	}
	if cid := c.Query("camera_id"); cid != "" {
		if uid, err := uuid.Parse(cid); err == nil {
			filter.CameraID = &uid
		}
	}
	if from := c.Query("from"); from != "" {
		if t, err := time.Parse(time.RFC3339, from); err == nil {
			filter.FromDate = &t
		}
	}
	if to := c.Query("to"); to != "" {
		if t, err := time.Parse(time.RFC3339, to); err == nil {
			filter.ToDate = &t
		}
	}
	if l := c.Query("limit"); l != "" {
		if val, err := strconv.Atoi(l); err == nil {
			filter.Limit = int32(val)
		}
	}

	logs, err := h.service.ListRecognitionLogs(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, logs)
}

// ListAttendance godoc
// @Summary List attendance records
// @Tags analytics
// @Accept json
// @Produce json
// @Param identity_id query string false "Identity ID"
// @Param status query string false "Status (late, on_time, absent, early_leave)"
// @Success 200 {object} AttendanceRecordResponse
// @Router /attendance/records [get]
func (h *AnalyticsHandler) ListAttendance(c *gin.Context) {
	filter := &ports.AttendanceFilter{
		Limit:  10,
		Offset: 0,
	}

	if id := c.Query("identity_id"); id != "" {
		if uid, err := uuid.Parse(id); err == nil {
			filter.IdentityID = &uid
		}
	}
	if st := c.Query("status"); st != "" {
		s := domain.AttendanceStatus(st)
		filter.Status = &s
	}

	records, err := h.service.ListAttendance(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, records)
}

// GetSummary godoc
// @Summary Get daily attendance summary
// @Tags analytics
// @Accept json
// @Produce json
// @Param date query string false "Date (YYYY-MM-DD)"
// @Success 200 {object} map[string]int64
// @Router /attendance/summary [get]
func (h *AnalyticsHandler) GetSummary(c *gin.Context) {
	dateStr := c.DefaultQuery("date", time.Now().Format("2006-01-02"))
	date, _ := time.Parse("2006-01-02", dateStr)

	summary, err := h.service.GetDailyAttendanceSummary(c.Request.Context(), date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, summary)
}
