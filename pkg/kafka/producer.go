package kafka

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(brokers []string) *Producer {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Balancer: &kafka.LeastBytes{},
	}
	return &Producer{
		writer: writer,
	}
}

// Close closes the Kafka writer
func (p *Producer) Close() error {
	return p.writer.Close()
}

// PublishEvent Helper function to publish an event to Kafka
func (p *Producer) PublishEvent(ctx context.Context, topic, key string, event any) error {
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
