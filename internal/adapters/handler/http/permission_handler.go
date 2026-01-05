package http

import (
	"net/http"

	"app/internal/core/ports"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PermissionHandler struct {
	service ports.PermissionService
}

func NewPermissionHandler(service ports.PermissionService) *PermissionHandler {
	return &PermissionHandler{service: service}
}

// UpdateCameraPermissions godoc
// @Summary Update user camera permissions
// @Tags permissions
// @Accept json
// @Produce json
// @Param userId path string true "User ID"
// @Param request body object true "Camera IDs array"
// @Success 204 "No Content"
// @Router /permissions/{userId}/cameras [post]
func (h *PermissionHandler) UpdateCameraPermissions(c *gin.Context) {
	userID, _ := uuid.Parse(c.Param("userId"))
	var req struct {
		CameraIDs []uuid.UUID `json:"camera_ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	if err := h.service.UpdateUserCameraPermissions(c.Request.Context(), userID, req.CameraIDs); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// UpdateZonePermissions godoc
// @Summary Update user zone permissions
// @Tags permissions
// @Accept json
// @Produce json
// @Param userId path string true "User ID"
// @Param request body object true "Zone IDs array"
// @Success 204 "No Content"
// @Router /permissions/{userId}/zones [post]
func (h *PermissionHandler) UpdateZonePermissions(c *gin.Context) {
	userID, _ := uuid.Parse(c.Param("userId"))
	var req struct {
		ZoneIDs []uuid.UUID `json:"zone_ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	if err := h.service.UpdateUserZonePermissions(c.Request.Context(), userID, req.ZoneIDs); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// GetPermissions godoc
// @Summary Get user permissions
// @Tags permissions
// @Accept json
// @Produce json
// @Param userId path string true "User ID"
// @Success 200 {object} ports.UserPermissions
// @Router /permissions/{userId} [get]
func (h *PermissionHandler) GetPermissions(c *gin.Context) {
	userID, _ := uuid.Parse(c.Param("userId"))
	perms, err := h.service.GetUserPermissions(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, perms)
}
