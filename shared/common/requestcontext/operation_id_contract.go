package requestcontext

import canonical "common/request"

const (
	OperationIDHeader      = "X-Operation-ID"
	OperationIDMetadataKey = "x-operation-id"
	MaxOperationIDLength   = canonical.MaxOperationIDLength
)

type OperationIDSource = canonical.OperationIDSource
type OperationIDReason = canonical.OperationIDReason
type OperationIDDecision = canonical.OperationIDDecision

const (
	OperationIDSourceCaller            = canonical.OperationIDSourceCaller
	OperationIDSourceGeneratedFallback = canonical.OperationIDSourceGeneratedFallback
	OperationIDReasonValid             = canonical.OperationIDReasonValid
	OperationIDReasonMissing           = canonical.OperationIDReasonMissing
)

var (
	ErrInvalidOperationID    = canonical.ErrInvalidOperationID
	ErrMultipleOperationIDs  = canonical.ErrMultipleOperationIDs
	ErrOperationIDGeneration = canonical.ErrOperationIDGeneration
)

func NormalizeOperationID(value string) string         { return canonical.NormalizeOperationID(value) }
func IsValidOperationID(value string) bool             { return canonical.IsValidOperationID(value) }
func AcceptOperationID(values []string) (string, bool) { return canonical.AcceptOperationID(values) }
func ResolveOperationID(values []string) (OperationIDDecision, error) {
	return canonical.ResolveOperationID(values)
}
