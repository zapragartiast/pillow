package models

import (
	"github.com/google/uuid"
)

type Permission struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description,omitempty" db:"description"`
	ScopeLevel  string    `json:"scope_level" db:"scope_level"`
}
