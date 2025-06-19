package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/distributedmarketplace/internal/billing/service"
	"github.com/segmentio/kafka-go"
	"log"
	"time"

	pb "github.com/distributedmarketplace/pkg/proto/kafka"
)

// BillingConsumer is a Kafka consumer for billing events
type BillingConsumer struct {
	readers        map[string]*kafka.Reader
	billingService service.BillingService
	stopChannels   map[string]chan struct{}
}

// NewBillingConsumer creates a new BillingConsumer
func NewBillingConsumer(brokers []string, groupID string, billingService service.BillingService) *BillingConsumer {
	consumer := &BillingConsumer{
		readers:        make(map[string]*kafka.Reader),
		billingService: billingService,
		stopChannels:   make(map[string]chan struct{}),
	}

	topics := []string{
		"task-completed",
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

func (c *BillingConsumer) Start(ctx context.Context) {
	for topic, reader := range c.readers {
		go c.consumeTopic(ctx, topic, reader, c.stopChannels[topic])
	}
}

func (c *BillingConsumer) Stop() {
	for topic, stopCh := range c.stopChannels {
		close(stopCh)
		if err := c.readers[topic].Close(); err != nil {
			log.Printf("Error closing Kafka reader for topic %s: %v", topic, err)
		}
	}
}

func (c *BillingConsumer) consumeTopic(ctx context.Context, topic string, reader *kafka.Reader, stopCh <-chan struct{}) {
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

func (c *BillingConsumer) handleMessage(ctx context.Context, topic string, message kafka.Message) error {
	switch topic {
	case "task-completed":
		return c.handleTaskCompleted(ctx, message)
	default:
		return fmt.Errorf("unknown topic: %s", topic)
	}
}

func (c *BillingConsumer) handleTaskCompleted(ctx context.Context, message kafka.Message) error {
	var event pb.TaskCompletedEvent
	if err := json.Unmarshal(message.Value, &event); err != nil {
		return fmt.Errorf("failed to unmarshal TaskCompletedEvent: %w", err)
	}
	record, err := c.billingService.ChargeTask(ctx, event.TaskId)
	if err != nil {
		return fmt.Errorf("failed to charge for task: %w", err)
	}

	log.Printf("Task %s charged successfully: %v", event.TaskId, record)
	return nil
}
