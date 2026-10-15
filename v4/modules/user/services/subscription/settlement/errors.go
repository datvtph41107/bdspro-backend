package settlement

import "errors"

var (
	ErrInvalidEffect        = errors.New("subscription settlement effect is invalid")
	ErrEffectConflict       = errors.New("subscription settlement effect conflicts with durable evidence")
	ErrCatalogTermsConflict = errors.New("subscription settlement commercial terms conflict with catalog")
	ErrSameTier             = errors.New("subscription settlement targets the current tier")
	ErrDowngrade            = errors.New("subscription settlement cannot apply a downgrade")
	ErrCurrentState         = errors.New("subscription state cannot accept settlement")
)
