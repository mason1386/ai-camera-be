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
	repo     ports.IdentityRepository
	faceRepo ports.IdentityFaceRepository
}

func NewIdentityService(repo ports.IdentityRepository, faceRepo ports.IdentityFaceRepository) ports.IdentityService {
	return &IdentityService{
		repo:     repo,
		faceRepo: faceRepo,
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
		Type:               req.Type,
		PhoneNumber:        req.PhoneNumber,
		IdentityCardNumber: req.IdentityCardNumber,
		Department:         req.Department,
		Metadata:           req.Metadata,
		Status:             domain.IdentityStatusPending,
		Note:               req.Note,
		CreatedBy:          req.CreatedBy,
	}

	// 3. Save to repo
	return s.repo.CreateIdentity(ctx, identity)
}

func (s *IdentityService) GetIdentity(ctx context.Context, id uuid.UUID) (*domain.Identity, error) {
	identity, err := s.repo.GetIdentity(ctx, id)
	if err != nil {
		return nil, err
	}
	if identity != nil {
		faces, _ := s.faceRepo.ListFaces(ctx, id)
		if faces != nil {
			identity.Faces = make([]domain.IdentityFace, len(faces))
			for i, f := range faces {
				identity.Faces[i] = *f
			}
		}
	}
	return identity, nil
}

func (s *IdentityService) ListIdentities(ctx context.Context, page, limit int, search string) ([]*domain.Identity, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	return s.repo.ListIdentities(ctx, page, limit, search)
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
	if req.Type != "" {
		current.Type = req.Type
	}
	if req.PhoneNumber != "" {
		current.PhoneNumber = req.PhoneNumber
	}
	if req.IdentityCardNumber != "" {
		current.IdentityCardNumber = req.IdentityCardNumber
	}
	if req.Department != "" {
		current.Department = req.Department
	}
	if req.Metadata != nil {
		current.Metadata = req.Metadata
	}
	if req.Note != "" {
		current.Note = req.Note
	}

	// 3. Save
	return s.repo.UpdateIdentity(ctx, current)
}

func (s *IdentityService) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.IdentityStatus) (*domain.Identity, error) {
	return s.repo.UpdateIdentityStatus(ctx, id, status)
}

func (s *IdentityService) DeleteIdentity(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteIdentity(ctx, id)
}

func (s *IdentityService) EnrollFace(ctx context.Context, identityID uuid.UUID, imageURL string, isPrimary bool) (*domain.IdentityFace, error) {
	face := &domain.IdentityFace{
		IdentityID: identityID,
		ImageURL:   imageURL,
		IsPrimary:  isPrimary,
	}
	if isPrimary {
		_ = s.faceRepo.SetPrimary(ctx, identityID, uuid.Nil) // Reset others if any (logic inside SetPrimary handled it)
	}
	return s.faceRepo.CreateFace(ctx, face)
}

func (s *IdentityService) DeleteFace(ctx context.Context, faceID uuid.UUID) error {
	return s.faceRepo.DeleteFace(ctx, faceID)
}
