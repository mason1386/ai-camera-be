package ports

import (
	"context"

	"app/internal/core/domain"

	"github.com/google/uuid"
)

type IdentityRepository interface {
	CreateIdentity(ctx context.Context, identity *domain.Identity) (*domain.Identity, error)
	GetIdentity(ctx context.Context, id uuid.UUID) (*domain.Identity, error)
	GetIdentityByCode(ctx context.Context, code string) (*domain.Identity, error)
	ListIdentities(ctx context.Context, limit, offset int32) ([]*domain.Identity, error)
	CountIdentities(ctx context.Context) (int64, error)
	UpdateIdentity(ctx context.Context, identity *domain.Identity) (*domain.Identity, error)
	UpdateIdentityStatus(ctx context.Context, id uuid.UUID, status domain.IdentityStatus, approvedBy *uuid.UUID) (*domain.Identity, error)
	DeleteIdentity(ctx context.Context, id uuid.UUID) error
}

type IdentityService interface {
	CreateIdentity(ctx context.Context, req *CreateIdentityRequest) (*domain.Identity, error)
	GetIdentity(ctx context.Context, id uuid.UUID) (*domain.Identity, error)
	ListIdentities(ctx context.Context, page, limit int) ([]*domain.Identity, int64, error)
	UpdateIdentity(ctx context.Context, id uuid.UUID, req *UpdateIdentityRequest) (*domain.Identity, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status domain.IdentityStatus, approvedBy uuid.UUID) (*domain.Identity, error)
	DeleteIdentity(ctx context.Context, id uuid.UUID) error
}

// DTOs
type CreateIdentityRequest struct {
	Code               string `json:"code" binding:"required"`
	FullName           string `json:"full_name" binding:"required"`
	PhoneNumber        string `json:"phone_number"`
	IdentityCardNumber string `json:"identity_card_number"`
	FaceImageURL       string `json:"face_image_url" binding:"required"`
	Type               string `json:"type" binding:"required"`
	Note               string `json:"note"`
	CreatedBy          *uuid.UUID
}

type UpdateIdentityRequest struct {
	FullName     string `json:"full_name"`
	PhoneNumber  string `json:"phone_number"`
	FaceImageURL string `json:"face_image_url"`
}
