package service

import (
	"context"
	"errors"
	"github.com/distributedmarketplace/internal/scheduler/grpc"
	"github.com/distributedmarketplace/internal/scheduler/kafka/producer"
	"github.com/distributedmarketplace/internal/task/model"
	"github.com/distributedmarketplace/internal/worker/model"
	"log"
	"sync"
)

// ErrTaskNotFound is returned when a task with the specified ID doesn't exist
var ErrTaskNotFound = errors.New("task not found")

// SchedulerServiceImpl implements the SchedulerService interface
type SchedulerServiceImpl struct {
	tasks             map[string]*model.Task
	taskMu            sync.RWMutex
	workerClient      *grpc.WorkerManagerClient
	kafkaProducer     *producer.SchedulerProducer
	taskServiceAddr   string
	workerServiceAddr string
}

// NewSchedulerServiceImpl creates a new instance of SchedulerServiceImpl
func NewSchedulerServiceImpl(taskServiceAddr, workerServiceAddr string, kafkaBrokers []string) (*SchedulerServiceImpl, error) {
	workerClient, err := grpc.NewWorkerManagerClient(workerServiceAddr)
	if err != nil {
		return nil, err
	}

	kafkaProducer := producer.NewSchedulerProducer(kafkaBrokers)

	return &SchedulerServiceImpl{
		tasks:             make(map[string]*model.Task),
		workerClient:      workerClient,
		kafkaProducer:     kafkaProducer,
		taskServiceAddr:   taskServiceAddr,
		workerServiceAddr: workerServiceAddr,
	}, nil
}

// Close closes all connections
func (s *SchedulerServiceImpl) Close() error {
	if err := s.workerClient.Close(); err != nil {
		return err
	}
	if err := s.kafkaProducer.Close(); err != nil {
		return err
	}
	return nil
}

// ScheduleTask assigns a task to appropriate workers
func (s *SchedulerServiceImpl) ScheduleTask(ctx context.Context, taskID string) error {
	log.Printf("Scheduling task %s", taskID)

	// In a real implementation, this would:
	// 1. Get the task details from the task service
	task, err := s.getTaskFromService(ctx, taskID)
	if err != nil {
		return err
	}

	// 2. Find available workers using the worker manager
	capabilities := []model.Capability{
		{Name: "default", Value: "1.0"},
	}
	resources := []model.Resource{
		{Type: "CPU", Value: 1},
	}

	workers, err := s.workerClient.FindAvailableWorkers(ctx, capabilities, resources)
	if err != nil {
		return err
	}

	if len(workers) == 0 {
		return errors.New("no available workers found")
	}

	// 3. Divide the task into subtasks if needed
	subtasks, err := s.DivideTask(ctx, taskID)
	if err != nil {
		return err
	}

	// 4. Assign the task or subtasks to workers
	workerIDs := make([]string, 0, len(workers))
	for i, worker := range workers {
		if i < len(subtasks) {
			// Assign subtask to worker
			if err := s.AssignTask(ctx, subtasks[i].ID, worker.ID); err != nil {
				return err
			}
		} else {
			// Assign main task to worker if no subtasks
			if len(subtasks) == 0 {
				if err := s.AssignTask(ctx, taskID, worker.ID); err != nil {
					return err
				}
			}
		}
		workerIDs = append(workerIDs, worker.ID)
	}

	// 5. Update the task status to SCHEDULED
	task.Status = model.StatusScheduled
	if err := s.updateTaskInService(ctx, task); err != nil {
		return err
	}

	// 6. Publish TaskScheduledEvent to Kafka
	if err := s.kafkaProducer.PublishTaskScheduled(ctx, taskID, workerIDs); err != nil {
		return err
	}

	return nil
}

// DivideTask splits a task into subtasks if needed
func (s *SchedulerServiceImpl) DivideTask(ctx context.Context, taskID string) ([]*model.SubTask, error) {
	log.Printf("Dividing task %s", taskID)
	// In a real implementation, this would:
	// 1. Get the task details from the task service
	// 2. Analyze the task to determine how to divide it
	// 3. Create subtasks
	// 4. Return the subtasks
	return nil, nil
}

// AssignTask assigns a task or subtask to a specific worker
func (s *SchedulerServiceImpl) AssignTask(ctx context.Context, taskID string, workerID string) error {
	log.Printf("Assigning task %s to worker %s", taskID, workerID)
	// In a real implementation, this would:
	// 1. Get the task details from the task service
	// 2. Verify the worker is available
	// 3. Assign the task to the worker
	// 4. Update the task status

	// 5. Publish TaskAssignedEvent to Kafka
	if err := s.kafkaProducer.PublishTaskAssigned(ctx, taskID, workerID); err != nil {
		return err
	}

	return nil
}

// GetTaskStatus retrieves the current status of a task
func (s *SchedulerServiceImpl) GetTaskStatus(ctx context.Context, taskID string) (model.Status, error) {
	log.Printf("Getting status for task %s", taskID)
	// In a real implementation, this would:
	// 1. Get the task details from the task service
	// 2. Return the task status
	return model.StatusPending, nil
}

// Helper function to get task from task service
func (s *SchedulerServiceImpl) getTaskFromService(ctx context.Context, taskID string) (*model.Task, error) {
	// In a real implementation, this would call the task service via gRPC
	// For now, just return a dummy task
	return &model.Task{
		ID:     taskID,
		Status: model.StatusPending,
	}, nil
}

// Helper function to update task in task service
func (s *SchedulerServiceImpl) updateTaskInService(ctx context.Context, task *model.Task) error {
	// In a real implementation, this would call the task service via gRPC
	return nil
}
