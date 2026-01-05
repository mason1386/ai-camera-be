package http

import (
	"net/http"

	"app/internal/core/domain"
	"app/internal/core/ports"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service ports.UserService
}

func NewUserHandler(service ports.UserService) *UserHandler {
	return &UserHandler{service: service}
}

// CreateUser godoc
// @Summary Create a new user
// @Tags users
// @Accept json
// @Produce json
// @Param request body domain.CreateWebUserRequest true "User Info"
// @Success 201 {object} domain.User
// @Failure 400 {object} map[string]string
// @Router /users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req domain.CreateWebUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	user, err := h.service.CreateUser(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}

// ListUsers godoc
// @Summary List all users
// @Tags users
// @Produce json
// @Param q query string false "Search query"
// @Success 200 {array} domain.User
// @Router /users [get]
func (h *UserHandler) ListUsers(c *gin.Context) {
	search := c.Query("q")
	users, err := h.service.ListUsers(c.Request.Context(), search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

// UpdateUser godoc
// @Summary Update a user
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param request body domain.UpdateWebUserRequest true "User Info"
// @Success 200 {object} domain.User
// @Router /users/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var req domain.UpdateWebUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	user, err := h.service.UpdateUser(c.Request.Context(), id, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// DeleteUser godoc
// @Summary Delete a user
// @Tags users
// @Param id path string true "User ID"
// @Success 204 "No Content"
// @Router /users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.DeleteUser(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// ResetPassword godoc
// @Summary Reset user password
// @Tags users
// @Accept json
// @Produce json
// @Param user_id path string true "User ID"
// @Param request body object true "New Password"
// @Success 200 {object} map[string]string
// @Router /users/{user_id}/reset-password [post]
func (h *UserHandler) ResetPassword(c *gin.Context) {
	userID := c.Param("user_id")
	var req struct {
		NewPassword string `json:"new_password" binding:"required,min=6"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	if err := h.service.ResetPassword(c.Request.Context(), userID, req.NewPassword); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password reset successfully"})
}
