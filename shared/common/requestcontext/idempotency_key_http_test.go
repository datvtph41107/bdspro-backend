package requestcontext

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestPrepareHTTPIdempotencyKeyMissingIsNoOp(t *testing.T) {
	t.Parallel()

	request := httptest.NewRequest(http.MethodPost, "/commands", nil)
	originalContext := request.Context()

	key, ok, err := PrepareHTTPIdempotencyKey(request)
	if err != nil || ok || key != "" {
		t.Fatalf("PrepareHTTPIdempotencyKey() = %q, %v, %v", key, ok, err)
	}
	if request.Context() != originalContext {
		t.Fatal("missing key changed the request context")
	}
	if values := request.Header.Values(IdempotencyKeyHeader); len(values) != 0 {
		t.Fatalf("missing key was originated: %v", values)
	}
}

func TestPrepareHTTPIdempotencyKeyBindsAndCanonicalizes(t *testing.T) {
	t.Parallel()

	request := httptest.NewRequest(http.MethodPost, "/commands", nil)
	request.Header.Add(IdempotencyKeyHeader, "  idem-command-123  ")

	key, ok, err := PrepareHTTPIdempotencyKey(request)
	if err != nil || !ok || key != "idem-command-123" {
		t.Fatalf("PrepareHTTPIdempotencyKey() = %q, %v, %v", key, ok, err)
	}
	if got, found := IdempotencyKeyFromContext(request.Context()); !found || got != key {
		t.Fatalf("context key = %q, %v; want %q, true", got, found, key)
	}
	if values := request.Header.Values(IdempotencyKeyHeader); !reflect.DeepEqual(values, []string{key}) {
		t.Fatalf("canonical header values = %v, want [%s]", values, key)
	}
}

func TestPrepareHTTPIdempotencyKeyRejectsInputAtomically(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name    string
		values  []string
		wantErr error
	}{
		{name: "invalid", values: []string{"unsafe key"}, wantErr: ErrInvalidIdempotencyKey},
		{name: "multiple", values: []string{"idem-first", "idem-second"}, wantErr: ErrMultipleIdempotencyKeys},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			request := httptest.NewRequest(http.MethodPost, "/commands", nil)
			for _, value := range testCase.values {
				request.Header.Add(IdempotencyKeyHeader, value)
			}
			originalContext := request.Context()
			originalHeaders := request.Header.Clone()

			_, _, err := PrepareHTTPIdempotencyKey(request)
			if !errors.Is(err, testCase.wantErr) {
				t.Fatalf("PrepareHTTPIdempotencyKey() error = %v, want %v", err, testCase.wantErr)
			}
			if request.Context() != originalContext || !reflect.DeepEqual(request.Header, originalHeaders) {
				t.Fatal("failed preparation partially mutated the request")
			}
		})
	}
}

func TestPrepareHTTPIdempotencyKeyRejectsContextConflictAtomically(t *testing.T) {
	t.Parallel()

	request := httptest.NewRequest(http.MethodPost, "/commands", nil)
	bound, err := BindIdempotencyKey(context.Background(), "idem-existing-123")
	if err != nil {
		t.Fatalf("BindIdempotencyKey() error = %v", err)
	}
	request = request.WithContext(bound)
	request.Header.Set(IdempotencyKeyHeader, "idem-different-456")
	originalHeaders := request.Header.Clone()

	_, _, err = PrepareHTTPIdempotencyKey(request)
	if !errors.Is(err, ErrIdempotencyKeyContextConflict) {
		t.Fatalf("PrepareHTTPIdempotencyKey() error = %v", err)
	}
	if !reflect.DeepEqual(request.Header, originalHeaders) {
		t.Fatal("context conflict partially mutated headers")
	}
	if key, ok := IdempotencyKeyFromContext(request.Context()); !ok || key != "idem-existing-123" {
		t.Fatalf("existing context key = %q, %v", key, ok)
	}
}
