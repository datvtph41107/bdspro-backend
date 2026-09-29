package helpers

import (
	"encoding/base64"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
)

var (
	ErrApplinkNotFound      = errors.New("applink: not found")
	ErrApplinkInvalidCode   = errors.New("applink: invalid base64 code")
	ErrApplinkCodeExhausted = errors.New("applink: failed to generate unique code after retries")
)

// ── Code generation ───────────────────────────────────────────────────────────
const (
	applinkCodeMin uint64 = 1_000_000      // 10^6
	applinkCodeMax uint64 = 10_000_000_000 // 10^10
)

// ApplinkGenerateCode sinh số random uint64 trong [10^6, 10^10]
func ApplinkGenerateCode() uint64 {
	rang := applinkCodeMax - applinkCodeMin + 1
	return applinkCodeMin + uint64(rand.Int63n(int64(rang)))
}

// ── Codec: number ↔ base64 ────────────────────────────────────────────────────

// ApplinkCodeToBase64 encode uint64 → base64 string (trả ra client / QR URL)
//
//	Ví dụ: 3_748_291_056 → "Mzc0ODI5MTA1Ng=="
func ApplinkCodeToBase64(code uint64) string {
	return base64.StdEncoding.EncodeToString(
		[]byte(strconv.FormatUint(code, 10)),
	)
}

// ApplinkBase64ToCode decode base64 string → uint64 (nhận từ client, query DB)
//
//	Ví dụ: "Mzc0ODI5MTA1Ng==" → 3_748_291_056
func ApplinkBase64ToCode(b64 string) (uint64, error) {
	raw, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return 0, fmt.Errorf("%w: %s", ErrApplinkInvalidCode, err.Error())
	}

	code, err := strconv.ParseUint(string(raw), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("%w: not a valid number", ErrApplinkInvalidCode)
	}

	return code, nil
}
