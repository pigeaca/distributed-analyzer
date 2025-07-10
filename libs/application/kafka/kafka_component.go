package app

import (
	"context"
	"github.com/pigeaca/DistributedMarketplace/libs/kafka"
)

type KafkaComponent struct {
	consumer *kafka.Consumer
}

func NewKafkaComponent(consumer *kafka.Consumer) *KafkaComponent {
	return &KafkaComponent{consumer: consumer}
}

func (k *KafkaComponent) Start(ctx context.Context) error {
	go k.consumer.Start(ctx)
	return nil
}

func (k *KafkaComponent) Stop(ctx context.Context) error {
	k.consumer.Stop()
	return nil
}

func (k *KafkaComponent) Name() string {
	return "KafkaComponent"
}
