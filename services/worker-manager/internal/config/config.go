package config

import (
	"distributed-analyzer/libs/config"
)

// Config holds the API gateway configuration
type Config struct {
	// Server settings
	configloader.ServerConfig `yaml:",inline"`

	// Services hold the configuration for all-services
	// WorkerManagement settings
	WorkerManagement WorkerManagementConfig `yaml:"worker_management"`

	// ShutdownTimeout specifies how long to wait for graceful shutdown
	// Can be set via SHUTDOWN_TIMEOUT environment variable
	// Default: 30s
	ShutdownTimeout string `envconfig:"SHUTDOWN_TIMEOUT" default:"30s" yaml:"shutdown_timeout"`
}

// WorkerManagementConfig holds worker management configuration
type WorkerManagementConfig struct {
	// HeartbeatInterval specifies the worker heartbeat interval
	// Can be set via WORKER_HEARTBEAT_INTERVAL environment variable
	// Default: 30s
	HeartbeatInterval string `envconfig:"WORKER_HEARTBEAT_INTERVAL" default:"30s" yaml:"heartbeat_interval"`

	// Timeout specifying the worker timeout
	// Can be set via WORKER_TIMEOUT environment variable
	// Default: 60s
	Timeout string `envconfig:"WORKER_TIMEOUT" default:"60s" yaml:"timeout"`
}
