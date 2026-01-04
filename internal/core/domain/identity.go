package domain

import (
	"time"

	"github.com/google/uuid"
)

type IdentityStatus string
type PersonGroup string

const (
	IdentityStatusPending  IdentityStatus = "pending"
	IdentityStatusActive   IdentityStatus = "active"
	IdentityStatusRejected IdentityStatus = "rejected"

	PersonGroupEmployee  PersonGroup = "employee"
	PersonGroupVIP       PersonGroup = "vip"
	PersonGroupBlacklist PersonGroup = "blacklist"
	PersonGroupVisitor   PersonGroup = "visitor"
	PersonGroupOther     PersonGroup = "other"
)

type Identity struct {
	ID                 uuid.UUID      `json:"id"`
	Code               string         `json:"code"`
	FullName           string         `json:"full_name"`
	Type               string         `json:"type"` // STAFF, STUDENT, VIP...
	PhoneNumber        string         `json:"phone_number"`
	IdentityCardNumber string         `json:"identity_card_number"`
	FaceImageURL       string         `json:"face_image_url"`
	Department         string         `json:"department"`
	Metadata           map[string]any `json:"metadata"`
	Status             IdentityStatus `json:"status"`
	Note               string         `json:"note"`
	CreatedBy          *uuid.UUID     `json:"created_by"`
	ApprovedBy         *uuid.UUID     `json:"approved_by"`
	UserAccountID      *uuid.UUID     `json:"user_account_id"`
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`
	DeletedAt          *time.Time     `json:"deleted_at,omitempty"`
	Faces              []IdentityFace `json:"faces,omitempty"`
}

type IdentityFace struct {
	ID           uuid.UUID `json:"id"`
	IdentityID   uuid.UUID `json:"identity_id"`
	ImageURL     string    `json:"image_url"`
	IsPrimary    bool      `json:"is_primary"`
	QualityScore float64   `json:"quality_score"`
	BlurScore    float64   `json:"blur_score"`
	CreatedAt    time.Time `json:"created_at"`
}
