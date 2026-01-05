package services

import (
	"context"

	"app/internal/core/domain"
	"app/internal/core/ports"
)

type RoleService struct {
	repo ports.RoleRepository
}

func NewRoleService(repo ports.RoleRepository) ports.RoleService {
	return &RoleService{repo: repo}
}

func (s *RoleService) CreateRole(ctx context.Context, role *domain.Role) error {
	return s.repo.Create(ctx, role)
}

func (s *RoleService) GetRole(ctx context.Context, id string) (*domain.Role, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *RoleService) ListRoles(ctx context.Context, search string) ([]*domain.Role, error) {
	return s.repo.List(ctx, search)
}

func (s *RoleService) UpdateRole(ctx context.Context, id string, role *domain.Role) error {
	role.ID = id
	return s.repo.Update(ctx, role)
}

func (s *RoleService) DeleteRole(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
