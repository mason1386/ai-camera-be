package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"app/internal/core/domain"
	"app/internal/core/ports"
)

type ZoneHandler struct {
	service ports.ZoneService
}

func NewZoneHandler(service ports.ZoneService) *ZoneHandler {
	return &ZoneHandler{service: service}
}

// CreateZone godoc
// @Summary Create a new zone
// @Tags zones
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param zone body domain.CreateZoneRequest true "Zone Info"
// @Success 201 {object} domain.Zone
// @Failure 400 {object} map[string]string
// @Router /zones [post]
func (h *ZoneHandler) CreateZone(c *gin.Context) {
	var req domain.CreateZoneRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	zone, err := h.service.CreateZone(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, zone)
}

// ListZones godoc
// @Summary List all zones
// @Tags zones
// @Produce json
// @Success 200 {array} domain.Zone
// @Security BearerAuth
// @Router /zones [get]
func (h *ZoneHandler) ListZones(c *gin.Context) {
	zones, err := h.service.ListZones(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, zones)
}

// GetZone godoc
// @Summary Get a zone by ID
// @Tags zones
// @Produce json
// @Param id path string true "Zone ID"
// @Success 200 {object} domain.Zone
// @Failure 404 {object} map[string]string
// @Security BearerAuth
// @Router /zones/{id} [get]
func (h *ZoneHandler) GetZone(c *gin.Context) {
	id := c.Param("id")
	zone, err := h.service.GetZone(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if zone == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Zone not found"})
		return
	}
	c.JSON(http.StatusOK, zone)
}

// UpdateZone godoc
// @Summary Update a zone
// @Tags zones
// @Accept json
// @Produce json
// @Param id path string true "Zone ID"
// @Param zone body domain.UpdateZoneRequest true "Update Info"
// @Success 200 {object} domain.Zone
// @Failure 404 {object} map[string]string
// @Security BearerAuth
// @Router /zones/{id} [put]
func (h *ZoneHandler) UpdateZone(c *gin.Context) {
	id := c.Param("id")
	var req domain.UpdateZoneRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	zone, err := h.service.UpdateZone(c.Request.Context(), id, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if zone == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Zone not found"})
		return
	}
	c.JSON(http.StatusOK, zone)
}

// DeleteZone godoc
// @Summary Delete a zone
// @Tags zones
// @Param id path string true "Zone ID"
// @Success 200 {object} map[string]string
// @Security BearerAuth
// @Router /zones/{id} [delete]
func (h *ZoneHandler) DeleteZone(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.DeleteZone(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Zone deleted"})
}
