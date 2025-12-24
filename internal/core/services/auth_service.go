package services

import (
	"context"
	"errors"

	"app/internal/core/domain"
	"app/internal/core/ports"
	"app/pkg/logger"
	"app/pkg/utils"
	"go.uber.org/zap"
)

type AuthService struct {
	userRepo ports.UserRepository
}

// TODO: Move secret key to config
const jwtSecret = "my_super_secret_key"

func NewAuthService(userRepo ports.UserRepository) ports.AuthService {
	return &AuthService{
		userRepo: userRepo,
	}
}

func (s *AuthService) Register(ctx context.Context, req *domain.RegisterRequest) (*domain.User, error) {
	// 1. Check if user exists (check both username and email)
	existingUser, err := s.userRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("username already exists")
	}

	existingUserEmail, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if existingUserEmail != nil {
		return nil, errors.New("email already exists")
	}

	// 2. Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// 3. Create user entity
	newUser := &domain.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hashedPassword,
		FullName:     req.FullName,
		Status:       domain.UserStatusActive,
	}

	// 4. Save to DB
	if err := s.userRepo.Save(ctx, newUser); err != nil {
		logger.Error("Failed to register user", zap.Error(err))
		return nil, err
	}

	return newUser, nil
}

func (s *AuthService) Login(ctx context.Context, req *domain.LoginRequest) (*domain.LoginResponse, error) {
	// 1. Find user by EMAIL
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("invalid credentials")
	}

	// 2. Check password
	if !utils.CheckPasswordHash(req.Password, user.PasswordHash) {
		return nil, errors.New("invalid credentials")
	}

	// 3. Generate Token
	token, err := utils.GenerateJWT(user.ID, user.Username, jwtSecret)
	if err != nil {
		return nil, err
	}

	return &domain.LoginResponse{
		AccessToken: token,
		User:        user,
	}, nil
}
