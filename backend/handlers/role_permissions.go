package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"pillow/models"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// AssignPermissionRequest represents the request payload for assigning a permission to a role
type AssignPermissionRequest struct {
	PermissionID uuid.UUID `json:"permission_id"`
}

// GetRolePermissions retrieves all permissions for a specific role
func GetRolePermissions(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		roleIDStr := vars["roleId"]

		roleID, err := uuid.Parse(roleIDStr)
		if err != nil {
			writeErrorResponse(w, "Invalid role ID format", http.StatusBadRequest, r)
			return
		}

		// Check if role exists
		var role models.Role
		err = db.QueryRow("SELECT id, name FROM \"roles\" WHERE id = $1", roleID).Scan(&role.ID, &role.Name)
		if err != nil {
			if err == sql.ErrNoRows {
				writeErrorResponse(w, "Role not found", http.StatusNotFound, r)
				return
			}
			writeErrorResponse(w, err.Error(), http.StatusInternalServerError, r)
			return
		}

		// Get permissions for this role
		rows, err := db.Query(`
			SELECT p.id, p.name, p.description, p.scope_level, p.created_at, p.updated_at
			FROM "permissions" p
			INNER JOIN "role_permissions" rp ON p.id = rp.permission_id
			WHERE rp.role_id = $1
			ORDER BY p.name
		`, roleID)

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

		response := models.RoleWithPermissions{
			Role:        role,
			Permissions: permissions,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

// AssignPermissionToRole assigns a permission to a role
func AssignPermissionToRole(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		roleIDStr := vars["roleId"]

		roleID, err := uuid.Parse(roleIDStr)
		if err != nil {
			writeErrorResponse(w, "Invalid role ID format", http.StatusBadRequest, r)
			return
		}

		var req AssignPermissionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeErrorResponse(w, "Invalid JSON format: "+err.Error(), http.StatusBadRequest, r)
			return
		}

		// Check if role exists
		var role models.Role
		err = db.QueryRow("SELECT id, name FROM \"roles\" WHERE id = $1", roleID).Scan(&role.ID, &role.Name)
		if err != nil {
			if err == sql.ErrNoRows {
				writeErrorResponse(w, "Role not found", http.StatusNotFound, r)
				return
			}
			writeErrorResponse(w, err.Error(), http.StatusInternalServerError, r)
			return
		}

		// Check if permission exists
		var permission models.Permission
		err = db.QueryRow("SELECT id, name FROM \"permissions\" WHERE id = $1", req.PermissionID).Scan(&permission.ID, &permission.Name)
		if err != nil {
			if err == sql.ErrNoRows {
				writeErrorResponse(w, "Permission not found", http.StatusNotFound, r)
				return
			}
			writeErrorResponse(w, err.Error(), http.StatusInternalServerError, r)
			return
		}

		// Check if assignment already exists
		var existingID uuid.UUID
		err = db.QueryRow("SELECT role_id FROM \"role_permissions\" WHERE role_id = $1 AND permission_id = $2",
			roleID, req.PermissionID).Scan(&existingID)

		if err == nil {
			writeErrorResponse(w, "Permission is already assigned to this role", http.StatusConflict, r)
			return
		} else if err != sql.ErrNoRows {
			writeErrorResponse(w, err.Error(), http.StatusInternalServerError, r)
			return
		}

		// Create the assignment
		_, err = db.Exec("INSERT INTO \"role_permissions\" (role_id, permission_id, created_at, updated_at) VALUES ($1, $2, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)",
			roleID, req.PermissionID)

		if err != nil {
			writeErrorResponse(w, "Failed to assign permission to role: "+err.Error(), http.StatusInternalServerError, r)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message":       "Permission assigned to role successfully",
			"role_id":       roleID,
			"permission_id": req.PermissionID,
		})
	}
}

// RemovePermissionFromRole removes a permission from a role
func RemovePermissionFromRole(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		roleIDStr := vars["roleId"]
		permissionIDStr := vars["permissionId"]

		roleID, err := uuid.Parse(roleIDStr)
		if err != nil {
			writeErrorResponse(w, "Invalid role ID format", http.StatusBadRequest, r)
			return
		}

		permissionID, err := uuid.Parse(permissionIDStr)
		if err != nil {
			writeErrorResponse(w, "Invalid permission ID format", http.StatusBadRequest, r)
			return
		}

		// Check if assignment exists
		var existingRoleID uuid.UUID
		err = db.QueryRow("SELECT role_id FROM \"role_permissions\" WHERE role_id = $1 AND permission_id = $2",
			roleID, permissionID).Scan(&existingRoleID)

		if err != nil {
			if err == sql.ErrNoRows {
				writeErrorResponse(w, "Permission is not assigned to this role", http.StatusNotFound, r)
				return
			}
			writeErrorResponse(w, err.Error(), http.StatusInternalServerError, r)
			return
		}

		// Remove the assignment
		_, err = db.Exec("DELETE FROM \"role_permissions\" WHERE role_id = $1 AND permission_id = $2",
			roleID, permissionID)

		if err != nil {
			writeErrorResponse(w, "Failed to remove permission from role: "+err.Error(), http.StatusInternalServerError, r)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message":       "Permission removed from role successfully",
			"role_id":       roleID,
			"permission_id": permissionID,
		})
	}
}

// GetPermissionRoles retrieves all roles that have a specific permission
func GetPermissionRoles(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		permissionIDStr := vars["permissionId"]

		permissionID, err := uuid.Parse(permissionIDStr)
		if err != nil {
			writeErrorResponse(w, "Invalid permission ID format", http.StatusBadRequest, r)
			return
		}

		// Check if permission exists
		var permission models.Permission
		err = db.QueryRow("SELECT id, name FROM \"permissions\" WHERE id = $1", permissionID).Scan(&permission.ID, &permission.Name)
		if err != nil {
			if err == sql.ErrNoRows {
				writeErrorResponse(w, "Permission not found", http.StatusNotFound, r)
				return
			}
			writeErrorResponse(w, err.Error(), http.StatusInternalServerError, r)
			return
		}

		// Get roles that have this permission
		rows, err := db.Query(`
			SELECT r.id, r.name, r.description, r.created_at, r.updated_at
			FROM "roles" r
			INNER JOIN "role_permissions" rp ON r.id = rp.role_id
			WHERE rp.permission_id = $1
			ORDER BY r.name
		`, permissionID)

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

		response := map[string]interface{}{
			"permission": permission,
			"roles":      roles,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}
