package services

import (
	"context"

	"app/internal/core/domain"
	"app/internal/core/ports"

	"github.com/google/uuid"
)

type AIService struct {
	repo ports.AIRepository
}

func NewAIService(repo ports.AIRepository) ports.AIService {
	return &AIService{repo: repo}
}

func (s *AIService) GetConfig(ctx context.Context, cameraID uuid.UUID) (*domain.AIConfig, error) {
	return s.repo.GetConfigByCamera(ctx, cameraID)
}

func (s *AIService) UpdateConfig(ctx context.Context, req *domain.AIConfig) error {
	return s.repo.SaveConfig(ctx, req)
}

func (s *AIService) CreateEvent(ctx context.Context, event *domain.AIEvent) (*domain.AIEvent, error) {
	if event.Status == "" {
		event.Status = domain.EventStatusNew
	}
	return s.repo.CreateEvent(ctx, event)
}

func (s *AIService) ListEvents(ctx context.Context, filter *ports.EventFilter) ([]*domain.AIEvent, error) {
	return s.repo.ListEvents(ctx, filter.CameraID, filter.EventType, filter.Status, filter.FromDate, filter.ToDate, filter.Limit, filter.Offset)
}

func (s *AIService) UpdateEventStatus(ctx context.Context, id uuid.UUID, status domain.EventStatus, resolvedBy uuid.UUID) (*domain.AIEvent, error) {
	return s.repo.UpdateEventStatus(ctx, id, status, &resolvedBy)
}

func (s *AIService) GetDashboardStats(ctx context.Context) (*domain.DashboardStats, error) {
	total, online, offline, maintenance, err := s.repo.GetDashboardStats(ctx)
	if err != nil {
		return nil, err
	}

	todayEvents, _ := s.repo.GetTodayEventsCount(ctx)

	// Fetch recent events (optional, can be separate call or limit filter)
	recent, _ := s.repo.ListEvents(ctx, nil, nil, nil, nil, nil, 5, 0)

	return &domain.DashboardStats{
		TotalCameras:       int(total),
		OnlineCameras:      int(online),
		OfflineCameras:     int(offline),
		MaintenanceCameras: int(maintenance),
		TodayEvents:        int(todayEvents),
		RecentEvents:       recent,
	}, nil
}
