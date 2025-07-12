package config

import (
	configloader "distributed-analyzer/libs/config"
)

type Config struct {
	configloader.ServerConfig `yaml:",inline"`

	Kafka           KafkaConfig            `yaml:"kafka"`
	Scheduling      SchedulingConfig       `yaml:"scheduling"`
	Log             configloader.LogConfig `yaml:"log"`
	ShutdownTimeout string                 `yaml:"shutdown_timeout" env:"SHUTDOWN_TIMEOUT" env-default:"30s"`
}

// KafkaConfig extends the common KafkaConfig with scheduler-specific topics
type KafkaConfig struct {
	configloader.KafkaConfig `yaml:",inline"`

	Topics KafkaTopicsConfig `yaml:"topics"`
}

type KafkaTopicsConfig struct {
	Tasks       string `yaml:"tasks"       env:"KAFKA_TOPIC_TASKS"       env-default:"tasks"`
	Assignments string `yaml:"assignments" env:"KAFKA_TOPIC_ASSIGNMENTS" env-default:"task_assignments"`
}

type SchedulingConfig struct {
	MaxRetries     int    `yaml:"max_retries"      env:"SCHEDULING_MAX_RETRIES"      env-default:"3"`
	RetryDelay     string `yaml:"retry_delay"      env:"SCHEDULING_RETRY_DELAY"      env-default:"5s"`
	DefaultTimeout string `yaml:"default_timeout"  env:"SCHEDULING_DEFAULT_TIMEOUT"  env-default:"300s"`
}
