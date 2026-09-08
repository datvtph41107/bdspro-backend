package application

import "errors"

var (
	ErrInvalidInput = errors.New(
		"report input is invalid",
	)

	ErrActorMissing = errors.New(
		"request actor is missing",
	)

	ErrProfileMissing = errors.New(
		"actor profile is missing",
	)

	ErrOperationMissing = errors.New(
		"operation is missing",
	)

	ErrOperationIDMissing = errors.New(
		"operation ID is missing",
	)

	ErrCommandKeyMissing = errors.New(
		"idempotency key is required",
	)

	ErrAccessUnavailable = errors.New(
		"entitlement access is unavailable",
	)

	ErrAcceptanceUnavailable = errors.New(
		"report acceptance lease is unavailable",
	)

	ErrProcessingUnavailable = errors.New(
		"generated report processing is unavailable",
	)

	ErrCommandConflict = errors.New(
		"idempotency key was already used with different report input",
	)

	ErrParcelNotFound = errors.New(
		"parcel not found",
	)

	ErrRegionNotFound = errors.New(
		"region not found",
	)
)
