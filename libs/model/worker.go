package model

import (
	"time"
)

// Capability represents a worker capability
type Capability struct {
	Name  string `json:"name"`  // e.g., "resnet50", "gpt-3"
	Value string `json:"value"` // Version or other details
}

// Worker represents a node that can execute tasks
type Worker struct {
	ID           string       `json:"id"`
	Name         string       `json:"name"`
	Status       string       `json:"status"` // Online, Offline, Busy
	Capabilities []Capability `json:"capabilities"`
	Resources    []Resource   `json:"resources"`
	LastSeen     time.Time    `json:"last_seen"`
}
