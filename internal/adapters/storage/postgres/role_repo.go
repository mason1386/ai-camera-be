package postgres

import (
	"context"

	"app/internal/core/domain"
	"app/internal/core/ports"

	"github.com/jackc/pgx/v5"
)

type RoleRepository struct {
	db *PostgresDB
}

func NewRoleRepository(db *PostgresDB) ports.RoleRepository {
	return &RoleRepository{db: db}
}

func (r *RoleRepository) Create(ctx context.Context, role *domain.Role) error {
	query := `INSERT INTO roles (name, description, permissions, is_system) VALUES ($1, $2, $3, $4) RETURNING id, created_at`
	return r.db.Pool.QueryRow(ctx, query, role.Name, role.Description, role.Permissions, role.IsSystem).
		Scan(&role.ID, &role.CreatedAt)
}

func (r *RoleRepository) GetByID(ctx context.Context, id string) (*domain.Role, error) {
	query := `SELECT id, name, description, permissions, is_system, created_at FROM roles WHERE id = $1`
	role := &domain.Role{}
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(&role.ID, &role.Name, &role.Description, &role.Permissions, &role.IsSystem, &role.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return role, nil
}

func (r *RoleRepository) List(ctx context.Context, search string) ([]*domain.Role, error) {
	query := `SELECT id, name, description, created_at, updated_at FROM roles`
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

	var roles []*domain.Role
	for rows.Next() {
		role := &domain.Role{}
		err := rows.Scan(&role.ID, &role.Name, &role.Description, &role.CreatedAt, &role.UpdatedAt)
		if err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}
	return roles, nil
}

func (r *RoleRepository) Update(ctx context.Context, role *domain.Role) error {
	query := `UPDATE roles SET name = $2, description = $3, permissions = $4, is_system = $5 WHERE id = $1`
	_, err := r.db.Pool.Exec(ctx, query, role.ID, role.Name, role.Description, role.Permissions, role.IsSystem)
	return err
}

func (r *RoleRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.Pool.Exec(ctx, "DELETE FROM roles WHERE id = $1 AND is_system = FALSE", id)
	return err
}
