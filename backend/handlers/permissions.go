package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"pillow/middleware"
	"pillow/models"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// CreatePermissionRequest represents the request payload for creating a permission
type CreatePermissionRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	ScopeLevel  string `json:"scope_level,omitempty"`
}

// UpdatePermissionRequest represents the request payload for updating a permission
type UpdatePermissionRequest struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	ScopeLevel  string `json:"scope_level,omitempty"`
}

// GetPermissions retrieves all permissions
func GetPermissions(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT id, name, description, scope_level, created_at, updated_at FROM \"Permissions\" ORDER BY name")
		if err != nil {
			writeErrorResponse(w, err.Error(), http.StatusInternalServerError, r)
			return
		}
		defer rows.Close()

		var permissions []models.Permission
		for rows.Next() {
			var permission models.Permission
			err := rows.Scan(&permission.ID, &permission.Name, &permission.Description, &permission.ScopeLevel, &permission.CreatedAt, &permission.UpdatedAt)
			if err != nil {
				writeErrorResponse(w, err.Error(), http.StatusInternalServerError, r)
				return
			}
			permissions = append(permissions, permission)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(permissions)
	}
}

// GetPermission retrieves a single permission by ID
func GetPermission(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idStr := vars["id"]

		permissionID, err := uuid.Parse(idStr)
		if err != nil {
			writeErrorResponse(w, "Invalid permission ID format", http.StatusBadRequest, r)
			return
		}

		var permission models.Permission
		err = db.QueryRow("SELECT id, name, description, scope_level, created_at, updated_at FROM \"Permissions\" WHERE id = $1",
			permissionID).Scan(&permission.ID, &permission.Name, &permission.Description, &permission.ScopeLevel, &permission.CreatedAt, &permission.UpdatedAt)

		if err != nil {
			if err == sql.ErrNoRows {
				writeErrorResponse(w, "Permission not found", http.StatusNotFound, r)
				return
			}
			writeErrorResponse(w, err.Error(), http.StatusInternalServerError, r)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(permission)
	}
}

// CreatePermission creates a new permission
func CreatePermission(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreatePermissionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeErrorResponse(w, "Invalid JSON format: "+err.Error(), http.StatusBadRequest, r)
			return
		}

		if strings.TrimSpace(req.Name) == "" {
			writeErrorResponse(w, "Permission name is required", http.StatusBadRequest, r)
			return
		}

		// Set default scope level if not provided
		scopeLevel := strings.TrimSpace(req.ScopeLevel)
		if scopeLevel == "" {
			scopeLevel = "user"
		}

		// Check if permission name already exists
		var existingID uuid.UUID
		err := db.QueryRow("SELECT id FROM \"Permissions\" WHERE name = $1", req.Name).Scan(&existingID)
		if err == nil {
			writeErrorResponse(w, "Permission with this name already exists", http.StatusConflict, r)
			return
		} else if err != sql.ErrNoRows {
			writeErrorResponse(w, err.Error(), http.StatusInternalServerError, r)
			return
		}

		// Create new permission
		permissionID := uuid.New()
		_, err = db.Exec("INSERT INTO \"Permissions\" (id, name, description, scope_level, created_at, updated_at) VALUES ($1, $2, $3, $4, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)",
			permissionID, strings.TrimSpace(req.Name), strings.TrimSpace(req.Description), scopeLevel)

		if err != nil {
			writeErrorResponse(w, "Failed to create permission: "+err.Error(), http.StatusInternalServerError, r)
			return
		}

		// Get the created permission
		var permission models.Permission
		err = db.QueryRow("SELECT id, name, description, scope_level, created_at, updated_at FROM \"Permissions\" WHERE id = $1",
			permissionID).Scan(&permission.ID, &permission.Name, &permission.Description, &permission.ScopeLevel, &permission.CreatedAt, &permission.UpdatedAt)

		if err != nil {
			writeErrorResponse(w, "Failed to retrieve created permission: "+err.Error(), http.StatusInternalServerError, r)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		// Prepare audit details and expose via response headers
		details := map[string]interface{}{
			"permission_after":  permission,
			"permission_before": nil,
			"action": map[string]interface{}{
				"method":     r.Method,
				"path":       r.URL.Path,
				"actor_id":   nil,
				"ip_address": r.RemoteAddr,
			},
		}
		if user, ok := middleware.GetUserFromContext(r.Context()); ok && user != nil {
			details["action"].(map[string]interface{})["actor_id"] = user.ID.String()
		}
		detBytes, _ := json.Marshal(details)
		w.Header().Set("X-Audit-Action", "PERMISSION_CREATED")
		w.Header().Set("X-Audit-Details", string(detBytes))

		json.NewEncoder(w).Encode(map[string]interface{}{
			"message":    "Permission created successfully",
			"permission": permission,
		})
	}
}

