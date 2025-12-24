package ports

import (
	"context"

	"app/internal/core/domain"
)

type CameraRepository interface {
	Save(ctx context.Context, camera *domain.Camera) error
	GetByID(ctx context.Context, id string) (*domain.Camera, error)
	List(ctx context.Context) ([]*domain.Camera, error)
	ListByZone(ctx context.Context, zoneID string) ([]*domain.Camera, error)
	Update(ctx context.Context, camera *domain.Camera) error
	Delete(ctx context.Context, id string) error
}

type CameraService interface {
	CreateCamera(ctx context.Context, req *domain.CreateCameraRequest) (*domain.Camera, error)
	GetCamera(ctx context.Context, id string) (*domain.Camera, error)
	ListCameras(ctx context.Context, zoneID string) ([]*domain.Camera, error)
	UpdateCamera(ctx context.Context, id string, req *domain.UpdateCameraRequest) (*domain.Camera, error)
	DeleteCamera(ctx context.Context, id string) error
}
