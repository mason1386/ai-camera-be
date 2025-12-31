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
// @Description Register a new face profile
// @Tags identities
// @Accept json
// @Produce json
// @Param request body ports.CreateIdentityRequest true "Identity Info"
// @Success 201 {object} domain.Identity
// @Failure 400 {object} ErrorResponse
// @Router /identities [post]
func (h *IdentityHandler) CreateIdentity(c *gin.Context) {
	var req ports.CreateIdentityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	// Get user ID from context (set by AuthMiddleware)
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
// @Description Get paginated list of identities
// @Tags identities
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Success 200 {object} PaginatedResponse
// @Router /identities [get]
func (h *IdentityHandler) ListIdentities(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	items, total, err := h.service.ListIdentities(c.Request.Context(), page, limit)
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
// @Summary Get identity details
// @Description Get info of a specific identity
// @Tags identities
// @Accept json
// @Produce json
// @Param id path string true "Identity ID"
// @Success 200 {object} domain.Identity
// @Failure 404 {object} ErrorResponse
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
// @Summary Update identity info
// @Description Update name, phone or face image
// @Tags identities
// @Accept json
// @Produce json
// @Param id path string true "Identity ID"
// @Param request body ports.UpdateIdentityRequest true "Update Info"
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
// @Summary Approve or Reject identity
// @Description Admin updates status
// @Tags identities
// @Accept json
// @Produce json
// @Param id path string true "Identity ID"
// @Param request body object{status=string} true "Status (active/rejected)"
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

	// Mocking approved by user ID from context
	approvedBy := uuid.Nil
	userIDStr, exists := c.Get("userID")
	if exists {
		if uid, err := uuid.Parse(userIDStr.(string)); err == nil {
			approvedBy = uid
		}
	}

	identity, err := h.service.UpdateStatus(c.Request.Context(), id, domain.IdentityStatus(req.Status), approvedBy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, identity)
}

// DeleteIdentity godoc
// @Summary Delete identity (Soft delete)
// @Description Remove identity from system
// @Tags identities
// @Accept json
// @Produce json
// @Param id path string true "Identity ID"
// @Success 204
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

type PaginatedResponse struct {
	Data  interface{} `json:"data"`
	Total int64       `json:"total"`
	Page  int         `json:"page"`
	Limit int         `json:"limit"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
