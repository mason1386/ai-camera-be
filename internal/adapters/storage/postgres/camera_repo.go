package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"app/internal/adapters/storage/postgres/generated"
	"app/internal/core/domain"
	"app/internal/core/ports"
)

type CameraRepository struct {
	db *PostgresDB
}

func NewCameraRepository(db *PostgresDB) ports.CameraRepository {
	return &CameraRepository{
		db: db,
	}
}

func (r *CameraRepository) Save(ctx context.Context, camera *domain.Camera) error {
	params := generated.CreateCameraParams{
		ZoneID:    pgtype.UUID{Valid: false},
		Name:      camera.Name,
		IpAddress: pgtype.Text{String: camera.IPAddress, Valid: camera.IPAddress != ""},
		RtspUrl:   camera.RTSPURL,
		Status:    pgtype.Text{String: string(camera.Status), Valid: true},
		AiEnabled: pgtype.Bool{Bool: camera.AIEnabled, Valid: true},
	}

	if camera.ZoneID != nil {
		var zoneUUID pgtype.UUID
		if err := zoneUUID.Scan(*camera.ZoneID); err == nil {
			params.ZoneID = zoneUUID
		}
	}

	result, err := r.db.Query.CreateCamera(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to insert camera: %w", err)
	}

	camera.ID = uuidToString(result.ID)
	camera.CreatedAt = result.CreatedAt.Time
	camera.UpdatedAt = result.UpdatedAt.Time
	return nil
}

func (r *CameraRepository) GetByID(ctx context.Context, id string) (*domain.Camera, error) {
	var uuid pgtype.UUID
	if err := uuid.Scan(id); err != nil {
		return nil, fmt.Errorf("invalid uuid: %w", err)
	}

	row, err := r.db.Query.GetCamera(ctx, uuid)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return mapCameraEntity(row), nil
}

func (r *CameraRepository) List(ctx context.Context) ([]*domain.Camera, error) {
	rows, err := r.db.Query.ListCameras(ctx)
	if err != nil {
		return nil, err
	}
	cameras := make([]*domain.Camera, len(rows))
	for i, row := range rows {
		cameras[i] = mapCameraEntity(row)
	}
	return cameras, nil
}

func (r *CameraRepository) ListByZone(ctx context.Context, zoneID string) ([]*domain.Camera, error) {
	var uuid pgtype.UUID
	if err := uuid.Scan(zoneID); err != nil {
		return nil, fmt.Errorf("invalid uuid: %w", err)
	}

	rows, err := r.db.Query.ListCamerasByZone(ctx, pgtype.UUID{Bytes: uuid.Bytes, Valid: true})
	if err != nil {
		return nil, err
	}
	cameras := make([]*domain.Camera, len(rows))
	for i, row := range rows {
		cameras[i] = mapCameraEntity(row)
	}
	return cameras, nil
}

func (r *CameraRepository) Update(ctx context.Context, camera *domain.Camera) error {
	var idUUID pgtype.UUID
	if err := idUUID.Scan(camera.ID); err != nil {
		return fmt.Errorf("invalid uuid: %w", err)
	}

	params := generated.UpdateCameraParams{
		ID:        idUUID,
		Name:      camera.Name,
		IpAddress: pgtype.Text{String: camera.IPAddress, Valid: camera.IPAddress != ""},
		RtspUrl:   camera.RTSPURL,
		Status:    pgtype.Text{String: string(camera.Status), Valid: true},
		AiEnabled: pgtype.Bool{Bool: camera.AIEnabled, Valid: true},
		ZoneID:    pgtype.UUID{Valid: false},
	}

	if camera.ZoneID != nil {
		var zoneUUID pgtype.UUID
		if err := zoneUUID.Scan(*camera.ZoneID); err == nil {
			params.ZoneID = zoneUUID
		}
	}

	_, err := r.db.Query.UpdateCamera(ctx, params)
	return err
}

func (r *CameraRepository) Delete(ctx context.Context, id string) error {
	var uuid pgtype.UUID
	if err := uuid.Scan(id); err != nil {
		return fmt.Errorf("invalid uuid: %w", err)
	}
	return r.db.Query.DeleteCamera(ctx, uuid)
}

func mapCameraEntity(row generated.Camera) *domain.Camera {
	var zoneIDPtr *string
	if row.ZoneID.Valid {
		s := uuidToString(row.ZoneID)
		zoneIDPtr = &s
	}

	return &domain.Camera{
		ID:        uuidToString(row.ID),
		ZoneID:    zoneIDPtr,
		Name:      row.Name,
		IPAddress: row.IpAddress.String,
		RTSPURL:   row.RtspUrl,
		Status:    domain.CameraStatus(row.Status.String),
		AIEnabled: row.AiEnabled.Bool,
		CreatedAt: row.CreatedAt.Time,
		UpdatedAt: row.UpdatedAt.Time,
	}
}

func uuidToString(uuid pgtype.UUID) string {
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid.Bytes[0:4], uuid.Bytes[4:6], uuid.Bytes[6:8], uuid.Bytes[8:10], uuid.Bytes[10:16])
}
