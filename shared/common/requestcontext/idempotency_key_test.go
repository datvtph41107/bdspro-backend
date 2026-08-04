package requestcontext

import (
	"context"
	"errors"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestIsValidIdempotencyKey(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		value string
		valid bool
	}{
		{name: "valid", value: "report-generation-123", valid: true},
		{name: "empty", value: "", valid: false},
		{name: "space", value: "unsafe key", valid: false},
		{name: "slash", value: "unsafe/key", valid: false},
		{name: "too long", value: strings.Repeat("a", MaxIdempotencyKeyLength+1), valid: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidIdempotencyKey(tt.value); got != tt.valid {
				t.Fatalf("IsValidIdempotencyKey(%q) = %v, want %v", tt.value, got, tt.valid)
			}
		})
	}
}

func TestIdempotencyKeyFromRequestMissing(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest("POST", "/reports", nil)
	key, ok, err := IdempotencyKeyFromRequest(req)
	if err != nil || ok || key != "" {
		t.Fatalf("IdempotencyKeyFromRequest() = %q, %v, %v", key, ok, err)
	}
}

func TestIdempotencyKeyFromRequestNormalizesValidValue(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest("POST", "/reports", nil)
	req.Header.Set(IdempotencyKeyHeader, "  report-generation-123  ")

	key, ok, err := IdempotencyKeyFromRequest(req)
	if err != nil || !ok || key != "report-generation-123" {
		t.Fatalf("IdempotencyKeyFromRequest() = %q, %v, %v", key, ok, err)
	}
	if got := req.Header.Get(IdempotencyKeyHeader); got != key {
		t.Fatalf("header = %q, want %q", got, key)
	}
}

func TestIdempotencyKeyFromRequestRejectsInvalidValue(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest("POST", "/reports", nil)
	req.Header.Set(IdempotencyKeyHeader, "unsafe key")

	_, _, err := IdempotencyKeyFromRequest(req)
	if !errors.Is(err, ErrInvalidIdempotencyKey) {
		t.Fatalf("error = %v, want ErrInvalidIdempotencyKey", err)
	}
}

func TestIdempotencyKeyFromRequestRejectsMultipleValues(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest("POST", "/reports", nil)
	req.Header.Add(IdempotencyKeyHeader, "idem-a")
	req.Header.Add(IdempotencyKeyHeader, "idem-b")

	_, _, err := IdempotencyKeyFromRequest(req)
	if !errors.Is(err, ErrInvalidIdempotencyKey) {
		t.Fatalf("error = %v, want ErrInvalidIdempotencyKey", err)
	}
}

func TestIdempotencyKeyContextRoundTrip(t *testing.T) {
	t.Parallel()

	ctx := WithIdempotencyKey(context.Background(), "idem-context-123")
	got, ok := IdempotencyKeyFromContext(ctx)
	if !ok || got != "idem-context-123" {
		t.Fatalf("IdempotencyKeyFromContext() = %q, %v", got, ok)
	}
}

func TestNewIdempotencyKey(t *testing.T) {
	t.Parallel()

	first := NewIdempotencyKey()
	second := NewIdempotencyKey()
	if !IsValidIdempotencyKey(first) || !IsValidIdempotencyKey(second) {
		t.Fatalf("generated keys are invalid: %q, %q", first, second)
	}
	if first == second {
		t.Fatalf("generated duplicate key: %q", first)
	}
}
