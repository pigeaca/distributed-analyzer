package config

import (
	configloader "distributed-analyzer/libs/config"
)

type Config struct {
	configloader.ServerConfig `yaml:",inline"`

	Kafka    KafkaConfig            `yaml:"kafka"`
	Services ServicesConfig         `yaml:"services"`
	Worker   WorkerConfig           `yaml:"worker"`
	Log      configloader.LogConfig `yaml:"log"`
}

// KafkaConfig extends the common KafkaConfig with worker-specific topics
type KafkaConfig struct {
	configloader.KafkaConfig `yaml:",inline"`
	Topics                   KafkaTopicsConfig `yaml:"topics"`
}

type KafkaTopicsConfig struct {
	Assignments string `yaml:"assignments" env:"KAFKA_TOPIC_ASSIGNMENTS" env-default:"task_assignments"`
	Results     string `yaml:"results"     env:"KAFKA_TOPIC_RESULTS"     env-default:"results"`
}

type ServicesConfig struct {
	Storage configloader.ServiceConnectionConfig `yaml:"storage"`
	Result  configloader.ServiceConnectionConfig `yaml:"result"`
}

type WorkerConfig struct {
	Capabilities       []string      `yaml:"capabilities"          env:"WORKER_CAPABILITIES"          env-default:"go_build,go_test,go_lint,go_benchmark,go_race"`
	MaxConcurrentTasks int           `yaml:"max_concurrent_tasks"  env:"WORKER_MAX_CONCURRENT_TASKS"  env-default:"5"`
	TaskTimeout        string        `yaml:"task_timeout"          env:"WORKER_TASK_TIMEOUT"          env-default:"300s"`
	Sandbox            SandboxConfig `yaml:"sandbox"`
}

type SandboxConfig struct {
	Enabled   bool            `yaml:"enabled" env:"SANDBOX_ENABLED" env-default:"true"`
	Type      string          `yaml:"type"    env:"SANDBOX_TYPE"    env-default:"docker"`
	Image     string          `yaml:"image"   env:"SANDBOX_IMAGE"   env-default:"golang:1.20-alpine"`
	Resources ResourcesConfig `yaml:"resources"`
}

type ResourcesConfig struct {
	CPULimit    int    `yaml:"cpu_limit"     env:"RESOURCES_CPU_LIMIT"     env-default:"1"`
	MemoryLimit string `yaml:"memory_limit"  env:"RESOURCES_MEMORY_LIMIT"  env-default:"512MB"`
}
