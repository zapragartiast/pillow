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

// CreateRoleRequest represents the request payload for creating a role
type CreateRoleRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

// UpdateRoleRequest represents the request payload for updating a role
type UpdateRoleRequest struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

// GetRoles retrieves all roles
func GetRoles(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT id, name, description, created_at, updated_at FROM \"Roles\" ORDER BY name")
		if err != nil {
			writeErrorResponse(w, err.Error(), http.StatusInternalServerError, r)
			return
		}
		defer rows.Close()

		var roles []models.Role
		for rows.Next() {
			var role models.Role
			err := rows.Scan(&role.ID, &role.Name, &role.Description, &role.CreatedAt, &role.UpdatedAt)
			if err != nil {
				writeErrorResponse(w, err.Error(), http.StatusInternalServerError, r)
				return
			}
			roles = append(roles, role)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(roles)
	}
}

// GetRole retrieves a single role by ID
func GetRole(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idStr := vars["id"]

		roleID, err := uuid.Parse(idStr)
		if err != nil {
			writeErrorResponse(w, "Invalid role ID format", http.StatusBadRequest, r)
			return
		}

		var role models.Role
		err = db.QueryRow("SELECT id, name, description, created_at, updated_at FROM \"Roles\" WHERE id = $1",
			roleID).Scan(&role.ID, &role.Name, &role.Description, &role.CreatedAt, &role.UpdatedAt)

		if err != nil {
			if err == sql.ErrNoRows {
				writeErrorResponse(w, "Role not found", http.StatusNotFound, r)
				return
			}
			writeErrorResponse(w, err.Error(), http.StatusInternalServerError, r)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(role)
	}
}

// CreateRole creates a new role
func CreateRole(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateRoleRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeErrorResponse(w, "Invalid JSON format: "+err.Error(), http.StatusBadRequest, r)
			return
		}

		if strings.TrimSpace(req.Name) == "" {
			writeErrorResponse(w, "Role name is required", http.StatusBadRequest, r)
			return
		}

		// Check if role name already exists
		var existingID uuid.UUID
		err := db.QueryRow("SELECT id FROM \"Roles\" WHERE name = $1", req.Name).Scan(&existingID)
		if err == nil {
			writeErrorResponse(w, "Role with this name already exists", http.StatusConflict, r)
			return
		} else if err != sql.ErrNoRows {
			writeErrorResponse(w, err.Error(), http.StatusInternalServerError, r)
			return
		}

		// Create new role
		roleID := uuid.New()
		_, err = db.Exec("INSERT INTO \"Roles\" (id, name, description, created_at, updated_at) VALUES ($1, $2, $3, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)",
			roleID, strings.TrimSpace(req.Name), strings.TrimSpace(req.Description))

		if err != nil {
			writeErrorResponse(w, "Failed to create role: "+err.Error(), http.StatusInternalServerError, r)
			return
		}

		// Get the created role
		var role models.Role
		err = db.QueryRow("SELECT id, name, description, created_at, updated_at FROM \"Roles\" WHERE id = $1",
			roleID).Scan(&role.ID, &role.Name, &role.Description, &role.CreatedAt, &role.UpdatedAt)

		if err != nil {
			writeErrorResponse(w, "Failed to retrieve created role: "+err.Error(), http.StatusInternalServerError, r)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		// Prepare audit details and expose via response headers so middleware records structured audit.
		details := map[string]interface{}{
			"role_after":  role,
			"role_before": nil,
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
		w.Header().Set("X-Audit-Action", "ROLE_CREATED")
		w.Header().Set("X-Audit-Details", string(detBytes))

		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Role created successfully",
			"role":    role,
		})
	}
}

