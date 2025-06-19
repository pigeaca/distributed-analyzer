package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/distributedmarketplace/internal/worker/service"
	"github.com/segmentio/kafka-go"
	"log"
	"time"

	pb "github.com/distributedmarketplace/pkg/proto/kafka"
)

// WorkerConsumer is a Kafka consumer for worker events
type WorkerConsumer struct {
	readers       map[string]*kafka.Reader
	workerService service.WorkerNodeService
	stopChannels  map[string]chan struct{}
}

// NewWorkerConsumer creates a new WorkerConsumer
func NewWorkerConsumer(brokers []string, groupID string, workerID string, workerService service.WorkerNodeService) *WorkerConsumer {
	consumer := &WorkerConsumer{
		readers:       make(map[string]*kafka.Reader),
		workerService: workerService,
		stopChannels:  make(map[string]chan struct{}),
	}

	// Create readers for each topic
	topics := []string{
		"task-assigned",
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
func (c *WorkerConsumer) Start(ctx context.Context) {
	for topic, reader := range c.readers {
		go c.consumeTopic(ctx, topic, reader, c.stopChannels[topic])
	}
}

// Stop stops consuming messages from Kafka
func (c *WorkerConsumer) Stop() {
	for topic, stopCh := range c.stopChannels {
		close(stopCh)
		if err := c.readers[topic].Close(); err != nil {
			log.Printf("Error closing Kafka reader for topic %s: %v", topic, err)
		}
	}
}

// consumeTopic consumes messages from a specific topic
func (c *WorkerConsumer) consumeTopic(ctx context.Context, topic string, reader *kafka.Reader, stopCh <-chan struct{}) {
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
func (c *WorkerConsumer) handleMessage(ctx context.Context, topic string, message kafka.Message) error {
	switch topic {
	case "task-assigned":
		return c.handleTaskAssigned(ctx, message)
	default:
		return fmt.Errorf("unknown topic: %s", topic)
	}
}

// handleTaskAssigned handles a TaskAssignedEvent
func (c *WorkerConsumer) handleTaskAssigned(ctx context.Context, message kafka.Message) error {
	var event pb.TaskAssignedEvent
	if err := json.Unmarshal(message.Value, &event); err != nil {
		return fmt.Errorf("failed to unmarshal TaskAssignedEvent: %w", err)
	}

	// Execute the task
	if err := c.workerService.ExecuteTask(ctx, event.TaskId); err != nil {
		return fmt.Errorf("failed to execute task: %w", err)
	}

	log.Printf("Task %s executed successfully", event.TaskId)
	return nil
}
