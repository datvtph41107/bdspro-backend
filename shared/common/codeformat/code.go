// Package codeformat validates the shared syntax of stable dotted identifiers.
// It owns syntax only; operation, feature, and meter semantics remain in their
// respective packages.
package codeformat

import "strings"

const MaxLength = 128

// IsValid reports whether value is a lowercase dotted identifier with at least
// two segments.
func IsValid(value string) bool {
	if value == "" ||
		len(value) > MaxLength ||
		strings.TrimSpace(value) != value {
		return false
	}

	parts := strings.Split(value, ".")
	if len(parts) < 2 {
		return false
	}

	for _, part := range parts {
		if !validPart(part) {
			return false
		}
	}

	return true
}

func validPart(value string) bool {
	if value == "" ||
		value[0] < 'a' ||
		value[0] > 'z' {
		return false
	}

	for i := 1; i < len(value); i++ {
		current := value[i]
		if (current >= 'a' && current <= 'z') ||
			(current >= '0' && current <= '9') ||
			current == '_' {
			continue
		}
		return false
	}

	return true
}
