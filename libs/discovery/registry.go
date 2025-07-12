package discovery

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"sync"
	"time"
)

var (
	// ErrServiceNotFound is returned when a service is not found
	ErrServiceNotFound = errors.New("service not found")

	// ErrInvalidServiceID is returned when an invalid service ID is provided
	ErrInvalidServiceID = errors.New("invalid service ID")
)

// generateUUID generates a random UUID
func generateUUID() string {
	uuid := make([]byte, 16)
	_, err := rand.Read(uuid)
	if err != nil {
		// If we can't generate a random UUID, use a timestamp-based one
		uuid = []byte(time.Now().String())
	}
	return hex.EncodeToString(uuid)
}

// ServiceRegistry is responsible for managing service registrations
type ServiceRegistry struct {
	// services is a map of service ID to service instance
	services map[string]*ServiceInstance

	// mu is a mutex to protect concurrent access to the services map
	mu sync.RWMutex

	// heartbeatTimeout is the duration after which a service is considered inactive if no heartbeat is received
	heartbeatTimeout time.Duration

	// checkInterval is the interval at which the registry checks for inactive services
	checkInterval time.Duration

	// stopCh is a channel used to signal the registry to stop
	stopCh chan struct{}
}

// NewServiceRegistry creates a new service registry
func NewServiceRegistry(heartbeatTimeout, checkInterval time.Duration) *ServiceRegistry {
	registry := &ServiceRegistry{
		services:         make(map[string]*ServiceInstance),
		heartbeatTimeout: heartbeatTimeout,
		checkInterval:    checkInterval,
		stopCh:           make(chan struct{}),
	}

	// Start a goroutine to check for inactive services
	go registry.checkInactiveServices()

	return registry
}

// Register registers a service
func (r *ServiceRegistry) Register(service *ServiceInstance) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Generate a new ID if one is not provided
	if service.ID == "" {
		service.ID = generateUUID()
	}

	// Set the last heartbeat to now
	service.LastHeartbeat = time.Now()

	// Add the service to the registry
	r.services[service.ID] = service

	return service.ID, nil
}

// Unregister removes a service from the registry
func (r *ServiceRegistry) Unregister(serviceID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check if the service exists
	if _, exists := r.services[serviceID]; !exists {
		return ErrServiceNotFound
	}

	// Remove the service from the registry
	delete(r.services, serviceID)

	return nil
}

// Heartbeat updates the last heartbeat timestamp for a service
func (r *ServiceRegistry) Heartbeat(serviceID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check if the service exists
	service, exists := r.services[serviceID]
	if !exists {
		return ErrServiceNotFound
	}

	// Update the last heartbeat
	service.LastHeartbeat = time.Now()

	return nil
}

// GetService retrieves a service by ID
func (r *ServiceRegistry) GetService(serviceID string) (*ServiceInstance, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Check if the service exists
	service, exists := r.services[serviceID]
	if !exists {
		return nil, ErrServiceNotFound
	}

	return service, nil
}

// FindService finds services by type and/or name
func (r *ServiceRegistry) FindService(serviceType ServiceType, serviceName string) []*ServiceInstance {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Create a slice to hold matching services
	services := make([]*ServiceInstance, 0)

	// Add matching services to the slice
	for _, service := range r.services {
		// If serviceType is not UNKNOWN, filter by type
		if serviceType != ServiceTypeUnknown && service.Type != serviceType {
			continue
		}

		// If serviceName is not empty, filter by name
		if serviceName != "" && service.Name != serviceName {
			continue
		}

		services = append(services, service)
	}

	return services
}

// ListServices lists all services
func (r *ServiceRegistry) ListServices(serviceType ServiceType) []*ServiceInstance {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Create a slice to hold all services
	services := make([]*ServiceInstance, 0, len(r.services))

	// Add all services to the slice
	for _, service := range r.services {
		// If serviceType is not UNKNOWN, filter by type
		if serviceType != ServiceTypeUnknown && service.Type != serviceType {
			continue
		}

		services = append(services, service)
	}

	return services
}

// Stop stops the service registry
func (r *ServiceRegistry) Stop() {
	close(r.stopCh)
}

// checkInactiveServices periodically checks for services that haven't sent a heartbeat recently
func (r *ServiceRegistry) checkInactiveServices() {
	ticker := time.NewTicker(r.checkInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			r.removeInactiveServices()
		case <-r.stopCh:
			return
		}
	}
}

// removeInactiveServices removes services that haven't sent a heartbeat recently
func (r *ServiceRegistry) removeInactiveServices() {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Get the current time
	now := time.Now()

	// Check each service
	for id, service := range r.services {
		// Calculate the time since the last heartbeat
		timeSinceLastHeartbeat := now.Sub(service.LastHeartbeat)

		// If the service hasn't sent a heartbeat recently, remove it
		if timeSinceLastHeartbeat > r.heartbeatTimeout {
			delete(r.services, id)
		}
	}
}
