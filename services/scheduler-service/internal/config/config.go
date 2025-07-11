package config

import "github.com/pigeaca/DistributedMarketplace/libs/common/config"

type Config struct {
	// Server settings
	config.ServerConfig `yaml:",inline"`

	// Kafka settings
	Kafka KafkaConfig `yaml:"kafka"`

	// Scheduling settings
	Scheduling SchedulingConfig `yaml:"scheduling"`

	// Log settings
	Log config.LogConfig `yaml:"log"`
}

// KafkaConfig extends the common KafkaConfig with scheduler-specific settings
type KafkaConfig struct {
	// Embed the common KafkaConfig
	config.KafkaConfig `yaml:",inline"`

	// Override Topics with the scheduler-specific topics
	Topics KafkaTopicsConfig `yaml:"topics"`
}

// SetDefaults sets default values for the KafkaConfig
func (c *KafkaConfig) SetDefaults() {
	if c.GroupID == "" {
		c.GroupID = "scheduler-service"
	}
}

// KafkaTopicsConfig holds Kafka topics configuration
type KafkaTopicsConfig struct {
	// Tasks specifies the task topic
	// Can be set via KAFKA_TOPIC_TASKS environment variable
	// Default: tasks
	Tasks string `envconfig:"KAFKA_TOPIC_TASKS" default:"tasks" yaml:"tasks"`

	// Assignments specify the task assignments topic
	// Can be set via KAFKA_TOPIC_ASSIGNMENTS environment variable
	// Default: task_assignments
	Assignments string `envconfig:"KAFKA_TOPIC_ASSIGNMENTS" default:"task_assignments" yaml:"assignments"`
}

// SchedulingConfig holds task scheduling configuration
type SchedulingConfig struct {
	// MaxRetries specifies the maximum number of retries for a task
	// Can be set via SCHEDULING_MAX_RETRIES environment variable
	// Default: 3
	MaxRetries int `envconfig:"SCHEDULING_MAX_RETRIES" default:"3" yaml:"max_retries"`

	// RetryDelay specifies the delay between retries
	// Can be set via SCHEDULING_RETRY_DELAY environment variable
	// Default: 5s
	RetryDelay string `envconfig:"SCHEDULING_RETRY_DELAY" default:"5s" yaml:"retry_delay"`

	// DefaultTimeout specifies the default task timeout
	// Can be set via SCHEDULING_DEFAULT_TIMEOUT environment variable
	// Default: 300s
	DefaultTimeout string `envconfig:"SCHEDULING_DEFAULT_TIMEOUT" default:"300s" yaml:"default_timeout"`
}
