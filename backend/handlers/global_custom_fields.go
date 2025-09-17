package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"pillow/middleware"
	"pillow/models"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// GetGlobalCustomFields retrieves all global custom fields
func GetGlobalCustomFields(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query(`
			SELECT id, name, label, type, required, options, validation, "order", is_active, created_at, updated_at
			FROM custom_fields
			WHERE is_active = 1
			ORDER BY "order" ASC
		`)
		if err != nil {
			writeErrorResponse(w, err.Error(), http.StatusInternalServerError, r)
			return
		}
		defer rows.Close()

		var fields []models.GlobalCustomField
		for rows.Next() {
			var field models.GlobalCustomField
			var optionsJSON, validationJSON []byte

			err := rows.Scan(
				&field.ID, &field.Name, &field.Label, &field.Type,
				&field.Required, &optionsJSON, &validationJSON,
				&field.Order, &field.IsActive, &field.CreatedAt, &field.UpdatedAt,
			)
			if err != nil {
				writeErrorResponse(w, err.Error(), http.StatusInternalServerError, r)
				return
			}

			// Parse JSON fields
			if len(optionsJSON) > 0 {
				field.OptionsFromJSON(optionsJSON)
			}
			if len(validationJSON) > 0 {
				field.ValidationFromJSON(validationJSON)
			}

			fields = append(fields, field)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(fields)
	}
}

// CreateGlobalCustomField creates a new global custom field
func CreateGlobalCustomField(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check if user has permission to manage custom fields
		currentUser, ok := middleware.GetUserFromContext(r.Context())
		if !ok {
			writeErrorResponse(w, "Unauthorized", http.StatusUnauthorized, r)
			return
		}

		// TODO: Add permission check for managing custom fields
		_ = currentUser

		var createReq models.CreateGlobalCustomFieldRequest
		if err := json.NewDecoder(r.Body).Decode(&createReq); err != nil {
			writeErrorResponse(w, "Invalid JSON format: "+err.Error(), http.StatusBadRequest, r)
			return
		}

		// Create field object
		field := models.GlobalCustomField{
			ID:         uuid.New(),
			Name:       createReq.Name,
			Label:      createReq.Label,
			Type:       createReq.Type,
			Required:   createReq.Required,
			Options:    createReq.Options,
			Validation: createReq.Validation,
			IsActive:   true,
		}

		// Validate field
		if err := models.ValidateGlobalCustomField(&field); err != nil {
			writeErrorResponse(w, err.Error(), http.StatusBadRequest, r)
			return
		}

		// Get next order index
		var maxOrder int
		err := db.QueryRow(`SELECT COALESCE(MAX("order"), 0) FROM custom_fields`).Scan(&maxOrder)
		if err != nil {
			writeErrorResponse(w, "Failed to get order index: "+err.Error(), http.StatusInternalServerError, r)
			return
		}
		field.Order = maxOrder + 1

		// Convert JSON fields
		optionsJSON, err := field.OptionsToJSON()
		if err != nil {
			writeErrorResponse(w, "Failed to serialize options: "+err.Error(), http.StatusInternalServerError, r)
			return
		}

		validationJSON, err := field.ValidationToJSON()
		if err != nil {
			writeErrorResponse(w, "Failed to serialize validation: "+err.Error(), http.StatusInternalServerError, r)
			return
		}

		// Convert []byte to string for database insertion
		// Handle empty JSON cases properly for PostgreSQL JSONB
		var optionsStr, validationStr sql.NullString

		if len(optionsJSON) > 0 {
			optionsStr = sql.NullString{String: string(optionsJSON), Valid: true}
		} else {
			optionsStr = sql.NullString{Valid: false} // NULL for empty options
		}
		if len(validationJSON) > 0 {
			validationStr = sql.NullString{String: string(validationJSON), Valid: true}
		} else {
			validationStr = sql.NullString{Valid: false} // NULL for empty validation
		}

		// Insert into database
		_, err = db.Exec(`
			INSERT INTO custom_fields (id, name, label, type, required, options, validation, "order", is_active, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		`, field.ID, field.Name, field.Label, field.Type, field.Required, optionsStr, validationStr, field.Order, field.IsActive)

		if err != nil {
			writeErrorResponse(w, "Failed to create field: "+err.Error(), http.StatusInternalServerError, r)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(field)
	}
}

