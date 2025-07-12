// Package health provides a health check service for all services
package health

import (
	"encoding/json"
	"net/http"
	"strings"
)

// HealthResponse represents the response from a health check
type HealthResponse struct {
	// Status is the overall health status of the service
	Status Status `json:"status"`

	// Services is a map of service name to health status
	Services map[string]Status `json:"services,omitempty"`
}

// Handler is an HTTP handler for health checks
type Handler struct {
	// server is the health check server
	server *Server

	// serviceName is the name of the service
	serviceName string
}

// NewHandler creates a new health check handler
func NewHandler(server *Server, serviceName string) *Handler {
	return &Handler{
		server:      server,
		serviceName: serviceName,
	}
}

// ServeHTTP implements the http.Handler interface
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Set the content type
	w.Header().Set("Content-Type", "application/json")

	// Parse the path to determine which service to check
	path := strings.TrimPrefix(r.URL.Path, "/health")
	path = strings.TrimPrefix(path, "/")

	// If no service is specified, check the current service
	if path == "" {
		h.handleServiceCheck(w, r, h.serviceName)
		return
	}

	// If the path is "all", check all services
	if path == "all" {
		h.handleAllServicesCheck(w, r)
		return
	}

	// Otherwise, check the specified service
	h.handleServiceCheck(w, r, path)
}

// handleServiceCheck handles a health check for a specific service
func (h *Handler) handleServiceCheck(w http.ResponseWriter, r *http.Request, service string) {
	// Check the service
	status, err := h.server.Check(r.Context(), service)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create the response
	response := HealthResponse{
		Status: status,
	}

	// Set the status code based on the health status
	statusCode := http.StatusOK
	if status == StatusNotServing {
		statusCode = http.StatusServiceUnavailable
	} else if status == StatusUnknown || status == StatusServiceUnknown {
		statusCode = http.StatusNotFound
	}

	// Write the response
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

// handleAllServicesCheck handles a health check for all services
func (h *Handler) handleAllServicesCheck(w http.ResponseWriter, r *http.Request) {
	// Check all services
	statuses := h.server.CheckAll(r.Context())

	// Determine the overall status
	overallStatus := StatusServing
	for _, status := range statuses {
		if status == StatusNotServing {
			overallStatus = StatusNotServing
			break
		} else if status == StatusUnknown {
			overallStatus = StatusUnknown
		}
	}

	// Create the response
	response := HealthResponse{
		Status:   overallStatus,
		Services: statuses,
	}

	// Set the status code based on the overall health status
	statusCode := http.StatusOK
	if overallStatus == StatusNotServing {
		statusCode = http.StatusServiceUnavailable
	} else if overallStatus == StatusUnknown {
		statusCode = http.StatusNotFound
	}

	// Write the response
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

// RegisterHandler registers the health check handler with the given mux
func RegisterHandler(mux *http.ServeMux, server *Server, serviceName string) {
	handler := NewHandler(server, serviceName)
	mux.Handle("/health", handler)
	mux.Handle("/health/", handler)
}
