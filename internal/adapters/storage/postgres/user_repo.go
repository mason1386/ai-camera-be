package postgres

import (
	"context"

	"app/internal/core/domain"
	"app/internal/core/ports"

	"github.com/jackc/pgx/v5"
)

type UserRepository struct {
	db *PostgresDB
}

func NewUserRepository(db *PostgresDB) ports.UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Save(ctx context.Context, user *domain.User) error {
	query := `INSERT INTO users (username, email, password_hash, full_name, phone, role_id, status, updated_at) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7, NOW()) RETURNING id, created_at, updated_at`
	return r.db.Pool.QueryRow(ctx, query, user.Username, user.Email, user.PasswordHash, user.FullName, user.Phone, user.RoleID, user.Status).
		Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `SELECT id, username, email, password_hash, COALESCE(full_name, ''), COALESCE(phone, ''), role_id, COALESCE(status, 'active'), last_login_at, created_at, updated_at FROM users WHERE email = $1`
	return r.scanUser(r.db.Pool.QueryRow(ctx, query, email))
}

func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	query := `SELECT id, username, email, password_hash, COALESCE(full_name, ''), COALESCE(phone, ''), role_id, COALESCE(status, 'active'), last_login_at, created_at, updated_at FROM users WHERE username = $1`
	return r.scanUser(r.db.Pool.QueryRow(ctx, query, username))
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	query := `SELECT id, username, email, password_hash, COALESCE(full_name, ''), COALESCE(phone, ''), role_id, COALESCE(status, 'active'), last_login_at, created_at, updated_at FROM users WHERE id = $1`
	return r.scanUser(r.db.Pool.QueryRow(ctx, query, id))
}

func (r *UserRepository) scanUser(row pgx.Row) (*domain.User, error) {
	user := &domain.User{}
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.FullName, &user.Phone, &user.RoleID, &user.Status, &user.LastLoginAt, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) List(ctx context.Context, search string) ([]*domain.User, error) {
	query := `SELECT id, username, email, password_hash, COALESCE(full_name, ''), COALESCE(phone, ''), role_id, COALESCE(status, 'active'), last_login_at, created_at, updated_at FROM users`
	var args []interface{}

	if search != "" {
		query += ` WHERE username ILIKE $1 OR full_name ILIKE $1 OR email ILIKE $1`
		args = append(args, "%"+search+"%")
	}

	query += ` ORDER BY created_at DESC`

	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []*domain.User{}
	for rows.Next() {
		user := &domain.User{}
		err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.FullName, &user.Phone, &user.RoleID, &user.Status, &user.LastLoginAt, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (r *UserRepository) Update(ctx context.Context, user *domain.User) error {
	query := `UPDATE users SET full_name = $2, phone = $3, role_id = $4, status = $5, updated_at = NOW(), password_hash = $6 WHERE id = $1`
	_, err := r.db.Pool.Exec(ctx, query, user.ID, user.FullName, user.Phone, user.RoleID, user.Status, user.PasswordHash)
	return err
}

func (r *UserRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.Pool.Exec(ctx, "DELETE FROM users WHERE id = $1", id)
	return err
}
