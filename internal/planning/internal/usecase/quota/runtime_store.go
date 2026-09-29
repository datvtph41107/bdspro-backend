package quota

import (
	commonmetering "common/metering"
	commonoperation "common/operation"
	"context"
	"time"
	"tqd/internal/access"
)

// StoreReserveInput là dữ liệu tối thiểu mà
// runtime quota store cần để Reserve atomically.
type StoreReserveInput struct {
	ReservationID string

	Subject   access.Subject
	Operation commonoperation.Code
	MeterCode commonmetering.Code

	OperationID    string
	IdempotencyKey string
	CommandKey     string

	Amount int64
	Limit  int64

	PeriodStart time.Time
	PeriodEnd   time.Time
	ExpiresAt   time.Time
}

// Store là dependency tối thiểu của quota fast path.
//
// Nó KHÔNG phải "mọi thứ Redis quota store làm được".
//
// Repair packages sở hữu interface riêng cho:
//   - expired reservation lookup
//   - usage projection
//   - compare-and-set repair
type Store interface {
	ReserveQuota(
		ctx context.Context,
		input StoreReserveInput,
	) (Reservation, error)

	CommitQuota(
		ctx context.Context,
		reservationID string,
	) (Reservation, error)

	CancelQuota(
		ctx context.Context,
		reservationID string,
	) (Reservation, error)
}
