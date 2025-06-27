package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/distributedmarketplace/internal/billing/service"
	pb "github.com/distributedmarketplace/pkg/proto/kafka"
	"github.com/segmentio/kafka-go"
	"log"
)

// BillingHandler is a Kafka consumer for billing events
type BillingHandler struct {
	billingService service.BillingService
}

// NewBillingHandler creates a new BillingHandler
func NewBillingHandler(billingService service.BillingService) *BillingHandler {
	return &BillingHandler{
		billingService: billingService,
	}
}

func (c *BillingHandler) HandleMessage(ctx context.Context, topic string, message kafka.Message) error {
	switch topic {
	case "task-completed":
		return c.handleTaskCompleted(ctx, message)
	default:
		return fmt.Errorf("unknown topic: %s", topic)
	}
}

func (c *BillingHandler) handleTaskCompleted(ctx context.Context, message kafka.Message) error {
	var event pb.TaskCompletedEvent
	if err := json.Unmarshal(message.Value, &event); err != nil {
		return fmt.Errorf("failed to unmarshal TaskCompletedEvent: %w", err)
	}
	record, err := c.billingService.ChargeTask(ctx, event.TaskId)
	if err != nil {
		return fmt.Errorf("failed to charge for task: %w", err)
	}

	log.Printf("Task %s charged successfully: %v", event.TaskId, record)
	return nil
}
