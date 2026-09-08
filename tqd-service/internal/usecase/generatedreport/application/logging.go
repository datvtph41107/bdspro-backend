package application

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

// commandRef is a safe correlation fingerprint for logs.
// Raw caller-owned command keys must not be logged.
func commandRef(commandKey string) string {
	commandKey = strings.TrimSpace(commandKey)
	if commandKey == "" {
		return ""
	}

	sum := sha256.Sum256([]byte(commandKey))
	return hex.EncodeToString(sum[:])
}
