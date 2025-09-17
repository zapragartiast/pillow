package models

import (
	"time"

	"github.com/google/uuid"
)

// CustomFieldType represents the type of custom field
type CustomFieldType string

const (
	FieldTypeText        CustomFieldType = "text"
	FieldTypeTextarea    CustomFieldType = "textarea"
	FieldTypeNumber      CustomFieldType = "number"
	FieldTypeEmail       CustomFieldType = "email"
	FieldTypePhone       CustomFieldType = "phone"
	FieldTypeDate        CustomFieldType = "date"
	FieldTypeBoolean     CustomFieldType = "boolean"
	FieldTypeSelect      CustomFieldType = "select"
	FieldTypeMultiselect CustomFieldType = "multiselect"
)

// CustomField represents a single custom field configuration
type CustomField struct {
	ID         string          `json:"id"`
	Name       string          `json:"name"`
	Label      string          `json:"label"`
	Type       CustomFieldType `json:"type"`
	Required   bool            `json:"required"`
	Value      any             `json:"value,omitempty"`
	Options    []string        `json:"options,omitempty"` // for select/multiselect
	Validation FieldValidation `json:"validation,omitempty"`
	Order      int             `json:"order"`
}

// FieldValidation represents validation rules for a field
type FieldValidation struct {
	MinLength *int     `json:"min_length,omitempty"`
	MaxLength *int     `json:"max_length,omitempty"`
	Min       *float64 `json:"min,omitempty"`
	Max       *float64 `json:"max,omitempty"`
	Pattern   *string  `json:"pattern,omitempty"`
}

// CustomFieldsData represents the structure of custom fields data stored in JSON
type CustomFieldsData struct {
	Fields   []CustomField  `json:"fields"`
	Metadata CustomMetadata `json:"metadata"`
}

// CustomMetadata contains metadata about custom fields
type CustomMetadata struct {
	Version     string    `json:"version"`
	LastUpdated time.Time `json:"last_updated"`
}

type User struct {
	ID           uuid.UUID `json:"id" db:"id"`
	Username     string    `json:"username" db:"username"`
	PasswordHash string    `json:"password_hash,omitempty" db:"password_hash"`
	Email        string    `json:"email" db:"email"`
	IsActive     bool      `json:"is_active" db:"is_active"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}
