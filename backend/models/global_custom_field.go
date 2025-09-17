package models

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// GlobalCustomField represents a unified custom field definition (supports both global and user-managed fields)
type GlobalCustomField struct {
	ID         uuid.UUID       `json:"id" db:"id"`
	Name       string          `json:"name" db:"name"`
	Label      string          `json:"label" db:"label"`
	Type       string          `json:"type" db:"type"`
	Required   bool            `json:"required" db:"required"`
	Options    []string        `json:"options,omitempty" db:"options"`
	Validation FieldValidation `json:"validation,omitempty" db:"validation"`
	Order      int             `json:"order_index" db:"order"`
	IsActive   bool            `json:"is_active" db:"is_active"`
	CreatedAt  time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time       `json:"updated_at" db:"updated_at"`
}

// UserCustomFieldValue represents a user's value for a global custom field
type UserCustomFieldValue struct {
	ID        uuid.UUID `json:"id" db:"id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	FieldID   uuid.UUID `json:"field_id" db:"field_id"`
	Value     string    `json:"value" db:"value"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// GlobalCustomFieldWithValue represents a global field with a user's value
type GlobalCustomFieldWithValue struct {
	GlobalCustomField
	Value *string `json:"value,omitempty"`
}

// CreateGlobalCustomFieldRequest represents the request to create a global custom field
type CreateGlobalCustomFieldRequest struct {
	Name       string          `json:"name"`
	Label      string          `json:"label"`
	Type       string          `json:"type"`
	Required   bool            `json:"required"`
	Options    []string        `json:"options,omitempty"`
	Validation FieldValidation `json:"validation,omitempty"`
}

// UpdateGlobalCustomFieldRequest represents the request to update a global custom field
type UpdateGlobalCustomFieldRequest struct {
	Label      string          `json:"label,omitempty"`
	Type       string          `json:"type,omitempty"`
	Required   *bool           `json:"required,omitempty"`
	Options    []string        `json:"options,omitempty"`
	Validation FieldValidation `json:"validation,omitempty"`
	IsActive   *bool           `json:"is_active,omitempty"`
}

// UserCustomFieldValuesRequest represents the request to update multiple field values for a user
type UserCustomFieldValuesRequest struct {
	FieldValues map[string]any `json:"field_values"` // field_id -> value
}

// ValidateGlobalCustomField validates a global custom field
func ValidateGlobalCustomField(field *GlobalCustomField) error {
	if field.Name == "" {
		return NewValidationError("name", "Field name is required")
	}

	if !isValidFieldName(field.Name) {
		return NewValidationError("name", "Field name must start with a letter and contain only letters, numbers, and underscores")
	}

	if field.Label == "" {
		return NewValidationError("label", "Field label is required")
	}

	if !isValidFieldType(field.Type) {
		return NewValidationError("type", "Invalid field type")
	}

	// Validate options for select/multiselect fields
	if (field.Type == "select" || field.Type == "multiselect") && len(field.Options) == 0 {
		return NewValidationError("options", "Options are required for select fields")
	}

	return nil
}

// ValidateUserCustomFieldValue validates a user's custom field value
func ValidateUserCustomFieldValue(field *GlobalCustomField, value any) error {
	if field.Required && (value == nil || value == "") {
		return NewValidationError("value", "This field is required")
	}

	if value == nil || value == "" {
		return nil // Optional field with no value is valid
	}

	// Convert value to string for validation
	var strValue string
	switch v := value.(type) {
	case string:
		strValue = v
	default:
		strValue = fmt.Sprintf("%v", v)
	}

	// Type-specific validation
	switch field.Type {
	case "text", "textarea":
		return validateTextValue(strValue, field.Validation)
	case "number":
		return validateNumberValue(strValue, field.Validation)
	case "email":
		return validateEmailValue(strValue)
	case "phone":
		return validatePhoneValue(strValue)
	case "date":
		return validateDateValue(strValue)
	case "boolean":
		return validateBooleanValue(strValue)
	case "select":
		return validateSelectValue(strValue, field.Options)
	case "multiselect":
		return validateMultiselectValue(value, field.Options)
	case "file":
		return validateFileValue(strValue, field.Validation)
	default:
		return NewValidationError("type", "Unknown field type")
	}
}

// Helper functions for validation
func isValidFieldName(name string) bool {
	if len(name) == 0 {
		return false
	}
	for i, r := range name {
		if i == 0 {
			if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')) {
				return false
			}
		} else {
			if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_') {
				return false
			}
		}
	}
	return true
}

func isValidFieldType(fieldType string) bool {
	validTypes := []string{"text", "textarea", "number", "email", "phone", "date", "boolean", "select", "multiselect", "file"}
	for _, t := range validTypes {
		if t == fieldType {
			return true
		}
	}
	return false
}

// ToJSON converts options to JSON bytes for database storage
func (gcf *GlobalCustomField) OptionsToJSON() ([]byte, error) {
	if len(gcf.Options) == 0 {
		return []byte{}, nil
	}
	return json.Marshal(gcf.Options)
}

// FromJSON parses options from JSON bytes
func (gcf *GlobalCustomField) OptionsFromJSON(jsonBytes []byte) error {
	if len(jsonBytes) == 0 {
		gcf.Options = nil
		return nil
	}
	return json.Unmarshal(jsonBytes, &gcf.Options)
}

// ValidationToJSON converts validation to JSON bytes for database storage
func (gcf *GlobalCustomField) ValidationToJSON() ([]byte, error) {
	if gcf.Validation == (FieldValidation{}) {
		return nil, nil
	}
	return json.Marshal(gcf.Validation)
}

// ValidationFromJSON parses validation from JSON bytes
func (gcf *GlobalCustomField) ValidationFromJSON(jsonBytes []byte) error {
	if len(jsonBytes) == 0 {
		gcf.Validation = FieldValidation{}
		return nil
	}
	return json.Unmarshal(jsonBytes, &gcf.Validation)
}
