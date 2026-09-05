package serviceauth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"strconv"
	"strings"
	"time"
)

const (
	defaultMaxAge     = 5 * time.Minute
	defaultFutureSkew = 1 * time.Minute
)

// Verifier validates the File-service internal-call authentication envelope.
//
// The shared secret proves caller possession; X-Service-Name is part of the
// signed material and therefore must be explicit. A raw secret is never a
// valid wire credential: every request is timestamped to bound replay.
type Verifier struct {
	key        []byte
	maxAge     time.Duration
	futureSkew time.Duration
}

func NewVerifier(secret string) (Verifier, bool) {
	secret = strings.TrimSpace(secret)
	if secret == "" {
		return Verifier{}, false
	}
	return Verifier{
		key:        []byte(secret),
		maxAge:     defaultMaxAge,
		futureSkew: defaultFutureSkew,
	}, true
}

// Verify checks "timestamp:signature" where signature is
// base64url(HMAC-SHA256(timestamp + ":" + serviceName, secret)).
func (v Verifier) Verify(header, serviceName string, now time.Time) bool {
	header = strings.TrimSpace(header)
	serviceName = strings.TrimSpace(serviceName)
	if header == "" || serviceName == "" || len(v.key) == 0 {
		return false
	}

	parts := strings.SplitN(header, ":", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return false
	}

	timestamp, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return false
	}

	requestTime := time.Unix(timestamp, 0)
	if now.Sub(requestTime) > v.maxAge || requestTime.Sub(now) > v.futureSkew {
		return false
	}

	expected := sign(parts[0]+":"+serviceName, v.key)
	return hmac.Equal([]byte(parts[1]), []byte(expected))
}

func sign(data string, key []byte) string {
	mac := hmac.New(sha256.New, key)
	_, _ = mac.Write([]byte(data))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}
