package services

import (
	"context"
	"time"

	"app/internal/core/domain"
	"app/internal/core/ports"
)

type AnalyticsService struct {
	repo ports.AnalyticsRepository
}

func NewAnalyticsService(repo ports.AnalyticsRepository) ports.AnalyticsService {
	return &AnalyticsService{repo: repo}
}

func (s *AnalyticsService) ListRecognitionLogs(ctx context.Context, filter *ports.RecognitionFilter) ([]*domain.RecognitionLog, error) {
	return s.repo.ListRecognitionLogs(ctx, filter.IdentityID, filter.CameraID, filter.FromDate, filter.ToDate, filter.Limit, filter.Offset)
}

func (s *AnalyticsService) ListAttendance(ctx context.Context, filter *ports.AttendanceFilter) ([]*domain.AttendanceRecord, error) {
	return s.repo.ListAttendanceRecords(ctx, filter.IdentityID, filter.FromDate, filter.ToDate, filter.Status, filter.Limit, filter.Offset)
}

func (s *AnalyticsService) GetDailyAttendanceSummary(ctx context.Context, date time.Time) (any, error) {
	return s.repo.GetAttendanceStats(ctx, date)
}
