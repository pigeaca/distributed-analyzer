package config

import (
	commonConfig "github.com/distributedmarketplace/internal/common/config"
)

type Config struct {
	// Server settings
	commonConfig.ServerConfig `yaml:",inline"`

	// Kafka settings
	Kafka KafkaConfig `yaml:"kafka"`

	// Database settings
	Database commonConfig.DatabaseConfig `yaml:"database"`

	// Log settings
	Log commonConfig.LogConfig `yaml:"log"`
}

// KafkaConfig extends the common KafkaConfig with task-specific settings
type KafkaConfig struct {
	// Embed the common KafkaConfig
	commonConfig.KafkaConfig `yaml:",inline"`

	// Override Topics with the task-specific topics
	Topics KafkaTopicsConfig `yaml:"topics"`
}

// SetDefaults sets default values for the KafkaConfig
func (c *KafkaConfig) SetDefaults() {
	if c.GroupID == "" {
		c.GroupID = "task-service"
	}
}

// KafkaTopicsConfig holds Kafka topics configuration
type KafkaTopicsConfig struct {
	// Tasks specifies the tasks topic
	// Can be set via KAFKA_TOPIC_TASKS environment variable
	// Default: tasks
	Tasks string `envconfig:"KAFKA_TOPIC_TASKS" default:"tasks" yaml:"tasks"`

	// Results specifies the results topic
	// Can be set via KAFKA_TOPIC_RESULTS environment variable
	// Default: results
	Results string `envconfig:"KAFKA_TOPIC_RESULTS" default:"results" yaml:"results"`
}
