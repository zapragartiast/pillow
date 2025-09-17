package models

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// ValidationError represents a validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (ve ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", ve.Field, ve.Message)
}

// NewValidationError creates a new validation error
func NewValidationError(field, message string) ValidationError {
	return ValidationError{
		Field:   field,
		Message: message,
	}
}

// validateTextValue validates text-based field values
func validateTextValue(value string, validation FieldValidation) error {
	if validation.MinLength != nil && len(value) < *validation.MinLength {
		return NewValidationError("value", fmt.Sprintf("Must be at least %d characters", *validation.MinLength))
	}

	if validation.MaxLength != nil && len(value) > *validation.MaxLength {
		return NewValidationError("value", fmt.Sprintf("Must be at most %d characters", *validation.MaxLength))
	}

	if validation.Pattern != nil {
		matched, err := regexp.MatchString(*validation.Pattern, value)
		if err != nil {
			return NewValidationError("value", "Invalid regex pattern")
		}
		if !matched {
			return NewValidationError("value", "Does not match required pattern")
		}
	}

	return nil
}

// validateNumberValue validates number field values
func validateNumberValue(value string, validation FieldValidation) error {
	num, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return NewValidationError("value", "Must be a valid number")
	}

	if validation.Min != nil && num < *validation.Min {
		return NewValidationError("value", fmt.Sprintf("Must be at least %g", *validation.Min))
	}

	if validation.Max != nil && num > *validation.Max {
		return NewValidationError("value", fmt.Sprintf("Must be at most %g", *validation.Max))
	}

	return nil
}

// validateEmailValue validates email field values
func validateEmailValue(value string) error {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(value) {
		return NewValidationError("value", "Must be a valid email address")
	}
	return nil
}

// validatePhoneValue validates phone number field values
func validatePhoneValue(value string) error {
	// Basic phone validation - allows various formats
	phoneRegex := regexp.MustCompile(`^[\+]?[1-9][\d]{0,15}$`)
	cleaned := strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(value, " ", ""), "-", ""), "+", "")
	if !phoneRegex.MatchString(cleaned) {
		return NewValidationError("value", "Must be a valid phone number")
	}
	return nil
}

// validateDateValue validates date field values
func validateDateValue(value string) error {
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

	return NewValidationError("value", "Must be a valid date")
}

// validateBooleanValue validates boolean field values
func validateBooleanValue(value string) error {
	lowerValue := strings.ToLower(value)
	if lowerValue != "true" && lowerValue != "false" && lowerValue != "1" && lowerValue != "0" {
		return NewValidationError("value", "Must be a boolean value")
	}
	return nil
}

// validateSelectValue validates select field values
func validateSelectValue(value string, options []string) error {
	if len(options) == 0 {
		return NewValidationError("value", "No options available")
	}

	for _, option := range options {
		if option == value {
			return nil
		}
	}

	return NewValidationError("value", "Invalid option selected")
}

// validateMultiselectValue validates multiselect field values
func validateMultiselectValue(value interface{}, options []string) error {
	if len(options) == 0 {
		return NewValidationError("value", "No options available")
	}

	var values []string
	switch v := value.(type) {
	case []string:
		values = v
	case []interface{}:
		for _, item := range v {
			if str, ok := item.(string); ok {
				values = append(values, str)
			} else {
				return NewValidationError("value", "All values must be strings")
			}
		}
	default:
		return NewValidationError("value", "Must be an array of values")
	}

	for _, val := range values {
		found := false
		for _, option := range options {
			if option == val {
				found = true
				break
			}
		}
		if !found {
			return NewValidationError("value", fmt.Sprintf("Invalid option: %s", val))
		}
	}

	return nil
}

// validateFileValue validates file field values
func validateFileValue(value string, validation FieldValidation) error {
	if value == "" {
		return nil // Empty value is valid
	}

	// Check file extension if pattern is specified
	if validation.Pattern != nil {
		matched, err := regexp.MatchString(*validation.Pattern, value)
		if err != nil {
			return NewValidationError("value", "Invalid file pattern")
		}
		if !matched {
			return NewValidationError("value", "File type not allowed")
		}
	}

	// Check file size if max length is specified (assuming max length represents max file size in bytes)
	if validation.MaxLength != nil && len(value) > *validation.MaxLength {
		return NewValidationError("value", fmt.Sprintf("File name too long (max %d characters)", *validation.MaxLength))
	}

	return nil
}
