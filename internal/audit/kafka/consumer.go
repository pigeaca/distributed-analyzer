package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/distributedmarketplace/internal/audit/model"
	"github.com/distributedmarketplace/internal/audit/service"
	pb "github.com/distributedmarketplace/pkg/proto/kafka"
	"github.com/segmentio/kafka-go"
)

// AuditHandler is a Kafka consumer for audit events
type AuditHandler struct {
	auditService service.AuditService
}

// NewAuditHandler creates a new AuditHandler
func NewAuditHandler(auditService service.AuditService) *AuditHandler {
	return &AuditHandler{
		auditService: auditService,
	}
}

// HandleMessage handles a message from Kafka
func (c *AuditHandler) HandleMessage(ctx context.Context, topic string, message kafka.Message) error {
	switch topic {
	case "audit-event":
		return c.handleAuditEvent(ctx, message)
	default:
		return fmt.Errorf("unknown topic: %s", topic)
	}
}

// handleAuditEvent handles an AuditEvent
func (c *AuditHandler) handleAuditEvent(ctx context.Context, message kafka.Message) error {
	var event pb.AuditEvent
	if err := json.Unmarshal(message.Value, &event); err != nil {
		return fmt.Errorf("failed to unmarshal AuditEvent: %w", err)
	}

	// Convert action string to model.AuditAction
	var action model.AuditAction
	switch event.Action {
	case "CREATE":
		action = model.ActionCreate
	case "READ":
		action = model.ActionRead
	case "UPDATE":
		action = model.ActionUpdate
	case "DELETE":
		action = model.ActionDelete
	default:
		action = model.ActionCreate // Default to CREATE if the action is not recognized
	}

	// Log the action
	auditLog, err := c.auditService.LogAction(ctx, event.UserId, action, event.Resource, event.ResourceId)
	if err != nil {
		return fmt.Errorf("failed to auditLog action: %w", err)
	}

	fmt.Printf("Audit auditLog created: %v\n", auditLog)
	return nil
}
