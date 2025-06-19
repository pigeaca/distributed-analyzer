package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/distributedmarketplace/internal/scheduler/service"
	pb "github.com/distributedmarketplace/pkg/proto/kafka"
	"github.com/segmentio/kafka-go"
	"log"
)

// SchedulerMessageHandler is a Kafka consumer for scheduler events
type SchedulerMessageHandler struct {
	schedulerService service.SchedulerService
}

func NewSchedulerHandler(schedulerService service.SchedulerService) *SchedulerMessageHandler {
	consumer := &SchedulerMessageHandler{
		schedulerService: schedulerService,
	}
	return consumer
}

// HandleMessage handleMessage handles a message from Kafka
func (c *SchedulerMessageHandler) HandleMessage(ctx context.Context, topic string, message kafka.Message) error {
	switch topic {
	case "task-created":
		return c.handleTaskCreated(ctx, message)
	default:
		return fmt.Errorf("unknown topic: %s", topic)
	}
}

// handleTaskCreated handles a TaskCreatedEvent
func (c *SchedulerMessageHandler) handleTaskCreated(ctx context.Context, message kafka.Message) error {
	var event pb.TaskCreatedEvent
	if err := json.Unmarshal(message.Value, &event); err != nil {
		return fmt.Errorf("failed to unmarshal TaskCreatedEvent: %w", err)
	}

	// Schedule the task
	if err := c.schedulerService.ScheduleTask(ctx, event.TaskId); err != nil {
		return fmt.Errorf("failed to schedule task: %w", err)
	}

	log.Printf("Task %s scheduled successfully", event.TaskId)
	return nil
}
