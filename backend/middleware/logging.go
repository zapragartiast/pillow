package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// LoggingMiddlewareMux logs incoming requests
func LoggingMiddlewareMux(logger *logrus.Logger, isEnabled bool) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !isEnabled || logger == nil {
				next.ServeHTTP(w, r)
				return
			}

			start := time.Now()

			// Capture request body if needed (but avoid sensitive data)
			var bodyBytes []byte
			if r.Body != nil {
				bodyBytes, _ = io.ReadAll(r.Body)
				r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			}

			// Create a response writer wrapper to capture status and size
			rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

			next.ServeHTTP(rw, r)

			duration := time.Since(start)

			// Extract important headers, excluding sensitive ones
			headers := make(map[string]string)
			for k, v := range r.Header {
				if !isSensitiveHeader(k) {
					headers[k] = strings.Join(v, ", ")
				}
			}

			// Prepare payload for logging (sanitize sensitive data)
			var payload interface{}
			if len(bodyBytes) > 0 {
				// Try to parse as JSON
				var parsed interface{}
				if err := json.Unmarshal(bodyBytes, &parsed); err == nil {
					payload = sanitizePayload(parsed)
				} else {
					// If not JSON, check if it's sensitive string data
					bodyStr := string(bodyBytes)
					if !containsSensitiveData(bodyStr) {
						payload = bodyStr
					} else {
						payload = "***MASKED***"
					}
				}
			}

			// Log the request
			logger.WithFields(logrus.Fields{
				"level":         "info",
				"type":          "request",
				"method":        r.Method,
				"url":           r.URL.String(),
				"headers":       headers,
				"payload":       payload,
				"request_size":  len(bodyBytes),
				"response_size": rw.size,
				"status_code":   rw.statusCode,
				"duration_ms":   duration.Milliseconds(),
				"user_agent":    r.Header.Get("User-Agent"),
				"remote_addr":   r.RemoteAddr,
			}).Info("API Request")
		})
	}
}

// responseWriter wraps http.ResponseWriter to capture status and size
type responseWriter struct {
	http.ResponseWriter
	statusCode int
	size       int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(data []byte) (int, error) {
	size, err := rw.ResponseWriter.Write(data)
	rw.size += size
	return size, err
}

// isSensitiveHeader checks if a header contains sensitive information
func isSensitiveHeader(header string) bool {
	sensitive := []string{
		"authorization",
		"cookie",
		"x-api-key",
		"x-auth-token",
	}
	headerLower := strings.ToLower(header)
	for _, s := range sensitive {
		if strings.Contains(headerLower, s) {
			return true
		}
	}
	return false
}

// containsSensitiveData checks if string might contain sensitive data
func containsSensitiveData(s string) bool {
	sensitive := []string{"password", "token", "secret", "key", "auth"}
	for _, word := range sensitive {
		if len(s) > len(word) && strings.Contains(strings.ToLower(s), word) {
			return true
		}
	}
	return false
}

// sanitizePayload recursively sanitizes sensitive data in JSON payloads
func sanitizePayload(data interface{}) interface{} {
	switch v := data.(type) {
	case map[string]interface{}:
		sanitized := make(map[string]interface{})
		for key, value := range v {
			if strings.Contains(strings.ToLower(key), "password") ||
				strings.Contains(strings.ToLower(key), "token") ||
				strings.Contains(strings.ToLower(key), "secret") ||
				strings.Contains(strings.ToLower(key), "key") {
				sanitized[key] = "***MASKED***"
			} else {
				sanitized[key] = sanitizePayload(value)
			}
		}
		return sanitized
	case []interface{}:
		sanitized := make([]interface{}, len(v))
		for i, item := range v {
			sanitized[i] = sanitizePayload(item)
		}
		return sanitized
	default:
		return v
	}
}
