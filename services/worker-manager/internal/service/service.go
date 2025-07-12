package service

import (
	"context"
	"distributed-analyzer/services/worker-manager/internal/config"
	"fmt"
	"log"
	"time"
)

// WorkerManager is the main service for managing workers
type WorkerManager struct {
	// registry is the worker registry
	registry *WorkerRegistry

	// heartbeatTracker is the heartbeat tracker
	heartbeatTracker *HeartbeatTracker

	// config is the service configuration
	config *config.Config
}

// NewWorkerManager creates a new worker manager
func NewWorkerManager(cfg *config.Config) (*WorkerManager, error) {
	// Parse heartbeat interval
	heartbeatInterval, err := time.ParseDuration(cfg.WorkerManagement.HeartbeatInterval)
	if err != nil {
		return nil, fmt.Errorf("invalid heartbeat interval: %w", err)
	}

	// Parse timeout
	timeout, err := time.ParseDuration(cfg.WorkerManagement.Timeout)
	if err != nil {
		return nil, fmt.Errorf("invalid timeout: %w", err)
	}

	// Create a worker registry
	registry := NewWorkerRegistry()

	// Create a heartbeat tracker
	heartbeatTracker := NewHeartbeatTracker(registry, heartbeatInterval, timeout)

	return &WorkerManager{
		registry:         registry,
		heartbeatTracker: heartbeatTracker,
		config:           cfg,
	}, nil
}

// Start starts the worker manager
func (m *WorkerManager) Start(ctx context.Context) error {
	log.Println("Starting worker manager")

	// Start the heartbeat tracker
	if err := m.heartbeatTracker.Start(ctx); err != nil {
		return fmt.Errorf("failed to start heartbeat tracker: %w", err)
	}

	return nil
}

// Stop stops the worker manager
func (m *WorkerManager) Stop(ctx context.Context) error {
	log.Println("Stopping worker manager")

	// Stop the heartbeat tracker
	if err := m.heartbeatTracker.Stop(ctx); err != nil {
		return fmt.Errorf("failed to stop heartbeat tracker: %w", err)
	}

	return nil
}

// RegisterWorker registers a new worker
func (m *WorkerManager) RegisterWorker(id, address string, capabilities []string) (*Worker, error) {
	log.Printf("Registering worker %s at %s with capabilities %v", id, address, capabilities)

	// Register the worker
	worker, err := m.registry.Register(id, address, capabilities)
	if err != nil {
		return nil, fmt.Errorf("failed to register worker: %w", err)
	}

	return worker, nil
}

// UnregisterWorker unregisters a worker
func (m *WorkerManager) UnregisterWorker(id string) error {
	log.Printf("Unregistering worker %s", id)

	// Unregister the worker
	if err := m.registry.Unregister(id); err != nil {
		return fmt.Errorf("failed to unregister worker: %w", err)
	}

	return nil
}

// RecordHeartbeat records a heartbeat for a worker
func (m *WorkerManager) RecordHeartbeat(id string) error {
	// Record the heartbeat
	if err := m.heartbeatTracker.RecordHeartbeat(id); err != nil {
		return fmt.Errorf("failed to record heartbeat: %w", err)
	}

	return nil
}

// GetWorker retrieves a worker by ID
func (m *WorkerManager) GetWorker(id string) (*Worker, error) {
	// Get the worker
	worker, err := m.registry.Get(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get worker: %w", err)
	}

	return worker, nil
}

// GetAllWorkers retrieves all workers
func (m *WorkerManager) GetAllWorkers() []*Worker {
	return m.registry.GetAll()
}

// GetWorkersByStatus retrieves all workers with the specified status
func (m *WorkerManager) GetWorkersByStatus(status WorkerStatus) []*Worker {
	return m.registry.GetByStatus(status)
}

// GetWorkersByCapability retrieves all workers with the specified capability
func (m *WorkerManager) GetWorkersByCapability(capability string) []*Worker {
	return m.registry.GetByCapability(capability)
}

// UpdateWorkerStatus updates the status of a worker
func (m *WorkerManager) UpdateWorkerStatus(id string, status WorkerStatus) error {
	// Get the worker
	worker, err := m.registry.Get(id)
	if err != nil {
		return fmt.Errorf("failed to get worker: %w", err)
	}

	// Update the worker's status
	worker.UpdateStatus(status)

	return nil
}

// UpdateWorkerLoad updates the load of a worker
func (m *WorkerManager) UpdateWorkerLoad(id string, load int) error {
	// Get the worker
	worker, err := m.registry.Get(id)
	if err != nil {
		return fmt.Errorf("failed to get worker: %w", err)
	}

	// Update the worker's load
	worker.UpdateLoad(load)

	return nil
}

// SetWorkerError sets an error for a worker
func (m *WorkerManager) SetWorkerError(id string, errorMsg string) error {
	// Get the worker
	worker, err := m.registry.Get(id)
	if err != nil {
		return fmt.Errorf("failed to get worker: %w", err)
	}

	// Set the worker's error
	worker.SetError(errorMsg)

	return nil
}

// ClearWorkerError clears the error for a worker
func (m *WorkerManager) ClearWorkerError(id string) error {
	// Get the worker
	worker, err := m.registry.Get(id)
	if err != nil {
		return fmt.Errorf("failed to get worker: %w", err)
	}

	// Clear the worker's error
	worker.ClearError()

	return nil
}
