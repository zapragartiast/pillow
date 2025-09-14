package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

// context key type to avoid collisions
type vcontextKey string

const validatedBodyKey vcontextKey = "validated_body"

var validate = validator.New()

// ValidateBody returns a middleware that decodes the JSON body into an instance
// produced by targetFactory(), runs struct-tag validation using go-playground/validator,
// and on success stores the validated instance in request context under key validatedBodyKey.
//
// Usage:
//
//	// route-level: ValidateBody(func() interface{} { return &CreateUserRequest{} })(CreateUserHandler)
//	// inside handler:
//	var req CreateUserRequest
//	if !GetValidatedBody(r.Context(), &req) { // handle missing
//	    http.Error(w, "validated body missing", http.StatusBadRequest)
//	    return
//	}
//	// use req (already validated)
func ValidateBody(targetFactory func() interface{}) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		if targetFactory == nil {
			// passthrough if no factory provided
			return next
		}
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Only attempt JSON decoding for requests with a body and JSON content-type
			if r.Body == nil {
				http.Error(w, "request body required", http.StatusBadRequest)
				return
			}
			ct := r.Header.Get("Content-Type")
			if ct == "" || !strings.Contains(ct, "application/json") {
				http.Error(w, "Content-Type must be application/json", http.StatusBadRequest)
				return
			}

			target := targetFactory()
			dec := json.NewDecoder(r.Body)
			dec.DisallowUnknownFields() // be strict by default
			if err := dec.Decode(target); err != nil {
				http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
				return
			}

			// run validation
			if err := validate.Struct(target); err != nil {
				// build structured errors
				var out []map[string]string
				if verrs, ok := err.(validator.ValidationErrors); ok {
					for _, fe := range verrs {
						out = append(out, map[string]string{
							"field":  fe.Field(),
							"tag":    fe.Tag(),
							"param":  fe.Param(),
							"reason": fe.Error(),
						})
					}
				} else {
					out = append(out, map[string]string{
						"field":  "",
						"tag":    "",
						"param":  "",
						"reason": err.Error(),
					})
				}
				resp := map[string]interface{}{
					"error":   "validation_failed",
					"details": out,
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				_ = json.NewEncoder(w).Encode(resp)
				return
			}

			// store validated value in context for handler to use
			ctx := context.WithValue(r.Context(), validatedBodyKey, target)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetValidatedBody copies the validated value found in ctx into out (out must be a pointer).
// Returns true if a value was found and assigned, false otherwise.
//
// Example:
//
//	var req CreateUserRequest
//	if !GetValidatedBody(r.Context(), &req) { ... }
func GetValidatedBody(ctx context.Context, out interface{}) bool {
	if ctx == nil || out == nil {
		return false
	}
	v := ctx.Value(validatedBodyKey)
	if v == nil {
		return false
	}
	// try marshal+unmarshal to copy into typed out (safe and simple)
	b, err := json.Marshal(v)
	if err != nil {
		return false
	}
	if err := json.Unmarshal(b, out); err != nil {
		return false
	}
	return true
}
