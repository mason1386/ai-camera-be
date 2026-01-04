package services

import (
	"context"

	"app/internal/core/ports"

	"github.com/google/uuid"
)

type PermissionService struct {
	repo ports.PermissionRepository
}

func NewPermissionService(repo ports.PermissionRepository) ports.PermissionService {
	return &PermissionService{repo: repo}
}

func (s *PermissionService) UpdateUserCameraPermissions(ctx context.Context, userID uuid.UUID, cameraIDs []uuid.UUID) error {
	// Simple strategy: Clear and re-grant (or more complex diff)
	// For now, let's assume we want to sync
	current, _ := s.repo.ListUserCameras(ctx, userID)

	// Poor man's sync
	for _, id := range current {
		_ = s.repo.RevokeCamera(ctx, userID, id)
	}
	for _, id := range cameraIDs {
		_ = s.repo.GrantCamera(ctx, userID, id)
	}
	return nil
}

func (s *PermissionService) UpdateUserZonePermissions(ctx context.Context, userID uuid.UUID, zoneIDs []uuid.UUID) error {
	current, _ := s.repo.ListUserZones(ctx, userID)
	for _, id := range current {
		_ = s.repo.RevokeZone(ctx, userID, id)
	}
	for _, id := range zoneIDs {
		_ = s.repo.GrantZone(ctx, userID, id)
	}
	return nil
}

func (s *PermissionService) GetUserPermissions(ctx context.Context, userID uuid.UUID) (*ports.UserPermissions, error) {
	cameras, _ := s.repo.ListUserCameras(ctx, userID)
	zones, _ := s.repo.ListUserZones(ctx, userID)
	return &ports.UserPermissions{
		CameraIDs: cameras,
		ZoneIDs:   zones,
	}, nil
}
