package models

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// ValidateCustomField validates a custom field based on its type and rules
func ValidateCustomField(field *CustomField) error {
	if field.Required && field.Value == nil {
		return fmt.Errorf("field '%s' is required", field.Label)
	}

	if field.Value == nil {
		return nil // Optional field with no value is valid
	}

	switch field.Type {
	case FieldTypeText, FieldTypeTextarea:
		return validateTextField(field)
	case FieldTypeNumber:
		return validateNumberField(field)
	case FieldTypeEmail:
		return validateEmailField(field)
	case FieldTypePhone:
		return validatePhoneField(field)
	case FieldTypeDate:
		return validateDateField(field)
	case FieldTypeBoolean:
		return validateBooleanField(field)
	case FieldTypeSelect:
		return validateSelectField(field)
	case FieldTypeMultiselect:
		return validateMultiselectField(field)
	default:
		return fmt.Errorf("unknown field type: %s", field.Type)
	}
}

// validateTextField validates text-based fields
func validateTextField(field *CustomField) error {
	value, ok := field.Value.(string)
	if !ok {
		return fmt.Errorf("field '%s' must be a string", field.Label)
	}

	if field.Validation.MinLength != nil && len(value) < *field.Validation.MinLength {
		return fmt.Errorf("field '%s' must be at least %d characters", field.Label, *field.Validation.MinLength)
	}

	if field.Validation.MaxLength != nil && len(value) > *field.Validation.MaxLength {
		return fmt.Errorf("field '%s' must be at most %d characters", field.Label, *field.Validation.MaxLength)
	}

	if field.Validation.Pattern != nil {
		matched, err := regexp.MatchString(*field.Validation.Pattern, value)
		if err != nil {
			return fmt.Errorf("invalid regex pattern for field '%s'", field.Label)
		}
		if !matched {
			return fmt.Errorf("field '%s' does not match required pattern", field.Label)
		}
	}

	return nil
}

// validateNumberField validates number fields
func validateNumberField(field *CustomField) error {
	var num float64
	switch v := field.Value.(type) {
	case float64:
		num = v
	case int:
		num = float64(v)
	case string:
		var err error
		num, err = strconv.ParseFloat(v, 64)
		if err != nil {
			return fmt.Errorf("field '%s' must be a valid number", field.Label)
		}
	default:
		return fmt.Errorf("field '%s' must be a number", field.Label)
	}

	if field.Validation.Min != nil && num < *field.Validation.Min {
		return fmt.Errorf("field '%s' must be at least %g", field.Label, *field.Validation.Min)
	}

	if field.Validation.Max != nil && num > *field.Validation.Max {
		return fmt.Errorf("field '%s' must be at most %g", field.Label, *field.Validation.Max)
	}

	return nil
}

// validateEmailField validates email fields
func validateEmailField(field *CustomField) error {
	value, ok := field.Value.(string)
	if !ok {
		return fmt.Errorf("field '%s' must be a string", field.Label)
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(value) {
		return fmt.Errorf("field '%s' must be a valid email address", field.Label)
	}

	return nil
}

// validatePhoneField validates phone number fields
func validatePhoneField(field *CustomField) error {
	value, ok := field.Value.(string)
	if !ok {
		return fmt.Errorf("field '%s' must be a string", field.Label)
	}

	// Basic phone validation - allows various formats
	phoneRegex := regexp.MustCompile(`^[\+]?[1-9][\d]{0,15}$`)
	cleaned := strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(value, " ", ""), "-", ""), "+", "")
	if !phoneRegex.MatchString(cleaned) {
		return fmt.Errorf("field '%s' must be a valid phone number", field.Label)
	}

	return nil
}

// validateDateField validates date fields
func validateDateField(field *CustomField) error {
	value, ok := field.Value.(string)
	if !ok {
		return fmt.Errorf("field '%s' must be a date string", field.Label)
	}

	// Try parsing with common date formats
	formats := []string{
		time.RFC3339,
		"2006-01-02",
		"02/01/2006",
		"01/02/2006",
		"2006/01/02",
	}

	for _, format := range formats {
		if _, err := time.Parse(format, value); err == nil {
			return nil
		}
	}

	return fmt.Errorf("field '%s' must be a valid date", field.Label)
}

