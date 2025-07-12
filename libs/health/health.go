// Package health provides a health check service for all services
package health

import (
	"context"
	"sync"
)

// Status represents the health status of a service
type Status string

const (
	// StatusUnknown indicates that the health status is unknown
	StatusUnknown Status = "UNKNOWN"

	// StatusServing indicates that the service is serving requests
	StatusServing Status = "SERVING"

	// StatusNotServing indicates that the service is not serving requests
	StatusNotServing Status = "NOT_SERVING"

	// StatusServiceUnknown indicates that the service is unknown
	StatusServiceUnknown Status = "SERVICE_UNKNOWN"
)

// Checker is an interface for checking the health of a service
type Checker interface {
	// Check checks the health of a service
	Check(ctx context.Context, service string) (Status, error)
}

// Server is a health check server
type Server struct {
	// mu protects the statuses map
	mu sync.RWMutex

	// statuses is a map of service name to health status
	statuses map[string]Status

	// checker is an optional health checker
	checker Checker
}

// NewServer creates a new health check server
func NewServer() *Server {
	return &Server{
		statuses: make(map[string]Status),
	}
}

// SetChecker sets the health checker
func (s *Server) SetChecker(checker Checker) {
	s.checker = checker
}

// SetStatus sets the health status of a service
func (s *Server) SetStatus(service string, status Status) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.statuses[service] = status
}

// Check checks the health of a service
func (s *Server) Check(ctx context.Context, service string) (Status, error) {
	// If a checker is set, use it to check the health of the service
	if s.checker != nil {
		return s.checker.Check(ctx, service)
	}

	// Otherwise, use the stored status
	s.mu.RLock()
	defer s.mu.RUnlock()

	status, ok := s.statuses[service]
	if !ok {
		return StatusServiceUnknown, nil
	}

	return status, nil
}

// CheckAll checks the health of all services
func (s *Server) CheckAll(ctx context.Context) map[string]Status {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Create a copy of the statuses map
	statuses := make(map[string]Status, len(s.statuses))
	for service, status := range s.statuses {
		statuses[service] = status
	}

	return statuses
}

// RegisterService registers a service with the health check server
func (s *Server) RegisterService(service string) {
	s.SetStatus(service, StatusUnknown)
}

// UnregisterService unregisters a service from the health check server
func (s *Server) UnregisterService(service string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.statuses, service)
}

// IsServing checks if a service is serving requests
func (s *Server) IsServing(ctx context.Context, service string) bool {
	status, _ := s.Check(ctx, service)
	return status == StatusServing
}

// IsNotServing checks if a service is not serving requests
func (s *Server) IsNotServing(ctx context.Context, service string) bool {
	status, _ := s.Check(ctx, service)
	return status == StatusNotServing
}

// IsUnknown checks if a service's health status is unknown
func (s *Server) IsUnknown(ctx context.Context, service string) bool {
	status, _ := s.Check(ctx, service)
	return status == StatusUnknown
}

// IsServiceUnknown checks if a service is unknown
func (s *Server) IsServiceUnknown(ctx context.Context, service string) bool {
	status, _ := s.Check(ctx, service)
	return status == StatusServiceUnknown
}
