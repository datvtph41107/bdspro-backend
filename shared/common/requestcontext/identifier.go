package requestcontext

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"sync/atomic"
	"time"
)

const maxContextIDLength = 128

var fallbackSequence atomic.Uint64

func ensureHTTPHeaderID(
	r *http.Request,
	header string,
	isValid func(string) bool,
	newID func() string,
) string {
	if r == nil {
		return newID()
	}

	if r.Header == nil {
		r.Header = make(http.Header)
	}

	value := strings.TrimSpace(r.Header.Get(header))
	if !isValid(value) {
		value = newID()
	}

	r.Header.Set(header, value)
	return value
}

func setHTTPResponseHeader(
	w http.ResponseWriter,
	header string,
	value string,
	isValid func(string) bool,
) {
	if w == nil || !isValid(value) {
		return
	}

	w.Header().Set(header, value)
}

func isValidContextID(value string, maxLength int) bool {
	if value == "" || len(value) > maxLength {
		return false
	}

	for i := 0; i < len(value); i++ {
		c := value[i]
		if (c >= 'a' && c <= 'z') ||
			(c >= 'A' && c <= 'Z') ||
			(c >= '0' && c <= '9') ||
			c == '-' ||
			c == '_' ||
			c == '.' ||
			c == ':' {
			continue
		}

		return false
	}

	return true
}

func newContextID(prefix string) string {
	var raw [16]byte
	if _, err := rand.Read(raw[:]); err == nil {
		return prefix + "_" + hex.EncodeToString(raw[:])
	}

	return fmt.Sprintf(
		"%s_%x_%x",
		prefix,
		time.Now().UnixNano(),
		fallbackSequence.Add(1),
	)
}
