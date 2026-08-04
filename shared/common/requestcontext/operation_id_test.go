package requestcontext

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestIsValidOperationID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		value string
		valid bool
	}{
		{name: "uuid", value: "550e8400-e29b-41d4-a716-446655440000", valid: true},
		{name: "ulid", value: "01J4Z7M3N8J4NQ9KQCYXGZP7A8", valid: true},
		{name: "server format", value: "op_0123456789abcdef", valid: true},
		{name: "empty", value: "", valid: false},
		{name: "space", value: "operation id", valid: false},
		{name: "newline", value: "operation\nid", valid: false},
		{name: "slash", value: "operation/id", valid: false},
		{name: "too long", value: strings.Repeat("a", MaxOperationIDLength+1), valid: false},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := IsValidOperationID(tt.value); got != tt.valid {
				t.Fatalf("IsValidOperationID(%q) = %v, want %v", tt.value, got, tt.valid)
			}
		})
	}
}

func TestEnsureOperationIDPreservesValidClientValue(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodPost, "http://example.test/report", nil)
	req.Header.Set(OperationIDHeader, "report-generation-abc123")

	got := EnsureOperationID(req)
	if got != "report-generation-abc123" {
		t.Fatalf("EnsureOperationID() = %q, want client value", got)
	}
	if header := req.Header.Get(OperationIDHeader); header != got {
		t.Fatalf("operation header = %q, want %q", header, got)
	}
}

func TestEnsureOperationIDGeneratesWhenMissing(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodPost, "http://example.test/report", nil)
	got := EnsureOperationID(req)

	if !IsValidOperationID(got) {
		t.Fatalf("generated operation ID is invalid: %q", got)
	}
	if !strings.HasPrefix(got, "op_") {
		t.Fatalf("generated operation ID = %q, want op_ prefix", got)
	}
	if header := req.Header.Get(OperationIDHeader); header != got {
		t.Fatalf("operation header = %q, want %q", header, got)
	}
}

func TestEnsureOperationIDReplacesInvalidClientValue(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodPost, "http://example.test/report", nil)
	req.Header.Set(OperationIDHeader, "unsafe operation id")

	got := EnsureOperationID(req)
	if got == "unsafe operation id" || !IsValidOperationID(got) {
		t.Fatalf("invalid operation ID was not replaced: %q", got)
	}
}

func TestNewOperationIDIsUniqueAcrossSample(t *testing.T) {
	t.Parallel()

	const count = 2000
	seen := make(map[string]struct{}, count)
	for i := 0; i < count; i++ {
		operationID := NewOperationID()
		if _, exists := seen[operationID]; exists {
			t.Fatalf("duplicate operation ID generated: %q", operationID)
		}
		seen[operationID] = struct{}{}
	}
}

func TestOperationIDContextRoundTrip(t *testing.T) {
	t.Parallel()

	ctx := WithOperationID(context.Background(), "operation-context-123")
	got, ok := OperationIDFromContext(ctx)
	if !ok || got != "operation-context-123" {
		t.Fatalf("OperationIDFromContext() = %q, %v", got, ok)
	}
}

func TestWithOperationIDIgnoresInvalidValue(t *testing.T) {
	t.Parallel()

	ctx := WithOperationID(context.Background(), "unsafe operation id")
	if got, ok := OperationIDFromContext(ctx); ok || got != "" {
		t.Fatalf("invalid operation ID reached context: %q", got)
	}
}

func TestSetOperationResponseHeaderEchoesCanonicalValue(t *testing.T) {
	t.Parallel()

	recorder := httptest.NewRecorder()
	SetOperationResponseHeader(recorder, "operation-response-123")

	if got := recorder.Header().Get(OperationIDHeader); got != "operation-response-123" {
		t.Fatalf("response operation ID = %q", got)
	}
}

func TestSetOperationResponseHeaderIgnoresInvalidValue(t *testing.T) {
	t.Parallel()

	recorder := httptest.NewRecorder()
	SetOperationResponseHeader(recorder, "unsafe operation id")

	if got := recorder.Header().Get(OperationIDHeader); got != "" {
		t.Fatalf("unsafe operation ID reflected in response: %q", got)
	}
}