// UpdatePermission updates an existing permission
func UpdatePermission(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idStr := vars["id"]

		permissionID, err := uuid.Parse(idStr)
		if err != nil {
			writeErrorResponse(w, "Invalid permission ID format", http.StatusBadRequest, r)
			return
		}

		var req UpdatePermissionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeErrorResponse(w, "Invalid JSON format: "+err.Error(), http.StatusBadRequest, r)
			return
		}

		// Check if permission exists
		var existingPermission models.Permission
		err = db.QueryRow("SELECT id, name, description, scope_level, created_at, updated_at FROM \"Permissions\" WHERE id = $1",
			permissionID).Scan(&existingPermission.ID, &existingPermission.Name, &existingPermission.Description, &existingPermission.ScopeLevel, &existingPermission.CreatedAt, &existingPermission.UpdatedAt)

		if err != nil {
			if err == sql.ErrNoRows {
				writeErrorResponse(w, "Permission not found", http.StatusNotFound, r)
				return
			}
			writeErrorResponse(w, err.Error(), http.StatusInternalServerError, r)
			return
		}

		// Check for name conflicts if name is being updated
		if req.Name != "" && req.Name != existingPermission.Name {
			var existingID uuid.UUID
			err := db.QueryRow("SELECT id FROM \"Permissions\" WHERE name = $1 AND id != $2",
				req.Name, permissionID).Scan(&existingID)
			if err == nil {
				writeErrorResponse(w, "Permission with this name already exists", http.StatusConflict, r)
				return
			} else if err != sql.ErrNoRows {
				writeErrorResponse(w, err.Error(), http.StatusInternalServerError, r)
				return
			}
		}

		// Build update query
		setParts := []string{}
		args := []interface{}{}
		argCount := 1

		if req.Name != "" {
			setParts = append(setParts, "name = $"+string(rune('0'+argCount)))
			args = append(args, strings.TrimSpace(req.Name))
			argCount++
		}

		if req.Description != "" {
			setParts = append(setParts, "description = $"+string(rune('0'+argCount)))
			args = append(args, strings.TrimSpace(req.Description))
			argCount++
		}

		if req.ScopeLevel != "" {
			setParts = append(setParts, "scope_level = $"+string(rune('0'+argCount)))
			args = append(args, strings.TrimSpace(req.ScopeLevel))
			argCount++
		}

		if len(setParts) == 0 {
			writeErrorResponse(w, "No fields to update", http.StatusBadRequest, r)
			return
		}

		setParts = append(setParts, "updated_at = CURRENT_TIMESTAMP")
		query := "UPDATE \"Permissions\" SET " + strings.Join(setParts, ", ") + " WHERE id = $" + string(rune('0'+argCount))
		args = append(args, permissionID)

		_, err = db.Exec(query, args...)
		if err != nil {
			writeErrorResponse(w, "Failed to update permission: "+err.Error(), http.StatusInternalServerError, r)
			return
		}

		// Get updated permission
		var updatedPermission models.Permission
		err = db.QueryRow("SELECT id, name, description, scope_level, created_at, updated_at FROM \"Permissions\" WHERE id = $1",
			permissionID).Scan(&updatedPermission.ID, &updatedPermission.Name, &updatedPermission.Description, &updatedPermission.ScopeLevel, &updatedPermission.CreatedAt, &updatedPermission.UpdatedAt)

		if err != nil {
			writeErrorResponse(w, "Failed to retrieve updated permission: "+err.Error(), http.StatusInternalServerError, r)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		// Prepare audit details and expose via response headers
		details := map[string]interface{}{
			"permission_after":  updatedPermission,
			"permission_before": existingPermission,
			"action": map[string]interface{}{
				"method":     r.Method,
				"path":       r.URL.Path,
				"actor_id":   nil,
				"ip_address": r.RemoteAddr,
			},
		}
		if user, ok := middleware.GetUserFromContext(r.Context()); ok && user != nil {
			details["action"].(map[string]interface{})["actor_id"] = user.ID.String()
		}
		detBytes, _ := json.Marshal(details)
		w.Header().Set("X-Audit-Action", "PERMISSION_UPDATED")
		w.Header().Set("X-Audit-Details", string(detBytes))

		json.NewEncoder(w).Encode(map[string]interface{}{
			"message":    "Permission updated successfully",
			"permission": updatedPermission,
		})
	}
}

// DeletePermission deletes a permission
func DeletePermission(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idStr := vars["id"]

		permissionID, err := uuid.Parse(idStr)
		if err != nil {
			writeErrorResponse(w, "Invalid permission ID format", http.StatusBadRequest, r)
			return
		}

		// Check if permission exists
		var permission models.Permission
		err = db.QueryRow("SELECT id, name FROM \"Permissions\" WHERE id = $1",
			permissionID).Scan(&permission.ID, &permission.Name)

		if err != nil {
			if err == sql.ErrNoRows {
				writeErrorResponse(w, "Permission not found", http.StatusNotFound, r)
				return
			}
			writeErrorResponse(w, err.Error(), http.StatusInternalServerError, r)
			return
		}

		// Check if permission is being used by any roles
		var roleCount int
		err = db.QueryRow("SELECT COUNT(*) FROM \"Role_Permissions\" WHERE permission_id = $1", permissionID).Scan(&roleCount)
		if err != nil {
			writeErrorResponse(w, err.Error(), http.StatusInternalServerError, r)
			return
		}

		if roleCount > 0 {
			writeErrorResponse(w, "Cannot delete permission that is assigned to roles", http.StatusConflict, r)
			return
		}

		// Delete permission
		_, err = db.Exec("DELETE FROM \"Permissions\" WHERE id = $1", permissionID)
		if err != nil {
			writeErrorResponse(w, "Failed to delete permission: "+err.Error(), http.StatusInternalServerError, r)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		// Prepare audit details and expose via response headers
		details := map[string]interface{}{
			"permission": nil,
			"action": map[string]interface{}{
				"method":     r.Method,
				"path":       r.URL.Path,
				"actor_id":   nil,
				"ip_address": r.RemoteAddr,
			},
		}
		if user, ok := middleware.GetUserFromContext(r.Context()); ok && user != nil {
			details["action"].(map[string]interface{})["actor_id"] = user.ID.String()
		}
		detBytes, _ := json.Marshal(details)
		w.Header().Set("X-Audit-Action", "PERMISSION_DELETED")
		w.Header().Set("X-Audit-Details", string(detBytes))

		json.NewEncoder(w).Encode(map[string]interface{}{
			"message":       "Permission deleted successfully",
			"permission_id": permissionID,
		})
	}
}
