package service

import (
	"context"
	"github.com/distributedmarketplace/internal/task/model"
)

// SchedulerService defines the interface for task scheduling operations
type SchedulerService interface {
	// ScheduleTask assigns a task to appropriate workers
	ScheduleTask(ctx context.Context, taskID string) error

	// DivideTask splits a task into subtasks if needed
	DivideTask(ctx context.Context, taskID string) ([]*model.SubTask, error)

	// AssignTask assigns a task or subtask to a specific worker
	AssignTask(ctx context.Context, taskID string, workerID string) error
}
