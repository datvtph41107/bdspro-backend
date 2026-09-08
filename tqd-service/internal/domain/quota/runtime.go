package quota

import (
	commonmetering "common/metering"
	"common/operation"
	"time"
	"tqd/internal/access"
)

/**
 * State cho biết reservation đang ở bước nào.
 */
type State string

const (
	StateReserved  State = "reserved"
	StateCommitted State = "committed"
	StateCanceled  State = "canceled"
)

/**
 * Reservation là phần quota đã được giữ cho một command.
 */
type Reservation struct {
	Required bool

	ID             string
	Subject        access.Subject
	Operation      operation.Code
	MeterCode      commonmetering.Code
	OperationID    string
	IdempotencyKey string
	CommandKey     string

	Amount int64
	Limit  int64

	PeriodStart time.Time
	PeriodEnd   time.Time

	Used      int64
	Reserved  int64
	Remaining int64

	State     State
	ExpiresAt time.Time
}

/**
 * Usage cho biết runtime usage hiện tại của một subject + meter + period.
 */
type Usage struct {
	Used     int64
	Reserved int64
}
