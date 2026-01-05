package ports

import (
	"context"

	"github.com/google/uuid"
)

type PermissionRepository interface {
	GrantCamera(ctx context.Context, userID, cameraID uuid.UUID) error
	RevokeCamera(ctx context.Context, userID, cameraID uuid.UUID) error
	ListUserCameras(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error)

	GrantZone(ctx context.Context, userID, zoneID uuid.UUID) error
	RevokeZone(ctx context.Context, userID, zoneID uuid.UUID) error
	ListUserZones(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error)
}

type UserPermissions struct {
	CameraIDs []uuid.UUID `json:"camera_ids"`
	ZoneIDs   []uuid.UUID `json:"zone_ids"`
}

type PermissionService interface {
	UpdateUserCameraPermissions(ctx context.Context, userID uuid.UUID, cameraIDs []uuid.UUID) error
	UpdateUserZonePermissions(ctx context.Context, userID uuid.UUID, zoneIDs []uuid.UUID) error
	GetUserPermissions(ctx context.Context, userID uuid.UUID) (*UserPermissions, error)
}
