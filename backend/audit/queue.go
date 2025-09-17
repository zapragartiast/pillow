package audit

import (
	"database/sql"
	"encoding/json"
	"log"
	"time"

	"github.com/google/uuid"
)

// AuditEvent is the structure enqueued for async processing
type AuditEvent struct {
	ID        uuid.UUID   `json:"id"`
	UserID    *uuid.UUID  `json:"user_id,omitempty"`
	Action    string      `json:"action"`
	Details   interface{} `json:"details,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
}

// Queue is a simple in-memory, bounded audit queue with a background worker
type Queue struct {
	ch   chan AuditEvent
	db   *sql.DB
	done chan struct{}
}

// NewQueue creates and starts a queue with a single worker
// bufferSize controls the channel buffer (e.g., 100)
func NewQueue(db *sql.DB, bufferSize int) *Queue {
	q := &Queue{
		ch:   make(chan AuditEvent, bufferSize),
		db:   db,
		done: make(chan struct{}),
	}

	go q.worker()
	return q
}

// Enqueue attempts to push an event onto the queue (non-blocking).
// Returns true if accepted, false if queue is full.
func (q *Queue) Enqueue(ev AuditEvent) bool {
	select {
	case q.ch <- ev:
		return true
	default:
		// Queue full, drop or fallback to synchronous insert
		return false
	}
}

// Shutdown gracefully stops the worker (drains channel)
func (q *Queue) Shutdown() {
	close(q.ch)
	<-q.done
}

// worker consumes events and writes them to the audit_log table
func (q *Queue) worker() {
	defer close(q.done)
	for ev := range q.ch {
		detailsBytes, err := json.Marshal(ev.Details)
		if err != nil {
			// fallback: stringified details
			detailsBytes = []byte(`"audit:marshal_error"`)
		}

		// best-effort insert; log error but continue
		_, err = q.db.Exec(
			`INSERT INTO "audit_log" (id, user_id, action, details, timestamp) VALUES ($1, $2, $3, $4, $5)`,
			ev.ID, ev.UserID, ev.Action, string(detailsBytes), ev.Timestamp,
		)
		if err != nil {
			log.Printf("audit: failed to insert audit log: %v (action=%s)\n", err, ev.Action)
		}
	}
}
