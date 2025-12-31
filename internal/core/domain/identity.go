package domain

import (
	"time"

	"github.com/google/uuid"
)

type IdentityStatus string

const (
	IdentityStatusPending  IdentityStatus = "pending"
	IdentityStatusActive   IdentityStatus = "active"
	IdentityStatusRejected IdentityStatus = "rejected"
)

type Identity struct {
	ID                 uuid.UUID      `json:"id"`
	Code               string         `json:"code"`
	FullName           string         `json:"full_name"`
	PhoneNumber        string         `json:"phone_number"`
	IdentityCardNumber string         `json:"identity_card_number"`
	FaceImageURL       string         `json:"face_image_url"`
	Type               string         `json:"type"`
	Status             IdentityStatus `json:"status"`
	Note               string         `json:"note"`
	CreatedBy          *uuid.UUID     `json:"created_by"`
	ApprovedBy         *uuid.UUID     `json:"approved_by"`
	UserAccountID      *uuid.UUID     `json:"user_account_id"`
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`
	DeletedAt          *time.Time     `json:"deleted_at,omitempty"`
}
