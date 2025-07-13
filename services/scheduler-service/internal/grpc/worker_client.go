package grpc

import (
	"context"
	"distributed-analyzer/libs/model"
	"distributed-analyzer/libs/network/client"
	pbW "distributed-analyzer/libs/proto/worker"
	"fmt"
	"google.golang.org/grpc"
	"log"
)

// WorkerManagerClient is a gRPC client for the WorkerManager service
type WorkerManagerClient struct {
	client pbW.WorkerManagerServiceClient
	conn   *grpc.ClientConn
}

// NewWorkerManagerClient creates a new WorkerManagerClient
func NewWorkerManagerClient(address string) (*WorkerManagerClient, error) {
	conn, err := client.NewGrpcResilientClient(nil, address)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to WorkerManager service: %w", err)
	}

	managerServiceClient := pbW.NewWorkerManagerServiceClient(conn)
	return &WorkerManagerClient{
		client: managerServiceClient,
		conn:   conn,
	}, nil
}

// Close closes the gRPC connection
func (c *WorkerManagerClient) Close() error {
	return c.conn.Close()
}

// FindAvailableWorkers finds workers that can handle a specific task
func (c *WorkerManagerClient) FindAvailableWorkers(ctx context.Context, capabilities []model.Capability, resources []model.Resource) ([]*model.Worker, error) {
	// Convert model.Capability to pb.Capability
	pbCapabilities := make([]*pbW.Capability, len(capabilities))
	for i, capability := range capabilities {
		pbCapabilities[i] = &pbW.Capability{
			Name:  capability.Name,
			Value: capability.Value,
		}
	}

	// Convert model.Resource to pb.Resource
	pbResources := make([]*pbW.Resource, len(resources))
	for i, resource := range resources {
		pbResources[i] = &pbW.Resource{
			Type:  resource.Type,
			Value: int32(resource.Value),
		}
	}

	req := &pbW.FindAvailableWorkersRequest{
		Capabilities: pbCapabilities,
		Resources:    pbResources,
	}

	resp, err := c.client.FindAvailableWorkers(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to find available workers: %w", err)
	}

	// Convert pb.Worker to model.Worker
	workers := make([]*model.Worker, len(resp.Workers))
	for i, pbWorker := range resp.Workers {
		// Convert pb.Capability to model.Capability
		capabilities := make([]model.Capability, len(pbWorker.Capabilities))
		for j, pbCapability := range pbWorker.Capabilities {
			capabilities[j] = model.Capability{
				Name:  pbCapability.Name,
				Value: pbCapability.Value,
			}
		}

		// Convert pb.Resource to model.Resource
		resources := make([]model.Resource, len(pbWorker.Resources))
		for j, pbResource := range pbWorker.Resources {
			resources[j] = model.Resource{
				Type:  pbResource.Type,
				Value: int(pbResource.Value),
			}
		}

		workers[i] = &model.Worker{
			ID:           pbWorker.Id,
			Name:         pbWorker.Name,
			Status:       pbWorker.Status,
			Capabilities: capabilities,
			Resources:    resources,
		}
	}

	log.Printf("Found %d available workers", len(workers))
	return workers, nil
}

// GetWorker retrieves a worker by its ID
func (c *WorkerManagerClient) GetWorker(ctx context.Context, id string) (*model.Worker, error) {
	req := &pbW.GetWorkerRequest{
		Id: id,
	}

	resp, err := c.client.GetWorker(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get worker: %w", err)
	}

	// Convert pb.Worker to model.Worker
	pbWorker := resp.Worker

	// Convert pb.Capability to model.Capability
	capabilities := make([]model.Capability, len(pbWorker.Capabilities))
	for j, pbCapability := range pbWorker.Capabilities {
		capabilities[j] = model.Capability{
			Name:  pbCapability.Name,
			Value: pbCapability.Value,
		}
	}

	// Convert pb.Resource to model.Resource
	resources := make([]model.Resource, len(pbWorker.Resources))
	for j, pbResource := range pbWorker.Resources {
		resources[j] = model.Resource{
			Type:  pbResource.Type,
			Value: int(pbResource.Value),
		}
	}

	worker := &model.Worker{
		ID:           pbWorker.Id,
		Name:         pbWorker.Name,
		Status:       pbWorker.Status,
		Capabilities: capabilities,
		Resources:    resources,
	}

	return worker, nil
}
