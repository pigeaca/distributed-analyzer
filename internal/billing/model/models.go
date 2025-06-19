package model

import (
	"time"
)

// BillingRecord represents a billing entry for task execution
type BillingRecord struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	TaskID    string    `json:"task_id"`
	Amount    float64   `json:"amount"`
	Currency  string    `json:"currency"`
	Timestamp time.Time `json:"timestamp"`
}

// UserBalance represents a user's current balance
type UserBalance struct {
	UserID    string    `json:"user_id"`
	Balance   float64   `json:"balance"`
	Currency  string    `json:"currency"`
	UpdatedAt time.Time `json:"updated_at"`
}
