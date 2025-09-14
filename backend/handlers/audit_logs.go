package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"pillow/models"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// GetAuditLogs retrieves paginated audit logs with full payload details
func GetAuditLogs(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse query parameters for pagination
		pageStr := r.URL.Query().Get("page")
		limitStr := r.URL.Query().Get("limit")

		page := 1
		limit := 50 // default limit

		if pageStr != "" {
			if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
				page = p
			}
		}

		if limitStr != "" {
			if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
				limit = l
			}
		}

		offset := (page - 1) * limit

		// Query audit logs with pagination
		query := `SELECT id, user_id, action, details, timestamp FROM "Audit_Log" ORDER BY timestamp DESC LIMIT $1 OFFSET $2`
		rows, err := db.Query(query, limit, offset)
		if err != nil {
			logrus.WithError(err).Error("Failed to query audit logs")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var auditLogs []models.AuditLog
		for rows.Next() {
			var log models.AuditLog
			err := rows.Scan(&log.ID, &log.UserID, &log.Action, &log.Details, &log.Timestamp)
			if err != nil {
				logrus.WithError(err).Error("Failed to scan audit log")
				continue
			}
			auditLogs = append(auditLogs, log)
		}

		if err = rows.Err(); err != nil {
			logrus.WithError(err).Error("Error iterating audit log rows")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// Get total count for pagination info
		var totalCount int
		countQuery := `SELECT COUNT(*) FROM "Audit_Log"`
		err = db.QueryRow(countQuery).Scan(&totalCount)
		if err != nil {
			logrus.WithError(err).Error("Failed to get audit log count")
			totalCount = 0
		}

		response := map[string]interface{}{
			"audit_logs": auditLogs,
			"pagination": map[string]interface{}{
				"page":       page,
				"limit":      limit,
				"total":      totalCount,
				"totalPages": (totalCount + limit - 1) / limit,
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

// GetAuditLog retrieves a specific audit log entry by ID
func GetAuditLog(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		if id == "" {
			http.Error(w, "Audit log ID is required", http.StatusBadRequest)
			return
		}

		var log models.AuditLog
		query := `SELECT id, user_id, action, details, timestamp FROM "Audit_Log" WHERE id = $1`
		err := db.QueryRow(query, id).Scan(&log.ID, &log.UserID, &log.Action, &log.Details, &log.Timestamp)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Audit log not found", http.StatusNotFound)
				return
			}
			logrus.WithError(err).Error("Failed to query audit log")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(log)
	}
}
