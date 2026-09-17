package wallet

import (
	"math"
	"testing"
)

func TestMinorFromWireRejectsRoundedInt64UpperBound(t *testing.T) {
	if _, err := MinorFromWire(float64(math.MaxInt64)); err == nil {
		t.Fatal("2^63 must not wrap into negative money")
	}
	value := math.Nextafter(float64(math.MaxInt64), 0)
	if got, err := MinorFromWire(value); err != nil || got <= 0 {
		t.Fatalf("valid positive value: got=%d err=%v", got, err)
	}
}
