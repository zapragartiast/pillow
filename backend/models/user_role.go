package models

import (
	"github.com/google/uuid"
)

type UserRole struct {
	UserID       uuid.UUID  `json:"user_id" db:"user_id"`
	RoleID       uuid.UUID  `json:"role_id" db:"role_id"`
	Scope        string     `json:"scope" db:"scope"`
	ParentRoleID *uuid.UUID `json:"parent_role_id,omitempty" db:"parent_role_id"`
}
