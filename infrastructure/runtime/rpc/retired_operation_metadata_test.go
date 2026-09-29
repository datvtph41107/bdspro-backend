package rpc

import (
	"testing"
	"time"

	"google.golang.org/grpc/metadata"
)

func TestRetiredOperationMetadataIsOutsideServiceAssertionContract(
	t *testing.T,
) {
	t.Parallel()

	now := time.Unix(1_700_000_000, 0)
	cfg := ServiceAssertionConfig{
		ServiceID: "gateway-service",
		Secret:    "test-secret",
		Now:       func() time.Time { return now },
	}

	md := metadata.Pairs(
		"x-qhpro-operation-code",
		"parcel.lookup",
		"x-qhpro-operation-binding",
		"parcel.lookup_transport",
		"x-qhpro-operation-source",
		"direct_binding",
	)

	if HasPrivilegedMetadata(md) {
		t.Fatal(
			"retired operation headers remain privileged",
		)
	}

	if err := SignServiceAssertion(
		"/test.Service/Call",
		md,
		cfg,
	); err != nil {
		t.Fatalf(
			"SignServiceAssertion() error = %v",
			err,
		)
	}

	// Retired headers are intentionally outside the assertion contract.
	md.Set(
		"x-qhpro-operation-code",
		"tampered.retired_value",
	)

	if _, err := VerifyServiceAssertion(
		"/test.Service/Call",
		md,
		cfg,
	); err != nil {
		t.Fatalf(
			"retired operation metadata changed assertion verification: %v",
			err,
		)
	}
}
