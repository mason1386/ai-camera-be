package domain

import (
	"time"
)

type UserStatus string

const (
	UserStatusActive UserStatus = "active"
	UserStatusLocked UserStatus = "locked"
	UserStatusBanned UserStatus = "banned"
)

type User struct {
	ID           string     `json:"id"` // Converted from UUID
	Username     string     `json:"username"`
	Email        string     `json:"email"`
	PasswordHash string     `json:"-"`
	FullName     string     `json:"full_name"`
	Phone        string     `json:"phone"`
	RoleID       *string    `json:"role_id"`
	Status       UserStatus `json:"status"`
	LastLoginAt  *time.Time `json:"last_login_at"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"` // Login by Email
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
	User        *User  `json:"user"`
}

type CreateWebUserRequest struct {
	Username string     `json:"username" binding:"required"`
	Email    string     `json:"email" binding:"required,email"`
	Password string     `json:"password" binding:"required,min=6"`
	FullName string     `json:"full_name"`
	Phone    string     `json:"phone"`
	RoleID   string     `json:"role_id"`
	Status   UserStatus `json:"status"`
}

type UpdateWebUserRequest struct {
	FullName *string     `json:"full_name"`
	Phone    *string     `json:"phone"`
	RoleID   *string     `json:"role_id"`
	Status   *UserStatus `json:"status"`
}
