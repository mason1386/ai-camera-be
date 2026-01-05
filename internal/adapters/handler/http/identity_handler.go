package http

import (
	"net/http"
	"strconv"

	"app/internal/core/domain"
	"app/internal/core/ports"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type IdentityHandler struct {
	service ports.IdentityService
}

func NewIdentityHandler(service ports.IdentityService) *IdentityHandler {
	return &IdentityHandler{
		service: service,
	}
}

// CreateIdentity godoc
// @Summary Create a new identity
// @Tags identities
// @Accept json
// @Produce json
// @Param request body ports.CreateIdentityRequest true "Identity Info"
// @Success 201 {object} domain.Identity
// @Router /identities [post]
func (h *IdentityHandler) CreateIdentity(c *gin.Context) {
	var req ports.CreateIdentityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	userIDStr, exists := c.Get("userID")
	if exists {
		if uid, err := uuid.Parse(userIDStr.(string)); err == nil {
			req.CreatedBy = &uid
		}
	}

	identity, err := h.service.CreateIdentity(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, identity)
}

// ListIdentities godoc
// @Summary List identities
// @Tags identities
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Limit"
// @Param q query string false "Search query"
// @Success 200 {object} PaginatedResponse
// @Router /identities [get]
func (h *IdentityHandler) ListIdentities(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.Query("q")

	items, total, err := h.service.ListIdentities(c.Request.Context(), page, limit, search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, PaginatedResponse{
		Data:  items,
		Total: total,
		Page:  page,
		Limit: limit,
	})
}

// GetIdentity godoc
// @Summary Get an identity by ID
// @Tags identities
// @Accept json
// @Produce json
// @Param id path string true "Identity ID"
// @Success 200 {object} domain.Identity
// @Router /identities/{id} [get]
func (h *IdentityHandler) GetIdentity(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid ID format"})
		return
	}

	identity, err := h.service.GetIdentity(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	if identity == nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Identity not found"})
		return
	}

	c.JSON(http.StatusOK, identity)
}

// UpdateIdentity godoc
// @Summary Update an identity
// @Tags identities
// @Accept json
// @Produce json
// @Param id path string true "Identity ID"
// @Param request body ports.UpdateIdentityRequest true "Identity Update Info"
// @Success 200 {object} domain.Identity
// @Router /identities/{id} [put]
func (h *IdentityHandler) UpdateIdentity(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid ID format"})
		return
	}

	var req ports.UpdateIdentityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	identity, err := h.service.UpdateIdentity(c.Request.Context(), id, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, identity)
}

// UpdateStatus godoc
// @Summary Update identity status
// @Tags identities
// @Accept json
// @Produce json
// @Param id path string true "Identity ID"
// @Param request body object true "Status object"
// @Success 200 {object} domain.Identity
// @Router /identities/{id}/status [patch]
func (h *IdentityHandler) UpdateStatus(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid ID format"})
		return
	}

	var req struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	identity, err := h.service.UpdateStatus(c.Request.Context(), id, domain.IdentityStatus(req.Status))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, identity)
}

// DeleteIdentity godoc
// @Summary Delete an identity
// @Tags identities
// @Param id path string true "Identity ID"
// @Success 204 "No Content"
// @Router /identities/{id} [delete]
func (h *IdentityHandler) DeleteIdentity(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid ID format"})
		return
	}

	if err := h.service.DeleteIdentity(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// EnrollFace godoc
// @Summary Enroll a new face for an identity
// @Tags identities
// @Accept json
// @Produce json
// @Param request body object true "Face Info"
// @Success 200 {object} domain.IdentityFace
// @Router /identities/enroll-face [post]
func (h *IdentityHandler) EnrollFace(c *gin.Context) {
	var req struct {
		IdentityID string `json:"identity_id" binding:"required"`
		ImageURL   string `json:"image_url" binding:"required"`
		IsPrimary  bool   `json:"is_primary"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	uid, _ := uuid.Parse(req.IdentityID)
	face, err := h.service.EnrollFace(c.Request.Context(), uid, req.ImageURL, req.IsPrimary)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, face)
}

// DeleteFace godoc
// @Summary Delete a face
// @Tags identities
// @Param face_id path string true "Face ID"
// @Success 204 "No Content"
// @Router /identities/faces/{face_id} [delete]
func (h *IdentityHandler) DeleteFace(c *gin.Context) {
	id, _ := uuid.Parse(c.Param("face_id"))
	if err := h.service.DeleteFace(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
