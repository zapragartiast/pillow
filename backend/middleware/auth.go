package middleware

import (
	"context"
	"database/sql"
	"net/http"
	"pillow/auth"
	"pillow/models"
	"strings"

	"github.com/gorilla/mux"
)

// Context keys
type contextKey string

const (
	UserContextKey contextKey = "user"
)

// AuthMiddleware validates JWT tokens and adds user info to request context
func AuthMiddleware(db *sql.DB) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Authorization header required", http.StatusUnauthorized)
				return
			}

			tokenString, err := auth.ExtractTokenFromHeader(authHeader)
			if err != nil {
				http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
				return
			}

			claims, err := auth.ValidateJWT(tokenString)
			if err != nil {
				http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
				return
			}

			// Get user from database to ensure they still exist and are active
			var user models.User
			err = db.QueryRow("SELECT id, username, email, is_active, created_at FROM \"users\" WHERE id = $1",
				claims.UserID).Scan(&user.ID, &user.Username, &user.Email, &user.IsActive, &user.CreatedAt)

			if err != nil {
				if err == sql.ErrNoRows {
					http.Error(w, "User not found", http.StatusUnauthorized)
					return
				}
				http.Error(w, "Database error", http.StatusInternalServerError)
				return
			}

			if !user.IsActive {
				http.Error(w, "Account is deactivated", http.StatusUnauthorized)
				return
			}

			// Add user to request context
			ctx := context.WithValue(r.Context(), UserContextKey, user)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		}
	}
}

// GetUserFromContext retrieves the user from request context
func GetUserFromContext(ctx context.Context) (*models.User, bool) {
	user, ok := ctx.Value(UserContextKey).(models.User)
	if !ok {
		return nil, false
	}
	return &user, true
}

// OptionalAuthMiddleware allows requests without authentication but adds user info if token is provided
func OptionalAuthMiddleware(db *sql.DB) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
				tokenString, err := auth.ExtractTokenFromHeader(authHeader)
				if err == nil {
					claims, err := auth.ValidateJWT(tokenString)
					if err == nil {
						// Get user from database
						var user models.User
						err = db.QueryRow("SELECT id, username, email, is_active, created_at FROM \"users\" WHERE id = $1",
							claims.UserID).Scan(&user.ID, &user.Username, &user.Email, &user.IsActive, &user.CreatedAt)

						if err == nil && user.IsActive {
							// Add user to request context
							ctx := context.WithValue(r.Context(), UserContextKey, user)
							r = r.WithContext(ctx)
						}
					}
				}
			}

			next.ServeHTTP(w, r)
		}
	}
}

// AuthMiddlewareMux creates Gorilla Mux compatible middleware
func AuthMiddlewareMux(db *sql.DB) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Authorization header required", http.StatusUnauthorized)
				return
			}

			tokenString, err := auth.ExtractTokenFromHeader(authHeader)
			if err != nil {
				http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
				return
			}

			claims, err := auth.ValidateJWT(tokenString)
			if err != nil {
				http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
				return
			}

			// Get user from database to ensure they still exist and are active
			var user models.User
			err = db.QueryRow("SELECT id, username, email, is_active, created_at FROM \"users\" WHERE id = $1",
				claims.UserID).Scan(&user.ID, &user.Username, &user.Email, &user.IsActive, &user.CreatedAt)

			if err != nil {
				if err == sql.ErrNoRows {
					http.Error(w, "User not found", http.StatusUnauthorized)
					return
				}
				http.Error(w, "Database error", http.StatusInternalServerError)
				return
			}

			if !user.IsActive {
				http.Error(w, "Account is deactivated", http.StatusUnauthorized)
				return
			}

			// Add user to request context
			ctx := context.WithValue(r.Context(), UserContextKey, user)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}
