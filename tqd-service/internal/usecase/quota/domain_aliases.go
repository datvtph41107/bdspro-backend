package quota

import quotadomain "tqd/internal/domain/quota"

type Reservation = quotadomain.Reservation
type Usage = quotadomain.Usage
type State = quotadomain.State
type ExhaustedError = quotadomain.ExhaustedError
type AccessDeniedError = quotadomain.AccessDeniedError

const (
	StateReserved  = quotadomain.StateReserved
	StateCommitted = quotadomain.StateCommitted
	StateCanceled  = quotadomain.StateCanceled
)

var (
	ErrAccessDenied         = quotadomain.ErrAccessDenied
	ErrQuotaExceeded        = quotadomain.ErrQuotaExceeded
	ErrInvalidInput         = quotadomain.ErrInvalidInput
	ErrReservationNotFound  = quotadomain.ErrReservationNotFound
	ErrReservationCommitted = quotadomain.ErrReservationCommitted
	ErrReservationCanceled  = quotadomain.ErrReservationCanceled
)