// validateBooleanField validates boolean fields
func validateBooleanField(field *CustomField) error {
	switch field.Value.(type) {
	case bool:
		return nil
	default:
		return fmt.Errorf("field '%s' must be a boolean", field.Label)
	}
}

// validateSelectField validates select fields
func validateSelectField(field *CustomField) error {
	value, ok := field.Value.(string)
	if !ok {
		return fmt.Errorf("field '%s' must be a string", field.Label)
	}

	if len(field.Options) == 0 {
		return fmt.Errorf("field '%s' has no options defined", field.Label)
	}

	for _, option := range field.Options {
		if option == value {
			return nil
		}
	}

	return fmt.Errorf("field '%s' value '%s' is not a valid option", field.Label, value)
}

// validateMultiselectField validates multi-select fields
func validateMultiselectField(field *CustomField) error {
	values, ok := field.Value.([]any)
	if !ok {
		// Try string array
		if strValues, ok := field.Value.([]string); ok {
			values = make([]any, len(strValues))
			for i, v := range strValues {
				values[i] = v
			}
		} else {
			return fmt.Errorf("field '%s' must be an array", field.Label)
		}
	}

	if len(field.Options) == 0 {
		return fmt.Errorf("field '%s' has no options defined", field.Label)
	}

	for _, value := range values {
		strValue, ok := value.(string)
		if !ok {
			return fmt.Errorf("field '%s' array elements must be strings", field.Label)
		}

		found := false
		for _, option := range field.Options {
			if option == strValue {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("field '%s' value '%s' is not a valid option", field.Label, strValue)
		}
	}

	return nil
}

// ValidateCustomFieldsData validates all custom fields in the data structure
func ValidateCustomFieldsData(data *CustomFieldsData) error {
	if data == nil {
		return nil
	}

	for _, field := range data.Fields {
		if err := ValidateCustomField(&field); err != nil {
			return err
		}
	}

	return nil
}

// NewCustomFieldsData creates a new CustomFieldsData with default metadata
func NewCustomFieldsData() *CustomFieldsData {
	return &CustomFieldsData{
		Fields: []CustomField{},
		Metadata: CustomMetadata{
			Version:     "1.0",
			LastUpdated: time.Now(),
		},
	}
}

// AddField adds a new field to the custom fields data
func (cfd *CustomFieldsData) AddField(field CustomField) {
	field.ID = generateFieldID()
	field.Order = len(cfd.Fields)
	cfd.Fields = append(cfd.Fields, field)
	cfd.Metadata.LastUpdated = time.Now()
}

// UpdateField updates an existing field
func (cfd *CustomFieldsData) UpdateField(fieldID string, updatedField CustomField) error {
	for i, field := range cfd.Fields {
		if field.ID == fieldID {
			updatedField.ID = fieldID
			updatedField.Order = field.Order
			cfd.Fields[i] = updatedField
			cfd.Metadata.LastUpdated = time.Now()
			return nil
		}
	}
	return fmt.Errorf("field with ID '%s' not found", fieldID)
}

// DeleteField removes a field from the custom fields data
func (cfd *CustomFieldsData) DeleteField(fieldID string) error {
	for i, field := range cfd.Fields {
		if field.ID == fieldID {
			cfd.Fields = append(cfd.Fields[:i], cfd.Fields[i+1:]...)
			// Reorder remaining fields
			for j := i; j < len(cfd.Fields); j++ {
				cfd.Fields[j].Order = j
			}
			cfd.Metadata.LastUpdated = time.Now()
			return nil
		}
	}
	return fmt.Errorf("field with ID '%s' not found", fieldID)
}

// GetField retrieves a field by ID
func (cfd *CustomFieldsData) GetField(fieldID string) (*CustomField, error) {
	for _, field := range cfd.Fields {
		if field.ID == fieldID {
			return &field, nil
		}
	}
	return nil, fmt.Errorf("field with ID '%s' not found", fieldID)
}

// ToJSON converts the custom fields data to JSON
func (cfd *CustomFieldsData) ToJSON() ([]byte, error) {
	return json.Marshal(cfd)
}

// FromJSON populates the custom fields data from JSON
func (cfd *CustomFieldsData) FromJSON(data []byte) error {
	return json.Unmarshal(data, cfd)
}

// generateFieldID generates a unique ID for a field
func generateFieldID() string {
	return fmt.Sprintf("field_%d", time.Now().UnixNano())
}
