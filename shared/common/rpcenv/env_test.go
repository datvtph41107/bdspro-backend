package rpcenv

import (
	"common/rpc"
	"errors"
	"testing"
	"time"

	"google.golang.org/grpc/metadata"
)

func TestLoadTransportConfigPreservesEnvironmentContract(
	t *testing.T,
) {
	t.Setenv("QHPRO_SERVICE_ID", "gateway-service")
	t.Setenv("SERVICE_NAME", "")
	t.Setenv("QHPRO_INTERNAL_METADATA_SECRET", "shared-secret")
	t.Setenv("QHPRO_TRUSTED_METADATA_MODE", "enforce")
	t.Setenv("QHPRO_INTERNAL_METADATA_MAX_AGE", "45s")
	t.Setenv("QHPRO_INTERNAL_METADATA_CLOCK_SKEW", "7s")
	t.Setenv(
		"QHPRO_INTERNAL_METADATA_VERIFICATION_KEYS",
		`{"gateway-service":"gateway-secret","invalid service":"ignored"}`,
	)

	cfg := LoadTransportConfig()

	if cfg.ServiceAssertion.ServiceID != "gateway-service" {
		t.Fatalf(
			"service ID = %q",
			cfg.ServiceAssertion.ServiceID,
		)
	}
	if cfg.ServiceAssertion.Secret != "shared-secret" {
		t.Fatal("shared secret was not preserved")
	}
	if cfg.ServiceAssertion.MaxAge != 45*time.Second ||
		cfg.ServiceAssertion.ClockSkew != 7*time.Second {
		t.Fatalf(
			"duration config = %+v",
			cfg.ServiceAssertion,
		)
	}
	if !cfg.ServiceAssertion.VerificationKeysConfigured {
		t.Fatal("verification keyring was not marked configured")
	}
	if got := cfg.ServiceAssertion.VerificationSecrets["gateway-service"]; got != "gateway-secret" {
		t.Fatalf("gateway verification secret = %q", got)
	}
	if _, exists := cfg.ServiceAssertion.VerificationSecrets["invalid service"]; exists {
		t.Fatalf(
			"invalid service ID was loaded: %#v",
			cfg.ServiceAssertion.VerificationSecrets,
		)
	}
	if !cfg.RequireServiceAssertion {
		t.Fatal("enforce mode did not require service assertion")
	}

	t.Setenv("QHPRO_TRUSTED_METADATA_MODE", "audit")
	if LoadTransportConfig().RequireServiceAssertion {
		t.Fatal("audit mode unexpectedly requires service assertion")
	}

	t.Setenv("QHPRO_TRUSTED_METADATA_MODE", "unknown")
	if LoadTransportConfig().RequireServiceAssertion {
		t.Fatal(
			"unknown compatibility mode unexpectedly requires service assertion",
		)
	}
}

func TestLoadTransportConfigUsesServiceNameFallback(t *testing.T) {
	t.Setenv("QHPRO_SERVICE_ID", "")
	t.Setenv("SERVICE_NAME", "worker-service")

	cfg := LoadTransportConfig()

	if cfg.ServiceAssertion.ServiceID != "worker-service" {
		t.Fatalf(
			"service ID = %q, want worker-service",
			cfg.ServiceAssertion.ServiceID,
		)
	}
}

func TestInvalidVerificationKeyringDoesNotFallBackToSharedSecret(
	t *testing.T,
) {
	now := time.Unix(1_700_000_000, 0)

	signer := rpc.ServiceAssertionConfig{
		ServiceID: "gateway-service",
		Secret:    "shared-secret",
		Now:       func() time.Time { return now },
	}

	md := metadata.Pairs(rpc.ProfileIDMetadataKey, "12")
	if err := rpc.SignServiceAssertion(
		"/test.Service/Call",
		md,
		signer,
	); err != nil {
		t.Fatal(err)
	}

	t.Setenv(
		"QHPRO_INTERNAL_METADATA_VERIFICATION_KEYS",
		"not-json",
	)
	t.Setenv(
		"QHPRO_INTERNAL_METADATA_SECRET",
		"shared-secret",
	)

	cfg := LoadTransportConfig()
	cfg.ServiceAssertion.Now = func() time.Time { return now }

	if _, err := rpc.VerifyServiceAssertion(
		"/test.Service/Call",
		md,
		cfg.ServiceAssertion,
	); !errors.Is(err, rpc.ErrServiceAssertionInvalid) {
		t.Fatalf(
			"VerifyServiceAssertion() error = %v, want ErrServiceAssertionInvalid",
			err,
		)
	}
}
