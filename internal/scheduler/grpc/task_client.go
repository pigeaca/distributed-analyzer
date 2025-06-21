package grpc

import (
	"context"
	"github.com/distributedmarketplace/internal/task/model"
	pb "github.com/distributedmarketplace/pkg/proto/task"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type TaskServiceGrpcClient struct {
	client pb.TaskServiceClient
	conn   *grpc.ClientConn
}

func NewTaskServiceGrpcClient(serverAddr string) (*TaskServiceGrpcClient, error) {
	conn, err := grpc.NewClient(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := pb.NewTaskServiceClient(conn)

	return &TaskServiceGrpcClient{
		client: client,
		conn:   conn,
	}, nil
}

func (t *TaskServiceGrpcClient) Close() error {
	return t.conn.Close()
}

// GetTask retrieves a task by its ID
func (t *TaskServiceGrpcClient) GetTask(ctx context.Context, id string) (*model.Task, error) {
	req := &pb.GetTaskRequest{
		Id: id,
	}

	resp, err := t.client.GetTask(ctx, req)
	if err != nil {
		return nil, err
	}

	return convertPbTaskToModelTask(resp.Task), nil
}

// UpdateTask updates an existing task
func (t *TaskServiceGrpcClient) UpdateTask(ctx context.Context, task *model.Task) (*model.Task, error) {
	pbTask := convertModelTaskToPbTask(task)

	req := &pb.UpdateTaskRequest{
		Task: pbTask,
	}

	resp, err := t.client.UpdateTask(ctx, req)
	if err != nil {
		return nil, err
	}

	return convertPbTaskToModelTask(resp.Task), nil
}

// ListTasks retrieves all tasks with optional filtering
func (t *TaskServiceGrpcClient) ListTasks(ctx context.Context) ([]*model.Task, error) {
	req := &pb.ListTasksRequest{}

	resp, err := t.client.ListTasks(ctx, req)
	if err != nil {
		return nil, err
	}

	tasks := make([]*model.Task, len(resp.Tasks))
	for i, pbTask := range resp.Tasks {
		tasks[i] = convertPbTaskToModelTask(pbTask)
	}

	return tasks, nil
}

// Helper functions to convert between model and protobuf types

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
