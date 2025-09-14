package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"pillow/auth"
	"pillow/middleware"
	"pillow/models"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq" // if needed
)

func GetUsers(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT id, username, email, is_active, created_at, updated_at FROM \"Users\" WHERE is_active = true")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var users []models.User
		for rows.Next() {
			var user models.User
			if err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.IsActive, &user.CreatedAt, &user.UpdatedAt); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			users = append(users, user)
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(users); err != nil {
			http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
		}
	}
}

func CreateUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, "Invalid JSON format: "+err.Error(), http.StatusBadRequest)
			return
		}

		// Debug: Log received data
		// fmt.Printf("Received user data: Username=%s, Email=%s, PasswordHash=%s\n", user.Username, user.Email, user.PasswordHash)

		if user.Username == "" {
			http.Error(w, "Username is required", http.StatusBadRequest)
			return
		}
		if user.Email == "" {
			http.Error(w, "Email is required", http.StatusBadRequest)
			return
		}
		if user.PasswordHash == "" {
			http.Error(w, "Password is required", http.StatusBadRequest)
			return
		}

		// Hash the password using bcrypt
		hashedPassword, err := auth.HashPassword(user.PasswordHash)
		if err != nil {
			http.Error(w, "Failed to hash password", http.StatusInternalServerError)
			return
		}

		// Generate UUID for new user
		user.ID = uuid.New()

		// Set user as active by default
		user.IsActive = true

		err = db.QueryRow("INSERT INTO \"Users\" (id, username, password_hash, email, is_active, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP) RETURNING id",
			user.ID, user.Username, hashedPassword, user.Email, user.IsActive).Scan(&user.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{"message": "User created", "user": user})

	}
}

// LoginRequest represents the login request payload
type LoginRequest struct {
	Identifier string `json:"identifier"` // Can be username or email
	Password   string `json:"password"`
}

// LoginResponse represents the login response
type LoginResponse struct {
	Token string      `json:"token"`
	User  models.User `json:"user"`
}

// UpdateUserRequest represents the update user request payload
type UpdateUserRequest struct {
	Username string `json:"username,omitempty"`
	Email    string `json:"email,omitempty"`
	IsActive *bool  `json:"is_active,omitempty"`
}

func Login(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var loginReq LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		if loginReq.Identifier == "" || loginReq.Password == "" {
			http.Error(w, "Identifier (username or email) and password required", http.StatusBadRequest)
			return
		}

		// Get user from database - check both username and email
		var user models.User
		err := db.QueryRow("SELECT id, username, password_hash, email, is_active, created_at, updated_at FROM \"Users\" WHERE username = $1 OR email = $1",
			loginReq.Identifier).Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Email, &user.IsActive, &user.CreatedAt, &user.UpdatedAt)

		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Invalid credentials", http.StatusUnauthorized)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Check if user is active
		if !user.IsActive {
			http.Error(w, "Account is deactivated", http.StatusUnauthorized)
			return
		}

		// Verify password
		if !auth.CheckPasswordHash(loginReq.Password, user.PasswordHash) {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		// Generate JWT token
		token, err := auth.GenerateJWT(user.ID, user.Username)
		if err != nil {
			http.Error(w, "Failed to generate token", http.StatusInternalServerError)
			return
		}

		// Clear password hash from response
		user.PasswordHash = ""

		response := LoginResponse{
			Token: token,
			User:  user,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

// GetUser retrieves a single user by ID
func GetUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idStr := vars["id"]

		userID, err := uuid.Parse(idStr)
		if err != nil {
			http.Error(w, "Invalid user ID format", http.StatusBadRequest)
			return
		}

		var user models.User
		err = db.QueryRow("SELECT id, username, email, is_active, created_at, updated_at FROM \"Users\" WHERE id = $1 AND is_active = true",
			userID).Scan(&user.ID, &user.Username, &user.Email, &user.IsActive, &user.CreatedAt, &user.UpdatedAt)

		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "User not found", http.StatusNotFound)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
	}
}

// GetUserProfile retrieves the current authenticated user's profile
func GetUserProfile(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := middleware.GetUserFromContext(r.Context())
		if !ok {
			http.Error(w, "User not found in context", http.StatusUnauthorized)
			return
		}

		// Get fresh user data from database
		var freshUser models.User
		err := db.QueryRow("SELECT id, username, email, is_active, created_at, updated_at FROM \"Users\" WHERE id = $1",
			user.ID).Scan(&freshUser.ID, &freshUser.Username, &freshUser.Email, &freshUser.IsActive, &freshUser.CreatedAt, &freshUser.UpdatedAt)

		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "User not found", http.StatusNotFound)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(freshUser)
	}
}

