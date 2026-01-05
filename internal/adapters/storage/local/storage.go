package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"app/internal/core/ports"
)

type LocalStorage struct {
	rootPath string
	baseURL  string
}

func NewLocalStorage(rootPath, baseURL string) ports.FileStorage {
	// Ensure directory exists
	if _, err := os.Stat(rootPath); os.IsNotExist(err) {
		os.MkdirAll(rootPath, 0755)
	}
	return &LocalStorage{
		rootPath: rootPath,
		baseURL:  strings.TrimSuffix(baseURL, "/"),
	}
}

func (s *LocalStorage) SaveFile(ctx context.Context, filename string, reader io.Reader) (string, error) {
	fullPath := filepath.Join(s.rootPath, filename)

	// Create subdirectories if needed
	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return "", err
	}

	out, err := os.Create(fullPath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	_, err = io.Copy(out, reader)
	if err != nil {
		return "", err
	}

	// Return URL
	return fmt.Sprintf("%s/%s", s.baseURL, filepath.ToSlash(filename)), nil
}

func (s *LocalStorage) DeleteFile(ctx context.Context, fileURL string) error {
	// Logic to extract relative path from URL and delete
	return nil
}
