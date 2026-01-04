package ports

import (
	"context"

	"app/internal/core/domain"
)

type UserRepository interface {
	Save(ctx context.Context, user *domain.User) error
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	GetByUsername(ctx context.Context, username string) (*domain.User, error)
	GetByID(ctx context.Context, id string) (*domain.User, error)
	List(ctx context.Context, search string) ([]*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id string) error
}

type AuthService interface {
	Register(ctx context.Context, req *domain.RegisterRequest) (*domain.User, error)
	Login(ctx context.Context, req *domain.LoginRequest) (*domain.LoginResponse, error)
}
