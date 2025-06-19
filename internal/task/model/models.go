package model

import (
	"time"
)

// Status represents the current state of a task
type Status string

// Task statuses
const (
	StatusPending   Status = "PENDING"
	StatusScheduled Status = "SCHEDULED"
	StatusRunning   Status = "RUNNING"
	StatusCompleted Status = "COMPLETED"
	StatusFailed    Status = "FAILED"
)

// Task represents a computational task in the system
type Task struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Status      Status            `json:"status"`
	Input       map[string]string `json:"input,omitempty"`
	Output      map[string]string `json:"output,omitempty"`
	Resources   []Resource        `json:"resources,omitempty"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
	CompletedAt time.Time         `json:"completed_at,omitempty"`
}

// SubTask represents a part of a larger task
type SubTask struct {
	ID        string            `json:"id"`
	ParentID  string            `json:"parent_id"`
	Name      string            `json:"name"`
	Status    Status            `json:"status"`
	Input     map[string]string `json:"input,omitempty"`
	Output    map[string]string `json:"output,omitempty"`
	WorkerID  string            `json:"worker_id,omitempty"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
}

// Resource represents a computational resource
type Resource struct {
	Type  string `json:"type"`  // CPU, GPU, Memory, etc.
	Value int    `json:"value"` // Amount of resource
}
