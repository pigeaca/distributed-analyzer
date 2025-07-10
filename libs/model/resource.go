package model

// Resource represents a computational resource
type Resource struct {
	Type  string `json:"type"`  // CPU, GPU, Memory, etc.
	Value int    `json:"value"` // Amount of resource
}
