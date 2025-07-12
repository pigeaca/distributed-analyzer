package producer

import (
	"context"
	"distributed-analyzer/libs/kafka"
	"distributed-analyzer/libs/model"
	pb "distributed-analyzer/libs/proto/kafka"
	taskpb "distributed-analyzer/libs/proto/task"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

type TaskProducer struct {
	*kafka.Producer
}

func NewTaskProducer(pr *kafka.Producer) *TaskProducer {
	return &TaskProducer{
		Producer: pr,
	}
}

// PublishTaskCreated publishes a TaskCreatedEvent to Kafka
func (p *TaskProducer) PublishTaskCreated(ctx context.Context, task *model.Task) error {
	event := &pb.TaskCreatedEvent{
		TaskId:    task.ID,
		CreatedAt: timestamppb.New(time.Now()),
	}

	return p.Producer.PublishEvent(ctx, "task-created", task.ID, event)
}

// PublishTaskStatusChanged publishes a TaskStatusChangedEvent to Kafka
func (p *TaskProducer) PublishTaskStatusChanged(ctx context.Context, taskID string, oldStatus, newStatus model.Status) error {
	event := &pb.TaskStatusChangedEvent{
		TaskId:    taskID,
		OldStatus: convertModelStatusToPbStatus(oldStatus),
		NewStatus: convertModelStatusToPbStatus(newStatus),
		ChangedAt: timestamppb.New(time.Now()),
	}

	return p.Producer.PublishEvent(ctx, "task-status-changed", taskID, event)
}

// PublishTaskCompleted publishes a TaskCompletedEvent to Kafka
func (p *TaskProducer) PublishTaskCompleted(ctx context.Context, task *model.Task) error {
	event := &pb.TaskCompletedEvent{
		TaskId:      task.ID,
		Result:      task.Output,
		CompletedAt: timestamppb.New(time.Now()),
	}

	return p.Producer.PublishEvent(ctx, "task-completed", task.ID, event)
}

// PublishTaskFailed publishes a TaskFailedEvent to Kafka
func (p *TaskProducer) PublishTaskFailed(ctx context.Context, taskID string, errorMsg string) error {
	event := &pb.TaskFailedEvent{
		TaskId:   taskID,
		Error:    errorMsg,
		FailedAt: timestamppb.New(time.Now()),
	}

	return p.Producer.PublishEvent(ctx, "task-failed", taskID, event)
}

// Helper function to convert model.Status to taskpb.Status
func convertModelStatusToPbStatus(status model.Status) taskpb.Status {
	switch status {
	case model.StatusPending:
		return taskpb.Status_STATUS_PENDING
	case model.StatusScheduled:
		return taskpb.Status_STATUS_SCHEDULED
	case model.StatusRunning:
		return taskpb.Status_STATUS_RUNNING
	case model.StatusCompleted:
		return taskpb.Status_STATUS_COMPLETED
	case model.StatusFailed:
		return taskpb.Status_STATUS_FAILED
	default:
		return taskpb.Status_STATUS_UNSPECIFIED
	}
}
