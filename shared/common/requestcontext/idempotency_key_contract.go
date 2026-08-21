package requestcontext

import canonical "common/request"

const (
	IdempotencyKeyHeader      = "Idempotency-Key"
	IdempotencyKeyMetadataKey = "idempotency-key"
	MaxIdempotencyKeyLength   = canonical.MaxIdempotencyKeyLength
)

var (
	ErrInvalidIdempotencyKey   = canonical.ErrInvalidIdempotencyKey
	ErrMultipleIdempotencyKeys = canonical.ErrMultipleIdempotencyKeys
)

func NormalizeIdempotencyKey(value string) string { return canonical.NormalizeIdempotencyKey(value) }
func IsValidIdempotencyKey(value string) bool     { return canonical.IsValidIdempotencyKey(value) }
func ParseIdempotencyKey(values []string) (string, bool, error) {
	return canonical.ParseIdempotencyKey(values)
}
