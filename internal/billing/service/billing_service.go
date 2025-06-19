package service

import (
	"context"
	"github.com/distributedmarketplace/internal/billing/model"
)

// BillingService defines the interface for billing operations
type BillingService interface {
	// ChargeTask charges a user for a completed task
	ChargeTask(ctx context.Context, taskID string) (*model.BillingRecord, error)

	// GetUserBalance retrieves a user's current balance
	GetUserBalance(ctx context.Context, userID string) (float64, error)

	// AddUserBalance adds to a user's balance
	AddUserBalance(ctx context.Context, userID string, amount float64) error

	// GetBillingHistory retrieves a user's billing history
	GetBillingHistory(ctx context.Context, userID string) ([]*model.BillingRecord, error)

	// CreateBillingRecord creates a new billing record
	CreateBillingRecord(ctx context.Context, record *model.BillingRecord) (*model.BillingRecord, error)
}
