package transporttrust

import (
	_callercontext "common/callercontext"
	"context"
	"errors"
	"testing"
	"time"

	"google.golang.org/grpc/metadata"
)

func TestBindCallerIsImmutable(t *testing.T) {
	t.Parallel()

	issuedAt := time.Unix(1_700_000_000, 0)
	caller := Caller{ServiceID: "gateway-service", Level: LevelVerified, IssuedAt: issuedAt}
	ctx, err := BindCaller(context.Background(), caller)
	if err != nil {
		t.Fatalf("first BindCaller() error = %v", err)
	}
	if rebound, err := BindCaller(ctx, caller); err != nil || rebound != ctx {
		t.Fatalf("same caller bind = %v, %v", rebound, err)
	}

	unchanged, err := BindCaller(ctx, Caller{
		ServiceID: "worker-service",
		Level:     LevelVerified,
		IssuedAt:  issuedAt,
	})
	if !errors.Is(err, ErrCallerConflict) {
		t.Fatalf("conflicting BindCaller() error = %v", err)
	}
	if unchanged != ctx {
		t.Fatal("conflicting transport caller changed context")
	}
}

func TestSignAndVerify(t *testing.T) {
	now := time.Unix(1_700_000_000, 0)
	cfg := Config{
		ServiceID: "gateway-service",
		Secret:    "test-secret",
		Mode:      ModeEnforce,
		MaxAge:    30 * time.Second,
		ClockSkew: 5 * time.Second,
		Now:       func() time.Time { return now },
	}
	md := metadata.Pairs("profileid", "12", "x-request-id", "req-123")

	if err := Sign("/test.Service/Call", md, cfg); err != nil {
		t.Fatalf("Sign() error = %v", err)
	}
	caller, err := Verify("/test.Service/Call", md, cfg)
	if err != nil {
		t.Fatalf("Verify() error = %v", err)
	}
	if caller.ServiceID != "gateway-service" || caller.Level != LevelVerified {
		t.Fatalf("caller = %+v", caller)
	}
}

func TestVerifyRejectsTamperedIdentity(t *testing.T) {
	now := time.Unix(1_700_000_000, 0)
	cfg := Config{ServiceID: "gateway-service", Secret: "test-secret", Now: func() time.Time { return now }}
	md := metadata.Pairs("profileid", "12")
	if err := Sign("/test.Service/Call", md, cfg); err != nil {
		t.Fatalf("Sign() error = %v", err)
	}
	md.Set("profileid", "99")

	_, err := Verify("/test.Service/Call", md, cfg)
	if !errors.Is(err, ErrInvalidAssertion) {
		t.Fatalf("Verify() error = %v, want ErrInvalidAssertion", err)
	}
}

func TestVerifyRejectsMethodReplay(t *testing.T) {
	now := time.Unix(1_700_000_000, 0)
	cfg := Config{ServiceID: "gateway-service", Secret: "test-secret", Now: func() time.Time { return now }}
	md := metadata.Pairs("profileid", "12")
	if err := Sign("/test.Service/One", md, cfg); err != nil {
		t.Fatalf("Sign() error = %v", err)
	}

	_, err := Verify("/test.Service/Two", md, cfg)
	if !errors.Is(err, ErrInvalidAssertion) {
		t.Fatalf("Verify() error = %v, want ErrInvalidAssertion", err)
	}
}

func TestVerifyRejectsExpiredAssertion(t *testing.T) {
	issued := time.Unix(1_700_000_000, 0)
	cfg := Config{ServiceID: "gateway-service", Secret: "test-secret", MaxAge: 30 * time.Second, Now: func() time.Time { return issued }}
	md := metadata.Pairs("profileid", "12")
	if err := Sign("/test.Service/Call", md, cfg); err != nil {
		t.Fatalf("Sign() error = %v", err)
	}
	cfg.Now = func() time.Time { return issued.Add(time.Minute) }

	_, err := Verify("/test.Service/Call", md, cfg)
	if !errors.Is(err, ErrExpiredAssertion) {
		t.Fatalf("Verify() error = %v, want ErrExpiredAssertion", err)
	}
}

