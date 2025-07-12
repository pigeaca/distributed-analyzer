package http

import (
	"context"
	"distributed-analyzer/services/task-service/internal/client"
	"distributed-analyzer/services/task-service/internal/model"
	taskService "distributed-analyzer/services/task-service/internal/service"
)

// TaskServiceHttpClient is a client implementation of the TaskService interface
type TaskServiceHttpClient struct {
	client *client.TaskClient
}

// Ensure TaskServiceHttpClient implements taskService.TaskService
var _ taskService.TaskService = (*TaskServiceHttpClient)(nil)

// NewTaskServiceClient creates a new TaskServiceHttpClient
func NewTaskServiceClient(baseURL string) *TaskServiceHttpClient {
	return &TaskServiceHttpClient{
		client: client.NewTaskClient(baseURL),
	}
}

// CreateTask creates a new task in the system
func (t *TaskServiceHttpClient) CreateTask(ctx context.Context, task *model.Task) (*model.Task, error) {
	return t.client.CreateTask(ctx, task)
}

// GetTask retrieves a task by its ID
func (t *TaskServiceHttpClient) GetTask(ctx context.Context, id string) (*model.Task, error) {
	return t.client.GetTask(ctx, id)
}

// UpdateTask updates an existing task
func (t *TaskServiceHttpClient) UpdateTask(ctx context.Context, task *model.Task) (*model.Task, error) {
	return t.client.UpdateTask(ctx, task)
}

// DeleteTask removes a task from the system
func (t *TaskServiceHttpClient) DeleteTask(ctx context.Context, id string) error {
	return t.client.DeleteTask(ctx, id)
}

// ListTasks retrieves all tasks with optional filtering
func (t *TaskServiceHttpClient) ListTasks(ctx context.Context) ([]*model.Task, error) {
	return t.client.ListTasks(ctx)
}
