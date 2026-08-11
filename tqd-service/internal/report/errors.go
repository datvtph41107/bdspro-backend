package report

import "errors"

var (
	ErrInvalidInput       = errors.New("report input is invalid")
	ErrActorMissing       = errors.New("request actor is missing")
	ErrProfileMissing     = errors.New("actor profile is missing")
	ErrOperationMissing   = errors.New("business operation is missing")
	ErrOperationMismatch  = errors.New("business operation does not match create report")
	ErrOperationIDMissing = errors.New("operation ID is missing")
	ErrCommandConflict    = errors.New("idempotency key was already used with different report input")
	ErrParcelNotFound     = errors.New("parcel not found")
	ErrRegionNotFound     = errors.New("region not found")
	ErrQuotaCommitFailed  = errors.New("report accepted but quota commit failed")
)
