package handlers

import (
	"encoding/json"
	"net/http"
	"time"
)

// ErrorResponse represents a consistent error response structure
type ErrorResponse struct {
	Error     string    `json:"error"`
	Code      int       `json:"code"`
	Message   string    `json:"message"`
	Path      string    `json:"path,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

// writeErrorResponse sends a consistent JSON error response
func writeErrorResponse(w http.ResponseWriter, message string, code int, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	errorResp := ErrorResponse{
		Error:     http.StatusText(code),
		Code:      code,
		Message:   message,
		Path:      r.URL.Path,
		Timestamp: time.Now(),
	}

	json.NewEncoder(w).Encode(errorResp)
}
