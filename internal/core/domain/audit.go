package domain

import (
	"time"

	"github.com/google/uuid"
)

type AuditLog struct {
	ID        int64          `json:"id"`
	UserID    *uuid.UUID     `json:"user_id"`
	Action    string         `json:"action"`
	TableName string         `json:"table_name"`
	RecordID  string         `json:"record_id"`
	OldValue  map[string]any `json:"old_value"`
	NewValue  map[string]any `json:"new_value"`
	IPAddress string         `json:"ip_address"`
	UserAgent string         `json:"user_agent"`
	CreatedAt time.Time      `json:"created_at"`

	// Join fields
	Username string `json:"username,omitempty"`
}
