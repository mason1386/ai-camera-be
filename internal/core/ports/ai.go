package ports

import (
	"context"
	"time"

	"app/internal/core/domain"

	"github.com/google/uuid"
)

type AIRepository interface {
	GetConfigByCamera(ctx context.Context, cameraID uuid.UUID) (*domain.AIConfig, error)
	SaveConfig(ctx context.Context, config *domain.AIConfig) error

	CreateEvent(ctx context.Context, event *domain.AIEvent) (*domain.AIEvent, error)
	ListEvents(ctx context.Context, cameraID *uuid.UUID, eventType *domain.EventType, status *domain.EventStatus, from, to *time.Time, limit, offset int32) ([]*domain.AIEvent, error)
	UpdateEventStatus(ctx context.Context, id uuid.UUID, status domain.EventStatus, resolvedBy *uuid.UUID) (*domain.AIEvent, error)

	GetDashboardStats(ctx context.Context) (total, online, offline, maintenance int64, err error)
	GetTodayEventsCount(ctx context.Context) (int64, error)
}

type AIService interface {
	GetConfig(ctx context.Context, cameraID uuid.UUID) (*domain.AIConfig, error)
	UpdateConfig(ctx context.Context, req *domain.AIConfig) error

	CreateEvent(ctx context.Context, event *domain.AIEvent) (*domain.AIEvent, error)
	ListEvents(ctx context.Context, filter *EventFilter) ([]*domain.AIEvent, error)
	UpdateEventStatus(ctx context.Context, id uuid.UUID, status domain.EventStatus, resolvedBy uuid.UUID) (*domain.AIEvent, error)

	GetDashboardStats(ctx context.Context) (*domain.DashboardStats, error)
}

type EventFilter struct {
	CameraID  *uuid.UUID
	EventType *domain.EventType
	Status    *domain.EventStatus
	FromDate  *time.Time
	ToDate    *time.Time
	Limit     int32
	Offset    int32
}
