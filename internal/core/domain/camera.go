package domain

import (
	"time"
)

type CameraStatus string

const (
	CameraStatusOnline      CameraStatus = "online"
	CameraStatusOffline     CameraStatus = "offline"
	CameraStatusMaintenance CameraStatus = "maintenance"
)

type Camera struct {
	ID        string       `json:"id"`
	ZoneID    *string      `json:"zone_id"`
	Name      string       `json:"name"`
	IPAddress string       `json:"ip_address"`
	RTSPURL   string       `json:"rtsp_url"`
	Status    CameraStatus `json:"status"`
	AIEnabled bool         `json:"ai_enabled"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
}

type CreateCameraRequest struct {
	ZoneID    *string `json:"zone_id"`
	Name      string  `json:"name" binding:"required"`
	IPAddress string  `json:"ip_address"`
	RTSPURL   string  `json:"rtsp_url" binding:"required"`
	AIEnabled bool    `json:"ai_enabled"`
}

type UpdateCameraRequest struct {
	ZoneID    *string      `json:"zone_id"`
	Name      string       `json:"name"`
	IPAddress string       `json:"ip_address"`
	RTSPURL   string       `json:"rtsp_url"`
	Status    CameraStatus `json:"status"`
	AIEnabled *bool        `json:"ai_enabled"` // Pointer to distinguish false vs nil
}
