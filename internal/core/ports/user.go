package ports

import (
	"context"

	"app/internal/core/domain"
)

type UserService interface {
	CreateUser(ctx context.Context, req *domain.CreateWebUserRequest) (*domain.User, error)
	ListUsers(ctx context.Context, search string) ([]*domain.User, error)
	UpdateUser(ctx context.Context, id string, req *domain.UpdateWebUserRequest) (*domain.User, error)
	DeleteUser(ctx context.Context, id string) error
	ResetPassword(ctx context.Context, userID string, newPassword string) error
}
