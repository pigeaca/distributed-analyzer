package service

import (
	"context"
	"github.com/distributedmarketplace/internal/user/model"
)

// UserService defines the interface for user operations
type UserService interface {
	// CreateUser creates a new user in the system
	CreateUser(ctx context.Context, user *model.User) (*model.User, error)

	// GetUser retrieves a user by their ID
	GetUser(ctx context.Context, id string) (*model.User, error)

	// UpdateUser updates an existing user
	UpdateUser(ctx context.Context, user *model.User) (*model.User, error)

	// DeleteUser removes a user from the system
	DeleteUser(ctx context.Context, id string) error

	// ListUsers retrieves all users with optional filtering
	ListUsers(ctx context.Context) ([]*model.User, error)

	// AssignRole assigns a role to a user
	AssignRole(ctx context.Context, userID string, role model.UserRole) error

	// GrantPermission grants a permission to a user
	GrantPermission(ctx context.Context, userID string, resource string, permission string) error
}