// UpdateRole updates an existing role
func UpdateRole(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idStr := vars["id"]

		roleID, err := uuid.Parse(idStr)
		if err != nil {
			writeErrorResponse(w, "Invalid role ID format", http.StatusBadRequest, r)
			return
		}

		var req UpdateRoleRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeErrorResponse(w, "Invalid JSON format: "+err.Error(), http.StatusBadRequest, r)
			return
		}

		// Check if role exists
		var existingRole models.Role
		err = db.QueryRow("SELECT id, name, description FROM \"Roles\" WHERE id = $1",
			roleID).Scan(&existingRole.ID, &existingRole.Name, &existingRole.Description)

		if err != nil {
			if err == sql.ErrNoRows {
				writeErrorResponse(w, "Role not found", http.StatusNotFound, r)
				return
			}
			writeErrorResponse(w, err.Error(), http.StatusInternalServerError, r)
			return
		}

		// Check for name conflicts if name is being updated
		if req.Name != "" && req.Name != existingRole.Name {
			var existingID uuid.UUID
			err := db.QueryRow("SELECT id FROM \"Roles\" WHERE name = $1 AND id != $2",
				req.Name, roleID).Scan(&existingID)
			if err == nil {
				writeErrorResponse(w, "Role with this name already exists", http.StatusConflict, r)
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

		if len(setParts) == 0 {
			writeErrorResponse(w, "No fields to update", http.StatusBadRequest, r)
			return
		}

		setParts = append(setParts, "updated_at = CURRENT_TIMESTAMP")
		query := "UPDATE \"Roles\" SET " + strings.Join(setParts, ", ") + " WHERE id = $" + string(rune('0'+argCount))
		args = append(args, roleID)

		_, err = db.Exec(query, args...)
		if err != nil {
			writeErrorResponse(w, "Failed to update role: "+err.Error(), http.StatusInternalServerError, r)
			return
		}

		// Get updated role
		var updatedRole models.Role
		err = db.QueryRow("SELECT id, name, description, created_at, updated_at FROM \"Roles\" WHERE id = $1",
			roleID).Scan(&updatedRole.ID, &updatedRole.Name, &updatedRole.Description, &updatedRole.CreatedAt, &updatedRole.UpdatedAt)

		if err != nil {
			writeErrorResponse(w, "Failed to retrieve updated role: "+err.Error(), http.StatusInternalServerError, r)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		// Prepare audit details and expose via response headers
		details := map[string]interface{}{
			"role_after":  updatedRole,
			"role_before": existingRole,
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
		w.Header().Set("X-Audit-Action", "ROLE_UPDATED")
		w.Header().Set("X-Audit-Details", string(detBytes))

		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Role updated successfully",
			"role":    updatedRole,
		})
	}
}

// DeleteRole deletes a role
func DeleteRole(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idStr := vars["id"]

		roleID, err := uuid.Parse(idStr)
		if err != nil {
			writeErrorResponse(w, "Invalid role ID format", http.StatusBadRequest, r)
			return
		}

		// Check if role exists
		var role models.Role
		err = db.QueryRow("SELECT id, name FROM \"Roles\" WHERE id = $1",
			roleID).Scan(&role.ID, &role.Name)

		if err != nil {
			if err == sql.ErrNoRows {
				writeErrorResponse(w, "Role not found", http.StatusNotFound, r)
				return
			}
			writeErrorResponse(w, err.Error(), http.StatusInternalServerError, r)
			return
		}

		// Check if role is being used by any users
		var userCount int
		err = db.QueryRow("SELECT COUNT(*) FROM \"User_Roles\" WHERE role_id = $1", roleID).Scan(&userCount)
		if err != nil {
			writeErrorResponse(w, err.Error(), http.StatusInternalServerError, r)
			return
		}

		if userCount > 0 {
			writeErrorResponse(w, "Cannot delete role that is assigned to users", http.StatusConflict, r)
			return
		}

		// Delete role
		_, err = db.Exec("DELETE FROM \"Roles\" WHERE id = $1", roleID)
		if err != nil {
			writeErrorResponse(w, "Failed to delete role: "+err.Error(), http.StatusInternalServerError, r)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		// Prepare audit details and expose via response headers
		details := map[string]interface{}{
			"role": role,
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
		w.Header().Set("X-Audit-Action", "ROLE_DELETED")
		w.Header().Set("X-Audit-Details", string(detBytes))

		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Role deleted successfully",
			"role_id": roleID,
		})
	}
}
