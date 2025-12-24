package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"app/internal/adapters/storage/postgres/generated"
	"app/internal/core/domain"
	"app/internal/core/ports"
)

type UserRepository struct {
	db *PostgresDB
}

func NewUserRepository(db *PostgresDB) ports.UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Save(ctx context.Context, user *domain.User) error {
	params := generated.CreateUserParams{
		Username:     user.Username,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		FullName:     pgtype.Text{String: user.FullName, Valid: user.FullName != ""},
		Status:       pgtype.Text{String: string(user.Status), Valid: true},
	}

	result, err := r.db.Query.CreateUser(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to save user: %w", err)
	}

	// UUID to String conversion
	user.ID = fmt.Sprintf("%x-%x-%x-%x-%x", result.ID.Bytes[0:4], result.ID.Bytes[4:6], result.ID.Bytes[6:8], result.ID.Bytes[8:10], result.ID.Bytes[10:16])
	user.CreatedAt = result.CreatedAt.Time
	user.UpdatedAt = result.UpdatedAt.Time
	return nil
}

func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	row, err := r.db.Query.GetUserByUsername(ctx, username)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user by username: %w", err)
	}
	return mapUserEntity(row), nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	row, err := r.db.Query.GetUserByEmail(ctx, email)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	return mapUserEntity(row), nil
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	var uuid pgtype.UUID
	if err := uuid.Scan(id); err != nil {
		return nil, fmt.Errorf("invalid uuid format: %w", err)
	}

	row, err := r.db.Query.GetUserByID(ctx, uuid)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}
	return mapUserEntity(row), nil
}

// Mapper function to convert DB model to Domain model
func mapUserEntity(row generated.User) *domain.User {
	var lastLogin *time.Time
	if row.LastLoginAt.Valid {
		t := row.LastLoginAt.Time
		lastLogin = &t
	}

	// UUID to String
	idStr := fmt.Sprintf("%x-%x-%x-%x-%x", row.ID.Bytes[0:4], row.ID.Bytes[4:6], row.ID.Bytes[6:8], row.ID.Bytes[8:10], row.ID.Bytes[10:16])

	var roleIDPtr *string
	// if row.RoleID.Valid { ... }

	return &domain.User{
		ID:           idStr,
		Username:     row.Username,
		Email:        row.Email,
		PasswordHash: row.PasswordHash,
		FullName:     row.FullName.String,
		Phone:        row.Phone.String,
		RoleID:       roleIDPtr,
		Status:       domain.UserStatus(row.Status.String),
		LastLoginAt:  lastLogin,
		CreatedAt:    row.CreatedAt.Time,
		UpdatedAt:    row.UpdatedAt.Time,
	}
}
