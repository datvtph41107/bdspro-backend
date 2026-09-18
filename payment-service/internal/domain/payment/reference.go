package domain

import (
	"encoding/hex"
	"regexp"
	"strings"
)

const OrderReferencePrefix = "QHP-"

var embeddedOrderReference = regexp.MustCompile(`(?i)QHP-?([0-9A-F]{16})`)

// FormatOrderReference is the only formatter for Payment-owned public order
// references. Callers provide cryptographically random entropy; the domain
// owns the stable wire and storage representation.
func FormatOrderReference(entropy [8]byte) string {
	return OrderReferencePrefix + strings.ToUpper(hex.EncodeToString(entropy[:]))
}

// ExtractOrderReference accepts the one known bank normalization: removal of
// the separator from transfer descriptions. It always returns the canonical
// uppercase storage representation.
func ExtractOrderReference(values ...string) (string, bool) {
	for _, value := range values {
		match := embeddedOrderReference.FindStringSubmatch(value)
		if len(match) == 2 {
			return OrderReferencePrefix + strings.ToUpper(match[1]), true
		}
	}
	return "", false
}

func IsCanonicalOrderReference(value string) bool {
	canonical, ok := ExtractOrderReference(value)
	return ok && value == canonical
}
