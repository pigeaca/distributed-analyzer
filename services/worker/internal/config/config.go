package config

import (
	"distributed-analyzer/libs/config"
)

type Config struct {
	// Server settings
	configloader.ServerConfig `yaml:",inline"`

	// Kafka settings
	Kafka KafkaConfig `yaml:"kafka"`

	// Services settings
	Services ServicesConfig `yaml:"services"`

	// Worker settings
	Worker WorkerConfig `yaml:"worker"`

	// Log settings
	Log configloader.LogConfig `yaml:"log"`
}

// KafkaConfig extends the common KafkaConfig with worker-specific settings
type KafkaConfig struct {
	// Embed the common KafkaConfig
	configloader.KafkaConfig `yaml:",inline"`

	// Override Topics with the worker-specific topics
	Topics KafkaTopicsConfig `yaml:"topics"`
}

// KafkaTopicsConfig holds Kafka topics configuration
type KafkaTopicsConfig struct {
	// Assignments specifies the task assignments topic
	// Can be set via KAFKA_TOPIC_ASSIGNMENTS environment variable
	// Default: task_assignments
	Assignments string `envconfig:"KAFKA_TOPIC_ASSIGNMENTS" default:"task_assignments" yaml:"assignments"`

	// Results specifies the results topic
	// Can be set via KAFKA_TOPIC_RESULTS environment variable
	// Default: results
	Results string `envconfig:"KAFKA_TOPIC_RESULTS" default:"results" yaml:"results"`
}

// ServicesConfig holds service connection configuration
type ServicesConfig struct {
	// Storage service configuration
	Storage configloader.ServiceConnectionConfig `yaml:"storage"`

	// Result service configuration
	Result configloader.ServiceConnectionConfig `yaml:"result"`
}

// WorkerConfig holds worker-specific configuration
type WorkerConfig struct {
	// Capabilities specifies the worker capabilities
	// Can be set via WORKER_CAPABILITIES environment variable
	// Default: go_build,go_test,go_lint, go_benchmark, go_race
	Capabilities []string `envconfig:"WORKER_CAPABILITIES" default:"go_build,go_test,go_lint,go_benchmark,go_race" yaml:"capabilities"`

	// MaxConcurrentTasks specifies the maximum number of concurrent tasks
	// Can be set via WORKER_MAX_CONCURRENT_TASKS environment variable
	// Default: 5
	MaxConcurrentTasks int `envconfig:"WORKER_MAX_CONCURRENT_TASKS" default:"5" yaml:"max_concurrent_tasks"`

	// TaskTimeout specifies the task timeout
	// Can be set via WORKER_TASK_TIMEOUT environment variable
	// Default: 300s
	TaskTimeout string `envconfig:"WORKER_TASK_TIMEOUT" default:"300s" yaml:"task_timeout"`

	// Sandbox settings
	Sandbox SandboxConfig `yaml:"sandbox"`
}

// SandboxConfig holds sandbox configuration
type SandboxConfig struct {
	// Enabled specifies whether sandboxing is enabled
	// Can be set via SANDBOX_ENABLED environment variable
	// Default: true
	Enabled bool `envconfig:"SANDBOX_ENABLED" default:"true" yaml:"enabled"`

	// Type specifies the sandbox type
	// Can be set via SANDBOX_TYPE environment variable
	// Default: docker
	Type string `envconfig:"SANDBOX_TYPE" default:"docker" yaml:"type"`

	// Image specifies the Docker image for the sandbox
	// Can be set via SANDBOX_IMAGE environment variable
	// Default: golang:1.20-alpine
	Image string `envconfig:"SANDBOX_IMAGE" default:"golang:1.20-alpine" yaml:"image"`

	// Resources specify the sandbox resource limits
	Resources ResourcesConfig `yaml:"resources"`
}

// ResourcesConfig holds resource limit configuration
type ResourcesConfig struct {
	// CPULimit specifies the CPU limit
	// Can be set via RESOURCES_CPU_LIMIT environment variable
	// Default: 1
	CPULimit int `envconfig:"RESOURCES_CPU_LIMIT" default:"1" yaml:"cpu_limit"`

	// MemoryLimit specifies the memory limit
	// Can be set via RESOURCES_MEMORY_LIMIT environment variable
	// Default: 512MB
	MemoryLimit string `envconfig:"RESOURCES_MEMORY_LIMIT" default:"512MB" yaml:"memory_limit"`
}
