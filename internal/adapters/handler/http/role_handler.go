package http

import (
	"net/http"

	"app/internal/core/domain"
	"app/internal/core/ports"

	"github.com/gin-gonic/gin"
)

type RoleHandler struct {
	service ports.RoleService
}

func NewRoleHandler(service ports.RoleService) *RoleHandler {
	return &RoleHandler{service: service}
}

// CreateRole godoc
// @Summary Create a new role
// @Tags roles
// @Accept json
// @Produce json
// @Param request body domain.Role true "Role Info"
// @Success 201 {object} domain.Role
// @Router /roles [post]
func (h *RoleHandler) CreateRole(c *gin.Context) {
	var role domain.Role
	if err := c.ShouldBindJSON(&role); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	if err := h.service.CreateRole(c.Request.Context(), &role); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, role)
}

// ListRoles godoc
// @Summary List all roles
// @Tags roles
// @Produce json
// @Param q query string false "Search query"
// @Success 200 {object} []domain.Role
// @Router /roles [get]
func (h *RoleHandler) ListRoles(c *gin.Context) {
	search := c.Query("q")
	roles, err := h.service.ListRoles(c.Request.Context(), search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, roles)
}

// GetRole godoc
// @Summary Get role by ID
// @Tags roles
// @Param id path string true "Role ID"
// @Success 200 {object} domain.Role
// @Router /roles/{id} [get]
func (h *RoleHandler) GetRole(c *gin.Context) {
	id := c.Param("id")
	role, err := h.service.GetRole(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	if role == nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Role not found"})
		return
	}
	c.JSON(http.StatusOK, role)
}

// UpdateRole godoc
// @Summary Update a role
// @Tags roles
// @Accept json
// @Produce json
// @Param id path string true "Role ID"
// @Param request body domain.Role true "Role Info"
// @Success 200 {object} domain.Role
// @Router /roles/{id} [put]
func (h *RoleHandler) UpdateRole(c *gin.Context) {
	id := c.Param("id")
	var role domain.Role
	if err := c.ShouldBindJSON(&role); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	if err := h.service.UpdateRole(c.Request.Context(), id, &role); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, role)
}

// DeleteRole godoc
// @Summary Delete a role
// @Tags roles
// @Param id path string true "Role ID"
// @Success 204 "No Content"
// @Router /roles/{id} [delete]
func (h *RoleHandler) DeleteRole(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.DeleteRole(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
