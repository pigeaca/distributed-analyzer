package grpc

import (
	"context"
	"distributed-analyzer/libs/proto/worker"
	"distributed-analyzer/services/worker-manager/internal/service"
	"fmt"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
)

// WorkerManagerServer implements the WorkerManagerServiceServer interface
type WorkerManagerServer struct {
	worker.UnimplementedWorkerManagerServiceServer

	// workerManager is the worker manager service
	workerManager *service.WorkerManager
}

// NewWorkerManagerServer creates a new worker manager server
func NewWorkerManagerServer(workerManager *service.WorkerManager) *WorkerManagerServer {
	return &WorkerManagerServer{
		workerManager: workerManager,
	}
}

// RegisterWorker registers a new worker
func (s *WorkerManagerServer) RegisterWorker(ctx context.Context, req *worker.RegisterWorkerRequest) (*worker.WorkerResponse, error) {
	log.Printf("Registering worker: %s", req.GetName())

	// Extract capabilities from the request
	capabilities := make([]string, 0, len(req.GetCapabilities()))
	for _, capability := range req.GetCapabilities() {
		capabilities = append(capabilities, capability.GetName())
	}

	// Register the worker
	w, err := s.workerManager.RegisterWorker(req.GetName(), req.GetName(), capabilities)
	if err != nil {
		return nil, fmt.Errorf("failed to register worker: %w", err)
	}

	// Convert the worker to a proto worker
	protoWorker := convertWorkerToProto(w)

	return &worker.WorkerResponse{
		Worker: protoWorker,
	}, nil
}

// GetWorker retrieves a worker by ID
func (s *WorkerManagerServer) GetWorker(ctx context.Context, req *worker.GetWorkerRequest) (*worker.WorkerResponse, error) {
	log.Printf("Getting worker: %s", req.GetId())

	// Get the worker
	w, err := s.workerManager.GetWorker(req.GetId())
	if err != nil {
		return nil, fmt.Errorf("failed to get worker: %w", err)
	}

	// Convert the worker to a proto worker
	protoWorker := convertWorkerToProto(w)

	return &worker.WorkerResponse{
		Worker: protoWorker,
	}, nil
}

// UpdateWorkerStatus updates the status of a worker
func (s *WorkerManagerServer) UpdateWorkerStatus(ctx context.Context, req *worker.UpdateWorkerStatusRequest) (*worker.UpdateWorkerStatusResponse, error) {
	log.Printf("Updating worker status: %s -> %s", req.GetId(), req.GetStatus())

	// Convert the status string to a WorkerStatus
	status := service.WorkerStatus(req.GetStatus())

	// Update the worker status
	err := s.workerManager.UpdateWorkerStatus(req.GetId(), status)
	if err != nil {
		return nil, fmt.Errorf("failed to update worker status: %w", err)
	}

	return &worker.UpdateWorkerStatusResponse{
		Success: true,
	}, nil
}

// ListWorkers retrieves all workers
func (s *WorkerManagerServer) ListWorkers(ctx context.Context, req *worker.ListWorkersRequest) (*worker.ListWorkersResponse, error) {
	log.Printf("Listing workers")

	// Get all workers
	workers := s.workerManager.GetAllWorkers()

	// Convert the workers to proto workers
	protoWorkers := make([]*worker.Worker, 0, len(workers))
	for _, w := range workers {
		protoWorkers = append(protoWorkers, convertWorkerToProto(w))
	}

	return &worker.ListWorkersResponse{
		Workers: protoWorkers,
	}, nil
}

// FindAvailableWorkers finds workers with the specified capabilities
func (s *WorkerManagerServer) FindAvailableWorkers(ctx context.Context, req *worker.FindAvailableWorkersRequest) (*worker.ListWorkersResponse, error) {
	log.Printf("Finding available workers")

	// Extract capabilities from the request
	capabilities := make([]string, 0, len(req.GetCapabilities()))
	for _, capability := range req.GetCapabilities() {
		capabilities = append(capabilities, capability.GetName())
	}

	// Get active workers
	activeWorkers := s.workerManager.GetWorkersByStatus(service.WorkerStatusActive)

	// Filter workers by capabilities
	availableWorkers := make([]*service.Worker, 0)
	for _, w := range activeWorkers {
		hasAllCapabilities := true
		for _, capability := range capabilities {
			if !w.HasCapability(capability) {
				hasAllCapabilities = false
				break
			}
		}
		if hasAllCapabilities {
			availableWorkers = append(availableWorkers, w)
		}
	}

	// Convert the workers to proto workers
	protoWorkers := make([]*worker.Worker, 0, len(availableWorkers))
	for _, w := range availableWorkers {
		protoWorkers = append(protoWorkers, convertWorkerToProto(w))
	}

	return &worker.ListWorkersResponse{
		Workers: protoWorkers,
	}, nil
}

// convertWorkerToProto converts a service.Worker to a worker.Worker
func convertWorkerToProto(w *service.Worker) *worker.Worker {
	// Convert capabilities to proto capabilities
	capabilities := make([]*worker.Capability, 0, len(w.Capabilities))
	for _, capability := range w.Capabilities {
		capabilities = append(capabilities, &worker.Capability{
			Name:  capability,
			Value: "",
		})
	}

	// Create a proto worker
	return &worker.Worker{
		Id:           w.ID,
		Name:         w.Address,
		Status:       string(w.Status),
		Capabilities: capabilities,
		Resources:    nil, // We don't track resources in our implementation
		LastSeen:     timestamppb.New(w.LastHeartbeat),
	}
}
