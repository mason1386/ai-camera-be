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
	ListIdentities(ctx context.Context, page, limit int, search string) ([]*domain.Identity, int64, error)
	CountIdentities(ctx context.Context) (int64, error)
	UpdateIdentity(ctx context.Context, identity *domain.Identity) (*domain.Identity, error)
	UpdateIdentityStatus(ctx context.Context, id uuid.UUID, status domain.IdentityStatus) (*domain.Identity, error)
	DeleteIdentity(ctx context.Context, id uuid.UUID) error
}

type IdentityFaceRepository interface {
	CreateFace(ctx context.Context, face *domain.IdentityFace) (*domain.IdentityFace, error)
	ListFaces(ctx context.Context, identityID uuid.UUID) ([]*domain.IdentityFace, error)
	DeleteFace(ctx context.Context, id uuid.UUID) error
	SetPrimary(ctx context.Context, identityID, faceID uuid.UUID) error
}

type IdentityService interface {
	CreateIdentity(ctx context.Context, req *CreateIdentityRequest) (*domain.Identity, error)
	GetIdentity(ctx context.Context, id uuid.UUID) (*domain.Identity, error)
	ListIdentities(ctx context.Context, page, limit int, search string) ([]*domain.Identity, int64, error)
	UpdateIdentity(ctx context.Context, id uuid.UUID, req *UpdateIdentityRequest) (*domain.Identity, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status domain.IdentityStatus) (*domain.Identity, error)
	DeleteIdentity(ctx context.Context, id uuid.UUID) error

	EnrollFace(ctx context.Context, identityID uuid.UUID, imageURL string, isPrimary bool) (*domain.IdentityFace, error)
	DeleteFace(ctx context.Context, faceID uuid.UUID) error
}

// DTOs
type CreateIdentityRequest struct {
	Code               string         `json:"code" binding:"required"`
	FullName           string         `json:"full_name" binding:"required"`
	Type               string         `json:"type"` // STAFF, STUDENT, VIP...
	PhoneNumber        string         `json:"phone_number"`
	IdentityCardNumber string         `json:"identity_card_number"`
	Department         string         `json:"department"`
	Metadata           map[string]any `json:"metadata"`
	Note               string         `json:"note"`
	CreatedBy          *uuid.UUID
}

type UpdateIdentityRequest struct {
	FullName           string         `json:"full_name"`
	Type               string         `json:"type"`
	PhoneNumber        string         `json:"phone_number"`
	IdentityCardNumber string         `json:"identity_card_number"`
	Department         string         `json:"department"`
	Metadata           map[string]any `json:"metadata"`
	Note               string         `json:"note"`
}
