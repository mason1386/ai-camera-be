package services

import (
	"context"
	"errors"
	"fmt"

	"app/internal/core/domain"
	"app/internal/core/ports"

	"github.com/google/uuid"
)

type IdentityService struct {
	repo ports.IdentityRepository
}

func NewIdentityService(repo ports.IdentityRepository) ports.IdentityService {
	return &IdentityService{
		repo: repo,
	}
}

func (s *IdentityService) CreateIdentity(ctx context.Context, req *ports.CreateIdentityRequest) (*domain.Identity, error) {
	// 1. Check if code exists
	existing, _ := s.repo.GetIdentityByCode(ctx, req.Code)
	if existing != nil {
		return nil, fmt.Errorf("identity with code %s already exists", req.Code)
	}

	// 2. Create domain entity
	identity := &domain.Identity{
		Code:               req.Code,
		FullName:           req.FullName,
		PhoneNumber:        req.PhoneNumber,
		IdentityCardNumber: req.IdentityCardNumber,
		FaceImageURL:       req.FaceImageURL,
		Type:               req.Type,
		Status:             domain.IdentityStatusPending, // Default Pending
		Note:               req.Note,
		CreatedBy:          req.CreatedBy,
	}

	// 3. Save to repo
	return s.repo.CreateIdentity(ctx, identity)
}

func (s *IdentityService) GetIdentity(ctx context.Context, id uuid.UUID) (*domain.Identity, error) {
	return s.repo.GetIdentity(ctx, id)
}

func (s *IdentityService) ListIdentities(ctx context.Context, page, limit int) ([]*domain.Identity, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	total, err := s.repo.CountIdentities(ctx)
	if err != nil {
		return nil, 0, err
	}

	identities, err := s.repo.ListIdentities(ctx, int32(limit), int32(offset))
	if err != nil {
		return nil, 0, err
	}

	return identities, total, nil
}

func (s *IdentityService) UpdateIdentity(ctx context.Context, id uuid.UUID, req *ports.UpdateIdentityRequest) (*domain.Identity, error) {
	// 1. Check existence
	current, err := s.repo.GetIdentity(ctx, id)
	if err != nil {
		return nil, err
	}
	if current == nil {
		return nil, errors.New("identity not found")
	}

	// 2. Update fields
	if req.FullName != "" {
		current.FullName = req.FullName
	}
	if req.PhoneNumber != "" {
		current.PhoneNumber = req.PhoneNumber
	}
	if req.FaceImageURL != "" {
		current.FaceImageURL = req.FaceImageURL
		// TODO: Trigger Kafka event for re-indexing face
	}

	// 3. Save
	return s.repo.UpdateIdentity(ctx, current)
}

func (s *IdentityService) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.IdentityStatus, approvedBy uuid.UUID) (*domain.Identity, error) {
	return s.repo.UpdateIdentityStatus(ctx, id, status, &approvedBy)
}

func (s *IdentityService) DeleteIdentity(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteIdentity(ctx, id)
}
