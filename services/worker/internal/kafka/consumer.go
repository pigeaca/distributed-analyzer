package kafka

import (
	"context"
	pb "distributed-analyzer/libs/proto/kafka"
	"distributed-analyzer/services/worker/internal/service"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
)

// WorkerHandler is a Kafka consumer for worker events
type WorkerHandler struct {
	workerService service.WorkerNodeService
}

// NewWorkerHandler creates a new WorkerHandler
func NewWorkerHandler(workerService service.WorkerNodeService) *WorkerHandler {
	return &WorkerHandler{
		workerService: workerService,
	}
}

func (c *WorkerHandler) HandleMessage(ctx context.Context, topic string, message kafka.Message) error {
	switch topic {
	case "task-assigned":
		return c.handleTaskAssigned(ctx, message)
	default:
		return fmt.Errorf("unknown topic: %s", topic)
	}
}

// handleTaskAssigned handles a TaskAssignedEvent
func (c *WorkerHandler) handleTaskAssigned(ctx context.Context, message kafka.Message) error {
	var event pb.TaskAssignedEvent
	if err := json.Unmarshal(message.Value, &event); err != nil {
		return fmt.Errorf("failed to unmarshal TaskAssignedEvent: %w", err)
	}

	// Execute the task
	if err := c.workerService.ExecuteTask(ctx, event.TaskId); err != nil {
		return fmt.Errorf("failed to execute task: %w", err)
	}

	log.Printf("Task %s executed successfully", event.TaskId)
	return nil
}
