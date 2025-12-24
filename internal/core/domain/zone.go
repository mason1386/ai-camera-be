package domain

import "time"

type Zone struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

type CreateZoneRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type UpdateZoneRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
