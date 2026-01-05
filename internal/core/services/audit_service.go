package services

import (
	"context"

	"app/internal/core/domain"
	"app/internal/core/ports"
)

type AuditService struct {
	repo ports.AuditRepository
}

func NewAuditService(repo ports.AuditRepository) ports.AuditService {
	return &AuditService{repo: repo}
}

func (s *AuditService) LogAction(ctx context.Context, log *domain.AuditLog) error {
	return s.repo.CreateLog(ctx, log)
}

func (s *AuditService) ListLogs(ctx context.Context, filter *ports.AuditFilter) ([]*domain.AuditLog, error) {
	return s.repo.ListLogs(ctx, filter.UserID, filter.Action, filter.TableName, filter.Limit, filter.Offset)
}
