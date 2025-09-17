package routes

import (
	"database/sql"
	"net/http"
	"pillow/database"
	"pillow/handlers"
	"pillow/middleware"

	cors "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// SetupRoutes configures all application routes
func SetupRoutes(db interface{}, logger *logrus.Logger, isLoggingEnabled bool) http.Handler {
	// Type assertion to get *sql.DB from either *sql.DB or *LoggingDB
	var sqlDB *sql.DB
	if loggingDB, ok := db.(*database.LoggingDB); ok {
		sqlDB = loggingDB.DB
	} else if regularDB, ok := db.(*sql.DB); ok {
		sqlDB = regularDB
	} else {
		panic("db must be *sql.DB or *database.LoggingDB")
	}

	r := mux.NewRouter()

	// Add logging middleware to all routes
	r.Use(middleware.LoggingMiddlewareMux(logger, isLoggingEnabled))

	// API routes group
	api := r.PathPrefix("/api").Subrouter()

	// Public authentication routes
	api.HandleFunc("/register", handlers.CreateUser(sqlDB)).Methods("POST")
	api.HandleFunc("/login", handlers.Login(sqlDB)).Methods("POST")

	// Protected routes - require authentication
	protected := api.PathPrefix("").Subrouter()
	protected.Use(middleware.AuthMiddlewareMux(sqlDB))
	// Audit middleware must run after authentication so actor is available in context
	protected.Use(middleware.AuditMiddlewareMux(sqlDB))

	// User routes (protected)
	protected.HandleFunc("/users", handlers.GetUsers(sqlDB)).Methods("GET")
	protected.HandleFunc("/users/profile", handlers.GetUserProfile(sqlDB)).Methods("GET")

	protected.HandleFunc("/users/{id}", handlers.GetUser(sqlDB)).Methods("GET")

	// User custom field values (protected)
	protected.HandleFunc("/users/{userId}/custom-field-values", handlers.GetUserCustomFieldValues(sqlDB)).Methods("GET")
	protected.HandleFunc("/users/{userId}/custom-field-values", handlers.UpdateUserCustomFieldValues(sqlDB)).Methods("PUT")
	// Admin-only routes - require specific permissions
	admin := protected.PathPrefix("").Subrouter()
	admin.Use(middleware.RequirePermissionMux(sqlDB, "manage_users"))

	admin.HandleFunc("/users/{id}", handlers.UpdateUser(sqlDB)).Methods("PUT")
	admin.HandleFunc("/users/{id}", handlers.DeleteUser(sqlDB)).Methods("DELETE")
	// Global custom fields management - require admin permission
	// admin.HandleFunc("/global-custom-fields", handlers.GetGlobalCustomFields(sqlDB)).Methods("GET")
	// admin.HandleFunc("/global-custom-fields", handlers.CreateGlobalCustomField(sqlDB)).Methods("POST")
	// admin.HandleFunc("/global-custom-fields/{fieldId}", handlers.UpdateGlobalCustomField(sqlDB)).Methods("PUT")
	// admin.HandleFunc("/global-custom-fields/{fieldId}", handlers.DeleteGlobalCustomField(sqlDB)).Methods("DELETE")

	// Audit log routes - require admin permission for viewing transaction logs
	admin.HandleFunc("/audit-logs", handlers.GetAuditLogs(sqlDB)).Methods("GET")
	admin.HandleFunc("/audit-logs/{id}", handlers.GetAuditLog(sqlDB)).Methods("GET")

	// Role management routes - require role management permission
	roleManager := protected.PathPrefix("").Subrouter()
	roleManager.Use(middleware.RequirePermissionMux(sqlDB, "manage_roles"))

	roleManager.HandleFunc("/roles", handlers.GetRoles(sqlDB)).Methods("GET")
	roleManager.HandleFunc("/roles", handlers.CreateRole(sqlDB)).Methods("POST")
	roleManager.HandleFunc("/roles/{id}", handlers.GetRole(sqlDB)).Methods("GET")
	roleManager.HandleFunc("/roles/{id}", handlers.UpdateRole(sqlDB)).Methods("PUT")
	roleManager.HandleFunc("/roles/{id}", handlers.DeleteRole(sqlDB)).Methods("DELETE")

	// Role-permission relationship management
	roleManager.HandleFunc("/roles/{roleId}/permissions", handlers.GetRolePermissions(sqlDB)).Methods("GET")
	roleManager.HandleFunc("/roles/{roleId}/permissions", handlers.AssignPermissionToRole(sqlDB)).Methods("POST")
	roleManager.HandleFunc("/roles/{roleId}/permissions/{permissionId}", handlers.RemovePermissionFromRole(sqlDB)).Methods("DELETE")

	// Permission management routes
	permissionManager := protected.PathPrefix("").Subrouter()
	permissionManager.Use(middleware.RequirePermissionMux(sqlDB, "manage_permissions"))

	permissionManager.HandleFunc("/permissions", handlers.GetPermissions(sqlDB)).Methods("GET")
	permissionManager.HandleFunc("/permissions", handlers.CreatePermission(sqlDB)).Methods("POST")
	permissionManager.HandleFunc("/permissions/{id}", handlers.GetPermission(sqlDB)).Methods("GET")
	permissionManager.HandleFunc("/permissions/{id}", handlers.UpdatePermission(sqlDB)).Methods("PUT")
	permissionManager.HandleFunc("/permissions/{id}", handlers.DeletePermission(sqlDB)).Methods("DELETE")

	// Permission-role relationship queries
	permissionManager.HandleFunc("/permissions/{permissionId}/roles", handlers.GetPermissionRoles(sqlDB)).Methods("GET")

	// Custom fields management routes - require custom fields management permission
	customFieldsManager := protected.PathPrefix("").Subrouter()
	customFieldsManager.Use(middleware.RequirePermissionMux(sqlDB, "manage_custom_fields"))

	customFieldsManager.HandleFunc("/global-custom-fields", handlers.GetGlobalCustomFields(sqlDB)).Methods("GET")
	customFieldsManager.HandleFunc("/global-custom-fields", handlers.CreateGlobalCustomField(sqlDB)).Methods("POST")
	customFieldsManager.HandleFunc("/global-custom-fields/{fieldId}", handlers.UpdateGlobalCustomField(sqlDB)).Methods("PUT")
	customFieldsManager.HandleFunc("/global-custom-fields/{fieldId}", handlers.DeleteGlobalCustomField(sqlDB)).Methods("DELETE")
	customFieldsManager.HandleFunc("/upload", handlers.UploadFile(sqlDB)).Methods("POST")

	// Organization management routes
	orgManager := protected.PathPrefix("").Subrouter()
	orgManager.Use(middleware.RequirePermissionMux(sqlDB, "manage_organizations"))

	orgManager.HandleFunc("/organizations", handlers.GetOrganizations(sqlDB)).Methods("GET")
	orgManager.HandleFunc("/organizations", handlers.CreateOrganization(sqlDB)).Methods("POST")
	orgManager.HandleFunc("/organizations/{id}", handlers.UpdateOrganization(sqlDB)).Methods("PUT")
	orgManager.HandleFunc("/organizations/{id}", handlers.DeleteOrganization(sqlDB)).Methods("DELETE")

	// Static file server for uploaded files
	r.PathPrefix("/uploads/").Handler(http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads/"))))

	return cors.CORS(
		cors.AllowedOrigins([]string{"http://localhost:3000"}),
		cors.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		cors.AllowedHeaders([]string{"Content-Type", "Authorization"}),
		cors.AllowCredentials(),
	)(r)
}
