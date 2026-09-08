package request

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sync/atomic"
	"time"
)

var requestIDFallbackSequence atomic.Uint64

type requestIDReadFunc func([]byte) (int, error)

/**
 * Description: NewRequestID tạo một request ID mới
 */
func NewRequestID() string {
	return newRequestID(
		rand.Read,
		time.Now,
		&requestIDFallbackSequence,
	)
}

func newRequestID(
	read requestIDReadFunc,
	now func() time.Time,
	sequence *atomic.Uint64,
) string {
	var raw [16]byte

	count, err := read(raw[:])

	// chỉ dùng random path khi đủ 16 bytes
	// happy path
	if err == nil && count == len(raw) {
		return "req_" + hex.EncodeToString(raw[:])
	}

	// fallback path
	return fmt.Sprintf(
		"req_%x_%x",
		now().UnixNano(),
		sequence.Add(1),
	)
	// -> đảm bảo tính duy nhất
}