// UpdateGlobalCustomField updates a global custom field
func UpdateGlobalCustomField(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		fieldIDStr := vars["fieldId"]

		fieldID, err := uuid.Parse(fieldIDStr)
		if err != nil {
			writeErrorResponse(w, "Invalid field ID format", http.StatusBadRequest, r)
			return
		}

		// Check permissions
		currentUser, ok := middleware.GetUserFromContext(r.Context())
		if !ok {
			writeErrorResponse(w, "Unauthorized", http.StatusUnauthorized, r)
			return
		}
		_ = currentUser

		var updateReq models.UpdateGlobalCustomFieldRequest
		if err := json.NewDecoder(r.Body).Decode(&updateReq); err != nil {
			writeErrorResponse(w, "Invalid JSON format: "+err.Error(), http.StatusBadRequest, r)
			return
		}

		// Get existing field
		var existingField models.GlobalCustomField
		var optionsJSON, validationJSON []byte

		err = db.QueryRow(`
			SELECT id, name, label, type, required, options, validation, "order", is_active, created_at, updated_at
			FROM custom_fields WHERE id = $1
		`, fieldID).Scan(
			&existingField.ID, &existingField.Name, &existingField.Label, &existingField.Type,
			&existingField.Required, &optionsJSON, &validationJSON,
			&existingField.Order, &existingField.IsActive, &existingField.CreatedAt, &existingField.UpdatedAt,
		)

		if err != nil {
			if err == sql.ErrNoRows {
				writeErrorResponse(w, "Field not found", http.StatusNotFound, r)
				return
			}
			writeErrorResponse(w, err.Error(), http.StatusInternalServerError, r)
			return
		}

		// Parse existing JSON fields
		if len(optionsJSON) > 0 {
			existingField.OptionsFromJSON(optionsJSON)
		}
		if len(validationJSON) > 0 {
			existingField.ValidationFromJSON(validationJSON)
		}

		// Update fields
		if updateReq.Label != "" {
			existingField.Label = updateReq.Label
		}
		if updateReq.Type != "" {
			existingField.Type = updateReq.Type
		}
		if updateReq.Required != nil {
			existingField.Required = *updateReq.Required
		}
		if len(updateReq.Options) > 0 {
			existingField.Options = updateReq.Options
		}
		if updateReq.Validation != (models.FieldValidation{}) {
			existingField.Validation = updateReq.Validation
		}
		if updateReq.IsActive != nil {
			existingField.IsActive = *updateReq.IsActive
		}

		// Validate updated field
		if err := models.ValidateGlobalCustomField(&existingField); err != nil {
			writeErrorResponse(w, err.Error(), http.StatusBadRequest, r)
			return
		}

		// Convert JSON fields
		optionsJSON, err = existingField.OptionsToJSON()
		if err != nil {
			writeErrorResponse(w, "Failed to serialize options: "+err.Error(), http.StatusInternalServerError, r)
			return
		}

		validationJSON, err = existingField.ValidationToJSON()
		if err != nil {
			writeErrorResponse(w, "Failed to serialize validation: "+err.Error(), http.StatusInternalServerError, r)
			return
		}

		// Convert []byte to string for database update
		var optionsUpdStr, validationUpdStr sql.NullString

		if len(optionsJSON) > 0 {
			optionsUpdStr = sql.NullString{String: string(optionsJSON), Valid: true}
		} else {
			optionsUpdStr = sql.NullString{Valid: false}
		}
		if len(validationJSON) > 0 {
			validationUpdStr = sql.NullString{String: string(validationJSON), Valid: true}
		} else {
			validationUpdStr = sql.NullString{Valid: false}
		}

		// Update database
		_, err = db.Exec(`
			UPDATE custom_fields
			SET label = $1, type = $2, required = $3, options = $4, validation = $5, is_active = $6, updated_at = CURRENT_TIMESTAMP
			WHERE id = $7
		`, existingField.Label, existingField.Type, existingField.Required, optionsUpdStr, validationUpdStr, existingField.IsActive, fieldID)

		if err != nil {
			writeErrorResponse(w, "Failed to update field: "+err.Error(), http.StatusInternalServerError, r)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(existingField)
	}
}

