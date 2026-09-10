// Package checkoutpostgres implements User Subscription Checkout persistence.
// It resolves authoritative Plan/current Subscription state and owns the
// durable frozen Checkout snapshot. It does not own Payment transport or
// Subscription mutation policy.
package checkoutpostgres

import (
	_db "common/db"
	"user/internal/usecase/subscription/checkout"
)

// @bind: user/internal/usecase/subscription/checkout.PlanResolver
// @bind: user/internal/usecase/subscription/checkout.SubscriptionReader
// @bind: user/internal/usecase/subscription/checkout.CommandStore
type CheckoutStore struct {
	db *_db.TransactionRepo
}

func NewCheckoutStore(db *_db.TransactionRepo) *CheckoutStore {
	return &CheckoutStore{db: db}
}

var (
	_ checkout.PlanResolver       = (*CheckoutStore)(nil)
	_ checkout.SubscriptionReader = (*CheckoutStore)(nil)
	_ checkout.CommandStore       = (*CheckoutStore)(nil)
)
