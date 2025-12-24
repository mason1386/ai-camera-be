package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"app/internal/core/domain"
	"app/internal/core/ports"
)

type CameraHandler struct {
	service ports.CameraService
}

func NewCameraHandler(service ports.CameraService) *CameraHandler {
	return &CameraHandler{
		service: service,
	}
}

// CreateCamera godoc
// @Summary Create a new camera
// @Description Create a new camera in the system
// @Tags cameras
// @Accept json
// @Produce json
// @Param camera body domain.CreateCameraRequest true "Camera Info"
// @Success 201 {object} domain.Camera
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /cameras [post]
func (h *CameraHandler) CreateCamera(c *gin.Context) {
	var req domain.CreateCameraRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	camera, err := h.service.CreateCamera(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create camera"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Camera created successfully",
		"data":    camera,
	})
}

// ListCameras godoc
// @Summary List cameras
// @Description List cameras, optionally filtered by zone_id
// @Tags cameras
// @Produce json
// @Param zone_id query string false "Filter by Zone ID"
// @Success 200 {array} domain.Camera
// @Security BearerAuth
// @Router /cameras [get]
func (h *CameraHandler) ListCameras(c *gin.Context) {
	zoneID := c.Query("zone_id")
	cameras, err := h.service.ListCameras(c.Request.Context(), zoneID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cameras)
}

// GetCamera godoc
// @Summary Get camera by ID
// @Tags cameras
// @Produce json
// @Param id path string true "Camera ID"
// @Success 200 {object} domain.Camera
// @Failure 404 {object} map[string]string
// @Security BearerAuth
// @Router /cameras/{id} [get]
func (h *CameraHandler) GetCamera(c *gin.Context) {
	id := c.Param("id")
	camera, err := h.service.GetCamera(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if camera == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Camera not found"})
		return
	}
	c.JSON(http.StatusOK, camera)
}

// UpdateCamera godoc
// @Summary Update camera
// @Tags cameras
// @Accept json
// @Produce json
// @Param id path string true "Camera ID"
// @Param camera body domain.UpdateCameraRequest true "Update Info"
// @Success 200 {object} domain.Camera
// @Failure 404 {object} map[string]string
// @Security BearerAuth
// @Router /cameras/{id} [put]
func (h *CameraHandler) UpdateCamera(c *gin.Context) {
	id := c.Param("id")
	var req domain.UpdateCameraRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	camera, err := h.service.UpdateCamera(c.Request.Context(), id, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if camera == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Camera not found"})
		return
	}
	c.JSON(http.StatusOK, camera)
}

// DeleteCamera godoc
// @Summary Delete camera
// @Tags cameras
// @Param id path string true "Camera ID"
// @Success 200 {object} map[string]string
// @Security BearerAuth
// @Router /cameras/{id} [delete]
func (h *CameraHandler) DeleteCamera(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.DeleteCamera(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Camera deleted"})
}