// DeleteGlobalCustomField deletes a global custom field
func DeleteGlobalCustomField(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		fieldIDStr := vars["fieldId"]

		fieldID, err := uuid.Parse(fieldIDStr)
		if err != nil {
			writeErrorResponse(w, "Invalid field ID format", http.StatusBadRequest, r)
			return
		}

		// Check permissions
		currentUser, ok := middleware.GetUserFromContext(r.Context())
		if !ok {
			writeErrorResponse(w, "Unauthorized", http.StatusUnauthorized, r)
			return
		}
		_ = currentUser

		// Soft delete by setting is_active to false
		result, err := db.Exec("UPDATE custom_fields SET is_active = 0, updated_at = CURRENT_TIMESTAMP WHERE id = $1", fieldID)
		if err != nil {
			writeErrorResponse(w, "Failed to delete field: "+err.Error(), http.StatusInternalServerError, r)
			return
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			writeErrorResponse(w, "Failed to check deletion result: "+err.Error(), http.StatusInternalServerError, r)
			return
		}

		if rowsAffected == 0 {
			writeErrorResponse(w, "Field not found", http.StatusNotFound, r)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"message":  "Field deleted successfully",
			"field_id": fieldID,
		})
	}
}

// GetUserCustomFieldValues gets all custom field values for a specific user
func GetUserCustomFieldValues(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// LOG: Function entry point
		fmt.Printf("[DEBUG] GetUserCustomFieldValues: Starting function execution\n")

		vars := mux.Vars(r)
		userIDStr := vars["userId"]

		// LOG: User ID parsing
		fmt.Printf("[DEBUG] GetUserCustomFieldValues: Raw userID from URL: %s\n", userIDStr)

		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			fmt.Printf("[ERROR] GetUserCustomFieldValues: Failed to parse userID '%s': %v\n", userIDStr, err)
			writeErrorResponse(w, "Invalid user ID format", http.StatusBadRequest, r)
			return
		}

		fmt.Printf("[DEBUG] GetUserCustomFieldValues: Parsed userID: %s\n", userID.String())

		// Check if requesting user can access this data
		currentUser, ok := middleware.GetUserFromContext(r.Context())
		if !ok {
			fmt.Printf("[ERROR] GetUserCustomFieldValues: No user found in context\n")
			writeErrorResponse(w, "Unauthorized", http.StatusUnauthorized, r)
			return
		}

		fmt.Printf("[DEBUG] GetUserCustomFieldValues: Current user ID: %s\n", currentUser.ID.String())

		// users can only access their own data, or admins can access any data
		if currentUser.ID != userID {
			fmt.Printf("[WARN] GetUserCustomFieldValues: Access denied - current user %s trying to access user %s\n",
				currentUser.ID.String(), userID.String())
			// TODO: Add admin permission check
			writeErrorResponse(w, "Access denied", http.StatusForbidden, r)
			return
		}

		fmt.Printf("[DEBUG] GetUserCustomFieldValues: Access granted for user %s\n", userID.String())

		// LOG: SQL Query preparation
		query := `
			SELECT cf.id, cf.name, cf.label, cf.type, cf.required, cf.options, cf.validation,
				   cf."order", cf.is_active, cf.created_at, cf.updated_at,
				   ucfv.value
			FROM custom_fields cf
			LEFT JOIN user_custom_field_values ucfv ON cf.id = ucfv.field_id AND ucfv.user_id = $1
			WHERE cf.is_active = true
			ORDER BY cf."order" ASC
		`

		fmt.Printf("[DEBUG] GetUserCustomFieldValues: Executing SQL query:\n%s\n", query)
		fmt.Printf("[DEBUG] GetUserCustomFieldValues: Query parameter userID: %s\n", userID.String())

		// Get all active global fields with user's values
		rows, err := db.Query(query, userID)

		if err != nil {
			fmt.Printf("[ERROR] GetUserCustomFieldValues: Database query failed: %v\n", err)
			fmt.Printf("[ERROR] GetUserCustomFieldValues: Query: %s\n", query)
			fmt.Printf("[ERROR] GetUserCustomFieldValues: Parameters: userID=%s\n", userID.String())
			writeErrorResponse(w, err.Error(), http.StatusInternalServerError, r)
			return
		}

		fmt.Printf("[DEBUG] GetUserCustomFieldValues: Query executed successfully\n")

		defer rows.Close()

		var fieldsWithValues []models.GlobalCustomFieldWithValue
		rowCount := 0

		for rows.Next() {
			rowCount++
			fmt.Printf("[DEBUG] GetUserCustomFieldValues: Processing row %d\n", rowCount)

			var field models.GlobalCustomFieldWithValue
			var optionsJSON, validationJSON []byte
			var value sql.NullString

			err := rows.Scan(
				&field.ID, &field.Name, &field.Label, &field.Type,
				&field.Required, &optionsJSON, &validationJSON,
				&field.Order, &field.IsActive, &field.CreatedAt, &field.UpdatedAt,
				&value,
			)
			if err != nil {
				fmt.Printf("[ERROR] GetUserCustomFieldValues: Failed to scan row %d: %v\n", rowCount, err)
				writeErrorResponse(w, err.Error(), http.StatusInternalServerError, r)
				return
			}

			fmt.Printf("[DEBUG] GetUserCustomFieldValues: Row %d scanned successfully - Field: %s (%s), is_active: %v\n",
				rowCount, field.Name, field.ID.String(), field.IsActive)

			// Parse JSON fields
			if len(optionsJSON) > 0 {
				fmt.Printf("[DEBUG] GetUserCustomFieldValues: Parsing options JSON for field %s\n", field.Name)
				field.OptionsFromJSON(optionsJSON)
			}
			if len(validationJSON) > 0 {
				fmt.Printf("[DEBUG] GetUserCustomFieldValues: Parsing validation JSON for field %s\n", field.Name)
				field.ValidationFromJSON(validationJSON)
			}

			// Set value if exists
			if value.Valid {
				field.Value = &value.String
				fmt.Printf("[DEBUG] GetUserCustomFieldValues: Field %s has value: %s\n", field.Name, value.String)
			} else {
				fmt.Printf("[DEBUG] GetUserCustomFieldValues: Field %s has no value (NULL)\n", field.Name)
			}

			fieldsWithValues = append(fieldsWithValues, field)
		}

		// Check for any iteration errors
		if err := rows.Err(); err != nil {
			fmt.Printf("[ERROR] GetUserCustomFieldValues: Error during row iteration: %v\n", err)
			writeErrorResponse(w, err.Error(), http.StatusInternalServerError, r)
			return
		}

		fmt.Printf("[DEBUG] GetUserCustomFieldValues: Successfully processed %d rows\n", rowCount)
		fmt.Printf("[DEBUG] GetUserCustomFieldValues: Returning %d fields with values\n", len(fieldsWithValues))

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(fieldsWithValues)
	}
}

