package providers

import (
	_db "common/db"
	"user/internal/usecase/plan/publish"
)

// NewPlanPublishTransaction adapts the shared transaction repo for plan publish.
func NewPlanPublishTransaction(db *_db.TransactionRepo) publish.Transaction {
	return db
}
