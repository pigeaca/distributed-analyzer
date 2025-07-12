package config

import configloader "distributed-analyzer/libs/config"

// Config holds the Scheduler Service configuration
type Config struct {
	configloader.ServerConfig `yaml:",inline"`

	// Worker management settings
	WorkerManagement WorkerManagementConfig `yaml:"worker_management"`

	// Graceful shutdown timeout
	ShutdownTimeout string `yaml:"shutdown_timeout" env:"SHUTDOWN_TIMEOUT" env-default:"30s"`
}

// WorkerManagementConfig holds worker management configuration
type WorkerManagementConfig struct {
	// Interval between worker heartbeats
	HeartbeatInterval string `yaml:"heartbeat_interval" env:"WORKER_HEARTBEAT_INTERVAL" env-default:"30s"`

	// Timeout after which a worker is considered dead
	Timeout string `yaml:"timeout" env:"WORKER_TIMEOUT" env-default:"60s"`
}
