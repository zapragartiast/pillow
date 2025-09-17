package middleware

import (
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"pillow/audit"
	"strings"
	"time"

	"github.com/google/uuid"
)

// statusCapturingResponseWriter is used by AuditMiddlewareMux to capture the HTTP
// status code written by handlers so the middleware can decide to persist audit
// rows only for successful (2xx) responses.
type statusCapturingResponseWriter struct {
	http.ResponseWriter
	status int
}

func (wc *statusCapturingResponseWriter) WriteHeader(code int) {
	wc.status = code
	wc.ResponseWriter.WriteHeader(code)
}

func (wc *statusCapturingResponseWriter) Write(b []byte) (int, error) {
	// If WriteHeader was not called explicitly, default to 200 OK.
	if wc.status == 0 {
		wc.status = http.StatusOK
	}
	return wc.ResponseWriter.Write(b)
}

// AuditEntry represents the payload stored in audit_log.details
type AuditEntry struct {
	Method  string      `json:"method"`
	Path    string      `json:"path"`
	ActorID *uuid.UUID  `json:"actor_id,omitempty"`
	Body    interface{} `json:"body,omitempty"`
	Note    string      `json:"note,omitempty"`
}

// WithAuditAction returns a new context containing audit action name and structured details.
// Handlers should call this before returning if they want the middleware to write a specific
// action name (e.g. "USER_UPDATED") and structured details JSON into the audit_log.
func WithAuditAction(ctx context.Context, actionName string, details interface{}) context.Context {
	if actionName != "" {
		ctx = context.WithValue(ctx, contextKey("audit_action_name"), actionName)
	}
	if details != nil {
		ctx = context.WithValue(ctx, contextKey("audit_action_details"), details)
	}
	return ctx
}

// AuditMiddlewareMux returns a mux-compatible middleware that records requests.
// It writes a row into "audit_log" for mutating methods (POST, PUT, DELETE).
// For safety it reads a copy of the request body (if present) but never modifies it.
func AuditMiddlewareMux(db *sql.DB) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Clone body if present and readable
			var bodyCopy interface{}
			if r.Body != nil && (r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodDelete) {
				contentType := r.Header.Get("Content-Type")
				// For file uploads (multipart/form-data), don't store the actual file content in audit log
				// to prevent database bloat. Store metadata only.
				if strings.Contains(contentType, "multipart/form-data") {
					bodyCopy = map[string]interface{}{
						"content_type": contentType,
						"upload_type":  "file",
						"note":         "File content omitted from audit log to prevent storage bloat",
					}
				} else {
					// read body
					b, err := io.ReadAll(r.Body)
					if err == nil && len(b) > 0 {
						// try to parse json, fallback to string
						var parsed interface{}
						if err := json.Unmarshal(b, &parsed); err == nil {
							bodyCopy = parsed
						} else {
							bodyCopy = string(b)
						}
						// restore r.Body so handler can read it
						r.Body = io.NopCloser(strings.NewReader(string(b)))
					}
				}
			}

			// try to get user from context if available so handlers can see it
			var actorID *uuid.UUID
			if user, ok := GetUserFromContext(r.Context()); ok && user != nil {
				actorID = &user.ID
			}
			// attach action_info to context for handlers that still want metadata
			actionInfo := map[string]interface{}{
				"method": r.Method,
				"path":   r.URL.Path,
				"actor":  actorID,
				"ip":     r.RemoteAddr,
			}
			ctx := context.WithValue(r.Context(), contextKey("audit_action_info"), actionInfo)
			r = r.WithContext(ctx)

			// call next with the request that carries the audit metadata
			// Wrap ResponseWriter to capture status code so we only write audit rows for successful responses (2xx).
			wrapper := &statusCapturingResponseWriter{ResponseWriter: w, status: 0}
			next.ServeHTTP(wrapper, r)
			// After handler returns, use wrapper.status when deciding whether to persist audit row.
			// Treat 2xx as success (200-299). If status is 0 (handler didn't write), default to 200.
			if wrapper.status == 0 {
				wrapper.status = http.StatusOK
			}
			// If response is not 2xx, skip writing the audit row.
			if wrapper.status < 200 || wrapper.status >= 300 {
				return
			}
			// Use wrapper for header checks below
			w = wrapper

			// Only log mutating methods to reduce noise
			if r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodDelete {
				// Debug: log detected actor id (helps verify actor propagated)
				if actorID != nil {
					log.Printf("audit-debug: detected actor id=%s for %s %s\n", actorID.String(), r.Method, r.URL.Path)
				} else {
					log.Printf("audit-debug: no actor detected for %s %s\n", r.Method, r.URL.Path)
				}

				// Prefer handler-supplied action name/details if present in response headers or context.
				var actionStr string
				var detailsBytes []byte

				// First check response headers set by handler (headers survive after next.ServeHTTP)
				if v := w.Header().Get("X-Audit-Action"); v != "" {
					actionStr = v
				}
				if v := w.Header().Get("X-Audit-Details"); v != "" {
					_ = json.Unmarshal([]byte(v), &detailsBytes) // validate but we'll re-use raw bytes below
					detailsBytes = []byte(v)
				}

				// Next, check context keys (backwards compatibility)
				if actionStr == "" {
					if v := r.Context().Value(contextKey("audit_action_name")); v != nil {
						if s, ok := v.(string); ok && s != "" {
							actionStr = s
						}
					}
				}
				if len(detailsBytes) == 0 {
					if v := r.Context().Value(contextKey("audit_action_details")); v != nil {
						detailsBytes, _ = json.Marshal(v)
					}
				}

				// fallback if handler didn't provide action or details
				if actionStr == "" {
					actionStr = strings.ToUpper(r.Method) + " " + r.URL.Path
				}
				if len(detailsBytes) == 0 {
					d := map[string]interface{}{
						"user_before": nil,
						"user_after":  nil,
						"action":      actionInfo,
						"body":        bodyCopy,
					}
					detailsBytes, _ = json.Marshal(d)
				}

				// Enqueue the audit event asynchronously. If the queue is not available
				// or enqueue fails, fall back to synchronous DB insert to avoid losing events.
				// Convert actorID (*uuid.UUID) to *string for audit.EnqueueEvent signature.
				enqueued := false
				if audit.Q != nil {
					ev := audit.AuditEvent{
						ID:        uuid.New(),
						UserID:    actorID,
						Action:    actionStr,
						Details:   json.RawMessage(detailsBytes),
						Timestamp: time.Now(),
					}
					enqueued = audit.Q.Enqueue(ev)
				}
				if !enqueued {
					_, _ = db.Exec(`INSERT INTO "audit_log" (id, user_id, action, details, timestamp) VALUES ($1, $2, $3, $4, $5)`,
						uuid.New(), actorID, actionStr, string(detailsBytes), time.Now())
				}

				_ = start // placeholder in case we want duration later
			}
		})
	}
}
