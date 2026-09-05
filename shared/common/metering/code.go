// Package metering owns stable identifiers for commercial usage meters.
package metering

import (
	"common/codeformat"
	"errors"
	"strings"
)

// ErrInvalidCode reports a meter identifier outside the shared code grammar.
var ErrInvalidCode = errors.New("meter code is invalid")

// Code identifies one allowance pool. It is distinct from the operation that caused usage.
type Code string

// ParseCode validates a meter code at a transport or persistence boundary.
// Surrounding whitespace is normalized for compatibility with the existing
// metering boundary; canonical storage remains the trimmed code.
func ParseCode(value string) (Code, error) {
	value = strings.TrimSpace(value)
	if !codeformat.IsValid(value) {
		return "", ErrInvalidCode
	}

	return Code(value), nil
}

// IsValid reports whether the meter code uses the stable shared grammar.
func (c Code) IsValid() bool {
	_, err := ParseCode(string(c))
	return err == nil
}
