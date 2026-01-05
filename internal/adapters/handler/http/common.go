package http

import "app/internal/core/domain"

type ErrorResponse struct {
	Error string `json:"error"`
}

type PaginatedResponse struct {
	Data  interface{} `json:"data"`
	Total int64       `json:"total"`
	Page  int         `json:"page"`
	Limit int         `json:"limit"`
}

type IdentityResponse struct {
	Data domain.Identity `json:"data"`
}

type AuditLogResponse struct {
	Data []domain.AuditLog `json:"data"`
}

type RecognitionLogResponse struct {
	Data []domain.RecognitionLog `json:"data"`
}

type AttendanceRecordResponse struct {
	Data []domain.AttendanceRecord `json:"data"`
}
