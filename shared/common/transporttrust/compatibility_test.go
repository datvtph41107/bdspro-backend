package transporttrust

import (
	_trust "common/trust"

	"context"
	"errors"
	"testing"
	"time"

	"google.golang.org/grpc/metadata"
)

func TestCompatibilityContextUsesCanonicalTrustAuthority(t *testing.T) {
	caller := Caller{
		ServiceID: "gateway-service",
		Level:     LevelVerified,
		IssuedAt:  time.Unix(1_700_000_000, 0),
	}

	oldBound, err := BindCaller(context.Background(), caller)
	if err != nil {
		t.Fatal(err)
	}
	if got, ok := _trust.FromContext(oldBound); !ok || got != caller {
		t.Fatalf("canonical trust cannot read compatibility caller: %+v %v", got, ok)
	}

	newBound, err := _trust.BindCaller(context.Background(), caller)
	if err != nil {
		t.Fatal(err)
	}
	if got, ok := FromContext(newBound); !ok || got != caller {
		t.Fatalf("compatibility trust cannot read canonical caller: %+v %v", got, ok)
	}
}

func TestCompatibilitySignVerifyPreservesLegacyMetadataContract(t *testing.T) {
	now := time.Unix(1_700_000_000, 0)
	cfg := Config{ServiceID: "gateway-service", Secret: "test-secret", Now: func() time.Time { return now }}
	md := metadata.Pairs(
		"profileid", "12",
		"x-qhpro-caller-kind", "user",
		"x-request-id", "req-123",
	)

	if err := Sign("/test.Service/Call", md, cfg); err != nil {
		t.Fatalf("Sign() error = %v", err)
	}
	if _, err := Verify("/test.Service/Call", md, cfg); err != nil {
		t.Fatalf("Verify() error = %v", err)
	}

	md.Set("profileid", "99")
	if _, err := Verify("/test.Service/Call", md, cfg); !errors.Is(err, ErrInvalidAssertion) {
		t.Fatalf("tampered profile Verify() error = %v", err)
	}
}

func TestConfigFromEnvPreservesCompatibilitySurface(t *testing.T) {
	t.Setenv("QHPRO_SERVICE_ID", "gateway-service")
	t.Setenv("QHPRO_INTERNAL_METADATA_SECRET", "shared-secret")
	t.Setenv("QHPRO_TRUSTED_METADATA_MODE", "enforce")
	t.Setenv("QHPRO_INTERNAL_METADATA_MAX_AGE", "45s")
	t.Setenv("QHPRO_INTERNAL_METADATA_CLOCK_SKEW", "7s")
	t.Setenv(
		"QHPRO_INTERNAL_METADATA_VERIFICATION_KEYS",
		`{"gateway-service":"gateway-secret","invalid service":"ignored"}`,
	)

	cfg := ConfigFromEnv()
	if cfg.ServiceID != "gateway-service" || cfg.Secret != "shared-secret" || cfg.Mode != ModeEnforce {
		t.Fatalf("ConfigFromEnv() = %+v", cfg)
	}
	if cfg.MaxAge != 45*time.Second || cfg.ClockSkew != 7*time.Second {
		t.Fatalf("duration config = %+v", cfg)
	}
	if !cfg.VerificationKeysConfigured || cfg.VerificationSecrets["gateway-service"] != "gateway-secret" {
		t.Fatalf("verification config = %+v", cfg)
	}
	if _, exists := cfg.VerificationSecrets["invalid service"]; exists {
		t.Fatalf("invalid service ID was loaded: %#v", cfg.VerificationSecrets)
	}
}

func TestInvalidVerificationKeyringDoesNotFallBackToSharedSecret(t *testing.T) {
	now := time.Unix(1_700_000_000, 0)
	signer := Config{ServiceID: "gateway-service", Secret: "shared-secret", Now: func() time.Time { return now }}
	md := metadata.Pairs("profileid", "12")
	if err := Sign("/test.Service/Call", md, signer); err != nil {
		t.Fatal(err)
	}

	t.Setenv("QHPRO_INTERNAL_METADATA_VERIFICATION_KEYS", "not-json")
	t.Setenv("QHPRO_INTERNAL_METADATA_SECRET", "shared-secret")
	cfg := ConfigFromEnv()
	cfg.Now = func() time.Time { return now }

	if _, err := Verify("/test.Service/Call", md, cfg); !errors.Is(err, ErrInvalidAssertion) {
		t.Fatalf("Verify() error = %v, want ErrInvalidAssertion", err)
	}
}
