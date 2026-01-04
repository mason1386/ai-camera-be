package postgres

import (
	"context"

	"app/internal/core/domain"
	"app/internal/core/ports"

	"github.com/jackc/pgx/v5"
)

type CameraRepository struct {
	db *PostgresDB
}

func NewCameraRepository(db *PostgresDB) ports.CameraRepository {
	return &CameraRepository{db: db}
}

func (r *CameraRepository) Save(ctx context.Context, camera *domain.Camera) error {
	query := `INSERT INTO cameras (zone_id, name, ip_address, rtsp_url, status, ai_enabled, updated_at) 
	          VALUES ($1, $2, $3, $4, $5, $6, NOW()) RETURNING id, created_at, updated_at`
	return r.db.Pool.QueryRow(ctx, query, camera.ZoneID, camera.Name, camera.IPAddress, camera.RTSPURL, camera.Status, camera.AIEnabled).
		Scan(&camera.ID, &camera.CreatedAt, &camera.UpdatedAt)
}

func (r *CameraRepository) GetByID(ctx context.Context, id string) (*domain.Camera, error) {
	query := `SELECT id, zone_id, name, ip_address, rtsp_url, status, ai_enabled, created_at, updated_at FROM cameras WHERE id = $1`
	camera := &domain.Camera{}
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&camera.ID, &camera.ZoneID, &camera.Name, &camera.IPAddress, &camera.RTSPURL, &camera.Status, &camera.AIEnabled, &camera.CreatedAt, &camera.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return camera, nil
}

func (r *CameraRepository) List(ctx context.Context, search string) ([]*domain.Camera, error) {
	query := `SELECT id, zone_id, name, ip_address, rtsp_url, status, ai_enabled, created_at, updated_at FROM cameras`
	var args []interface{}

	if search != "" {
		query += ` WHERE name ILIKE $1 OR ip_address ILIKE $1`
		args = append(args, "%"+search+"%")
	}

	query += ` ORDER BY created_at DESC`

	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cameras := []*domain.Camera{}
	for rows.Next() {
		camera := &domain.Camera{}
		err := rows.Scan(&camera.ID, &camera.ZoneID, &camera.Name, &camera.IPAddress, &camera.RTSPURL, &camera.Status, &camera.AIEnabled, &camera.CreatedAt, &camera.UpdatedAt)
		if err != nil {
			return nil, err
		}
		cameras = append(cameras, camera)
	}
	return cameras, nil
}

func (r *CameraRepository) ListByZone(ctx context.Context, zoneID string, search string) ([]*domain.Camera, error) {
	query := `SELECT id, zone_id, name, ip_address, rtsp_url, status, ai_enabled, created_at, updated_at FROM cameras WHERE zone_id = $1`
	args := []interface{}{zoneID}

	if search != "" {
		query += ` AND (name ILIKE $2 OR ip_address ILIKE $2)`
		args = append(args, "%"+search+"%")
	}

	query += ` ORDER BY created_at DESC`

	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cameras := []*domain.Camera{}
	for rows.Next() {
		camera := &domain.Camera{}
		err := rows.Scan(&camera.ID, &camera.ZoneID, &camera.Name, &camera.IPAddress, &camera.RTSPURL, &camera.Status, &camera.AIEnabled, &camera.CreatedAt, &camera.UpdatedAt)
		if err != nil {
			return nil, err
		}
		cameras = append(cameras, camera)
	}
	return cameras, nil
}

func (r *CameraRepository) Update(ctx context.Context, camera *domain.Camera) error {
	query := `UPDATE cameras SET zone_id = $2, name = $3, ip_address = $4, rtsp_url = $5, status = $6, ai_enabled = $7, updated_at = NOW() WHERE id = $1`
	_, err := r.db.Pool.Exec(ctx, query, camera.ID, camera.ZoneID, camera.Name, camera.IPAddress, camera.RTSPURL, camera.Status, camera.AIEnabled)
	return err
}

func (r *CameraRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.Pool.Exec(ctx, "DELETE FROM cameras WHERE id = $1", id)
	return err
}
