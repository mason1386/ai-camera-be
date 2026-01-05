package postgres

import (
	"context"

	"app/internal/core/ports"

	"github.com/google/uuid"
)

type PermissionRepository struct {
	db *PostgresDB
}

func NewPermissionRepository(db *PostgresDB) ports.PermissionRepository {
	return &PermissionRepository{db: db}
}

func (r *PermissionRepository) GrantCamera(ctx context.Context, userID, cameraID uuid.UUID) error {
	query := `INSERT INTO user_camera_permissions (user_id, camera_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`
	_, err := r.db.Pool.Exec(ctx, query, userID, cameraID)
	return err
}

func (r *PermissionRepository) RevokeCamera(ctx context.Context, userID, cameraID uuid.UUID) error {
	query := `DELETE FROM user_camera_permissions WHERE user_id = $1 AND camera_id = $2`
	_, err := r.db.Pool.Exec(ctx, query, userID, cameraID)
	return err
}

func (r *PermissionRepository) ListUserCameras(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error) {
	query := `SELECT camera_id FROM user_camera_permissions WHERE user_id = $1`
	rows, err := r.db.Pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []uuid.UUID
	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func (r *PermissionRepository) GrantZone(ctx context.Context, userID, zoneID uuid.UUID) error {
	query := `INSERT INTO user_zone_permissions (user_id, zone_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`
	_, err := r.db.Pool.Exec(ctx, query, userID, zoneID)
	return err
}

func (r *PermissionRepository) RevokeZone(ctx context.Context, userID, zoneID uuid.UUID) error {
	query := `DELETE FROM user_zone_permissions WHERE user_id = $1 AND zone_id = $2`
	_, err := r.db.Pool.Exec(ctx, query, userID, zoneID)
	return err
}

func (r *PermissionRepository) ListUserZones(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error) {
	query := `SELECT zone_id FROM user_zone_permissions WHERE user_id = $1`
	rows, err := r.db.Pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []uuid.UUID
	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}
