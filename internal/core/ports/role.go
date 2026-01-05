package ports

import (
	"context"

	"app/internal/core/domain"
)

type RoleRepository interface {
	Create(ctx context.Context, role *domain.Role) error
	GetByID(ctx context.Context, id string) (*domain.Role, error)
	List(ctx context.Context, search string) ([]*domain.Role, error)
	Update(ctx context.Context, role *domain.Role) error
	Delete(ctx context.Context, id string) error
}

type RoleService interface {
	CreateRole(ctx context.Context, role *domain.Role) error
	GetRole(ctx context.Context, id string) (*domain.Role, error)
	ListRoles(ctx context.Context, search string) ([]*domain.Role, error)
	UpdateRole(ctx context.Context, id string, role *domain.Role) error
	DeleteRole(ctx context.Context, id string) error
}
