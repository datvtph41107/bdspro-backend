package grpcadapter

import (
	"common/requestcontext"
	"context"
	"testing"
	"tqd/internal/workspace/report/application"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestBindRequestCommandKeyCompatibility(t *testing.T) {
	ctx, err := bindRequestCommandKeyCompatibility(context.Background(), " idem_legacy_1 ")
	if err != nil {
		t.Fatal(err)
	}
	value, ok := requestcontext.IdempotencyKeyFromContext(ctx)
	if !ok || value != "idem_legacy_1" {
		t.Fatalf("bound key = %q, %v", value, ok)
	}
}

func TestBindRequestCommandKeyCompatibilityRequiresTransportAgreement(t *testing.T) {
	ctx, err := requestcontext.BindIdempotencyKey(context.Background(), "idem_header")
	if err != nil {
		t.Fatal(err)
	}

	_, err = bindRequestCommandKeyCompatibility(ctx, "idem_body")
	if status.Code(err) != codes.InvalidArgument {
		t.Fatalf("conflict error = %v, want InvalidArgument", err)
	}

	if _, err := bindRequestCommandKeyCompatibility(ctx, "idem_header"); err != nil {
		t.Fatalf("matching keys must be accepted: %v", err)
	}
}

func TestBindRequestCommandKeyCompatibilityRejectsInvalidValue(t *testing.T) {
	_, err := bindRequestCommandKeyCompatibility(context.Background(), "unsafe key")
	if status.Code(err) != codes.InvalidArgument {
		t.Fatalf("invalid key error = %v, want InvalidArgument", err)
	}
}

func TestToStatusTreatsMissingOperationAsServerFailure(t *testing.T) {
	err := toStatus(application.ErrOperationMissing)
	if status.Code(err) != codes.Internal {
		t.Fatalf(
			"missing operation status = %v, want Internal; error = %v",
			status.Code(err),
			err,
		)
	}
}
