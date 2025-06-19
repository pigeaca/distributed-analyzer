package kafka

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"

	pb "github.com/distributedmarketplace/pkg/proto/kafka"
)

// ResultProducer is a Kafka producer for result events
type ResultProducer struct {
	writer *kafka.Writer
}

// NewResultProducer creates a new ResultProducer
func NewResultProducer(brokers []string) *ResultProducer {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Balancer: &kafka.LeastBytes{},
	}
	return &ResultProducer{
		writer: writer,
	}
}

// Close closes the Kafka writer
func (p *ResultProducer) Close() error {
	return p.writer.Close()
}

// PublishTaskCompleted publishes a TaskCompletedEvent to Kafka
func (p *ResultProducer) PublishTaskCompleted(ctx context.Context, taskID string, result map[string]string) error {
	event := &pb.TaskCompletedEvent{
		TaskId:      taskID,
		Result:      result,
		CompletedAt: timestamppb.New(time.Now()),
	}

	return p.publishEvent(ctx, "task-completed", taskID, event)
}

// Helper function to publish an event to Kafka
func (p *ResultProducer) publishEvent(ctx context.Context, topic, key string, event interface{}) error {
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