func TestVerifyRejectsTamperedCallerClassification(t *testing.T) {
	now := time.Unix(1_700_000_000, 0)
	cfg := Config{ServiceID: "gateway-service", Secret: "test-secret", Now: func() time.Time { return now }}
	md := metadata.Pairs(_callercontext.CallerKindMetadataKey, string(_callercontext.KindAnonymous))
	if err := Sign("/test.Service/Call", md, cfg); err != nil {
		t.Fatalf("Sign() error = %v", err)
	}
	md.Set(_callercontext.CallerKindMetadataKey, string(_callercontext.KindUser))

	_, err := Verify("/test.Service/Call", md, cfg)
	if !errors.Is(err, ErrInvalidAssertion) {
		t.Fatalf("Verify() error = %v, want ErrInvalidAssertion", err)
	}
}

func TestVerifyRejectsTamperedOperationSelection(t *testing.T) {
	now := time.Unix(1_700_000_000, 0)
	cfg := Config{ServiceID: "gateway-service", Secret: "test-secret", Now: func() time.Time { return now }}
	md := metadata.Pairs(
		"x-qhpro-operation-code", "parcel.lookup",
		"x-qhpro-operation-binding", "parcel.lookup_transport",
		"x-qhpro-operation-source", "direct_binding",
	)
	if err := Sign("/test.Service/Call", md, cfg); err != nil {
		t.Fatalf("Sign() error = %v", err)
	}
	md.Set("x-qhpro-operation-code", "parcel.report_generate")

	_, err := Verify("/test.Service/Call", md, cfg)
	if !errors.Is(err, ErrInvalidAssertion) {
		t.Fatalf("Verify() error = %v, want ErrInvalidAssertion", err)
	}
}

func TestVerifyUsesCallerSpecificSecret(t *testing.T) {
	now := time.Unix(1_700_000_000, 0)
	signer := Config{
		ServiceID: "gateway-service",
		Secret:    "gateway-secret",
		Now:       func() time.Time { return now },
	}
	md := metadata.Pairs("profileid", "12")
	if err := Sign("/test.Service/Call", md, signer); err != nil {
		t.Fatalf("Sign() error = %v", err)
	}

	verifier := Config{
		VerificationSecrets: map[string]string{
			"gateway-service": "gateway-secret",
			"worker-service":  "worker-secret",
		},
		Now: func() time.Time { return now },
	}
	if _, err := Verify("/test.Service/Call", md, verifier); err != nil {
		t.Fatalf("Verify() error = %v", err)
	}

	md.Set(CallerMetadataKey, "worker-service")
	if _, err := Verify("/test.Service/Call", md, verifier); !errors.Is(err, ErrInvalidAssertion) {
		t.Fatalf("impersonation error = %v, want ErrInvalidAssertion", err)
	}
}

func TestConfigLoadsPerServiceVerificationKeys(t *testing.T) {
	t.Setenv(
		"QHPRO_INTERNAL_METADATA_VERIFICATION_KEYS",
		`{"gateway-service":"gateway-secret","invalid service":"ignored"}`,
	)

	cfg := ConfigFromEnv()
	if !cfg.VerificationKeysConfigured {
		t.Fatal("verification key configuration was not detected")
	}
	if cfg.VerificationSecrets["gateway-service"] != "gateway-secret" {
		t.Fatalf("verification keys = %#v", cfg.VerificationSecrets)
	}
	if _, exists := cfg.VerificationSecrets["invalid service"]; exists {
		t.Fatalf("invalid service ID was loaded: %#v", cfg.VerificationSecrets)
	}
}

func TestInvalidVerificationKeyringDoesNotFallBackToSharedSecret(t *testing.T) {
	t.Setenv("QHPRO_INTERNAL_METADATA_VERIFICATION_KEYS", "not-json")
	t.Setenv("QHPRO_INTERNAL_METADATA_SECRET", "shared-secret")

	cfg := ConfigFromEnv()
	_, err := secretForCaller(cfg, "gateway-service")
	if !errors.Is(err, ErrInvalidAssertion) {
		t.Fatalf("secretForCaller() error = %v, want ErrInvalidAssertion", err)
	}
}