// UpdateUserCustomFieldValues updates multiple custom field values for a user
func UpdateUserCustomFieldValues(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		userIDStr := vars["userId"]

		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			writeErrorResponse(w, "Invalid user ID format", http.StatusBadRequest, r)
			return
		}

		// Check if requesting user can modify this data
		currentUser, ok := middleware.GetUserFromContext(r.Context())
		if !ok {
			writeErrorResponse(w, "Unauthorized", http.StatusUnauthorized, r)
			return
		}

		// users can only modify their own data, or admins can modify any data
		if currentUser.ID != userID {
			// TODO: Add admin permission check
			writeErrorResponse(w, "Access denied", http.StatusForbidden, r)
			return
		}

		var updateReq models.UserCustomFieldValuesRequest
		if err := json.NewDecoder(r.Body).Decode(&updateReq); err != nil {
			writeErrorResponse(w, "Invalid JSON format: "+err.Error(), http.StatusBadRequest, r)
			return
		}

		// Process each field value
		for fieldIDStr, value := range updateReq.FieldValues {
			fieldID, err := uuid.Parse(fieldIDStr)
			if err != nil {
				writeErrorResponse(w, "Invalid field ID format: "+fieldIDStr, http.StatusBadRequest, r)
				return
			}

			// Get field definition to validate
			var field models.GlobalCustomField
			var optionsJSON, validationJSON []byte

			err = db.QueryRow(`
				SELECT id, name, label, type, required, options, validation, "order", is_active
				FROM custom_fields WHERE id = $1 AND is_active = true
			`, fieldID).Scan(
				&field.ID, &field.Name, &field.Label, &field.Type,
				&field.Required, &optionsJSON, &validationJSON,
				&field.Order, &field.IsActive,
			)

			if err != nil {
				if err == sql.ErrNoRows {
					writeErrorResponse(w, "Field not found: "+fieldIDStr, http.StatusBadRequest, r)
					return
				}
				writeErrorResponse(w, err.Error(), http.StatusInternalServerError, r)
				return
			}

			// Parse JSON fields
			if len(optionsJSON) > 0 {
				field.OptionsFromJSON(optionsJSON)
			}
			if len(validationJSON) > 0 {
				field.ValidationFromJSON(validationJSON)
			}

			// RE-ENABLE VALIDATION AFTER DEBUGGING
			_ = value // Prevent unused variable error

			// Convert value to string for storage
			var valueStr string
			if value != nil {
				switch v := value.(type) {
				case string:
					valueStr = v
				case []any:
					// For multiselect, convert to JSON string
					jsonBytes, _ := json.Marshal(v)
					valueStr = string(jsonBytes)
				case bool:
					// Handle boolean values - FORCE to string representation
					valueStr = strconv.FormatBool(v) // This will give "true" or "false"
				case float64:
					// Handle number values
					if v == float64(int(v)) {
						// If it's a whole number, store as integer string
						valueStr = strconv.Itoa(int(v))
					} else {
						// If it's a decimal, convert to string
						valueStr = strconv.FormatFloat(v, 'f', -1, 64)
					}
				default:
					// Fallback for any other unexpected types
					// Try to convert to string using fmt.Sprint
					valueStr = fmt.Sprintf("%v", v)
				}
			}

			// Upsert value
			_, err = db.Exec(`
				INSERT INTO user_custom_field_values (id, user_id, field_id, value, created_at, updated_at)
				VALUES ($1, $2, $3, $4, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
				ON CONFLICT(user_id, field_id) DO UPDATE SET
				value = excluded.value,
				updated_at = CURRENT_TIMESTAMP
			`, uuid.New(), userID, fieldID, valueStr)

			if err != nil {
				writeErrorResponse(w, "Failed to update field value: "+err.Error(), http.StatusInternalServerError, r)
				return
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"message": "Custom field values updated successfully",
		})
	}
}

