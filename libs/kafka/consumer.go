package kafka

import (
	"context"
	"github.com/segmentio/kafka-go"
	"log"
	"time"
)

type MessageHandler interface {
	HandleMessage(ctx context.Context, topic string, message kafka.Message) error
}

// Consumer is a Kafka consumer for billing events
type Consumer struct {
	readers      map[string]*kafka.Reader
	handler      MessageHandler
	stopChannels map[string]chan struct{}
}

func NewConsumer(topics []string, brokers []string, groupID string, msgConsumer MessageHandler) *Consumer {
	consumer := &Consumer{
		readers:      make(map[string]*kafka.Reader),
		handler:      msgConsumer,
		stopChannels: make(map[string]chan struct{}),
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

func (c *Consumer) Start(ctx context.Context) {
	for topic, reader := range c.readers {
		go c.consumeTopic(ctx, topic, reader, c.stopChannels[topic])
	}
}

func (c *Consumer) Stop(ctx context.Context) {
	// Create a channel to signal when all readers are closed
	done := make(chan struct{})

	go func() {
		for topic, stopCh := range c.stopChannels {
			close(stopCh)
			if err := c.readers[topic].Close(); err != nil {
				log.Printf("Error closing Kafka reader for topic %s: %v", topic, err)
			}
		}
		close(done)
	}()

	// Wait for either all readers to close or context to be canceled
	select {
	case <-done:
		return
	case <-ctx.Done():
		log.Printf("Context canceled while stopping Kafka consumer: %v", ctx.Err())
		return
	}
}

func (c *Consumer) consumeTopic(ctx context.Context, topic string, reader *kafka.Reader, stopCh <-chan struct{}) {
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
			if err := c.handler.HandleMessage(ctx, topic, message); err != nil {
				log.Printf("Error handling message from topic %s: %v", topic, err)
			}
		}
	}
}