// UpdateUser updates a user's information
func UpdateUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idStr := vars["id"]

		userID, err := uuid.Parse(idStr)
		if err != nil {
			http.Error(w, "Invalid user ID format", http.StatusBadRequest)
			return
		}

		var updateReq UpdateUserRequest
		if err := json.NewDecoder(r.Body).Decode(&updateReq); err != nil {
			http.Error(w, "Invalid JSON format: "+err.Error(), http.StatusBadRequest)
			return
		}

		// Check if user exists
		var existingUser models.User
		err = db.QueryRow("SELECT id, username, email, is_active, created_at, updated_at FROM \"Users\" WHERE id = $1",
			userID).Scan(&existingUser.ID, &existingUser.Username, &existingUser.Email, &existingUser.IsActive, &existingUser.CreatedAt, &existingUser.UpdatedAt)

		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "User not found", http.StatusNotFound)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Build update query dynamically
		setParts := []string{}
		args := []interface{}{}
		argCount := 1

		if updateReq.Username != "" {
			setParts = append(setParts, "username = $"+strconv.Itoa(argCount))
			args = append(args, updateReq.Username)
			argCount++
		}

		if updateReq.Email != "" {
			setParts = append(setParts, "email = $"+strconv.Itoa(argCount))
			args = append(args, updateReq.Email)
			argCount++
		}

		if updateReq.IsActive != nil {
			setParts = append(setParts, "is_active = $"+strconv.Itoa(argCount))
			args = append(args, *updateReq.IsActive)
			argCount++
		}

		if len(setParts) == 0 {
			http.Error(w, "No fields to update", http.StatusBadRequest)
			return
		}

		query := "UPDATE \"Users\" SET " + strings.Join(setParts, ", ") + " WHERE id = $" + strconv.Itoa(argCount)
		args = append(args, userID)

		_, err = db.Exec(query, args...)
		if err != nil {
			http.Error(w, "Failed to update user: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Get updated user data
		var updatedUser models.User
		err = db.QueryRow("SELECT id, username, email, is_active, created_at, updated_at FROM \"Users\" WHERE id = $1",
			userID).Scan(&updatedUser.ID, &updatedUser.Username, &updatedUser.Email, &updatedUser.IsActive, &updatedUser.CreatedAt, &updatedUser.UpdatedAt)

		if err != nil {
			http.Error(w, "Failed to retrieve updated user: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		// Prepare audit details in the exact required structure and attach to context
		actorIDPtr := (*uuid.UUID)(nil)
		if actor, ok := middleware.GetUserFromContext(r.Context()); ok && actor != nil {
			actorIDPtr = &actor.ID
		}
		actionObj := map[string]interface{}{
			"method":     r.Method,
			"path":       r.URL.Path,
			"actor_id":   nil,
			"ip_address": r.RemoteAddr,
		}
		if actorIDPtr != nil {
			actionObj["actor_id"] = actorIDPtr.String()
		}

		details := map[string]interface{}{
			"user_after": map[string]interface{}{
				"id":         updatedUser.ID,
				"username":   updatedUser.Username,
				"email":      updatedUser.Email,
				"is_active":  updatedUser.IsActive,
				"created_at": updatedUser.CreatedAt,
			},
			"user_before": map[string]interface{}{
				"id":         existingUser.ID,
				"username":   existingUser.Username,
				"email":      existingUser.Email,
				"is_active":  existingUser.IsActive,
				"created_at": existingUser.CreatedAt,
			},
			"action": actionObj,
		}

		// Attach desired action name and details so middleware will write the DB row with these exact values.
		// Use response headers so middleware (which runs after handlers) can read them reliably.
		detBytes, _ := json.Marshal(details)
		w.Header().Set("X-Audit-Action", "USER_UPDATED")
		w.Header().Set("X-Audit-Details", string(detBytes))

		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "User updated successfully",
			"user":    updatedUser,
		})

		// Audit: handled by middleware (middleware will write the audit record).
	}
}

// DeleteUser performs a soft delete by setting is_active to false
func DeleteUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idStr := vars["id"]

		userID, err := uuid.Parse(idStr)
		if err != nil {
			http.Error(w, "Invalid user ID format", http.StatusBadRequest)
			return
		}

		// Check if user exists and is active
		var user models.User
		err = db.QueryRow("SELECT id, username, email FROM \"Users\" WHERE id = $1 AND is_active = true",
			userID).Scan(&user.ID, &user.Username, &user.Email)

		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "User not found or already deactivated", http.StatusNotFound)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Soft delete by setting is_active to false
		_, err = db.Exec("UPDATE \"Users\" SET is_active = false, updated_at = CURRENT_TIMESTAMP WHERE id = $1", userID)
		if err != nil {
			http.Error(w, "Failed to delete user: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "User deleted successfully",
			"user_id": userID,
		})

		// Audit: handled by middleware (middleware will write the audit record).
	}
}
