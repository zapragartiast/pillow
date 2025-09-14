package middleware

import (
	"database/sql"
	"net/http"
	"pillow/models"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// GetUserRoles retrieves all roles for a given user
func GetUserRoles(db *sql.DB, userID uuid.UUID) ([]models.Role, error) {
	query := `
		SELECT r.id, r.name, r.description, r.created_at
		FROM "Roles" r
		INNER JOIN "User_Roles" ur ON r.id = ur.role_id
		WHERE ur.user_id = $1
	`

	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []models.Role
	for rows.Next() {
		var role models.Role
		err := rows.Scan(&role.ID, &role.Name, &role.Description, &role.CreatedAt)
		if err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}

	return roles, nil
}

// GetUserPermissions retrieves all permissions for a given user through their roles
func GetUserPermissions(db *sql.DB, userID uuid.UUID) ([]models.Permission, error) {
	query := `
		SELECT DISTINCT p.id, p.name, p.description, p.scope_level
		FROM "Permissions" p
		INNER JOIN "Role_Permissions" rp ON p.id = rp.permission_id
		INNER JOIN "User_Roles" ur ON rp.role_id = ur.role_id
		WHERE ur.user_id = $1
	`

	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []models.Permission
	for rows.Next() {
		var permission models.Permission
		err := rows.Scan(&permission.ID, &permission.Name, &permission.Description, &permission.ScopeLevel)
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, permission)
	}

	return permissions, nil
}

// HasPermission checks if a user has a specific permission
func HasPermission(db *sql.DB, userID uuid.UUID, permissionName string) (bool, error) {
	query := `
		SELECT COUNT(*) > 0
		FROM "Permissions" p
		INNER JOIN "Role_Permissions" rp ON p.id = rp.permission_id
		INNER JOIN "User_Roles" ur ON rp.role_id = ur.role_id
		WHERE ur.user_id = $1 AND p.name = $2
	`

	var hasPermission bool
	err := db.QueryRow(query, userID, permissionName).Scan(&hasPermission)
	if err != nil {
		return false, err
	}

	return hasPermission, nil
}

// HasRole checks if a user has a specific role
func HasRole(db *sql.DB, userID uuid.UUID, roleName string) (bool, error) {
	query := `
		SELECT COUNT(*) > 0
		FROM "Roles" r
		INNER JOIN "User_Roles" ur ON r.id = ur.role_id
		WHERE ur.user_id = $1 AND r.name = $2
	`

	var hasRole bool
	err := db.QueryRow(query, userID, roleName).Scan(&hasRole)
	if err != nil {
		return false, err
	}

	return hasRole, nil
}

// RequirePermission creates middleware that requires a specific permission
func RequirePermission(db *sql.DB, permissionName string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			user, ok := GetUserFromContext(r.Context())
			if !ok {
				http.Error(w, "Authentication required", http.StatusUnauthorized)
				return
			}

			hasPermission, err := HasPermission(db, user.ID, permissionName)
			if err != nil {
				http.Error(w, "Error checking permissions", http.StatusInternalServerError)
				return
			}

			if !hasPermission {
				http.Error(w, "Insufficient permissions", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		}
	}
}

// RequireRole creates middleware that requires a specific role
func RequireRole(db *sql.DB, roleName string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			user, ok := GetUserFromContext(r.Context())
			if !ok {
				http.Error(w, "Authentication required", http.StatusUnauthorized)
				return
			}

			hasRole, err := HasRole(db, user.ID, roleName)
			if err != nil {
				http.Error(w, "Error checking roles", http.StatusInternalServerError)
				return
			}

			if !hasRole {
				http.Error(w, "Insufficient role permissions", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		}
	}
}

// RequireAnyPermission creates middleware that requires any of the specified permissions
func RequireAnyPermission(db *sql.DB, permissionNames ...string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			user, ok := GetUserFromContext(r.Context())
			if !ok {
				http.Error(w, "Authentication required", http.StatusUnauthorized)
				return
			}

			hasAnyPermission := false
			for _, permissionName := range permissionNames {
				hasPermission, err := HasPermission(db, user.ID, permissionName)
				if err != nil {
					http.Error(w, "Error checking permissions", http.StatusInternalServerError)
					return
				}
				if hasPermission {
					hasAnyPermission = true
					break
				}
			}

			if !hasAnyPermission {
				http.Error(w, "Insufficient permissions", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		}
	}
}

// RequirePermissionMux creates Gorilla Mux compatible middleware that requires a specific permission
func RequirePermissionMux(db *sql.DB, permissionName string) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, ok := GetUserFromContext(r.Context())
			if !ok {
				http.Error(w, "Authentication required", http.StatusUnauthorized)
				return
			}

			hasPermission, err := HasPermission(db, user.ID, permissionName)
			if err != nil {
				http.Error(w, "Error checking permissions", http.StatusInternalServerError)
				return
			}

			if !hasPermission {
				http.Error(w, "Insufficient permissions", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RequireRoleMux creates Gorilla Mux compatible middleware that requires a specific role
func RequireRoleMux(db *sql.DB, roleName string) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, ok := GetUserFromContext(r.Context())
			if !ok {
				http.Error(w, "Authentication required", http.StatusUnauthorized)
				return
			}

			hasRole, err := HasRole(db, user.ID, roleName)
			if err != nil {
				http.Error(w, "Error checking roles", http.StatusInternalServerError)
				return
			}

			if !hasRole {
				http.Error(w, "Insufficient role permissions", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RequireAnyRole creates middleware that requires any of the specified roles
func RequireAnyRole(db *sql.DB, roleNames ...string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			user, ok := GetUserFromContext(r.Context())
			if !ok {
				http.Error(w, "Authentication required", http.StatusUnauthorized)
				return
			}

			hasAnyRole := false
			for _, roleName := range roleNames {
				hasRole, err := HasRole(db, user.ID, roleName)
				if err != nil {
					http.Error(w, "Error checking roles", http.StatusInternalServerError)
					return
				}
				if hasRole {
					hasAnyRole = true
					break
				}
			}

			if !hasAnyRole {
				http.Error(w, "Insufficient role permissions", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		}
	}
}
