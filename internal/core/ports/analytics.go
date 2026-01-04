package ports

import (
	"context"
	"time"

	"app/internal/core/domain"

	"github.com/google/uuid"
)

type AnalyticsRepository interface {
	CreateRecognitionLog(ctx context.Context, log *domain.RecognitionLog) error
	ListRecognitionLogs(ctx context.Context, identityID *uuid.UUID, cameraID *uuid.UUID, from, to *time.Time, limit, offset int32) ([]*domain.RecognitionLog, error)

	ListAttendanceRecords(ctx context.Context, identityID *uuid.UUID, from, to *time.Time, status *domain.AttendanceStatus, limit, offset int32) ([]*domain.AttendanceRecord, error)
	GetAttendanceStats(ctx context.Context, date time.Time) (map[string]int64, error)
}

type AnalyticsService interface {
	ListRecognitionLogs(ctx context.Context, filter *RecognitionFilter) ([]*domain.RecognitionLog, error)
	ListAttendance(ctx context.Context, filter *AttendanceFilter) ([]*domain.AttendanceRecord, error)
	GetDailyAttendanceSummary(ctx context.Context, date time.Time) (any, error)
}

type RecognitionFilter struct {
	IdentityID *uuid.UUID
	CameraID   *uuid.UUID
	FromDate   *time.Time
	ToDate     *time.Time
	Limit      int32
	Offset     int32
}

type AttendanceFilter struct {
	IdentityID *uuid.UUID
	FromDate   *time.Time
	ToDate     *time.Time
	Status     *domain.AttendanceStatus
	Limit      int32
	Offset     int32
}
