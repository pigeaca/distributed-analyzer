package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/distributedmarketplace/internal/task/model"
	"github.com/distributedmarketplace/internal/task/service"
	kf "github.com/distributedmarketplace/pkg/kafka"
	pb "github.com/distributedmarketplace/pkg/proto/kafka"
	taskpb "github.com/distributedmarketplace/pkg/proto/task"
	"github.com/segmentio/kafka-go"
	"time"
)

var _ kf.MessageHandler = (*TaskMessageHandler)(nil)

type TaskMessageHandler struct {
	taskService service.TaskService
}

func NewTaskMessageHandler(taskService service.TaskService) *TaskMessageHandler {
	return &TaskMessageHandler{taskService: taskService}
}

func (c *TaskMessageHandler) HandleMessage(ctx context.Context, topic string, message kafka.Message) error {
	switch topic {
	case "task-status-changed":
		return c.handleTaskStatusChanged(ctx, message)
	case "task-completed":
		return c.handleTaskCompleted(ctx, message)
	case "task-failed":
		return c.handleTaskFailed(ctx, message)
	default:
		return fmt.Errorf("unknown topic: %s", topic)
	}
}

// handleTaskStatusChanged handles a TaskStatusChangedEvent
func (c *TaskMessageHandler) handleTaskStatusChanged(ctx context.Context, message kafka.Message) error {
	var event pb.TaskStatusChangedEvent
	if err := json.Unmarshal(message.Value, &event); err != nil {
		return fmt.Errorf("failed to unmarshal TaskStatusChangedEvent: %w", err)
	}

	task, err := c.taskService.GetTask(ctx, event.TaskId)
	if err != nil {
		return fmt.Errorf("failed to get task: %w", err)
	}

	task.Status = convertPbStatusToModelStatus(event.NewStatus)
	task.UpdatedAt = time.Now()

	_, err = c.taskService.UpdateTask(ctx, task)
	if err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}

	return nil
}

// handleTaskCompleted handles a TaskCompletedEvent
func (c *TaskMessageHandler) handleTaskCompleted(ctx context.Context, message kafka.Message) error {
	var event pb.TaskCompletedEvent
	if err := json.Unmarshal(message.Value, &event); err != nil {
		return fmt.Errorf("failed to unmarshal TaskCompletedEvent: %w", err)
	}

	task, err := c.taskService.GetTask(ctx, event.TaskId)
	if err != nil {
		return fmt.Errorf("failed to get task: %w", err)
	}

	task.Status = model.StatusCompleted
	task.Output = event.Result
	task.UpdatedAt = time.Now()
	task.CompletedAt = time.Now()

	_, err = c.taskService.UpdateTask(ctx, task)
	if err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}

	return nil
}

// handleTaskFailed handles a TaskFailedEvent
func (c *TaskMessageHandler) handleTaskFailed(ctx context.Context, message kafka.Message) error {
	var event pb.TaskFailedEvent
	if err := json.Unmarshal(message.Value, &event); err != nil {
		return fmt.Errorf("failed to unmarshal TaskFailedEvent: %w", err)
	}

	task, err := c.taskService.GetTask(ctx, event.TaskId)
	if err != nil {
		return fmt.Errorf("failed to get task: %w", err)
	}

	task.Status = model.StatusFailed
	task.UpdatedAt = time.Now()

	_, err = c.taskService.UpdateTask(ctx, task)
	if err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}

	return nil
}

// Helper function to convert taskpb.Status to model.Status
func convertPbStatusToModelStatus(status taskpb.Status) model.Status {
	switch status {
	case taskpb.Status_STATUS_PENDING:
		return model.StatusPending
	case taskpb.Status_STATUS_SCHEDULED:
		return model.StatusScheduled
	case taskpb.Status_STATUS_RUNNING:
		return model.StatusRunning
	case taskpb.Status_STATUS_COMPLETED:
		return model.StatusCompleted
	case taskpb.Status_STATUS_FAILED:
		return model.StatusFailed
	default:
		return model.StatusPending
	}
}
