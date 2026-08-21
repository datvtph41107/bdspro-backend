package requestcontext

import canonical "common/request"

const (
	RequestIDHeader      = "X-Request-ID"
	RequestIDMetadataKey = "x-request-id"
	MaxRequestIDLength   = canonical.MaxRequestIDLength
)

type RequestIDSource = canonical.RequestIDSource
type RequestIDReason = canonical.RequestIDReason
type RequestIDDecision = canonical.RequestIDDecision

const (
	RequestIDSourceClient    = canonical.RequestIDSourceClient
	RequestIDSourceGenerated = canonical.RequestIDSourceGenerated

	RequestIDReasonValid    = canonical.RequestIDReasonValid
	RequestIDReasonMissing  = canonical.RequestIDReasonMissing
	RequestIDReasonInvalid  = canonical.RequestIDReasonInvalid
	RequestIDReasonMultiple = canonical.RequestIDReasonMultiple
)

func NormalizeRequestID(value string) string             { return canonical.NormalizeRequestID(value) }
func IsValidRequestID(value string) bool                 { return canonical.IsValidRequestID(value) }
func ResolveRequestID(values []string) RequestIDDecision { return canonical.ResolveRequestID(values) }
