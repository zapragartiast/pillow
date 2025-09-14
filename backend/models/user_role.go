package models

import (
	"time"

	"github.com/google/uuid"
)

type UserRole struct {
	UserID       uuid.UUID  `json:"user_id" db:"user_id"`
	RoleID       uuid.UUID  `json:"role_id" db:"role_id"`
	Scope        string     `json:"scope" db:"scope"`
	ParentRoleID *uuid.UUID `json:"parent_role_id,omitempty" db:"parent_role_id"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
}
