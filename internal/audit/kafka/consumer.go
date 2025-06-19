package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/distributedmarketplace/internal/audit/model"
	"github.com/distributedmarketplace/internal/audit/service"
	"github.com/segmentio/kafka-go"
	"log"
	"time"

	pb "github.com/distributedmarketplace/pkg/proto/kafka"
)

// AuditConsumer is a Kafka consumer for audit events
type AuditConsumer struct {
	readers      map[string]*kafka.Reader
	auditService service.AuditService
	stopChannels map[string]chan struct{}
}

// NewAuditConsumer creates a new AuditConsumer
func NewAuditConsumer(brokers []string, groupID string, auditService service.AuditService) *AuditConsumer {
	consumer := &AuditConsumer{
		readers:      make(map[string]*kafka.Reader),
		auditService: auditService,
		stopChannels: make(map[string]chan struct{}),
	}

	// Create readers for each topic
	topics := []string{
		"audit-event",
	}

	for _, topic := range topics {
		consumer.readers[topic] = kafka.NewReader(kafka.ReaderConfig{
			Brokers:     brokers,
			Topic:       topic,
			GroupID:     groupID,
			MinBytes:    10e3, // 10KB
			MaxBytes:    10e6, // 10MB
			MaxWait:     1 * time.Second,
			StartOffset: kafka.FirstOffset,
		})
		consumer.stopChannels[topic] = make(chan struct{})
	}

	return consumer
}

// Start starts consuming messages from Kafka
func (c *AuditConsumer) Start(ctx context.Context) {
	for topic, reader := range c.readers {
		go c.consumeTopic(ctx, topic, reader, c.stopChannels[topic])
	}
}

// Stop stops consuming messages from Kafka
func (c *AuditConsumer) Stop() {
	for topic, stopCh := range c.stopChannels {
		close(stopCh)
		if err := c.readers[topic].Close(); err != nil {
			log.Printf("Error closing Kafka reader for topic %s: %v", topic, err)
		}
	}
}

// consumeTopic consumes messages from a specific topic
func (c *AuditConsumer) consumeTopic(ctx context.Context, topic string, reader *kafka.Reader, stopCh <-chan struct{}) {
	for {
		select {
		case <-stopCh:
			return
		case <-ctx.Done():
			return
		default:
			message, err := reader.ReadMessage(ctx)
			if err != nil {
				log.Printf("Error reading message from topic %s: %v", topic, err)
				continue
			}

			if err := c.handleMessage(ctx, topic, message); err != nil {
				log.Printf("Error handling message from topic %s: %v", topic, err)
			}
		}
	}
}

// handleMessage handles a message from Kafka
func (c *AuditConsumer) handleMessage(ctx context.Context, topic string, message kafka.Message) error {
	switch topic {
	case "audit-event":
		return c.handleAuditEvent(ctx, message)
	default:
		return fmt.Errorf("unknown topic: %s", topic)
	}
}

// handleAuditEvent handles an AuditEvent
func (c *AuditConsumer) handleAuditEvent(ctx context.Context, message kafka.Message) error {
	var event pb.AuditEvent
	if err := json.Unmarshal(message.Value, &event); err != nil {
		return fmt.Errorf("failed to unmarshal AuditEvent: %w", err)
	}

	// Convert action string to model.AuditAction
	var action model.AuditAction
	switch event.Action {
	case "CREATE":
		action = model.ActionCreate
	case "READ":
		action = model.ActionRead
	case "UPDATE":
		action = model.ActionUpdate
	case "DELETE":
		action = model.ActionDelete
	default:
		action = model.ActionCreate // Default to CREATE if action is not recognized
	}

	// Log the action
	log, err := c.auditService.LogAction(ctx, event.UserId, action, event.Resource, event.ResourceId)
	if err != nil {
		return fmt.Errorf("failed to log action: %w", err)
	}

	fmt.Printf("Audit log created: %v\n", log)
	return nil
}
