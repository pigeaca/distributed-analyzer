package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/distributedmarketplace/internal/result/service"
	"github.com/segmentio/kafka-go"
	"log"
	"time"

	pb "github.com/distributedmarketplace/pkg/proto/kafka"
)

// ResultConsumer is a Kafka consumer for result events
type ResultConsumer struct {
	readers       map[string]*kafka.Reader
	resultService service.ResultAggregatorService
	stopChannels  map[string]chan struct{}
}

// NewResultConsumer creates a new ResultConsumer
func NewResultConsumer(brokers []string, groupID string, resultService service.ResultAggregatorService) *ResultConsumer {
	consumer := &ResultConsumer{
		readers:       make(map[string]*kafka.Reader),
		resultService: resultService,
		stopChannels:  make(map[string]chan struct{}),
	}

	// Create readers for each topic
	topics := []string{
		"subtask-completed",
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
func (c *ResultConsumer) Start(ctx context.Context) {
	for topic, reader := range c.readers {
		go c.consumeTopic(ctx, topic, reader, c.stopChannels[topic])
	}
}

// Stop stops consuming messages from Kafka
func (c *ResultConsumer) Stop() {
	for topic, stopCh := range c.stopChannels {
		close(stopCh)
		if err := c.readers[topic].Close(); err != nil {
			log.Printf("Error closing Kafka reader for topic %s: %v", topic, err)
		}
	}
}

// consumeTopic consumes messages from a specific topic
func (c *ResultConsumer) consumeTopic(ctx context.Context, topic string, reader *kafka.Reader, stopCh <-chan struct{}) {
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
func (c *ResultConsumer) handleMessage(ctx context.Context, topic string, message kafka.Message) error {
	switch topic {
	case "subtask-completed":
		return c.handleSubTaskCompleted(ctx, message)
	default:
		return fmt.Errorf("unknown topic: %s", topic)
	}
}

// handleSubTaskCompleted handles a SubTaskCompletedEvent
func (c *ResultConsumer) handleSubTaskCompleted(ctx context.Context, message kafka.Message) error {
	var event pb.SubTaskCompletedEvent
	if err := json.Unmarshal(message.Value, &event); err != nil {
		return fmt.Errorf("failed to unmarshal SubTaskCompletedEvent: %w", err)
	}

	// Save the partial result
	if err := c.resultService.SavePartialResult(ctx, event.TaskId, event.SubtaskId, event.Result); err != nil {
		return fmt.Errorf("failed to save partial result: %w", err)
	}

	// Try to finalize the result (this will check if all subtasks are completed)
	if err := c.resultService.FinalizeResult(ctx, event.TaskId); err != nil {
		log.Printf("Not finalizing result for task %s yet: %v", event.TaskId, err)
	} else {
		log.Printf("Result for task %s finalized successfully", event.TaskId)
	}

	return nil
}
