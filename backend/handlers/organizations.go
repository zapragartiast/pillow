package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"pillow/middleware"
	"pillow/models"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// CreateOrganizationRequest represents payload to create an organization
type CreateOrganizationRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Domain      string `json:"domain,omitempty"`
	ManagedBy   string `json:"managed_by,omitempty"` // user id
	ParentOrgID string `json:"parent_org_id,omitempty"`
}

// UpdateOrganizationRequest represents payload to update an organization
type UpdateOrganizationRequest struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Domain      string `json:"domain,omitempty"`
	ManagedBy   string `json:"managed_by,omitempty"`
	ParentOrgID string `json:"parent_org_id,omitempty"`
}

// GetOrganizations returns all organizations
func GetOrganizations(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT id, name, description, domain, managed_by, created_at, updated_at, parent_org_id FROM \"organizations\" ORDER BY name")
		if err != nil {
			writeErrorResponse(w, err.Error(), http.StatusInternalServerError, r)
			return
		}
		defer rows.Close()

		var orgs []models.Organization
		for rows.Next() {
			var o models.Organization
			var managedBy sql.NullString
			var parentOrg sql.NullString
			if err := rows.Scan(&o.ID, &o.Name, &o.Description, &o.Domain, &managedBy, &o.CreatedAt, &o.UpdatedAt, &parentOrg); err != nil {
				writeErrorResponse(w, err.Error(), http.StatusInternalServerError, r)
				return
			}
			if managedBy.Valid {
				id, _ := uuid.Parse(managedBy.String)
				o.ManagedBy = &id
			}
			if parentOrg.Valid {
				id, _ := uuid.Parse(parentOrg.String)
				o.ParentOrgID = &id
			}
			orgs = append(orgs, o)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(orgs)
	}
}

// GetOrganization retrieves a single organization by ID
func GetOrganization(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idStr := vars["id"]

		orgID, err := uuid.Parse(idStr)
		if err != nil {
			writeErrorResponse(w, "Invalid organization ID format", http.StatusBadRequest, r)
			return
		}

		var o models.Organization
		var managedBy sql.NullString
		var parentOrg sql.NullString
		err = db.QueryRow("SELECT id, name, description, domain, managed_by, created_at, updated_at, parent_org_id FROM \"organizations\" WHERE id = $1",
			orgID).Scan(&o.ID, &o.Name, &o.Description, &o.Domain, &managedBy, &o.CreatedAt, &o.UpdatedAt, &parentOrg)

		if err != nil {
			if err == sql.ErrNoRows {
				writeErrorResponse(w, "Organization not found", http.StatusNotFound, r)
				return
			}
			writeErrorResponse(w, err.Error(), http.StatusInternalServerError, r)
			return
		}
		if managedBy.Valid {
			id, _ := uuid.Parse(managedBy.String)
			o.ManagedBy = &id
		}
		if parentOrg.Valid {
			id, _ := uuid.Parse(parentOrg.String)
			o.ParentOrgID = &id
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(o)
	}
}

