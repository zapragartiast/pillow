package audit

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

// Global queue instance (initialized by Start)
var Q *Queue

// StartAuditQueue initializes the global queue with given bufferSize
func StartAuditQueue(db *sql.DB, bufferSize int) {
	if Q == nil {
		Q = NewQueue(db, bufferSize)
	}
}

// StopAuditQueue shuts down the global queue gracefully
func StopAuditQueue() {
	if Q != nil {
		Q.Shutdown()
		Q = nil
	}
}

// EnqueueEvent tries to enqueue and falls back to false if full
func EnqueueEvent(action string, userID *string, details interface{}) bool {
	if Q == nil {
		return false
	}

	var uidPtr *time.Time // placeholder not used, keep signature simple
	_ = uidPtr

	// convert userID string to uuid pointer is handled by caller when needed
	ev := AuditEvent{
		ID:        uuid.New(),
		UserID:    nil,
		Action:    action,
		Details:   details,
		Timestamp: time.Now(),
	}
	return Q.Enqueue(ev)
}
