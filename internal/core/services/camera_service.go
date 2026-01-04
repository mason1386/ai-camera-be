package services

import (
	"context"

	"app/internal/core/domain"
	"app/internal/core/ports"
	"app/pkg/logger"

	"go.uber.org/zap"
)

type CameraService struct {
	repo ports.CameraRepository
}

func NewCameraService(repo ports.CameraRepository) ports.CameraService {
	return &CameraService{
		repo: repo,
	}
}

func (s *CameraService) CreateCamera(ctx context.Context, req *domain.CreateCameraRequest) (*domain.Camera, error) {
	camera := &domain.Camera{
		ZoneID:    req.ZoneID,
		Name:      req.Name,
		IPAddress: req.IPAddress,
		RTSPURL:   req.RTSPURL,
		Status:    domain.CameraStatusOnline,
		AIEnabled: req.AIEnabled,
	}

	err := s.repo.Save(ctx, camera)
	if err != nil {
		logger.Error("Failed to create camera", zap.Error(err))
		return nil, err
	}

	logger.Info("Camera created successfully", zap.String("id", camera.ID), zap.String("name", camera.Name))
	return camera, nil
}

func (s *CameraService) GetCamera(ctx context.Context, id string) (*domain.Camera, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *CameraService) ListCameras(ctx context.Context, zoneID string, search string) ([]*domain.Camera, error) {
	if zoneID != "" {
		return s.repo.ListByZone(ctx, zoneID, search)
	}
	return s.repo.List(ctx, search)
}

func (s *CameraService) UpdateCamera(ctx context.Context, id string, req *domain.UpdateCameraRequest) (*domain.Camera, error) {
	camera, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if camera == nil {
		return nil, nil // Or NotFound error
	}

	if req.ZoneID != nil {
		camera.ZoneID = req.ZoneID
	}
	if req.Name != "" {
		camera.Name = req.Name
	}
	if req.IPAddress != "" {
		camera.IPAddress = req.IPAddress
	}
	if req.RTSPURL != "" {
		camera.RTSPURL = req.RTSPURL
	}
	if req.Status != "" {
		camera.Status = req.Status
	}
	if req.AIEnabled != nil {
		camera.AIEnabled = *req.AIEnabled
	}

	if err := s.repo.Update(ctx, camera); err != nil {
		return nil, err
	}
	return camera, nil
}

func (s *CameraService) DeleteCamera(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