// CreateOrganization creates a new organization
func CreateOrganization(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateOrganizationRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeErrorResponse(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest, r)
			return
		}

		if strings.TrimSpace(req.Name) == "" {
			writeErrorResponse(w, "Organization name is required", http.StatusBadRequest, r)
			return
		}

		// generate id
		orgID := uuid.New()
		var managedBy *uuid.UUID
		var parentOrg *uuid.UUID

		if strings.TrimSpace(req.ManagedBy) != "" {
			if id, err := uuid.Parse(req.ManagedBy); err == nil {
				managedBy = &id
			} else {
				writeErrorResponse(w, "Invalid managed_by UUID", http.StatusBadRequest, r)
				return
			}
		}

		if strings.TrimSpace(req.ParentOrgID) != "" {
			if id, err := uuid.Parse(req.ParentOrgID); err == nil {
				parentOrg = &id
			} else {
				writeErrorResponse(w, "Invalid parent_org_id UUID", http.StatusBadRequest, r)
				return
			}
		}

		_, err := db.Exec("INSERT INTO \"organizations\" (id, name, description, domain, managed_by, created_at, updated_at, parent_org_id) VALUES ($1, $2, $3, $4, $5, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, $6)",
			orgID, strings.TrimSpace(req.Name), strings.TrimSpace(req.Description), strings.TrimSpace(req.Domain), managedBy, parentOrg)
		if err != nil {
			writeErrorResponse(w, "Failed to create organization: "+err.Error(), http.StatusInternalServerError, r)
			return
		}

		var o models.Organization
		err = db.QueryRow("SELECT id, name, description, domain, managed_by, created_at, updated_at, parent_org_id FROM \"organizations\" WHERE id = $1",
			orgID).Scan(&o.ID, &o.Name, &o.Description, &o.Domain, &managedBy, &o.CreatedAt, &o.UpdatedAt, &parentOrg)
		if err != nil {
			writeErrorResponse(w, "Failed to retrieve created organization: "+err.Error(), http.StatusInternalServerError, r)
			return
		}
		if managedBy != nil {
			o.ManagedBy = managedBy
		}
		if parentOrg != nil {
			o.ParentOrgID = parentOrg
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		// Prepare audit details and expose via response headers so middleware records structured audit.
		detailsMap := map[string]interface{}{
			"user_before": nil,
			"user_after":  o,
			"action": map[string]interface{}{
				"method":     r.Method,
				"path":       r.URL.Path,
				"actor_id":   nil,
				"ip_address": r.RemoteAddr,
			},
		}
		if user, ok := middleware.GetUserFromContext(r.Context()); ok && user != nil {
			detailsMap["action"].(map[string]interface{})["actor_id"] = user.ID.String()
		}
		dBytes, _ := json.Marshal(detailsMap)
		w.Header().Set("X-Audit-Action", "ORGANIZATION_CREATED")
		w.Header().Set("X-Audit-Details", string(dBytes))

		json.NewEncoder(w).Encode(map[string]interface{}{
			"message":      "Organization created successfully",
			"organization": o,
		})
	}
}

