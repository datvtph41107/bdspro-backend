package requestgrpc

import (
	"common/requestcontext"
	"context"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func TestRestoreFromMetadataBindsOperationAndIdempotency(t *testing.T) {
	ctx, err := RestoreFromMetadata(context.Background(), metadata.Pairs(
		requestcontext.RequestIDMetadataKey, "req_123",
		requestcontext.OperationIDMetadataKey, "op_123",
		requestcontext.IdempotencyKeyMetadataKey, "idem_123",
	))
	if err != nil {
		t.Fatal(err)
	}
	if got, ok := requestcontext.OperationIDFromContext(ctx); !ok || got != "op_123" {
		t.Fatalf("operation id = %q, %v", got, ok)
	}
	if got, ok := requestcontext.IdempotencyKeyFromContext(ctx); !ok || got != "idem_123" {
		t.Fatalf("idempotency key = %q, %v", got, ok)
	}
}

func TestRestoreFromMetadataAllowsMissingOperationID(t *testing.T) {
	if _, err := RestoreFromMetadata(context.Background(), metadata.MD{}); err != nil {
		t.Fatal(err)
	}
}

func TestRestoreFromMetadataRejectsDuplicateOperationID(t *testing.T) {
	_, err := RestoreFromMetadata(context.Background(), metadata.MD{
		requestcontext.OperationIDMetadataKey: []string{"op_a", "op_b"},
	})
	if status.Code(err) != codes.InvalidArgument {
		t.Fatalf("code = %v, want InvalidArgument", status.Code(err))
	}
}

func TestRestoreFromMetadataRejectsOperationConflict(t *testing.T) {
	ctx, err := requestcontext.BindOperationID(context.Background(), "op_a")
	if err != nil {
		t.Fatal(err)
	}
	_, err = RestoreFromMetadata(ctx, metadata.Pairs(requestcontext.OperationIDMetadataKey, "op_b"))
	if status.Code(err) != codes.Internal {
		t.Fatalf("code = %v, want Internal", status.Code(err))
	}
}
