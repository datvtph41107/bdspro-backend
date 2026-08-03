package requestcontext

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestIsValidRequestID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		value string
		valid bool
	}{
		{
			name:  "uuid",
			value: "550e8400-e29b-41d4-a716-446655440000",
			valid: true,
		},
		{
			name:  "ulid",
			value: "01J4Z7M3N8J4NQ9KQCYXGZP7A8",
			valid: true,
		},
		{
			name:  "server format",
			value: "req_0123456789abcdef",
			valid: true,
		},
		{
			name:  "empty",
			value: "",
			valid: false,
		},
		{
			name:  "space",
			value: "request id",
			valid: false,
		},
		{
			name:  "newline",
			value: "request\nid",
			valid: false,
		},
		{
			name:  "slash",
			value: "request/id",
			valid: false,
		},
		{
			name: "too long",
			value: strings.Repeat(
				"a",
				MaxRequestIDLength+1,
			),
			valid: false,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := IsValidRequestID(tt.value)

			if got != tt.valid {
				t.Fatalf(
					"IsValidRequestID(%q) = %v, want %v",
					tt.value,
					got,
					tt.valid,
				)
			}
		})
	}
}

func TestEnsureRequestIDPreservesValidClientValue(
	t *testing.T,
) {
	t.Parallel()

	req := httptest.NewRequest(
		http.MethodGet,
		"http://example.test/v2/tqd/parcels",
		nil,
	)

	req.Header.Set(
		RequestIDHeader,
		"client-request-123",
	)

	got := EnsureRequestID(req)

	if got != "client-request-123" {
		t.Fatalf(
			"EnsureRequestID() = %q, want client value",
			got,
		)
	}

	if header := req.Header.Get(RequestIDHeader); header != got {
		t.Fatalf(
			"request header = %q, want %q",
			header,
			got,
		)
	}
}

func TestEnsureRequestIDGeneratesWhenMissing(
	t *testing.T,
) {
	t.Parallel()

	req := httptest.NewRequest(
		http.MethodGet,
		"http://example.test/v2/tqd/parcels",
		nil,
	)

	got := EnsureRequestID(req)

	if !IsValidRequestID(got) {
		t.Fatalf(
			"generated request ID is invalid: %q",
			got,
		)
	}

	if header := req.Header.Get(RequestIDHeader); header != got {
		t.Fatalf(
			"request header = %q, want %q",
			header,
			got,
		)
	}
}

func TestEnsureRequestIDReplacesInvalidClientValue(
	t *testing.T,
) {
	t.Parallel()

	req := httptest.NewRequest(
		http.MethodGet,
		"http://example.test/v2/tqd/parcels",
		nil,
	)

	req.Header.Set(
		RequestIDHeader,
		"unsafe request id",
	)

	got := EnsureRequestID(req)

	if got == "unsafe request id" {
		t.Fatal(
			"invalid client request ID was preserved",
		)
	}

	if !IsValidRequestID(got) {
		t.Fatalf(
			"replacement request ID is invalid: %q",
			got,
		)
	}
}

func TestNewRequestIDIsUniqueAcrossSample(
	t *testing.T,
) {
	t.Parallel()

	const count = 2000

	seen := make(map[string]struct{}, count)

	for i := 0; i < count; i++ {
		requestID := NewRequestID()

		if _, exists := seen[requestID]; exists {
			t.Fatalf(
				"duplicate request ID generated: %q",
				requestID,
			)
		}

		seen[requestID] = struct{}{}
	}
}

func TestEnsureRequestIDInitializesNilHeader(
	t *testing.T,
) {
	t.Parallel()

	req := &http.Request{}

	got := EnsureRequestID(req)

	if !IsValidRequestID(got) {
		t.Fatalf(
			"generated request ID is invalid: %q",
			got,
		)
	}

	if req.Header.Get(RequestIDHeader) != got {
		t.Fatalf(
			"request header was not initialized with %q",
			got,
		)
	}
}

func TestRequestIDContextRoundTrip(
	t *testing.T,
) {
	t.Parallel()

	ctx := WithRequestID(
		context.Background(),
		"request-context-123",
	)

	got, ok := FromContext(ctx)

	if !ok || got != "request-context-123" {
		t.Fatalf(
			"FromContext() = %q, %v",
			got,
			ok,
		)
	}
}

func TestWithRequestIDIgnoresInvalidValue(
	t *testing.T,
) {
	t.Parallel()

	ctx := WithRequestID(
		context.Background(),
		"unsafe request id",
	)

	if got, ok := FromContext(ctx); ok || got != "" {
		t.Fatalf(
			"invalid request ID reached context: %q",
			got,
		)
	}
}

func TestSetResponseHeaderEchoesCanonicalValue(
	t *testing.T,
) {
	t.Parallel()

	recorder := httptest.NewRecorder()

	SetResponseHeader(
		recorder,
		"request-response-123",
	)

	got := recorder.Header().Get(RequestIDHeader)

	if got != "request-response-123" {
		t.Fatalf(
			"response request ID = %q",
			got,
		)
	}
}

func TestSetResponseHeaderIgnoresInvalidValue(
	t *testing.T,
) {
	t.Parallel()

	recorder := httptest.NewRecorder()

	SetResponseHeader(
		recorder,
		"unsafe request id",
	)

	got := recorder.Header().Get(RequestIDHeader)

	if got != "" {
		t.Fatalf(
			"unsafe request ID reflected in response: %q",
			got,
		)
	}
}
