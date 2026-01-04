package postgres

import (
	"context"

	"app/internal/core/domain"
	"app/internal/core/ports"

	"github.com/google/uuid"
)

type AuditRepository struct {
	db *PostgresDB
}

func NewAuditRepository(db *PostgresDB) ports.AuditRepository {
	return &AuditRepository{db: db}
}

func (r *AuditRepository) CreateLog(ctx context.Context, log *domain.AuditLog) error {
	query := `INSERT INTO audit_logs (user_id, action, table_name, record_id, old_value, new_value, ip_address, user_agent, created_at) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW()) RETURNING id`
	return r.db.Pool.QueryRow(ctx, query, log.UserID, log.Action, log.TableName, log.RecordID, log.OldValue, log.NewValue, log.IPAddress, log.UserAgent).
		Scan(&log.ID)
}

func (r *AuditRepository) ListLogs(ctx context.Context, userID *uuid.UUID, action *string, tableName *string, limit, offset int32) ([]*domain.AuditLog, error) {
	query := `SELECT al.id, al.user_id, al.action, al.table_name, al.record_id, al.old_value, al.new_value, al.ip_address, al.user_agent, al.created_at, u.username
	          FROM audit_logs al
	          LEFT JOIN users u ON al.user_id = u.id
	          WHERE ($1::uuid IS NULL OR al.user_id = $1)
	            AND ($2::text IS NULL OR al.action = $2)
	            AND ($3::text IS NULL OR al.table_name = $3)
	          ORDER BY al.created_at DESC
	          LIMIT $4 OFFSET $5`

	rows, err := r.db.Pool.Query(ctx, query, userID, action, tableName, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []*domain.AuditLog
	for rows.Next() {
		log := &domain.AuditLog{}
		err := rows.Scan(
			&log.ID, &log.UserID, &log.Action, &log.TableName,
			&log.RecordID, &log.OldValue, &log.NewValue,
			&log.IPAddress, &log.UserAgent, &log.CreatedAt, &log.Username,
		)
		if err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}
	return logs, nil
}
