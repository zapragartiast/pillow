package models

import (
	"time"

	"github.com/google/uuid"
)

type AuditLog struct {
	ID        uuid.UUID  `json:"id" db:"id"`
	UserID    *uuid.UUID `json:"user_id,omitempty" db:"user_id"`
	Action    string     `json:"action" db:"action"`
	Timestamp time.Time  `json:"timestamp" db:"timestamp"`
	Details   string     `json:"details,omitempty" db:"details"`
}
