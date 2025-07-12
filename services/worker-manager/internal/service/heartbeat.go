package service

import (
	"context"
	"log"
	"time"
)

// HeartbeatTracker is responsible for tracking worker heartbeats and updating worker statuses
type HeartbeatTracker struct {
	// registry is the worker registry
	registry *WorkerRegistry

	// heartbeatInterval is the interval at which workers should send heartbeats
	heartbeatInterval time.Duration

	// timeout is the duration after which a worker is considered inactive if no heartbeat is received
	timeout time.Duration

	// checkInterval is the interval at which the tracker checks for inactive workers
	checkInterval time.Duration

	// stopCh is a channel used to signal the tracker to stop
	stopCh chan struct{}
}

// NewHeartbeatTracker creates a new heartbeat tracker
func NewHeartbeatTracker(registry *WorkerRegistry, heartbeatInterval, timeout time.Duration) *HeartbeatTracker {
	// Set check interval to half the heartbeat interval
	checkInterval := heartbeatInterval / 2
	if checkInterval < time.Second {
		checkInterval = time.Second
	}

	return &HeartbeatTracker{
		registry:          registry,
		heartbeatInterval: heartbeatInterval,
		timeout:           timeout,
		checkInterval:     checkInterval,
		stopCh:            make(chan struct{}),
	}
}

// Start starts the heartbeat tracker
func (t *HeartbeatTracker) Start(ctx context.Context) error {
	log.Println("Starting heartbeat tracker")

	// Start a goroutine to check for inactive workers
	go t.checkInactiveWorkers(ctx)

	return nil
}

// Stop stops the heartbeat tracker
func (t *HeartbeatTracker) Stop(ctx context.Context) error {
	log.Println("Stopping heartbeat tracker")

	// Signal the tracker to stop
	close(t.stopCh)

	return nil
}

// RecordHeartbeat records a heartbeat for the specified worker
func (t *HeartbeatTracker) RecordHeartbeat(workerID string) error {
	// Get the worker from the registry
	worker, err := t.registry.Get(workerID)
	if err != nil {
		return err
	}

	// Update the worker's heartbeat
	worker.UpdateHeartbeat()

	return nil
}

// checkInactiveWorkers periodically checks for workers that haven't sent a heartbeat recently
func (t *HeartbeatTracker) checkInactiveWorkers(ctx context.Context) {
	ticker := time.NewTicker(t.checkInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			t.markInactiveWorkers()
		case <-t.stopCh:
			return
		case <-ctx.Done():
			return
		}
	}
}

// markInactiveWorkers marks workers as inactive if they haven't sent a heartbeat recently
func (t *HeartbeatTracker) markInactiveWorkers() {
	// Get all workers
	workers := t.registry.GetAll()

	// Check each worker
	now := time.Now()
	for _, worker := range workers {
		// Skip workers that are already inactive
		if worker.Status == WorkerStatusInactive {
			continue
		}

		// Check if the worker has timed out
		if now.Sub(worker.LastHeartbeat) > t.timeout {
			log.Printf("Worker %s has timed out (last heartbeat: %s)", worker.ID, worker.LastHeartbeat)
			worker.UpdateStatus(WorkerStatusInactive)
		}
	}
}
