package config

import (
	commonConfig "distributed-analyzer/libs/config"
)

type Config struct {
	// Server settings
	commonConfig.ServerConfig `yaml:",inline"`

	// Kafka settings
	Kafka KafkaConfig `yaml:"kafka"`

	// Storage settings
	Storage StorageConfig `yaml:"storage"`

	// Database settings
	commonConfig.DatabaseConfig `yaml:"database"`

	// Aggregation settings
	Aggregation AggregationConfig `yaml:"aggregation"`

	// Log settings
	commonConfig.LogConfig `yaml:"log"`
}

type KafkaConfig struct {
	commonConfig.KafkaConfig `yaml:",inline"`

	Topics KafkaTopicsConfig `yaml:"topics"`
}

type KafkaTopicsConfig struct {
	Results   string `yaml:"results"   env:"KAFKA_TOPIC_RESULTS"   env-default:"results"`
	Completed string `yaml:"completed" env:"KAFKA_TOPIC_COMPLETED" env-default:"task_completed"`
}

type StorageConfig struct {
	ServiceGrpcAddr string `yaml:"service_grpc_addr" env:"STORAGE_SERVICE_GRPC_ADDR" env-default:"localhost:9085"`
	ResultTTL       string `yaml:"result_ttl"        env:"STORAGE_RESULT_TTL"        env-default:"30d"`
}

type AggregationConfig struct {
	BatchSize     int    `yaml:"batch_size"     env:"AGGREGATION_BATCH_SIZE"     env-default:"100"`
	FlushInterval string `yaml:"flush_interval" env:"AGGREGATION_FLUSH_INTERVAL" env-default:"5s"`
}
