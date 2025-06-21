package producer

import (
	"context"
	"github.com/distributedmarketplace/pkg/kafka"
	pb "github.com/distributedmarketplace/pkg/proto/kafka"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

// SchedulerProducer is a Kafka producer for scheduler events
type SchedulerProducer struct {
	*kafka.Producer
}

func NewSchedulerProducer(pr *kafka.Producer) *SchedulerProducer {
	return &SchedulerProducer{
		Producer: pr,
	}
}

// PublishTaskAssigned publishes a TaskAssignedEvent to Kafka
func (p *SchedulerProducer) PublishTaskAssigned(ctx context.Context, taskID string, workerID string) error {
	event := &pb.TaskAssignedEvent{
		TaskId:     taskID,
		WorkerId:   workerID,
		AssignedAt: timestamppb.New(time.Now()),
	}

	return p.Producer.PublishEvent(ctx, "task-assigned", taskID, event)
}

// PublishTaskScheduled publishes a TaskScheduledEvent to Kafka
func (p *SchedulerProducer) PublishTaskScheduled(ctx context.Context, taskID string, workerIDs []string) error {
	event := &pb.TaskScheduledEvent{
		TaskId:      taskID,
		WorkerIds:   workerIDs,
		ScheduledAt: timestamppb.New(time.Now()),
	}

	return p.Producer.PublishEvent(ctx, "task-scheduled", taskID, event)
}
