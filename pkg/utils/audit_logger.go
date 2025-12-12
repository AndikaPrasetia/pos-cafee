package utils

import (
	"time"

	"github.com/gin-gonic/gin"
)

// AuditEvent represents an auditable event in the system
type AuditEvent struct {
	ID          string                 `json:"id"`
	Timestamp   time.Time              `json:"timestamp"`
	UserID      string                 `json:"user_id,omitempty"`
	Username    string                 `json:"username,omitempty"`
	Action      string                 `json:"action"`
	Resource    string                 `json:"resource"`
	Details     map[string]interface{} `json:"details,omitempty"`
	IPAddress   string                 `json:"ip_address,omitempty"`
	UserAgent   string                 `json:"user_agent,omitempty"`
	Status      string                 `json:"status"` // success, failed, warning
	SessionID   string                 `json:"session_id,omitempty"`
	RequestID   string                 `json:"request_id,omitempty"`
	Environment string                 `json:"environment"`
}

// AuditLogger provides audit logging functionality
type AuditLogger struct{}

// NewAuditLogger creates a new audit logger instance
func NewAuditLogger() *AuditLogger {
	return &AuditLogger{}
}

// LogUserLogin logs a user login event
func (al *AuditLogger) LogUserLogin(ctx *gin.Context, userID, username, ip, userAgent string, success bool) {
	status := "success"
	if !success {
		status = "failed"
	}

	event := AuditEvent{
		ID:        GenerateUUID(),
		Timestamp: time.Now(),
		UserID:    userID,
		Username:  username,
		Action:    "user_login",
		Resource:  "auth",
		IPAddress: ip,
		UserAgent: userAgent,
		Status:    status,
		SessionID: ctx.GetString("session_id"),
		RequestID: ctx.GetString("request_id"),
	}

	fields := map[string]interface{}{
		"event_id":    event.ID,
		"user_id":     event.UserID,
		"username":    event.Username,
		"action":      event.Action,
		"resource":    event.Resource,
		"ip_address":  event.IPAddress,
		"user_agent":  event.UserAgent,
		"status":      event.Status,
		"session_id":  event.SessionID,
		"request_id":  event.RequestID,
		"timestamp":   event.Timestamp.Format(time.RFC3339),
		"environment": event.Environment,
	}

	if success {
		LogInfo("User login event", fields)
	} else {
		LogWarn("User login failed", fields)
	}
}

// LogUserLogout logs a user logout event
func (al *AuditLogger) LogUserLogout(ctx *gin.Context, userID, username, ip string) {
	event := AuditEvent{
		ID:        GenerateUUID(),
		Timestamp: time.Now(),
		UserID:    userID,
		Username:  username,
		Action:    "user_logout",
		Resource:  "auth",
		IPAddress: ip,
		Status:    "success",
		SessionID: ctx.GetString("session_id"),
		RequestID: ctx.GetString("request_id"),
	}

	fields := map[string]interface{}{
		"event_id":    event.ID,
		"user_id":     event.UserID,
		"username":    event.Username,
		"action":      event.Action,
		"resource":    event.Resource,
		"ip_address":  event.IPAddress,
		"status":      event.Status,
		"session_id":  event.SessionID,
		"request_id":  event.RequestID,
		"timestamp":   event.Timestamp.Format(time.RFC3339),
		"environment": event.Environment,
	}

	LogInfo("User logout event", fields)
}

// LogSensitiveAction logs sensitive operations like order creation, payments, etc.
func (al *AuditLogger) LogSensitiveAction(ctx *gin.Context, userID, username, action, resource string, details map[string]interface{}, status string) {
	event := AuditEvent{
		ID:        GenerateUUID(),
		Timestamp: time.Now(),
		UserID:    userID,
		Username:  username,
		Action:    action,
		Resource:  resource,
		Details:   details,
		IPAddress: ctx.GetString("ip_address"),
		Status:    status,
		SessionID: ctx.GetString("session_id"),
		RequestID: ctx.GetString("request_id"),
	}

	fields := map[string]interface{}{
		"event_id":    event.ID,
		"user_id":     event.UserID,
		"username":    event.Username,
		"action":      event.Action,
		"resource":    event.Resource,
		"details":     event.Details,
		"ip_address":  event.IPAddress,
		"status":      event.Status,
		"session_id":  event.SessionID,
		"request_id":  event.RequestID,
		"timestamp":   event.Timestamp.Format(time.RFC3339),
		"environment": event.Environment,
	}

	switch status {
	case "success":
		LogInfo("Sensitive action performed", fields)
	case "failed":
		LogError("Sensitive action failed", fields)
	case "warning":
		LogWarn("Sensitive action warning", fields)
	default:
		LogInfo("Sensitive action performed", fields)
	}
}

// LogDataModification logs data modification events
func (al *AuditLogger) LogDataModification(ctx *gin.Context, userID, username, action, resource, recordID string, oldValues, newValues map[string]interface{}) {
	details := map[string]interface{}{
		"record_id":  recordID,
		"old_values": oldValues,
		"new_values": newValues,
	}

	event := AuditEvent{
		ID:        GenerateUUID(),
		Timestamp: time.Now(),
		UserID:    userID,
		Username:  username,
		Action:    action,
		Resource:  resource,
		Details:   details,
		IPAddress: ctx.GetString("ip_address"),
		Status:    "success",
		SessionID: ctx.GetString("session_id"),
		RequestID: ctx.GetString("request_id"),
	}

	fields := map[string]interface{}{
		"event_id":    event.ID,
		"user_id":     event.UserID,
		"username":    event.Username,
		"action":      event.Action,
		"resource":    event.Resource,
		"details":     event.Details,
		"ip_address":  event.IPAddress,
		"status":      event.Status,
		"session_id":  event.SessionID,
		"request_id":  event.RequestID,
		"timestamp":   event.Timestamp.Format(time.RFC3339),
		"environment": event.Environment,
	}

	LogInfo("Data modification event", fields)
}