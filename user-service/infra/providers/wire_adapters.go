package providers

import (
	"time"

	_db "common/db"
	"user/internal/usecase/plan/publish"
	"user/internal/usecase/subscription/checkout"
)

// NewPlanPublishTransaction adapts the shared transaction repo for plan publish.
func NewPlanPublishTransaction(db *_db.TransactionRepo) publish.Transaction {
	return db
}

// NewCancelClock supplies wall-clock time for subscription cancellation.
func NewCancelClock() func() time.Time {
	return time.Now
}

// NewCheckoutClock supplies wall-clock time for subscription checkout.
func NewCheckoutClock() checkout.Clock {
	return time.Now
}
