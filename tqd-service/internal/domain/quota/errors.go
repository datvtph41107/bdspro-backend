package quota

import (
	"errors"
	"time"

	commonmetering "common/metering"
	commonoperation "common/operation"
	"tqd/internal/access"
)

var (
	ErrAccessDenied         = errors.New("operation is not allowed")
	ErrQuotaExceeded        = errors.New("quota exceeded")
	ErrInvalidInput         = errors.New("quota input is invalid")
	ErrReservationNotFound  = errors.New("quota reservation not found")
	ErrReservationCommitted = errors.New("quota reservation is already committed")
	ErrReservationCanceled  = errors.New("quota reservation is canceled")
)

// ExhaustedError carries the owner evidence needed by a transport presenter.
// It describes the failed admission snapshot, not durable usage accounting;
// PostgreSQL remains the authority for accepted consumption.
type ExhaustedError struct {
	Subject     access.Subject
	Operation   commonoperation.Code
	MeterCode   commonmetering.Code
	Limit       int64
	Used        int64
	Reserved    int64
	Remaining   int64
	PeriodStart time.Time
	PeriodEnd   time.Time
}

func (e *ExhaustedError) Error() string { return ErrQuotaExceeded.Error() }

func (e *ExhaustedError) Unwrap() error { return ErrQuotaExceeded }

// AccessDeniedError distinguishes a missing commercial capability from an
// exhausted allowance. The client may offer pricing for this reason, while a
// quota exhaustion presenter can show the current period capacity instead.
type AccessDeniedError struct {
	Subject   access.Subject
	Operation commonoperation.Code
}

func (e *AccessDeniedError) Error() string { return ErrAccessDenied.Error() }

func (e *AccessDeniedError) Unwrap() error { return ErrAccessDenied }
