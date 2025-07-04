package grpc

import (
	"context"
	"fmt"
	"github.com/distributedmarketplace/internal/worker/model"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"

	pb "github.com/distributedmarketplace/pkg/proto/worker"
)

// WorkerManagerClient is a gRPC client for the WorkerManager service
type WorkerManagerClient struct {
	client pb.WorkerManagerServiceClient
	conn   *grpc.ClientConn
}

// NewWorkerManagerClient creates a new WorkerManagerClient
func NewWorkerManagerClient(address string) (*WorkerManagerClient, error) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to WorkerManager service: %w", err)
	}

	client := pb.NewWorkerManagerServiceClient(conn)
	return &WorkerManagerClient{
		client: client,
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
	pbCapabilities := make([]*pb.Capability, len(capabilities))
	for i, capability := range capabilities {
		pbCapabilities[i] = &pb.Capability{
			Name:  capability.Name,
			Value: capability.Value,
		}
	}

	// Convert model.Resource to pb.Resource
	pbResources := make([]*pb.Resource, len(resources))
	for i, resource := range resources {
		pbResources[i] = &pb.Resource{
			Type:  resource.Type,
			Value: int32(resource.Value),
		}
	}

	req := &pb.FindAvailableWorkersRequest{
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
	req := &pb.GetWorkerRequest{
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
