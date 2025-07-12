// Package kafka provides application components for Kafka integration.
package kafka

import (
	"context"
	"fmt"
	"log"

	"distributed-analyzer/libs/kafka"
)

// ConsumerComponent wraps a Kafka consumer as an application component.
type ConsumerComponent struct {
	consumer *kafka.Consumer
}

// NewKafkaComponent creates a new KafkaConsumerComponent with the provided consumer.
// It returns an error if the consumer is nil.
func NewKafkaComponent(consumer *kafka.Consumer) *ConsumerComponent {
	if consumer == nil {
		log.Fatal("kafka consumer cannot be nil")
	}
	return &ConsumerComponent{consumer: consumer}
}

// Start begins consuming messages from Kafka topics.
// It starts the consumer in a separate goroutine to avoid blocking.
func (k *ConsumerComponent) Start(ctx context.Context) error {
	if k.consumer == nil {
		return fmt.Errorf("kafka consumer is nil")
	}
	go k.consumer.Start(ctx)
	return nil
}

// Stop gracefully stops the Kafka consumer.
// It respects the provided context for cancellation.
func (k *ConsumerComponent) Stop(ctx context.Context) error {
	if k.consumer == nil {
		return nil // Nothing to stop
	}
	k.consumer.Stop(ctx)
	return nil
}

// Name returns the component name for logging and identification.
func (k *ConsumerComponent) Name() string {
	return "KafkaConsumerComponent"
}

// ProducerComponent wraps a Kafka producer as an application component.
type ProducerComponent struct {
	producer *kafka.Producer
}

// NewKafkaProducerComponent creates a new KafkaProducerComponent with the provided producer.
// It returns an error if the producer is nil.
func NewKafkaProducerComponent(producer *kafka.Producer) *ProducerComponent {
	if producer == nil {
		log.Fatal("kafka producer cannot be nil")
	}
	return &ProducerComponent{producer: producer}
}

// Producer returns the underlying Kafka producer.
// This method allows other components to access the producer for sending messages.
func (k *ProducerComponent) Producer() *kafka.Producer {
	return k.producer
}

// Start initializes the Kafka producer component.
// For producers, this is typically a no-op as the producer is ready upon creation.
func (k *ProducerComponent) Start(ctx context.Context) error {
	if k.producer == nil {
		return fmt.Errorf("kafka producer is nil")
	}
	return nil
}

// Stop gracefully closes the Kafka producer.
// It ensures all pending messages are sent before closing.
func (k *ProducerComponent) Stop(ctx context.Context) error {
	if k.producer == nil {
		return nil // Nothing to close
	}

	if err := k.producer.Close(); err != nil {
		return fmt.Errorf("failed to close kafka producer: %w", err)
	}
	return nil
}

// Name returns the component name for logging and identification.
func (k *ProducerComponent) Name() string {
	return "KafkaProducerComponent"
}
