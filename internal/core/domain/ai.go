package domain

import (
	"time"

	"github.com/google/uuid"
)

type EventType string
type EventStatus string

const (
	EventTypePerson    EventType = "person"
	EventTypeVehicle   EventType = "vehicle"
	EventTypeFace      EventType = "face"
	EventTypeIntrusion EventType = "intrusion"
	EventTypeLoitering EventType = "loitering"
	EventTypeCrowd     EventType = "crowd"
	EventTypeFire      EventType = "fire"
	EventTypeOther     EventType = "other"

	EventStatusNew        EventStatus = "new"
	EventStatusProcessing EventStatus = "processing"
	EventStatusResolved   EventStatus = "resolved"
	EventStatusIgnored    EventStatus = "ignored"
)

type AIConfig struct {
	ID            uuid.UUID   `json:"id"`
	CameraID      uuid.UUID   `json:"camera_id"`
	AIEnabled     bool        `json:"ai_enabled"`
	AITypes       []EventType `json:"ai_types"`
	ROIZones      any         `json:"roi_zones"`    // JSONB
	ActiveHours   any         `json:"active_hours"` // JSONB
	Sensitivity   int         `json:"sensitivity"`
	MinConfidence int         `json:"min_confidence"`
	CreatedAt     time.Time   `json:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at"`
}

type AIEvent struct {
	ID          uuid.UUID      `json:"id"`
	CameraID    uuid.UUID      `json:"camera_id"`
	EventType   EventType      `json:"event_type"`
	Confidence  float64        `json:"confidence"`
	SnapshotURL string         `json:"snapshot_url"`
	Metadata    map[string]any `json:"metadata"`
	Status      EventStatus    `json:"status"`
	ResolvedBy  *uuid.UUID     `json:"resolved_by"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

type DashboardStats struct {
	TotalCameras       int        `json:"total_cameras"`
	OnlineCameras      int        `json:"online_cameras"`
	OfflineCameras     int        `json:"offline_cameras"`
	MaintenanceCameras int        `json:"maintenance_cameras"`
	TodayEvents        int        `json:"today_events"`
	UnresolvedEvents   int        `json:"unresolved_events"`
	RecentEvents       []*AIEvent `json:"recent_events"`
}
