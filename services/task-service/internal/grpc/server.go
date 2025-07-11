package grpc

import (
	"context"
	"github.com/pigeaca/DistributedMarketplace/libs/model"
	pb "github.com/pigeaca/DistributedMarketplace/libs/proto/task"
	"github.com/pigeaca/DistributedMarketplace/services/task-service/internal/kafka/producer"
	"github.com/pigeaca/DistributedMarketplace/services/task-service/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
)

type TaskServer struct {
	pb.UnimplementedTaskServiceServer
	taskService service.TaskService
	producer    *producer.TaskProducer
}

func NewTaskServer(taskService service.TaskService, producer *producer.TaskProducer) *TaskServer {
	return &TaskServer{
		taskService: taskService,
		producer:    producer,
	}
}

// CreateTask creates a new task in the system
func (s *TaskServer) CreateTask(ctx context.Context, req *pb.CreateTaskRequest) (*pb.TaskResponse, error) {
	task := &model.Task{
		Name:        req.Name,
		Description: req.Description,
		Input:       req.Input,
	}

	// Convert resources if any
	if len(req.Resources) > 0 {
		task.Resources = make([]model.Resource, len(req.Resources))
		for i, r := range req.Resources {
			task.Resources[i] = model.Resource{
				Type:  r.Type,
				Value: int(r.Value),
			}
		}
	}

	createdTask, err := s.taskService.CreateTask(ctx, task)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create task: %v", err)
	}

	// Publish TaskCreatedEvent to Kafka
	if s.producer != nil {
		if err := s.producer.PublishTaskCreated(ctx, createdTask); err != nil {
			log.Printf("Failed to publish TaskCreatedEvent: %v", err)
			// Continue even if publishing fails
		}
	}

	return &pb.TaskResponse{
		Task: convertModelTaskToPbTask(createdTask),
	}, nil
}

// GetTask retrieves a task by its ID
func (s *TaskServer) GetTask(ctx context.Context, req *pb.GetTaskRequest) (*pb.TaskResponse, error) {
	task, err := s.taskService.GetTask(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "task not found: %v", err)
	}

	return &pb.TaskResponse{
		Task: convertModelTaskToPbTask(task),
	}, nil
}

// UpdateTask updates an existing task
func (s *TaskServer) UpdateTask(ctx context.Context, req *pb.UpdateTaskRequest) (*pb.TaskResponse, error) {
	if req.Task == nil {
		return nil, status.Errorf(codes.InvalidArgument, "task is required")
	}

	modelTask := convertPbTaskToModelTask(req.Task)
	updatedTask, err := s.taskService.UpdateTask(ctx, modelTask)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update task: %v", err)
	}

	return &pb.TaskResponse{
		Task: convertModelTaskToPbTask(updatedTask),
	}, nil
}

// DeleteTask removes a task from the system
func (s *TaskServer) DeleteTask(ctx context.Context, req *pb.DeleteTaskRequest) (*pb.DeleteTaskResponse, error) {
	err := s.taskService.DeleteTask(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete task: %v", err)
	}

	return &pb.DeleteTaskResponse{
		Success: true,
	}, nil
}

// ListTasks retrieves all tasks with optional filtering
func (s *TaskServer) ListTasks(ctx context.Context, req *pb.ListTasksRequest) (*pb.ListTasksResponse, error) {
	tasks, err := s.taskService.ListTasks(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list tasks: %v", err)
	}

	pbTasks := make([]*pb.Task, len(tasks))
	for i, task := range tasks {
		pbTasks[i] = convertModelTaskToPbTask(task)
	}

	return &pb.ListTasksResponse{
		Tasks: pbTasks,
	}, nil
}

// convertModelTaskToPbTask converts a model.Task to a pb.Task
func convertModelTaskToPbTask(task *model.Task) *pb.Task {
	pbTask := &pb.Task{
		Id:          task.ID,
		Name:        task.Name,
		Description: task.Description,
		Status:      convertModelStatusToPbStatus(task.Status),
		Input:       task.Input,
		Output:      task.Output,
		CreatedAt:   timestamppb.New(task.CreatedAt),
		UpdatedAt:   timestamppb.New(task.UpdatedAt),
	}

	if !task.CompletedAt.IsZero() {
		pbTask.CompletedAt = timestamppb.New(task.CompletedAt)
	}

	if len(task.Resources) > 0 {
		pbTask.Resources = make([]*pb.Resource, len(task.Resources))
		for i, r := range task.Resources {
			pbTask.Resources[i] = &pb.Resource{
				Type:  r.Type,
				Value: int32(r.Value),
			}
		}
	}

	return pbTask
}

// convertPbTaskToModelTask converts a pb.Task to a model.Task
func convertPbTaskToModelTask(pbTask *pb.Task) *model.Task {
	task := &model.Task{
		ID:          pbTask.Id,
		Name:        pbTask.Name,
		Description: pbTask.Description,
		Status:      convertPbStatusToModelStatus(pbTask.Status),
		Input:       pbTask.Input,
		Output:      pbTask.Output,
	}

	if pbTask.CreatedAt != nil {
		task.CreatedAt = pbTask.CreatedAt.AsTime()
	}

	if pbTask.UpdatedAt != nil {
		task.UpdatedAt = pbTask.UpdatedAt.AsTime()
	}

	if pbTask.CompletedAt != nil {
		task.CompletedAt = pbTask.CompletedAt.AsTime()
	}

	if len(pbTask.Resources) > 0 {
		task.Resources = make([]model.Resource, len(pbTask.Resources))
		for i, r := range pbTask.Resources {
			task.Resources[i] = model.Resource{
				Type:  r.Type,
				Value: int(r.Value),
			}
		}
	}

	return task
}

// convertModelStatusToPbStatus converts a model.Status to a pb.Status
func convertModelStatusToPbStatus(status model.Status) pb.Status {
	switch status {
	case model.StatusPending:
		return pb.Status_STATUS_PENDING
	case model.StatusScheduled:
		return pb.Status_STATUS_SCHEDULED
	case model.StatusRunning:
		return pb.Status_STATUS_RUNNING
	case model.StatusCompleted:
		return pb.Status_STATUS_COMPLETED
	case model.StatusFailed:
		return pb.Status_STATUS_FAILED
	default:
		return pb.Status_STATUS_UNSPECIFIED
	}
}

// convertPbStatusToModelStatus converts a pb.Status to a model.Status
func convertPbStatusToModelStatus(status pb.Status) model.Status {
	switch status {
	case pb.Status_STATUS_PENDING:
		return model.StatusPending
	case pb.Status_STATUS_SCHEDULED:
		return model.StatusScheduled
	case pb.Status_STATUS_RUNNING:
		return model.StatusRunning
	case pb.Status_STATUS_COMPLETED:
		return model.StatusCompleted
	case pb.Status_STATUS_FAILED:
		return model.StatusFailed
	default:
		return model.StatusPending
	}
}
