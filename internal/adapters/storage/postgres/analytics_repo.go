package postgres

import (
	"context"
	"time"

	"app/internal/core/domain"
	"app/internal/core/ports"

	"github.com/google/uuid"
)

type AnalyticsRepository struct {
	db *PostgresDB
}

func NewAnalyticsRepository(db *PostgresDB) ports.AnalyticsRepository {
	return &AnalyticsRepository{db: db}
}

func (r *AnalyticsRepository) CreateRecognitionLog(ctx context.Context, log *domain.RecognitionLog) error {
	query := `INSERT INTO recognition_logs (camera_id, identity_id, snapshot_url, face_crop_url, confidence, label, occurred_at) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id, created_at`
	return r.db.Pool.QueryRow(ctx, query, log.CameraID, log.IdentityID, log.SnapshotURL, log.FaceCropURL, log.Confidence, log.Label, log.OccurredAt).
		Scan(&log.ID, &log.CreatedAt)
}

func (r *AnalyticsRepository) ListRecognitionLogs(ctx context.Context, identityID *uuid.UUID, cameraID *uuid.UUID, from, to *time.Time, limit, offset int32) ([]*domain.RecognitionLog, error) {
	query := `SELECT rl.id, rl.camera_id, rl.identity_id, rl.snapshot_url, rl.face_crop_url, rl.confidence, rl.label, rl.occurred_at, rl.created_at, i.full_name as identity_name, c.name as camera_name
	          FROM recognition_logs rl
	          JOIN identities i ON rl.identity_id = i.id
	          JOIN cameras c ON rl.camera_id = c.id
	          WHERE ($1::uuid IS NULL OR rl.identity_id = $1)
	            AND ($2::uuid IS NULL OR rl.camera_id = $2)
	            AND ($3::timestamp IS NULL OR rl.occurred_at >= $3)
	            AND ($4::timestamp IS NULL OR rl.occurred_at <= $4)
	          ORDER BY rl.occurred_at DESC
	          LIMIT $5 OFFSET $6`

	rows, err := r.db.Pool.Query(ctx, query, identityID, cameraID, from, to, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []*domain.RecognitionLog
	for rows.Next() {
		log := &domain.RecognitionLog{}
		err := rows.Scan(
			&log.ID, &log.CameraID, &log.IdentityID,
			&log.SnapshotURL, &log.FaceCropURL, &log.Confidence, &log.Label,
			&log.OccurredAt, &log.CreatedAt,
			&log.IdentityName, &log.CameraName,
		)
		if err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}
	return logs, nil
}

func (r *AnalyticsRepository) ListAttendanceRecords(ctx context.Context, identityID *uuid.UUID, from, to *time.Time, status *domain.AttendanceStatus, limit, offset int32) ([]*domain.AttendanceRecord, error) {
	query := `SELECT ar.id, ar.identity_id, ar.date, ar.check_in, ar.check_out, ar.work_hours, ar.status, ar.created_at, ar.updated_at, i.full_name as identity_name
	          FROM attendance_records ar
	          JOIN identities i ON ar.identity_id = i.id
	          WHERE ($1::uuid IS NULL OR ar.identity_id = $1)
	            AND ($2::date IS NULL OR ar.date >= $2)
	            AND ($3::date IS NULL OR ar.date <= $3)
	            AND ($4::attendance_status IS NULL OR ar.status = $4)
	          ORDER BY ar.date DESC, ar.check_in DESC
	          LIMIT $5 OFFSET $6`

	rows, err := r.db.Pool.Query(ctx, query, identityID, from, to, status, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []*domain.AttendanceRecord
	for rows.Next() {
		record := &domain.AttendanceRecord{}
		err := rows.Scan(
			&record.ID, &record.IdentityID, &record.Date, &record.CheckIn,
			&record.CheckOut, &record.WorkHours, &record.Status,
			&record.CreatedAt, &record.UpdatedAt, &record.IdentityName,
		)
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}
	return records, nil
}

func (r *AnalyticsRepository) GetAttendanceStats(ctx context.Context, date time.Time) (map[string]int64, error) {
	query := `SELECT status, COUNT(*) as count FROM attendance_records WHERE date = $1 GROUP BY status`
	rows, err := r.db.Pool.Query(ctx, query, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stats := make(map[string]int64)
	for rows.Next() {
		var status string
		var count int64
		if err := rows.Scan(&status, &count); err != nil {
			return nil, err
		}
		stats[status] = count
	}
	return stats, nil
}
