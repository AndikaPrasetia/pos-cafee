package models

import "time"

// AuditLog represents an audit trail entry for sensitive operations
type AuditLog struct {
	ID          string    `json:"id" db:"id"`
	UserID      *string   `json:"user_id,omitempty" db:"user_id"`
	Username    *string   `json:"username,omitempty" db:"username"`
	Action      string    `json:"action" db:"action"`
	Resource    string    `json:"resource" db:"resource"`
	ResourceID  *string   `json:"resource_id,omitempty" db:"resource_id"`
	IPAddress   string    `json:"ip_address" db:"ip_address"`
	UserAgent   string    `json:"user_agent" db:"user_agent"`
	Details     string    `json:"details" db:"details"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

// AuditTrailFilter represents filter options for audit logs
type AuditTrailFilter struct {
	UserID       *string    `json:"user_id,omitempty"`
	Action       *string    `json:"action,omitempty"`
	Resource     *string    `json:"resource,omitempty"`
	StartDate    *time.Time `json:"start_date,omitempty"`
	EndDate      *time.Time `json:"end_date,omitempty"`
	Limit        int        `json:"limit"`
	Offset       int        `json:"offset"`
}