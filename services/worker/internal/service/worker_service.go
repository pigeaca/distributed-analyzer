package service

import (
	"context"
	"distributed-analyzer/libs/model"
)

// WorkerManagerService defines the interface for worker management operations
type WorkerManagerService interface {
	// RegisterWorker registers a new worker in the system
	RegisterWorker(ctx context.Context, worker *model.Worker) (*model.Worker, error)

	// GetWorker retrieves a worker by its ID
	GetWorker(ctx context.Context, id string) (*model.Worker, error)

	// UpdateWorkerStatus updates the status of a worker
	UpdateWorkerStatus(ctx context.Context, id string, status string) error

	// ListWorkers retrieves all workers with optional filtering
	ListWorkers(ctx context.Context) ([]*model.Worker, error)

	// FindAvailableWorkers finds workers that can handle a specific task
	FindAvailableWorkers(ctx context.Context, capabilities []model.Capability, resources []model.Resource) ([]*model.Worker, error)
}

// WorkerNodeService defines the interface for worker node operations
type WorkerNodeService interface {
	// ExecuteTask executes a task on the worker
	ExecuteTask(ctx context.Context, taskID string) error

	// LoadModel loads a model required for task execution
	LoadModel(ctx context.Context, modelName string) error

	// ReportStatus reports the worker's current status
	ReportStatus(ctx context.Context, status string) error

	// SendResult sends the result of a completed task
	SendResult(ctx context.Context, taskID string, result map[string]string) error
}
