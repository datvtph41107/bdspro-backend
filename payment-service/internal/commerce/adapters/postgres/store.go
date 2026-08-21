// Package postgres implements Payment Commerce persistence only.
// Business policy remains in domain/application; this adapter owns PostgreSQL
// transaction, uniqueness, row-lock and stale-claim mechanics.
package postgres

import (
	"payment/internal/commerce/application"

	"gorm.io/gorm"
)

type Store struct {
	db *gorm.DB
}

func NewStore(db *gorm.DB) *Store { return &Store{db: db} }

var (
	_ application.OrderStore          = (*Store)(nil)
	_ application.SettlementStore     = (*Store)(nil)
	_ application.FulfillmentStore    = (*Store)(nil)
	_ application.RecoveryStore       = (*Store)(nil)
	_ application.PaymentAttemptStore = (*Store)(nil)
	_ application.OutboxStore         = (*Store)(nil)
)
