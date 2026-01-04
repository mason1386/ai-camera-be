package domain

import (
	"time"

	"github.com/google/uuid"
)

type AttendanceStatus string

const (
	AttendanceLate       AttendanceStatus = "late"
	AttendanceOnTime     AttendanceStatus = "on_time"
	AttendanceAbsent     AttendanceStatus = "absent"
	AttendanceEarlyLeave AttendanceStatus = "early_leave"
)

type RecognitionLog struct {
	ID          uuid.UUID `json:"id"`
	CameraID    uuid.UUID `json:"camera_id"`
	IdentityID  uuid.UUID `json:"identity_id"`
	SnapshotURL string    `json:"snapshot_url"`
	FaceCropURL string    `json:"face_crop_url"`
	Confidence  float64   `json:"confidence"`
	Label       string    `json:"label"`
	OccurredAt  time.Time `json:"occurred_at"`
	CreatedAt   time.Time `json:"created_at"`

	// Join fields
	IdentityName string `json:"identity_name,omitempty"`
	CameraName   string `json:"camera_name,omitempty"`
}

type AttendanceRecord struct {
	ID         uuid.UUID        `json:"id"`
	IdentityID uuid.UUID        `json:"identity_id"`
	Date       time.Time        `json:"date"`
	CheckIn    *time.Time       `json:"check_in"`
	CheckOut   *time.Time       `json:"check_out"`
	WorkHours  float64          `json:"work_hours"`
	Status     AttendanceStatus `json:"status"`
	CreatedAt  time.Time        `json:"created_at"`
	UpdatedAt  time.Time        `json:"updated_at"`

	// Join fields
	IdentityName string `json:"identity_name,omitempty"`
}
