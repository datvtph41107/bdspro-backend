package quota

import (
	commonaccess "common/access"
	"common/businessoperation"
	"context"
	"time"
)

/**
 * StoreReserveInput là dữ liệu tối thiểu để storage giữ quota atomically.
 */
type StoreReserveInput struct {
	ReservationID string
	Subject       commonaccess.Subject
	Operation     businessoperation.Code

	OperationID    string
	IdempotencyKey string
	CommandKey     string

	Amount int64
	Limit  int64

	PeriodStart time.Time
	PeriodEnd   time.Time
	ExpiresAt   time.Time
}

/**
 * Store giữ usage và reservation runtime.
 *
 * Business code không cần biết implementation bên dưới là Redis hay memory.
 */
type Store interface {
	ReserveQuota(ctx context.Context, input StoreReserveInput) (Reservation, error)
	CommitQuota(ctx context.Context, reservationID string) (Reservation, error)
	CancelQuota(ctx context.Context, reservationID string) (Reservation, error)

	GetUsage(
		ctx context.Context,
		subject commonaccess.Subject,
		operation businessoperation.Code,
		periodStart time.Time,
		periodEnd time.Time,
	) (Usage, error)

	SetUsed(
		ctx context.Context,
		subject commonaccess.Subject,
		operation businessoperation.Code,
		periodStart time.Time,
		periodEnd time.Time,
		used int64,
	) error

	FindExpiredReservations(ctx context.Context, before time.Time, limit int) ([]Reservation, error)
}
