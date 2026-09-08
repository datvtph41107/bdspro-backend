package domain

import "errors"

var (
	ErrInvalidCommercialTerms      = errors.New("invalid commercial terms")
	ErrInvalidCommand              = errors.New("invalid commerce command")
	ErrOrderNotFound               = errors.New("commerce order not found")
	ErrOrderCommandConflict        = errors.New("commerce order command conflict")
	ErrProviderEvidenceConflict    = errors.New("provider evidence conflict")
	ErrUnsupportedProviderEvent    = errors.New("unsupported provider event")
	ErrFulfillmentNotFound         = errors.New("commerce fulfillment not found")
	ErrClaimLost                   = errors.New("commerce fulfillment claim lost")
	ErrRedriveNotAllowed           = errors.New("commerce fulfillment redrive not allowed")
	ErrSubscriptionTransient       = errors.New("subscription dependency transient failure")
	ErrSubscriptionPermanent       = errors.New("subscription effect permanent failure")
	ErrSubscriptionEffectConflict  = errors.New("subscription effect conflict")
	ErrAttemptNotFound             = errors.New("payment attempt not found")
	ErrAttemptCommandConflict      = errors.New("payment attempt command conflict")
	ErrAttemptNotAllowed           = errors.New("payment attempt not allowed")
	ErrFundsConfirmationNotAllowed = errors.New("manual funds confirmation not allowed")
	ErrCommandEvidenceConflict     = errors.New("commerce command evidence conflict")
	ErrUnsupportedPaymentMethod    = errors.New("unsupported payment method")
	ErrProviderUnavailable         = errors.New("payment provider unavailable")
)
