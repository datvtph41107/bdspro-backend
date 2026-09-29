package tilesession

import (
	"testing"
	"time"
)

func TestKeyOwnsTileSessionNamespace(t *testing.T) {
	if got, want := key(42), "ss:k:42"; got != want {
		t.Fatalf("key(42) = %q, want %q", got, want)
	}
}

func TestExpirationTTLPreservesPositiveExpiry(t *testing.T) {
	now := time.Unix(100, 0)
	expiresAt := now.Add(37 * time.Minute)
	if got := expirationTTL(expiresAt, now); got != 37*time.Minute {
		t.Fatalf("expirationTTL() = %s, want %s", got, 37*time.Minute)
	}
}

func TestExpirationTTLFallsBackToCanonicalTTL(t *testing.T) {
	now := time.Unix(100, 0)
	for _, expiresAt := range []time.Time{now, now.Add(-time.Second)} {
		if got := expirationTTL(expiresAt, now); got != TTL {
			t.Fatalf("expirationTTL(%s) = %s, want %s", expiresAt, got, TTL)
		}
	}
}
