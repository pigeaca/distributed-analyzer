package kafka

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"

	pb "github.com/distributedmarketplace/pkg/proto/kafka"
)

// AuditProducer is a Kafka producer for audit events
type AuditProducer struct {
	writer *kafka.Writer
}

// NewAuditProducer creates a new AuditProducer
func NewAuditProducer(brokers []string) *AuditProducer {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Balancer: &kafka.LeastBytes{},
	}
	return &AuditProducer{
		writer: writer,
	}
}

// Close closes the Kafka writer
func (p *AuditProducer) Close() error {
	return p.writer.Close()
}

// PublishAuditEvent publishes an AuditEvent to Kafka
func (p *AuditProducer) PublishAuditEvent(ctx context.Context, userID string, action string, resource string, resourceID string) error {
	event := &pb.AuditEvent{
		UserId:     userID,
		Action:     action,
		Resource:   resource,
		ResourceId: resourceID,
		Timestamp:  timestamppb.New(time.Now()),
	}

	return p.publishEvent(ctx, "audit-event", resourceID, event)
}

// Helper function to publish an event to Kafka
func (p *AuditProducer) publishEvent(ctx context.Context, topic, key string, event interface{}) error {
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