// UploadFile handles file uploads for custom fields
func UploadFile(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check if user has permission
		currentUser, ok := middleware.GetUserFromContext(r.Context())
		if !ok {
			writeErrorResponse(w, "Unauthorized", http.StatusUnauthorized, r)
			return
		}
		_ = currentUser // TODO: Add permission check for file uploads

		// Parse multipart form
		err := r.ParseMultipartForm(32 << 20) // 32 MB max
		if err != nil {
			writeErrorResponse(w, "Failed to parse multipart form: "+err.Error(), http.StatusBadRequest, r)
			return
		}

		file, header, err := r.FormFile("file")
		if err != nil {
			writeErrorResponse(w, "Failed to get file from form: "+err.Error(), http.StatusBadRequest, r)
			return
		}
		defer file.Close()

		// Validate file size (max 10MB)
		if header.Size > 10<<20 {
			writeErrorResponse(w, "File too large. Maximum size is 10MB", http.StatusBadRequest, r)
			return
		}

		// Generate unique filename
		fileID := uuid.New()
		fileExt := filepath.Ext(header.Filename)
		filename := fileID.String() + fileExt

		// Create uploads directory if it doesn't exist
		uploadDir := "./uploads"
		if err := os.MkdirAll(uploadDir, 0755); err != nil {
			writeErrorResponse(w, "Failed to create upload directory: "+err.Error(), http.StatusInternalServerError, r)
			return
		}

		// Save file to disk
		filePath := filepath.Join(uploadDir, filename)
		dst, err := os.Create(filePath)
		if err != nil {
			writeErrorResponse(w, "Failed to create file: "+err.Error(), http.StatusInternalServerError, r)
			return
		}
		defer dst.Close()

		if _, err := io.Copy(dst, file); err != nil {
			writeErrorResponse(w, "Failed to save file: "+err.Error(), http.StatusInternalServerError, r)
			return
		}

		// Return file information
		fileInfo := map[string]any{
			"id":   fileID.String(),
			"name": header.Filename,
			"size": header.Size,
			"type": header.Header.Get("Content-Type"),
			"path": filePath,
			"url":  "/uploads/" + filename,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(fileInfo)
	}
}
