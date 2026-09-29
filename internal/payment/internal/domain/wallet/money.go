package wallet

import (
	"fmt"
	"math"
)

// MinorFromWire converts the legacy protobuf double amount into the exact
// integer representation used by Wallet persistence. Wallet RPC values are
// interpreted as minor units; fractional values are rejected instead of being
// rounded silently.
func MinorFromWire(value float64) (int64, error) {
	if math.IsNaN(value) || math.IsInf(value, 0) {
		return 0, fmt.Errorf("amount must be finite")
	}
	if value != math.Trunc(value) {
		return 0, fmt.Errorf("amount must be expressed in whole minor units")
	}
	if value > float64(math.MaxInt64) || value < float64(math.MinInt64) {
		return 0, fmt.Errorf("amount is outside int64 range")
	}
	return int64(value), nil
}

func MinorToWire(value int64) float64 { return float64(value) }
