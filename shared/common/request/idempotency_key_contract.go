package request

import (
	"errors"
	"strings"
)

const (
	IdempotencyKeyHeader      = "Idempotency-Key"
	IdempotencyKeyMetadataKey = "idempotency-key"
	MaxIdempotencyKeyLength   = 128
)

var (
	ErrInvalidIdempotencyKey = errors.New(
		"invalid idempotency key",
	)

	ErrMultipleIdempotencyKeys = errors.New(
		"multiple idempotency keys",
	)
)

/**
 * NormalizeIdempotencyKey normalizes one caller-owned command key.
 */
func NormalizeIdempotencyKey(value string) string {
	return strings.TrimSpace(value)
}

/**
 * IsValidIdempotencyKey validates one normalized transport-safe key.
 */
func IsValidIdempotencyKey(value string) bool {
	if value == "" || len(value) > MaxIdempotencyKeyLength {
		return false
	}

	for index := 0; index < len(value); index++ {
		character := value[index]
		if (character >= 'a' && character <= 'z') ||
			(character >= 'A' && character <= 'Z') ||
			(character >= '0' && character <= '9') ||
			character == '-' ||
			character == '_' ||
			character == '.' ||
			character == ':' {
			continue
		}

		return false
	}

	return true
}

/**
 * ParseIdempotencyKey accepts zero or exactly one valid optional key.
 * It never originates or replaces a key.
 */
func ParseIdempotencyKey(
	values []string,
) (
	string,
	bool,
	error,
) {
	switch len(values) {
	case 0:
		return "", false, nil

	case 1:
		key := NormalizeIdempotencyKey(values[0])
		if !IsValidIdempotencyKey(key) {
			return "", false, ErrInvalidIdempotencyKey
		}

		return key, true, nil

	default:
		return "", false, ErrMultipleIdempotencyKeys
	}
}
