package config

import (
	configloader "distributed-analyzer/libs/config"
)

type Config struct {
	ServerConfig    configloader.ServerConfig   `yaml:",inline"`
	Kafka           KafkaConfig                 `yaml:"kafka"`
	Database        configloader.DatabaseConfig `yaml:"database"`
	Log             configloader.LogConfig      `yaml:"log"`
	ShutdownTimeout string                      `yaml:"shutdown_timeout" env:"SHUTDOWN_TIMEOUT" env-default:"30s"`
}

// KafkaConfig extends the common KafkaConfig with task-specific topics
type KafkaConfig struct {
	configloader.KafkaConfig `yaml:",inline"`
	Topics                   KafkaTopicsConfig `yaml:"topics"`
}

type KafkaTopicsConfig struct {
	Tasks   string `yaml:"tasks"   env:"KAFKA_TOPIC_TASKS"   env-default:"tasks"`
	Results string `yaml:"results" env:"KAFKA_TOPIC_RESULTS" env-default:"results"`
}
