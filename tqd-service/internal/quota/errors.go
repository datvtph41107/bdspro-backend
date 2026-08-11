package quota

import "errors"

var (
	ErrAccessDenied         = errors.New("operation is not allowed")
	ErrQuotaExceeded        = errors.New("quota exceeded")
	ErrInvalidInput         = errors.New("quota input is invalid")
	ErrReservationNotFound  = errors.New("quota reservation not found")
	ErrReservationCommitted = errors.New("quota reservation is already committed")
	ErrReservationCanceled  = errors.New("quota reservation is canceled")
)
