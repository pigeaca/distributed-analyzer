package discovery

import (
	"time"
)

// ServiceType represents the type of service
type ServiceType int

const (
	// ServiceTypeUnknown represents an unknown service type
	ServiceTypeUnknown ServiceType = iota

	// ServiceTypeAPI represents an API service
	ServiceTypeAPI

	// ServiceTypeWorker represents a worker service
	ServiceTypeWorker

	// ServiceTypeStorage represents a storage service
	ServiceTypeStorage

	// ServiceTypeTask represents a task service
	ServiceTypeTask

	// ServiceTypeResult represents a result service
	ServiceTypeResult

	// ServiceTypeScheduler represents a scheduler service
	ServiceTypeScheduler

	// ServiceTypeWorkerManager represents a worker manager service
	ServiceTypeWorkerManager
)

// String returns the string representation of the service type
func (t ServiceType) String() string {
	switch t {
	case ServiceTypeAPI:
		return "api"
	case ServiceTypeWorker:
		return "worker"
	case ServiceTypeStorage:
		return "storage"
	case ServiceTypeTask:
		return "task"
	case ServiceTypeResult:
		return "result"
	case ServiceTypeScheduler:
		return "scheduler"
	case ServiceTypeWorkerManager:
		return "worker-manager"
	default:
		return "unknown"
	}
}

// ServiceInstance represents a registered service instance
type ServiceInstance struct {
	// ID is the unique identifier for the service instance
	ID string

	// Name is the name of the service
	Name string

	// Type is the type of service
	Type ServiceType

	// Host is the hostname or IP address of the service
	Host string

	// Port is the port number of the service
	Port int

	// GrpcPort is the gRPC port number of the service
	GrpcPort int

	// Metadata is additional metadata about the service
	Metadata map[string]string

	// LastHeartbeat is the timestamp of the last heartbeat received from the service
	LastHeartbeat time.Time

	// RegisteredAt is the timestamp when the service was registered
	RegisteredAt time.Time
}

// ServiceRegistration represents a service registration request
type ServiceRegistration struct {
	// Name is the name of the service
	Name string

	// Type is the type of service
	Type ServiceType

	// Host is the hostname or IP address of the service
	Host string

	// Port is the port number of the service
	Port int

	// GrpcPort is the gRPC port number of the service
	GrpcPort int

	// Metadata is additional metadata about the service
	Metadata map[string]string
}

// ServiceHeartbeat represents a service heartbeat request
type ServiceHeartbeat struct {
	// ID is the unique identifier for the service instance
	ID string

	// Metadata is additional metadata about the service that can be updated
	Metadata map[string]string
}

// ServiceQuery represents a service discovery query
type ServiceQuery struct {
	// Type is the type of service to find
	Type ServiceType

	// Name is the name of the service to find
	Name string
}
