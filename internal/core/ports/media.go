package ports

import (
	"context"
	"io"
)

type FileStorage interface {
	SaveFile(ctx context.Context, filename string, reader io.Reader) (string, error)
	DeleteFile(ctx context.Context, fileURL string) error
}

type MediaService interface {
	UploadImage(ctx context.Context, folder string, filename string, reader io.Reader) (string, error)
}
