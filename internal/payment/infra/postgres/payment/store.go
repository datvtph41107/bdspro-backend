// Package postgres implements Payment durable persistence. Business policy is
// owned by the consumer usecases; this adapter owns PostgreSQL transaction,
// uniqueness, row-lock and stale-claim mechanics.
package postgres

import (
	"payment/internal/usecase/attempt"
	"payment/internal/usecase/fulfillment"
	"payment/internal/usecase/order"
	"payment/internal/usecase/outbox"
	"payment/internal/usecase/recovery"
	"payment/internal/usecase/settlement"

	"gorm.io/gorm"
)

type Store struct {
	db *gorm.DB
}

func NewStore(db *gorm.DB) *Store { return &Store{db: db} }

var (
	_ order.Store           = (*Store)(nil)
	_ attempt.OrderStore    = (*Store)(nil)
	_ attempt.Store         = (*Store)(nil)
	_ settlement.OrderStore = (*Store)(nil)
	_ settlement.Store      = (*Store)(nil)
	_ fulfillment.Store     = (*Store)(nil)
	_ recovery.Store        = (*Store)(nil)
	_ outbox.Store          = (*Store)(nil)
)
