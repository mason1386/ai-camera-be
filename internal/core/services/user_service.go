package services

import (
	"context"
	"errors"
	"time"

	"app/internal/core/domain"
	"app/internal/core/ports"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo ports.UserRepository
}

func NewUserService(repo ports.UserRepository) ports.UserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) CreateUser(ctx context.Context, req *domain.CreateWebUserRequest) (*domain.User, error) {
	// Check if username or email exists
	existingAuth, _ := s.repo.GetByUsername(ctx, req.Username)
	if existingAuth != nil {
		return nil, errors.New("username already exists")
	}
	existingEmail, _ := s.repo.GetByEmail(ctx, req.Email)
	if existingEmail != nil {
		return nil, errors.New("email already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		FullName:     req.FullName,
		Phone:        req.Phone,
		Status:       req.Status,
	}
	if req.RoleID != "" {
		user.RoleID = &req.RoleID
	}

	if err := s.repo.Save(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) ListUsers(ctx context.Context, search string) ([]*domain.User, error) {
	return s.repo.List(ctx, search)
}

func (s *UserService) UpdateUser(ctx context.Context, id string, req *domain.UpdateWebUserRequest) (*domain.User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	if req.FullName != nil {
		user.FullName = *req.FullName
	}
	if req.Phone != nil {
		user.Phone = *req.Phone
	}
	if req.RoleID != nil {
		user.RoleID = req.RoleID // Can be nil or changed
	}
	if req.Status != nil {
		user.Status = *req.Status
	}
	user.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) DeleteUser(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *UserService) ResetPassword(ctx context.Context, userID string, newPassword string) error {
	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hashedPassword)
	user.UpdatedAt = time.Now()

	return s.repo.Update(ctx, user)
}