// UpdateOrganization updates an existing organization
func UpdateOrganization(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idStr := vars["id"]

		orgID, err := uuid.Parse(idStr)
		if err != nil {
			writeErrorResponse(w, "Invalid organization ID format", http.StatusBadRequest, r)
			return
		}

		var req UpdateOrganizationRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeErrorResponse(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest, r)
			return
		}

		// check exists
		var existing models.Organization
		var mby sql.NullString
		var porg sql.NullString
		err = db.QueryRow("SELECT id, name, description, domain, managed_by, created_at, updated_at FROM \"organizations\" WHERE id = $1",
			orgID).Scan(&existing.ID, &existing.Name, &existing.Description, &existing.Domain, &mby, &existing.CreatedAt, &existing.UpdatedAt)
		if err != nil {
			if err == sql.ErrNoRows {
				writeErrorResponse(w, "Organization not found", http.StatusNotFound, r)
				return
			}
			writeErrorResponse(w, err.Error(), http.StatusInternalServerError, r)
			return
		}

		setParts := []string{}
		args := []interface{}{}
		argCnt := 1

		if strings.TrimSpace(req.Name) != "" {
			setParts = append(setParts, "name = $"+strconv.Itoa(argCnt))
			args = append(args, strings.TrimSpace(req.Name))
			argCnt++
		}
		if req.Description != "" {
			setParts = append(setParts, "description = $"+strconv.Itoa(argCnt))
			args = append(args, strings.TrimSpace(req.Description))
			argCnt++
		}
		if req.Domain != "" {
			setParts = append(setParts, "domain = $"+strconv.Itoa(argCnt))
			args = append(args, strings.TrimSpace(req.Domain))
			argCnt++
		}
		if req.ManagedBy != "" {
			if id, err := uuid.Parse(req.ManagedBy); err == nil {
				setParts = append(setParts, "managed_by = $"+strconv.Itoa(argCnt))
				args = append(args, id)
				argCnt++
			} else {
				writeErrorResponse(w, "Invalid managed_by UUID", http.StatusBadRequest, r)
				return
			}
		}
		if req.ParentOrgID != "" {
			if id, err := uuid.Parse(req.ParentOrgID); err == nil {
				setParts = append(setParts, "parent_org_id = $"+strconv.Itoa(argCnt))
				args = append(args, id)
				argCnt++
			} else {
				writeErrorResponse(w, "Invalid parent_org_id UUID", http.StatusBadRequest, r)
				return
			}
		}

		if len(setParts) == 0 {
			writeErrorResponse(w, "No fields to update", http.StatusBadRequest, r)
			return
		}

		setParts = append(setParts, "updated_at = CURRENT_TIMESTAMP")
		// build query (note: use strconv to be safe, but keep simple here)
		query := "UPDATE \"organizations\" SET " + strings.Join(setParts, ", ") + " WHERE id = $" + strconv.Itoa(argCnt)
		args = append(args, orgID)

		_, err = db.Exec(query, args...)
		if err != nil {
			writeErrorResponse(w, "Failed to update organization: "+err.Error(), http.StatusInternalServerError, r)
			return
		}

		var updated models.Organization
		err = db.QueryRow("SELECT id, name, description, domain, managed_by, created_at, updated_at, parent_org_id FROM \"organizations\" WHERE id = $1",
			orgID).Scan(&updated.ID, &updated.Name, &updated.Description, &updated.Domain, &mby, &updated.CreatedAt, &updated.UpdatedAt, &porg)
		if err != nil {
			writeErrorResponse(w, "Failed to retrieve updated organization: "+err.Error(), http.StatusInternalServerError, r)
			return
		}
		if mby.Valid {
			id, _ := uuid.Parse(mby.String)
			updated.ManagedBy = &id
		}
		if porg.Valid {
			id, _ := uuid.Parse(porg.String)
			updated.ParentOrgID = &id
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message":      "Organization updated successfully",
			"organization": updated,
		})

		// Prepare audit details and expose via response headers so middleware records structured audit.
		detailsMap := map[string]interface{}{
			"user_before": existing,
			"user_after":  updated,
			"action": map[string]interface{}{
				"method":     r.Method,
				"path":       r.URL.Path,
				"actor_id":   nil,
				"ip_address": r.RemoteAddr,
			},
		}
		if user, ok := middleware.GetUserFromContext(r.Context()); ok && user != nil {
			detailsMap["action"].(map[string]interface{})["actor_id"] = user.ID.String()
		}
		dBytes, _ := json.Marshal(detailsMap)
		w.Header().Set("X-Audit-Action", "ORGANIZATION_UPDATED")
		w.Header().Set("X-Audit-Details", string(dBytes))
	}
}

// DeleteOrganization deletes an organization (soft delete not implemented here)
func DeleteOrganization(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idStr := vars["id"]

		orgID, err := uuid.Parse(idStr)
		if err != nil {
			writeErrorResponse(w, "Invalid organization ID format", http.StatusBadRequest, r)
			return
		}

		// prevent deletion if child organizations or user memberships exist
		var childCount int
		err = db.QueryRow("SELECT COUNT(*) FROM \"organizations\" WHERE parent_org_id = $1", orgID).Scan(&childCount)
		if err != nil {
			writeErrorResponse(w, "Error checking children: "+err.Error(), http.StatusInternalServerError, r)
			return
		}
		if childCount > 0 {
			writeErrorResponse(w, "Cannot delete organization that has child organizations", http.StatusConflict, r)
			return
		}

		var memberCount int
		err = db.QueryRow("SELECT COUNT(*) FROM \"user_organizations\" WHERE org_id = $1", orgID).Scan(&memberCount)
		if err != nil {
			writeErrorResponse(w, "Error checking memberships: "+err.Error(), http.StatusInternalServerError, r)
			return
		}
		if memberCount > 0 {
			writeErrorResponse(w, "Cannot delete organization that has user memberships", http.StatusConflict, r)
			return
		}

		_, err = db.Exec("DELETE FROM \"organizations\" WHERE id = $1", orgID)
		if err != nil {
			writeErrorResponse(w, "Failed to delete organization: "+err.Error(), http.StatusInternalServerError, r)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message":   "Organization deleted successfully",
			"org_id":    orgID,
			"deletedAt": time.Now(),
		})

		// Audit: handled by middleware (middleware will write the audit record).
	}
}
