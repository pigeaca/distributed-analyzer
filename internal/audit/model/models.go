package model

import (
	"time"
)

// AuditLog represents a system audit entry
type AuditLog struct {
	ID         string    `json:"id"`
	UserID     string    `json:"user_id,omitempty"`
	Action     string    `json:"action"`
	Resource   string    `json:"resource"`
	ResourceID string    `json:"resource_id"`
	Timestamp  time.Time `json:"timestamp"`
}

// AuditAction represents the type of action performed
type AuditAction string

// Audit action types
const (
	ActionCreate AuditAction = "CREATE"
	ActionRead   AuditAction = "READ"
	ActionUpdate AuditAction = "UPDATE"
	ActionDelete AuditAction = "DELETE"
)
