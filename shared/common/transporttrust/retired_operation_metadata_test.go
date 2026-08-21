package transporttrust

import (
	"testing"
	"time"

	"google.golang.org/grpc/metadata"
)

func TestRetiredOperationMetadataIsOutsideTrustContract(
	t *testing.T,
) {
	t.Parallel()

	now := time.Unix(1_700_000_000, 0)
	cfg := Config{
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

	if err := Sign(
		"/test.Service/Call",
		md,
		cfg,
	); err != nil {
		t.Fatalf("Sign() error = %v", err)
	}

	// Retired headers are intentionally outside the signature contract.
	// No runtime component is allowed to consume them after W9-E.
	md.Set(
		"x-qhpro-operation-code",
		"tampered.retired_value",
	)

	if _, err := Verify(
		"/test.Service/Call",
		md,
		cfg,
	); err != nil {
		t.Fatalf(
			"retired operation metadata changed trust verification: %v",
			err,
		)
	}
}
