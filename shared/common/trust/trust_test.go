package trust

import (
	"context"
	"errors"
	"testing"
	"time"

	"google.golang.org/grpc/metadata"
)

func testMetadataContract() MetadataContract {
	return NewMetadataContract(
		[]string{"profileid", "x-test-caller-kind"},
		[]string{"profileid", "x-test-caller-kind", "x-request-id"},
		[]string{"/grpc.health.v1.Health/Check", "/grpc.health.v1.Health/Watch"},
	)
}

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

	if err := Sign("/test.Service/Call", md, cfg, testMetadataContract()); err != nil {
		t.Fatalf("Sign() error = %v", err)
	}
	caller, err := Verify("/test.Service/Call", md, cfg, testMetadataContract())
	if err != nil {
		t.Fatalf("Verify() error = %v", err)
	}
	if caller.ServiceID != "gateway-service" || caller.Level != LevelVerified {
		t.Fatalf("caller = %+v", caller)
	}
}

func TestVerifyRejectsTamperedSignedMetadata(t *testing.T) {
	now := time.Unix(1_700_000_000, 0)
	cfg := Config{ServiceID: "gateway-service", Secret: "test-secret", Now: func() time.Time { return now }}
	md := metadata.Pairs("profileid", "12")
	if err := Sign("/test.Service/Call", md, cfg, testMetadataContract()); err != nil {
		t.Fatalf("Sign() error = %v", err)
	}
	md.Set("profileid", "99")

	_, err := Verify("/test.Service/Call", md, cfg, testMetadataContract())
	if !errors.Is(err, ErrInvalidAssertion) {
		t.Fatalf("Verify() error = %v, want ErrInvalidAssertion", err)
	}
}

func TestVerifyRejectsMethodReplay(t *testing.T) {
	now := time.Unix(1_700_000_000, 0)
	cfg := Config{ServiceID: "gateway-service", Secret: "test-secret", Now: func() time.Time { return now }}
	md := metadata.Pairs("profileid", "12")
	if err := Sign("/test.Service/One", md, cfg, testMetadataContract()); err != nil {
		t.Fatalf("Sign() error = %v", err)
	}

	_, err := Verify("/test.Service/Two", md, cfg, testMetadataContract())
	if !errors.Is(err, ErrInvalidAssertion) {
		t.Fatalf("Verify() error = %v, want ErrInvalidAssertion", err)
	}
}

func TestVerifyRejectsExpiredAndFutureAssertions(t *testing.T) {
	issued := time.Unix(1_700_000_000, 0)
	cfg := Config{
		ServiceID: "gateway-service",
		Secret:    "test-secret",
		MaxAge:    30 * time.Second,
		ClockSkew: 5 * time.Second,
		Now:       func() time.Time { return issued },
	}

	for name, verifyNow := range map[string]time.Time{
		"expired": issued.Add(time.Minute),
		"future":  issued.Add(-10 * time.Second),
	} {
		t.Run(name, func(t *testing.T) {
			md := metadata.Pairs("profileid", "12")
			cfg.Now = func() time.Time { return issued }
			if err := Sign("/test.Service/Call", md, cfg, testMetadataContract()); err != nil {
				t.Fatalf("Sign() error = %v", err)
			}
			cfg.Now = func() time.Time { return verifyNow }
			if _, err := Verify("/test.Service/Call", md, cfg, testMetadataContract()); !errors.Is(err, ErrExpiredAssertion) {
				t.Fatalf("Verify() error = %v, want ErrExpiredAssertion", err)
			}
		})
	}
}

func TestVerifyUsesCallerSpecificSecret(t *testing.T) {
	now := time.Unix(1_700_000_000, 0)
	signer := Config{ServiceID: "gateway-service", Secret: "gateway-secret", Now: func() time.Time { return now }}
	md := metadata.Pairs("profileid", "12")
	if err := Sign("/test.Service/Call", md, signer, testMetadataContract()); err != nil {
		t.Fatalf("Sign() error = %v", err)
	}

	verifier := Config{
		VerificationSecrets: map[string]string{
			"gateway-service": "gateway-secret",
			"worker-service":  "worker-secret",
		},
		Now: func() time.Time { return now },
	}
	if _, err := Verify("/test.Service/Call", md, verifier, testMetadataContract()); err != nil {
		t.Fatalf("Verify() error = %v", err)
	}

	md.Set(CallerMetadataKey, "worker-service")
	if _, err := Verify("/test.Service/Call", md, verifier, testMetadataContract()); !errors.Is(err, ErrInvalidAssertion) {
		t.Fatalf("impersonation error = %v, want ErrInvalidAssertion", err)
	}
}

func TestSameMethodReplayWithinWindowIsKnownLimitation(t *testing.T) {
	now := time.Unix(1_700_000_000, 0)
	cfg := Config{ServiceID: "gateway-service", Secret: "test-secret", Now: func() time.Time { return now }}
	md := metadata.Pairs("profileid", "42", "x-request-id", "req-123")
	if err := Sign("/test.Service/Call", md, cfg, testMetadataContract()); err != nil {
		t.Fatalf("Sign() error = %v", err)
	}

	for attempt := 0; attempt < 2; attempt++ {
		if _, err := Verify("/test.Service/Call", md, cfg, testMetadataContract()); err != nil {
			t.Fatalf("attempt %d Verify() error = %v", attempt+1, err)
		}
	}
}
