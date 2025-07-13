package service

import (
	"sync"
	"time"
)

// WorkerStatus represents the status of a worker
type WorkerStatus string

const (
	// WorkerStatusActive indicates that the worker is active and available for tasks
	WorkerStatusActive WorkerStatus = "active"

	// WorkerStatusInactive indicates that the worker is inactive (hasn't sent a heartbeat recently)
	WorkerStatusInactive WorkerStatus = "inactive"

	// WorkerStatusBusy indicates that the worker is currently processing a task
	WorkerStatusBusy WorkerStatus = "busy"

	// WorkerStatusError indicates that the worker has encountered an error
	WorkerStatusError WorkerStatus = "error"
)

// Worker represents a worker node in the distributed system
type Worker struct {
	// ID is the unique identifier of the worker
	ID string

	// Address is the network address of the worker
	Address string

	// Status is the current status of the worker
	Status WorkerStatus

	// Capabilities is a list of capabilities that the worker supports
	Capabilities []string

	// LastHeartbeat is the timestamp of the last heartbeat received from the worker
	LastHeartbeat time.Time

	// RegisteredAt is the timestamp when the worker was registered
	RegisteredAt time.Time

	// CurrentLoad is the current load of the worker (0-100)
	CurrentLoad int

	// Error is the last error reported by the worker
	Error string

	// mu is a mutex to protect concurrent access to the worker
	mu sync.RWMutex
}

func NewWorker(id, address string, capabilities []string) *Worker {
	return &Worker{
		ID:            id,
		Address:       address,
		Capabilities:  capabilities,
		Status:        WorkerStatusActive,
		RegisteredAt:  time.Now(),
		CurrentLoad:   0,
		LastHeartbeat: time.Now(),
	}
}

// UpdateStatus updates the status of the worker
func (w *Worker) UpdateStatus(status WorkerStatus) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.Status = status
}

// UpdateHeartbeat updates the last heartbeat timestamp of the worker
func (w *Worker) UpdateHeartbeat() {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.LastHeartbeat = time.Now()

	// If the worker was inactive, set it to active
	if w.Status == WorkerStatusInactive {
		w.Status = WorkerStatusActive
	}
}

// UpdateLoad updates the current load of the worker
func (w *Worker) UpdateLoad(load int) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.CurrentLoad = load
}

// SetError sets the error message and updates the status to error
func (w *Worker) SetError(err string) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.Error = err
	w.Status = WorkerStatusError
}

// ClearError clears the error message and updates the status to active
func (w *Worker) ClearError() {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.Error = ""
	w.Status = WorkerStatusActive
}

// IsActive returns true if the worker is active
func (w *Worker) IsActive() bool {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.Status == WorkerStatusActive
}

// IsBusy returns true if the worker is busy
func (w *Worker) IsBusy() bool {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.Status == WorkerStatusBusy
}

// HasCapability returns true if the worker has the specified capability
func (w *Worker) HasCapability(capability string) bool {
	w.mu.RLock()
	defer w.mu.RUnlock()
	for _, cap := range w.Capabilities {
		if cap == capability {
			return true
		}
	}
	return false
}
