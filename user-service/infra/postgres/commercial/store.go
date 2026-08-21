// Package commercialpostgres implements User commercial persistence.
// It resolves authoritative Catalog/current Subscription state and owns the
// durable frozen Checkout command. It does not own Payment transport or
// Subscription mutation policy.
package commercialpostgres

import (
	_db "common/db"
	"user/internal/modules/commercial/checkout"
)

type CheckoutStore struct {
	db *_db.TransactionRepo
}

func NewCheckoutStore(db *_db.TransactionRepo) *CheckoutStore {
	return &CheckoutStore{db: db}
}

var (
	_ checkout.CatalogResolver    = (*CheckoutStore)(nil)
	_ checkout.SubscriptionReader = (*CheckoutStore)(nil)
	_ checkout.CommandStore       = (*CheckoutStore)(nil)
)
