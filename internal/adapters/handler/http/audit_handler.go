package http

import (
	"net/http"
	"strconv"

	"app/internal/core/ports"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuditHandler struct {
	service ports.AuditService
}

func NewAuditHandler(service ports.AuditService) *AuditHandler {
	return &AuditHandler{service: service}
}

// ListLogs godoc
// @Summary List audit logs
// @Tags audit
// @Accept json
// @Produce json
// @Param user_id query string false "User ID"
// @Param action query string false "Action"
// @Param table query string false "Table Name"
// @Param limit query int false "Limit"
// @Success 200 {object} AuditLogResponse
// @Router /audit-logs [get]
func (h *AuditHandler) ListLogs(c *gin.Context) {
	filter := &ports.AuditFilter{
		Limit:  20,
		Offset: 0,
	}

	if uid := c.Query("user_id"); uid != "" {
		if id, err := uuid.Parse(uid); err == nil {
			filter.UserID = &id
		}
	}
	if act := c.Query("action"); act != "" {
		filter.Action = &act
	}
	if tbl := c.Query("table"); tbl != "" {
		filter.TableName = &tbl
	}
	if l := c.Query("limit"); l != "" {
		if val, err := strconv.Atoi(l); err == nil {
			filter.Limit = int32(val)
		}
	}

	logs, err := h.service.ListLogs(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, logs)
}
