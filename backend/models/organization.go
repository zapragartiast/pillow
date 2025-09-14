package models

import (
	"time"

	"github.com/google/uuid"
)

type Organization struct {
	ID          uuid.UUID  `json:"id" db:"id"`
	Name        string     `json:"name" db:"name"`
	Description string     `json:"description,omitempty" db:"description"`
	Domain      string     `json:"domain,omitempty" db:"domain"`
	ManagedBy   *uuid.UUID `json:"managed_by,omitempty" db:"managed_by"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	ParentOrgID *uuid.UUID `json:"parent_org_id,omitempty" db:"parent_org_id"`
}

// OrganizationWithUsers represents an organization with its associated users
type OrganizationWithUsers struct {
	Organization Organization `json:"organization"`
	Users        []User       `json:"users"`
}

// UserOrganization represents the many-to-many relationship between users and organizations
type UserOrganization struct {
	ID        uuid.UUID  `json:"id" db:"id"`
	UserID    uuid.UUID  `json:"user_id" db:"user_id"`
	OrgID     uuid.UUID  `json:"org_id" db:"org_id"`
	RoleID    *uuid.UUID `json:"role_id,omitempty" db:"role_id"`
	InvitedBy *uuid.UUID `json:"invited_by,omitempty" db:"invited_by"`
	JoinedAt  time.Time  `json:"joined_at" db:"joined_at"`
}
