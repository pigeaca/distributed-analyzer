package service

import (
	"context"
	"github.com/distributedmarketplace/internal/result/model"
)

// ResultAggregatorService defines the interface for result aggregation operations
type ResultAggregatorService interface {
	// SavePartialResult saves a partial result for a task
	SavePartialResult(ctx context.Context, taskID string, subtaskID string, result map[string]string) error

	// FinalizeResult finalizes the result when all subtasks are completed
	FinalizeResult(ctx context.Context, taskID string) error

	// GetResult retrieves the result of a completed task
	GetResult(ctx context.Context, taskID string) (map[string]string, error)

	// GetTaskResult retrieves the full task result object
	GetTaskResult(ctx context.Context, taskID string) (*model.TaskResult, error)

	// GetSubTaskResults retrieves all subtask results for a task
	GetSubTaskResults(ctx context.Context, taskID string) ([]*model.SubTaskResult, error)
}
