package domain

import (
	"time"

	"github.com/google/uuid"
)

type UserCameraPermission struct {
	UserID    uuid.UUID `json:"user_id"`
	CameraID  uuid.UUID `json:"camera_id"`
	CreatedAt time.Time `json:"created_at"`
}

type UserZonePermission struct {
	UserID    uuid.UUID `json:"user_id"`
	ZoneID    uuid.UUID `json:"zone_id"`
	CreatedAt time.Time `json:"created_at"`
}
