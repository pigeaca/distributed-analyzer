package config

import (
	"distributed-analyzer/libs/common/config"
)

type Config struct {
	// Server settings
	ServerConfig config.ServerConfig `yaml:",inline"`

	// Kafka settings
	Kafka KafkaConfig `yaml:"kafka"`

	// Database settings
	Database config.DatabaseConfig `yaml:"database"`

	// Log settings
	Log config.LogConfig `yaml:"log"`

	// ShutdownTimeout specifies how long to wait for graceful shutdown
	// Can be set via SHUTDOWN_TIMEOUT environment variable
	// Default: 30s
	ShutdownTimeout string `envconfig:"SHUTDOWN_TIMEOUT" default:"30s" yaml:"shutdown_timeout"`
}

// KafkaConfig extends the common KafkaConfig with task-specific settings
type KafkaConfig struct {
	// Embed the common KafkaConfig
	config.KafkaConfig `yaml:",inline"`

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
	// Tasks specifies the task topic
	// Can be set via KAFKA_TOPIC_TASKS environment variable
	// Default: tasks
	Tasks string `envconfig:"KAFKA_TOPIC_TASKS" default:"tasks" yaml:"tasks"`

	// Results specify the result topic
	// Can be set via KAFKA_TOPIC_RESULTS environment variable
	// Default: results
	Results string `envconfig:"KAFKA_TOPIC_RESULTS" default:"results" yaml:"results"`
}
