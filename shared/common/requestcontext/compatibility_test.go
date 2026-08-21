package requestcontext

import (
	canonical "common/request"
	"context"
	"testing"
)

func TestCompatibilityRequestIDSharesCanonicalContext(t *testing.T) {
	ctx := WithRequestID(context.Background(), "req_bridge_123")
	if got, ok := canonical.RequestIDFromContext(ctx); !ok || got != "req_bridge_123" {
		t.Fatalf("canonical request ID = %q, %v", got, ok)
	}

	ctx = canonical.WithRequestID(context.Background(), "req_target_456")
	if got, ok := FromContext(ctx); !ok || got != "req_target_456" {
		t.Fatalf("legacy request ID = %q, %v", got, ok)
	}
}

func TestCompatibilityOperationIDSharesCanonicalContext(t *testing.T) {
	ctx, err := BindOperationID(context.Background(), "op_bridge_123")
	if err != nil {
		t.Fatal(err)
	}
	if got, ok := canonical.OperationIDFromContext(ctx); !ok || got != "op_bridge_123" {
		t.Fatalf("canonical operation ID = %q, %v", got, ok)
	}

	ctx, err = canonical.BindOperationID(context.Background(), "op_target_456")
	if err != nil {
		t.Fatal(err)
	}
	if got, ok := OperationIDFromContext(ctx); !ok || got != "op_target_456" {
		t.Fatalf("legacy operation ID = %q, %v", got, ok)
	}
}

func TestCompatibilityIdempotencyKeySharesCanonicalContext(t *testing.T) {
	ctx, err := BindIdempotencyKey(context.Background(), "idem_bridge_123")
	if err != nil {
		t.Fatal(err)
	}
	if got, ok := canonical.IdempotencyKeyFromContext(ctx); !ok || got != "idem_bridge_123" {
		t.Fatalf("canonical idempotency key = %q, %v", got, ok)
	}

	ctx, err = canonical.BindIdempotencyKey(context.Background(), "idem_target_456")
	if err != nil {
		t.Fatal(err)
	}
	if got, ok := IdempotencyKeyFromContext(ctx); !ok || got != "idem_target_456" {
		t.Fatalf("legacy idempotency key = %q, %v", got, ok)
	}
}
