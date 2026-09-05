package application

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
)

/**
 * BuildRequestHash tạo dấu vết ổn định cho input đi cùng Idempotency-Key.
 *
 * Cùng key nhưng input khác sẽ bị reject thay vì trả nhầm report cũ.
 */
func BuildRequestHash(input Input) (string, error) {
	payload, err := json.Marshal(input)
	if err != nil {
		return "", ErrInvalidInput
	}
	sum := sha256.Sum256(payload)
	return hex.EncodeToString(sum[:]), nil
}
