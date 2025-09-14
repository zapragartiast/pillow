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
		rows, err := db.Query("SELECT id, name, description, domain, managed_by, created_at, parent_org_id FROM \"Organizations\" ORDER BY name")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var orgs []models.Organization
		for rows.Next() {
			var o models.Organization
			var managedBy sql.NullString
			var parentOrg sql.NullString
			if err := rows.Scan(&o.ID, &o.Name, &o.Description, &o.Domain, &managedBy, &o.CreatedAt, &parentOrg); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
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
			http.Error(w, "Invalid organization ID format", http.StatusBadRequest)
			return
		}

		var o models.Organization
		var managedBy sql.NullString
		var parentOrg sql.NullString
		err = db.QueryRow("SELECT id, name, description, domain, managed_by, created_at, parent_org_id FROM \"Organizations\" WHERE id = $1",
			orgID).Scan(&o.ID, &o.Name, &o.Description, &o.Domain, &managedBy, &o.CreatedAt, &parentOrg)

		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Organization not found", http.StatusNotFound)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
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
			http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
			return
		}

		if strings.TrimSpace(req.Name) == "" {
			http.Error(w, "Organization name is required", http.StatusBadRequest)
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
				http.Error(w, "Invalid managed_by UUID", http.StatusBadRequest)
				return
			}
		}

		if strings.TrimSpace(req.ParentOrgID) != "" {
			if id, err := uuid.Parse(req.ParentOrgID); err == nil {
				parentOrg = &id
			} else {
				http.Error(w, "Invalid parent_org_id UUID", http.StatusBadRequest)
				return
			}
		}

		_, err := db.Exec("INSERT INTO \"Organizations\" (id, name, description, domain, managed_by, created_at, parent_org_id) VALUES ($1, $2, $3, $4, $5, CURRENT_TIMESTAMP, $6)",
			orgID, strings.TrimSpace(req.Name), strings.TrimSpace(req.Description), strings.TrimSpace(req.Domain), managedBy, parentOrg)
		if err != nil {
			http.Error(w, "Failed to create organization: "+err.Error(), http.StatusInternalServerError)
			return
		}

		var o models.Organization
		err = db.QueryRow("SELECT id, name, description, domain, managed_by, created_at, parent_org_id FROM \"Organizations\" WHERE id = $1",
			orgID).Scan(&o.ID, &o.Name, &o.Description, &o.Domain, &managedBy, &o.CreatedAt, &parentOrg)
		if err != nil {
			http.Error(w, "Failed to retrieve created organization: "+err.Error(), http.StatusInternalServerError)
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
			http.Error(w, "Invalid organization ID format", http.StatusBadRequest)
			return
		}

		var req UpdateOrganizationRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
			return
		}

		// check exists
		var existing models.Organization
		var mby sql.NullString
		var porg sql.NullString
		err = db.QueryRow("SELECT id, name, description, domain, managed_by FROM \"Organizations\" WHERE id = $1",
			orgID).Scan(&existing.ID, &existing.Name, &existing.Description, &existing.Domain, &mby)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Organization not found", http.StatusNotFound)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
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
				http.Error(w, "Invalid managed_by UUID", http.StatusBadRequest)
				return
			}
		}
		if req.ParentOrgID != "" {
			if id, err := uuid.Parse(req.ParentOrgID); err == nil {
				setParts = append(setParts, "parent_org_id = $"+strconv.Itoa(argCnt))
				args = append(args, id)
				argCnt++
			} else {
				http.Error(w, "Invalid parent_org_id UUID", http.StatusBadRequest)
				return
			}
		}

		if len(setParts) == 0 {
			http.Error(w, "No fields to update", http.StatusBadRequest)
			return
		}

		// build query (note: use strconv to be safe, but keep simple here)
		query := "UPDATE \"Organizations\" SET " + strings.Join(setParts, ", ") + " WHERE id = $" + strconv.Itoa(argCnt)
		args = append(args, orgID)

		_, err = db.Exec(query, args...)
		if err != nil {
			http.Error(w, "Failed to update organization: "+err.Error(), http.StatusInternalServerError)
			return
		}

		var updated models.Organization
		err = db.QueryRow("SELECT id, name, description, domain, managed_by, created_at, parent_org_id FROM \"Organizations\" WHERE id = $1",
			orgID).Scan(&updated.ID, &updated.Name, &updated.Description, &updated.Domain, &mby, &updated.CreatedAt, &porg)
		if err != nil {
			http.Error(w, "Failed to retrieve updated organization: "+err.Error(), http.StatusInternalServerError)
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

		// Audit: handled by middleware (middleware will write the audit record).
	}
}

// DeleteOrganization deletes an organization (soft delete not implemented here)
func DeleteOrganization(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idStr := vars["id"]

		orgID, err := uuid.Parse(idStr)
		if err != nil {
			http.Error(w, "Invalid organization ID format", http.StatusBadRequest)
			return
		}

		// prevent deletion if child organizations or user memberships exist
		var childCount int
		err = db.QueryRow("SELECT COUNT(*) FROM \"Organizations\" WHERE parent_org_id = $1", orgID).Scan(&childCount)
		if err != nil {
			http.Error(w, "Error checking children: "+err.Error(), http.StatusInternalServerError)
			return
		}
		if childCount > 0 {
			http.Error(w, "Cannot delete organization that has child organizations", http.StatusConflict)
			return
		}

		var memberCount int
		err = db.QueryRow("SELECT COUNT(*) FROM \"User_Organizations\" WHERE org_id = $1", orgID).Scan(&memberCount)
		if err != nil {
			http.Error(w, "Error checking memberships: "+err.Error(), http.StatusInternalServerError)
			return
		}
		if memberCount > 0 {
			http.Error(w, "Cannot delete organization that has user memberships", http.StatusConflict)
			return
		}

		_, err = db.Exec("DELETE FROM \"Organizations\" WHERE id = $1", orgID)
		if err != nil {
			http.Error(w, "Failed to delete organization: "+err.Error(), http.StatusInternalServerError)
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
