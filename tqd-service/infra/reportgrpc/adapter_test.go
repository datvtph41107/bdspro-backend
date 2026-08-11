package reportgrpc

import (
	"common/requestcontext"
	"context"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestBindLegacyIdempotencyKey(t *testing.T) {
	ctx, err := bindLegacyIdempotencyKey(context.Background(), " idem_legacy_1 ")
	if err != nil {
		t.Fatal(err)
	}
	value, ok := requestcontext.IdempotencyKeyFromContext(ctx)
	if !ok || value != "idem_legacy_1" {
		t.Fatalf("bound key = %q, %v", value, ok)
	}
}

func TestBindLegacyIdempotencyKeyRequiresTransportAgreement(t *testing.T) {
	ctx, err := requestcontext.BindIdempotencyKey(context.Background(), "idem_header")
	if err != nil {
		t.Fatal(err)
	}

	_, err = bindLegacyIdempotencyKey(ctx, "idem_body")
	if status.Code(err) != codes.InvalidArgument {
		t.Fatalf("conflict error = %v, want InvalidArgument", err)
	}

	if _, err := bindLegacyIdempotencyKey(ctx, "idem_header"); err != nil {
		t.Fatalf("matching keys must be accepted: %v", err)
	}
}

func TestBindLegacyIdempotencyKeyRejectsInvalidValue(t *testing.T) {
	_, err := bindLegacyIdempotencyKey(context.Background(), "unsafe key")
	if status.Code(err) != codes.InvalidArgument {
		t.Fatalf("invalid key error = %v, want InvalidArgument", err)
	}
}
