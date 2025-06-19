package service

import (
	"context"
	"github.com/distributedmarketplace/internal/audit/model"
)

// AuditService defines the interface for audit operations
type AuditService interface {
	// LogAction logs an action in the audit system
	LogAction(ctx context.Context, userID string, action model.AuditAction, resource string, resourceID string) (*model.AuditLog, error)

	// GetAuditLogs retrieves audit logs with optional filtering
	GetAuditLogs(ctx context.Context, userID string, resource string, fromTime string, toTime string) ([]*model.AuditLog, error)

	// GetResourceAuditLogs retrieves audit logs for a specific resource
	GetResourceAuditLogs(ctx context.Context, resource string, resourceID string) ([]*model.AuditLog, error)

	// GetUserAuditLogs retrieves audit logs for a specific user
	GetUserAuditLogs(ctx context.Context, userID string) ([]*model.AuditLog, error)
}
