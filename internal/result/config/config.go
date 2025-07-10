package config

import (
	commonConfig "github.com/distributedmarketplace/internal/common/config"
)

type Config struct {
	// Server settings
	commonConfig.ServerConfig `yaml:",inline"`

	// Kafka settings
	Kafka KafkaConfig `yaml:"kafka"`

	// Storage settings
	Storage StorageConfig `yaml:"storage"`

	// Database settings
	Database commonConfig.DatabaseConfig `yaml:"database"`

	// Aggregation settings
	Aggregation AggregationConfig `yaml:"aggregation"`

	// Log settings
	Log commonConfig.LogConfig `yaml:"log"`
}

// KafkaConfig extends the common KafkaConfig with result-specific settings
type KafkaConfig struct {
	// Embed the common KafkaConfig
	commonConfig.KafkaConfig `yaml:",inline"`

	// Override Topics with the result-specific topics
	Topics KafkaTopicsConfig `yaml:"topics"`
}

// SetDefaults sets default values for the KafkaConfig
func (c *KafkaConfig) SetDefaults() {
	if c.GroupID == "" {
		c.GroupID = "result-service"
	}
}

// KafkaTopicsConfig holds Kafka topics configuration
type KafkaTopicsConfig struct {
	// Results specify the result topic
	// Can be set via KAFKA_TOPIC_RESULTS environment variable
	// Default: results
	Results string `envconfig:"KAFKA_TOPIC_RESULTS" default:"results" yaml:"results"`

	// Completed specifies the task-completed topic
	// Can be set via KAFKA_TOPIC_COMPLETED environment variable
	// Default: task_completed
	Completed string `envconfig:"KAFKA_TOPIC_COMPLETED" default:"task_completed" yaml:"completed"`
}

// StorageConfig holds storage-related configuration
type StorageConfig struct {
	// ServiceGrpcAddr specifies the gRPC address of the storage service
	// Can be set via STORAGE_SERVICE_GRPC_ADDR environment variable
	// Default: localhost:9085
	ServiceGrpcAddr string `envconfig:"STORAGE_SERVICE_GRPC_ADDR" default:"localhost:9085" yaml:"service_grpc_addr"`

	// ResultTTL specifies the time-to-live for results
	// Can be set via STORAGE_RESULT_TTL environment variable
	// Default: 30d
	ResultTTL string `envconfig:"STORAGE_RESULT_TTL" default:"30d" yaml:"result_ttl"`
}

// AggregationConfig holds result aggregation configuration
type AggregationConfig struct {
	// BatchSize specifies the batch size for result aggregation
	// Can be set via AGGREGATION_BATCH_SIZE environment variable
	// Default: 100
	BatchSize int `envconfig:"AGGREGATION_BATCH_SIZE" default:"100" yaml:"batch_size"`

	// FlushInterval specifies the flush interval for result aggregation
	// Can be set via AGGREGATION_FLUSH_INTERVAL environment variable
	// Default: 5s
	FlushInterval string `envconfig:"AGGREGATION_FLUSH_INTERVAL" default:"5s" yaml:"flush_interval"`
}
