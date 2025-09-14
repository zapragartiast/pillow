package routes

import (
	"database/sql"
	"net/http"
	"pillow/handlers"
	"pillow/middleware"

	cors "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// SetupRoutes configures all application routes
func SetupRoutes(db *sql.DB) http.Handler {
	r := mux.NewRouter()

	// API routes group
	api := r.PathPrefix("/api").Subrouter()

	// Public authentication routes
	api.HandleFunc("/register", handlers.CreateUser(db)).Methods("POST")
	api.HandleFunc("/login", handlers.Login(db)).Methods("POST")

	// Protected routes - require authentication
	protected := api.PathPrefix("").Subrouter()
	protected.Use(middleware.AuthMiddlewareMux(db))
	// Audit middleware must run after authentication so actor is available in context
	protected.Use(middleware.AuditMiddlewareMux(db))

	// User routes (protected)
	protected.HandleFunc("/users", handlers.GetUsers(db)).Methods("GET")
	protected.HandleFunc("/users/profile", handlers.GetUserProfile(db)).Methods("GET")
	protected.HandleFunc("/users/{id}", handlers.GetUser(db)).Methods("GET")

	// Admin-only routes - require specific permissions
	admin := protected.PathPrefix("").Subrouter()
	admin.Use(middleware.RequirePermissionMux(db, "manage_users"))

	admin.HandleFunc("/users/{id}", handlers.UpdateUser(db)).Methods("PUT")
	admin.HandleFunc("/users/{id}", handlers.DeleteUser(db)).Methods("DELETE")

	// Role management routes - require role management permission
	roleManager := protected.PathPrefix("").Subrouter()
	roleManager.Use(middleware.RequirePermissionMux(db, "manage_roles"))

	roleManager.HandleFunc("/roles", handlers.GetRoles(db)).Methods("GET")
	roleManager.HandleFunc("/roles", handlers.CreateRole(db)).Methods("POST")
	roleManager.HandleFunc("/roles/{id}", handlers.GetRole(db)).Methods("GET")
	roleManager.HandleFunc("/roles/{id}", handlers.UpdateRole(db)).Methods("PUT")
	roleManager.HandleFunc("/roles/{id}", handlers.DeleteRole(db)).Methods("DELETE")

	// Role-permission relationship management
	roleManager.HandleFunc("/roles/{roleId}/permissions", handlers.GetRolePermissions(db)).Methods("GET")
	roleManager.HandleFunc("/roles/{roleId}/permissions", handlers.AssignPermissionToRole(db)).Methods("POST")
	roleManager.HandleFunc("/roles/{roleId}/permissions/{permissionId}", handlers.RemovePermissionFromRole(db)).Methods("DELETE")

	// Permission management routes
	permissionManager := protected.PathPrefix("").Subrouter()
	permissionManager.Use(middleware.RequirePermissionMux(db, "manage_permissions"))

	permissionManager.HandleFunc("/permissions", handlers.GetPermissions(db)).Methods("GET")
	permissionManager.HandleFunc("/permissions", handlers.CreatePermission(db)).Methods("POST")
	permissionManager.HandleFunc("/permissions/{id}", handlers.GetPermission(db)).Methods("GET")
	permissionManager.HandleFunc("/permissions/{id}", handlers.UpdatePermission(db)).Methods("PUT")
	permissionManager.HandleFunc("/permissions/{id}", handlers.DeletePermission(db)).Methods("DELETE")

	// Permission-role relationship queries
	permissionManager.HandleFunc("/permissions/{permissionId}/roles", handlers.GetPermissionRoles(db)).Methods("GET")

	// Organization management routes
	orgManager := protected.PathPrefix("").Subrouter()
	orgManager.Use(middleware.RequirePermissionMux(db, "manage_organizations"))

	orgManager.HandleFunc("/organizations", handlers.GetOrganizations(db)).Methods("GET")
	orgManager.HandleFunc("/organizations", handlers.CreateOrganization(db)).Methods("POST")
	orgManager.HandleFunc("/organizations/{id}", handlers.UpdateOrganization(db)).Methods("PUT")
	orgManager.HandleFunc("/organizations/{id}", handlers.DeleteOrganization(db)).Methods("DELETE")

	return cors.CORS(
		cors.AllowedOrigins([]string{"http://localhost:3000"}),
		cors.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		cors.AllowedHeaders([]string{"Content-Type", "Authorization"}),
		cors.AllowCredentials(),
	)(r)
}
