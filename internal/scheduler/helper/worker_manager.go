package helper

import (
	"context"
	"fmt"
	"github.com/distributedmarketplace/internal/scheduler/grpc"
	"github.com/distributedmarketplace/internal/scheduler/kafka/producer"
	workermodel "github.com/distributedmarketplace/internal/worker/model"
	"log"
)

// WorkerManagerHelper provides helper functions for interacting with the WorkerManager service
type WorkerManagerHelper struct {
	workerClient  *grpc.WorkerManagerClient
	kafkaProducer *producer.SchedulerProducer
}

// NewWorkerManagerHelper creates a new WorkerManagerHelper
func NewWorkerManagerHelper(workerServiceAddr string, kafkaBrokers []string) (*WorkerManagerHelper, error) {
	workerClient, err := grpc.NewWorkerManagerClient(workerServiceAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to create WorkerManager client: %w", err)
	}

	kafkaProducer := producer.NewSchedulerProducer(kafkaBrokers)

	return &WorkerManagerHelper{
		workerClient:  workerClient,
		kafkaProducer: kafkaProducer,
	}, nil
}

// Close closes all connections
func (h *WorkerManagerHelper) Close() error {
	if err := h.workerClient.Close(); err != nil {
		return err
	}
	if err := h.kafkaProducer.Close(); err != nil {
		return err
	}
	return nil
}

// FindAvailableWorkers finds workers that can handle a specific task
func (h *WorkerManagerHelper) FindAvailableWorkers(ctx context.Context, capabilities []workermodel.Capability, resources []workermodel.Resource) ([]*workermodel.Worker, error) {
	workers, err := h.workerClient.FindAvailableWorkers(ctx, capabilities, resources)
	if err != nil {
		return nil, fmt.Errorf("failed to find available workers: %w", err)
	}

	log.Printf("Found %d available workers", len(workers))
	return workers, nil
}

// GetWorker retrieves a worker by its ID
func (h *WorkerManagerHelper) GetWorker(ctx context.Context, id string) (*workermodel.Worker, error) {
	worker, err := h.workerClient.GetWorker(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get worker: %w", err)
	}

	return worker, nil
}

// PublishTaskAssigned publishes a TaskAssignedEvent to Kafka
func (h *WorkerManagerHelper) PublishTaskAssigned(ctx context.Context, taskID string, workerID string) error {
	if err := h.kafkaProducer.PublishTaskAssigned(ctx, taskID, workerID); err != nil {
		return fmt.Errorf("failed to publish TaskAssignedEvent: %w", err)
	}
	log.Printf("Published TaskAssignedEvent for task %s, worker %s", taskID, workerID)
	return nil
}

// PublishTaskScheduled publishes a TaskScheduledEvent to Kafka
func (h *WorkerManagerHelper) PublishTaskScheduled(ctx context.Context, taskID string, workerIDs []string) error {
	if err := h.kafkaProducer.PublishTaskScheduled(ctx, taskID, workerIDs); err != nil {
		return fmt.Errorf("failed to publish TaskScheduledEvent: %w", err)
	}

	log.Printf("Published TaskScheduledEvent for task %s, workers %v", taskID, workerIDs)
	return nil
}
