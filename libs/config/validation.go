// Package configloader Package config loader provides a centralized configuration management system for all services.
package configloader

import (
	"fmt"
)

// Validator is an interface that configuration structs can implement to provide custom validation.
type Validator interface {
	Validate() error
}

// DefaultSetter is an interface that configuration structs can implement to set default values.
type DefaultSetter interface {
	SetDefaults()
}

// ValidateConfig validates the provided configuration if it implements the Validator interface.
// It returns an error if validation fails, or nil if validation succeeds, or the config doesn't
// implement the Validator interface.
func ValidateConfig(config interface{}) error {
	if validator, ok := config.(Validator); ok {
		if err := validator.Validate(); err != nil {
			return fmt.Errorf("config validation failed: %w", err)
		}
	}
	return nil
}
