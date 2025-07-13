package service

import (
	"context"
	"distributed-analyzer/services/worker-manager/internal/discovery"
	k8s2 "distributed-analyzer/services/worker-manager/internal/discovery/k8s"
	"errors"
	"log"
	"sync"
	"time"
)

const WorkerServiceName = "worker"

var (
	// ErrWorkerNotFound is returned when a worker is not found in the registry
	ErrWorkerNotFound = errors.New("worker not found")

	// ErrWorkerAlreadyRegistered is returned when a worker with the same ID is already registered
	ErrWorkerAlreadyRegistered = errors.New("worker already registered")
)

// WorkerRegistry is responsible for managing worker registrations
type WorkerRegistry struct {
	// workers is a map of worker ID to worker
	workers map[string]*Worker

	discovery discovery.ServiceDiscovery

	// mu is a mutex to protect concurrent access to the worker's map
	mu sync.RWMutex
}

// NewWorkerRegistry creates a new worker registry
func NewWorkerRegistry() *WorkerRegistry {
	return &WorkerRegistry{
		workers:   make(map[string]*Worker),
		discovery: k8s2.NewServiceDiscovery(),
	}
}

func (r *WorkerRegistry) StartRegistration(ctx context.Context) {
	go r.runAutoDiscovery(ctx, 10*time.Second)
}

func (r *WorkerRegistry) runAutoDiscovery(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("Stopping worker discovery")
			return
		case <-ticker.C:
			r.runDiscoveryOnce()
		}
	}
}

func (r *WorkerRegistry) runDiscoveryOnce() {
	services, err := r.discovery.DiscoverServices(WorkerServiceName)
	if err != nil {
		log.Printf("Discovery error: %v", err)
		return
	}
	r.RegisterServices(services)
}

// RegisterServices registers a new worker service and removes any stale workers
func (r *WorkerRegistry) RegisterServices(services []*discovery.Service) {
	r.mu.Lock()
	defer r.mu.Unlock()

	active := make(map[string]struct{})
	for _, s := range services {
		active[s.Addr] = struct{}{}

		if _, exists := r.workers[s.Addr]; !exists {
			r.workers[s.Addr] = NewWorker(s.Addr, s.Addr, []string{})
			log.Printf("✔ Registered new worker: %s", s.Addr)
		}
	}

	for id := range r.workers {
		if _, stillActive := active[id]; !stillActive {
			delete(r.workers, id)
			log.Printf("✖ Unregistered stale worker: %s", id)
		}
	}
}

// Register registers a new worker
func (r *WorkerRegistry) Register(id, address string, capabilities []string) (*Worker, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check if the worker is already registered
	if _, exists := r.workers[id]; exists {
		return nil, ErrWorkerAlreadyRegistered
	}

	// Create a new worker
	worker := NewWorker(id, address, capabilities)

	// Add the worker to the registry
	r.workers[id] = worker

	return worker, nil
}

// Unregister removes a worker from the registry
func (r *WorkerRegistry) Unregister(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check if the worker exists
	if _, exists := r.workers[id]; !exists {
		return ErrWorkerNotFound
	}

	// Remove the worker from the registry
	delete(r.workers, id)

	return nil
}

// Get retrieves a worker by ID
func (r *WorkerRegistry) Get(id string) (*Worker, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Check if the worker exists
	worker, exists := r.workers[id]
	if !exists {
		return nil, ErrWorkerNotFound
	}

	return worker, nil
}

// GetAll retrieves all workers
func (r *WorkerRegistry) GetAll() []*Worker {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Create a slice to hold all workers
	workers := make([]*Worker, 0, len(r.workers))

	// Add all workers to the slice
	for _, worker := range r.workers {
		workers = append(workers, worker)
	}

	return workers
}

// GetByStatus retrieves all workers with the specified status
func (r *WorkerRegistry) GetByStatus(status WorkerStatus) []*Worker {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Create a slice to hold matching workers
	workers := make([]*Worker, 0)

	// Add matching workers to the slice
	for _, worker := range r.workers {
		if worker.Status == status {
			workers = append(workers, worker)
		}
	}

	return workers
}

// GetByCapability retrieves all workers with the specified capability
func (r *WorkerRegistry) GetByCapability(capability string) []*Worker {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Create a slice to hold matching workers
	workers := make([]*Worker, 0)

	// Add matching workers to the slice
	for _, worker := range r.workers {
		if worker.HasCapability(capability) {
			workers = append(workers, worker)
		}
	}

	return workers
}

// Count returns the total number of registered workers
func (r *WorkerRegistry) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return len(r.workers)
}

// CountByStatus returns the number of workers with the specified status
func (r *WorkerRegistry) CountByStatus(status WorkerStatus) int {
	r.mu.RLock()
	defer r.mu.RUnlock()

	count := 0
	for _, worker := range r.workers {
		if worker.Status == status {
			count++
		}
	}

	return count
}
