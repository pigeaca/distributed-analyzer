package service

import (
	"context"
	"github.com/distributedmarketplace/internal/task/model"
)

// TaskService defines the interface for task-related operations
type TaskService interface {
	// CreateTask creates a new task in the system
	CreateTask(ctx context.Context, task *model.Task) (*model.Task, error)

	// GetTask retrieves a task by its ID
	GetTask(ctx context.Context, id string) (*model.Task, error)

	// UpdateTask updates an existing task
	UpdateTask(ctx context.Context, task *model.Task) (*model.Task, error)

	// DeleteTask removes a task from the system
	DeleteTask(ctx context.Context, id string) error

	// ListTasks retrieves all tasks with optional filtering
	ListTasks(ctx context.Context) ([]*model.Task, error)
}
