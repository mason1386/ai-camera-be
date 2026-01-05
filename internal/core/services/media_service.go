package services

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"time"

	"app/internal/core/ports"

	"github.com/google/uuid"
)

type MediaService struct {
	storage ports.FileStorage
}

func NewMediaService(storage ports.FileStorage) ports.MediaService {
	return &MediaService{storage: storage}
}

func (s *MediaService) UploadImage(ctx context.Context, folder string, filename string, reader io.Reader) (string, error) {
	// Generate unique filename to avoid collision
	ext := filepath.Ext(filename)
	newFilename := fmt.Sprintf("%s/%d_%s%s", folder, time.Now().Unix(), uuid.New().String()[:8], ext)

	return s.storage.SaveFile(ctx, newFilename, reader)
}
