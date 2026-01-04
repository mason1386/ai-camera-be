package postgres

import (
	"context"

	"app/internal/core/domain"
	"app/internal/core/ports"

	"github.com/jackc/pgx/v5"
)

type ZoneRepository struct {
	db *PostgresDB
}

func NewZoneRepository(db *PostgresDB) ports.ZoneRepository {
	return &ZoneRepository{db: db}
}

func (r *ZoneRepository) Save(ctx context.Context, zone *domain.Zone) error {
	query := `INSERT INTO zones (name, description) VALUES ($1, $2) RETURNING id, created_at`
	return r.db.Pool.QueryRow(ctx, query, zone.Name, zone.Description).Scan(&zone.ID, &zone.CreatedAt)
}

func (r *ZoneRepository) GetByID(ctx context.Context, id string) (*domain.Zone, error) {
	query := `SELECT id, name, COALESCE(description, ''), created_at FROM zones WHERE id = $1`
	zone := &domain.Zone{}
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(&zone.ID, &zone.Name, &zone.Description, &zone.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return zone, nil
}

func (r *ZoneRepository) List(ctx context.Context, search string) ([]*domain.Zone, error) {
	query := `SELECT id, name, COALESCE(description, ''), created_at FROM zones`
	var args []interface{}

	if search != "" {
		query += ` WHERE name ILIKE $1 OR description ILIKE $1`
		args = append(args, "%"+search+"%")
	}

	query += ` ORDER BY created_at DESC`

	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	zones := []*domain.Zone{} // Initialize as empty slice
	for rows.Next() {
		zone := &domain.Zone{}
		err := rows.Scan(&zone.ID, &zone.Name, &zone.Description, &zone.CreatedAt)
		if err != nil {
			return nil, err
		}
		zones = append(zones, zone)
	}
	return zones, nil
}

func (r *ZoneRepository) Update(ctx context.Context, zone *domain.Zone) error {
	query := `UPDATE zones SET name = $2, description = $3 WHERE id = $1`
	_, err := r.db.Pool.Exec(ctx, query, zone.ID, zone.Name, zone.Description)
	return err
}

func (r *ZoneRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.Pool.Exec(ctx, "DELETE FROM zones WHERE id = $1", id)
	return err
}
