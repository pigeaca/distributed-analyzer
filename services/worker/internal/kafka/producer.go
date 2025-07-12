package kafka

import (
	"context"
	pb "distributed-analyzer/libs/proto/kafka"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

// WorkerProducer is a Kafka producer for worker events
type WorkerProducer struct {
	writer *kafka.Writer
}

// NewWorkerProducer creates a new WorkerProducer
func NewWorkerProducer(brokers []string) *WorkerProducer {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Balancer: &kafka.LeastBytes{},
	}
	return &WorkerProducer{
		writer: writer,
	}
}

// Close closes the Kafka writer
func (p *WorkerProducer) Close() error {
	return p.writer.Close()
}

// PublishSubTaskCompleted publishes a SubTaskCompletedEvent to Kafka
func (p *WorkerProducer) PublishSubTaskCompleted(ctx context.Context, subtaskID string, taskID string, workerID string, result map[string]string) error {
	event := &pb.SubTaskCompletedEvent{
		SubtaskId:   subtaskID,
		TaskId:      taskID,
		WorkerId:    workerID,
		Result:      result,
		CompletedAt: timestamppb.New(time.Now()),
	}

	return p.publishEvent(ctx, "subtask-completed", subtaskID, event)
}

// PublishWorkerStatusChanged publishes a WorkerStatusChangedEvent to Kafka
func (p *WorkerProducer) PublishWorkerStatusChanged(ctx context.Context, workerID string, oldStatus string, newStatus string) error {
	event := &pb.WorkerStatusChangedEvent{
		WorkerId:  workerID,
		OldStatus: oldStatus,
		NewStatus: newStatus,
		ChangedAt: timestamppb.New(time.Now()),
	}

	return p.publishEvent(ctx, "worker-status-changed", workerID, event)
}

// Helper function to publish an event to Kafka
func (p *WorkerProducer) publishEvent(ctx context.Context, topic, key string, event interface{}) error {
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
