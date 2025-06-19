package model

import (
	"time"
)

// User represents a system user
type User struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

// UserRole represents a user's role in the system
type UserRole string

// User roles
const (
	RoleAdmin  UserRole = "ADMIN"
	RoleUser   UserRole = "USER"
	RoleWorker UserRole = "WORKER"
)

// UserPermission represents a user's permission in the system
type UserPermission struct {
	UserID     string    `json:"user_id"`
	Resource   string    `json:"resource"`
	Permission string    `json:"permission"`
	CreatedAt  time.Time `json:"created_at"`
}
