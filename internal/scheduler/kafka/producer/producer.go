package producer

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"

	pb "github.com/distributedmarketplace/pkg/proto/kafka"
)

// SchedulerProducer is a Kafka producer for scheduler events
type SchedulerProducer struct {
	writer *kafka.Writer
}

// NewSchedulerProducer creates a new SchedulerProducer
func NewSchedulerProducer(brokers []string) *SchedulerProducer {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Balancer: &kafka.LeastBytes{},
	}
	return &SchedulerProducer{
		writer: writer,
	}
}

// Close closes the Kafka writer
func (p *SchedulerProducer) Close() error {
	return p.writer.Close()
}

// PublishTaskAssigned publishes a TaskAssignedEvent to Kafka
func (p *SchedulerProducer) PublishTaskAssigned(ctx context.Context, taskID string, workerID string) error {
	event := &pb.TaskAssignedEvent{
		TaskId:     taskID,
		WorkerId:   workerID,
		AssignedAt: timestamppb.New(time.Now()),
	}

	return p.publishEvent(ctx, "task-assigned", taskID, event)
}

// PublishTaskScheduled publishes a TaskScheduledEvent to Kafka
func (p *SchedulerProducer) PublishTaskScheduled(ctx context.Context, taskID string, workerIDs []string) error {
	event := &pb.TaskScheduledEvent{
		TaskId:      taskID,
		WorkerIds:   workerIDs,
		ScheduledAt: timestamppb.New(time.Now()),
	}

	return p.publishEvent(ctx, "task-scheduled", taskID, event)
}

// Helper function to publish an event to Kafka
func (p *SchedulerProducer) publishEvent(ctx context.Context, topic, key string, event interface{}) error {
	value, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return p.writer.WriteMessages(ctx, kafka.Message{
		Topic: topic,
		Key:   []byte(key),
		Value: value,
	})
}
