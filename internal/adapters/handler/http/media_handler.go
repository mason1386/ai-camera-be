package http

import (
	"net/http"

	"app/internal/core/ports"

	"github.com/gin-gonic/gin"
)

type MediaHandler struct {
	service ports.MediaService
}

func NewMediaHandler(service ports.MediaService) *MediaHandler {
	return &MediaHandler{service: service}
}

// UploadImage godoc
// @Summary Upload an image
// @Tags media
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Image file"
// @Param folder query string false "Subfolder name"
// @Success 200 {object} map[string]string
// @Router /media/upload [post]
func (h *MediaHandler) UploadImage(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "No file provided"})
		return
	}
	defer file.Close()

	// Folder from query or default
	folder := c.DefaultQuery("folder", "identities")

	url, err := h.service.UploadImage(c.Request.Context(), folder, header.Filename, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"url":      url,
		"filename": header.Filename,
	})
}
