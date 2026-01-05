package ports

import (
	"context"

	"app/internal/core/domain"

	"github.com/google/uuid"
)

type AuditRepository interface {
	CreateLog(ctx context.Context, log *domain.AuditLog) error
	ListLogs(ctx context.Context, userID *uuid.UUID, action *string, tableName *string, limit, offset int32) ([]*domain.AuditLog, error)
}

type AuditService interface {
	LogAction(ctx context.Context, log *domain.AuditLog) error
	ListLogs(ctx context.Context, filter *AuditFilter) ([]*domain.AuditLog, error)
}

type AuditFilter struct {
	UserID    *uuid.UUID
	Action    *string
	TableName *string
	Limit     int32
	Offset    int32
}
