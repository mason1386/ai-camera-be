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

type ZoneRepository struct {
	db *PostgresDB
}

func NewZoneRepository(db *PostgresDB) ports.ZoneRepository {
	return &ZoneRepository{db: db}
}

func (r *ZoneRepository) Save(ctx context.Context, zone *domain.Zone) error {
	params := generated.CreateZoneParams{
		Name:        zone.Name,
		Description: pgtype.Text{String: zone.Description, Valid: zone.Description != ""},
	}

	result, err := r.db.Query.CreateZone(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to create zone: %w", err)
	}

	zone.ID = fmt.Sprintf("%x-%x-%x-%x-%x", result.ID.Bytes[0:4], result.ID.Bytes[4:6], result.ID.Bytes[6:8], result.ID.Bytes[8:10], result.ID.Bytes[10:16])
	zone.CreatedAt = result.CreatedAt.Time
	return nil
}

func (r *ZoneRepository) GetByID(ctx context.Context, id string) (*domain.Zone, error) {
	var uuid pgtype.UUID
	if err := uuid.Scan(id); err != nil {
		return nil, fmt.Errorf("invalid uuid: %w", err)
	}

	row, err := r.db.Query.GetZone(ctx, uuid)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return mapZoneEntity(row), nil
}

func (r *ZoneRepository) List(ctx context.Context) ([]*domain.Zone, error) {
	rows, err := r.db.Query.ListZones(ctx)
	if err != nil {
		return nil, err
	}

	zones := make([]*domain.Zone, len(rows))
	for i, row := range rows {
		zones[i] = mapZoneEntity(row)
	}
	return zones, nil
}

func (r *ZoneRepository) Update(ctx context.Context, zone *domain.Zone) error {
	var uuid pgtype.UUID
	if err := uuid.Scan(zone.ID); err != nil {
		return fmt.Errorf("invalid uuid: %w", err)
	}

	params := generated.UpdateZoneParams{
		ID:          uuid,
		Name:        zone.Name,
		Description: pgtype.Text{String: zone.Description, Valid: zone.Description != ""},
	}

	_, err := r.db.Query.UpdateZone(ctx, params)
	return err
}

func (r *ZoneRepository) Delete(ctx context.Context, id string) error {
	var uuid pgtype.UUID
	if err := uuid.Scan(id); err != nil {
		return fmt.Errorf("invalid uuid: %w", err)
	}
	return r.db.Query.DeleteZone(ctx, uuid)
}

func mapZoneEntity(row generated.Zone) *domain.Zone {
	idStr := fmt.Sprintf("%x-%x-%x-%x-%x", row.ID.Bytes[0:4], row.ID.Bytes[4:6], row.ID.Bytes[6:8], row.ID.Bytes[8:10], row.ID.Bytes[10:16])
	return &domain.Zone{
		ID:          idStr,
		Name:        row.Name,
		Description: row.Description.String,
		CreatedAt:   row.CreatedAt.Time,
	}
}
