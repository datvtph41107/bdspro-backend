package rpc

import (
	"errors"
	"testing"
	"time"

	"google.golang.org/grpc/metadata"
)

func TestServiceAssertionSignAndVerify(t *testing.T) {
	now := time.Unix(1_700_000_000, 0)
	cfg := ServiceAssertionConfig{
		ServiceID: "gateway-service",
		Secret:    "test-secret",
		MaxAge:    30 * time.Second,
		ClockSkew: 5 * time.Second,
		Now:       func() time.Time { return now },
	}
	md := metadata.Pairs(
		ProfileIDMetadataKey, "12",
		RequestIDMetadataKey, "req-123",
	)

	if err := SignServiceAssertion(
		"/test.Service/Call",
		md,
		cfg,
	); err != nil {
		t.Fatalf("SignServiceAssertion() error = %v", err)
	}

	caller, err := VerifyServiceAssertion(
		"/test.Service/Call",
		md,
		cfg,
	)
	if err != nil {
		t.Fatalf("VerifyServiceAssertion() error = %v", err)
	}
	if caller.ServiceID != "gateway-service" {
		t.Fatalf("caller = %+v", caller)
	}
}

func TestServiceAssertionRejectsTamperedSignedMetadata(t *testing.T) {
	now := time.Unix(1_700_000_000, 0)
	cfg := ServiceAssertionConfig{
		ServiceID: "gateway-service",
		Secret:    "test-secret",
		Now:       func() time.Time { return now },
	}

	md := metadata.Pairs(ProfileIDMetadataKey, "12")
	if err := SignServiceAssertion("/test.Service/Call", md, cfg); err != nil {
		t.Fatal(err)
	}

	md.Set(ProfileIDMetadataKey, "99")

	if _, err := VerifyServiceAssertion(
		"/test.Service/Call",
		md,
		cfg,
	); !errors.Is(err, ErrServiceAssertionInvalid) {
		t.Fatalf(
			"VerifyServiceAssertion() error = %v, want ErrServiceAssertionInvalid",
			err,
		)
	}
}

func TestServiceAssertionRejectsMethodReplay(t *testing.T) {
	now := time.Unix(1_700_000_000, 0)
	cfg := ServiceAssertionConfig{
		ServiceID: "gateway-service",
		Secret:    "test-secret",
		Now:       func() time.Time { return now },
	}

	md := metadata.Pairs(ProfileIDMetadataKey, "12")
	if err := SignServiceAssertion("/test.Service/One", md, cfg); err != nil {
		t.Fatal(err)
	}

	if _, err := VerifyServiceAssertion(
		"/test.Service/Two",
		md,
		cfg,
	); !errors.Is(err, ErrServiceAssertionInvalid) {
		t.Fatalf(
			"VerifyServiceAssertion() error = %v, want ErrServiceAssertionInvalid",
			err,
		)
	}
}

func TestServiceAssertionRejectsExpiredAndFutureAssertions(t *testing.T) {
	issued := time.Unix(1_700_000_000, 0)

	cfg := ServiceAssertionConfig{
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
			md := metadata.Pairs(ProfileIDMetadataKey, "12")

			cfg.Now = func() time.Time { return issued }
			if err := SignServiceAssertion(
				"/test.Service/Call",
				md,
				cfg,
			); err != nil {
				t.Fatal(err)
			}

			cfg.Now = func() time.Time { return verifyNow }
			if _, err := VerifyServiceAssertion(
				"/test.Service/Call",
				md,
				cfg,
			); !errors.Is(err, ErrServiceAssertionExpired) {
				t.Fatalf(
					"VerifyServiceAssertion() error = %v, want ErrServiceAssertionExpired",
					err,
				)
			}
		})
	}
}

func TestServiceAssertionUsesCallerSpecificSecret(t *testing.T) {
	now := time.Unix(1_700_000_000, 0)

	signer := ServiceAssertionConfig{
		ServiceID: "gateway-service",
		Secret:    "gateway-secret",
		Now:       func() time.Time { return now },
	}

	md := metadata.Pairs(ProfileIDMetadataKey, "12")
	if err := SignServiceAssertion(
		"/test.Service/Call",
		md,
		signer,
	); err != nil {
		t.Fatal(err)
	}

	verifier := ServiceAssertionConfig{
		VerificationSecrets: map[string]string{
			"gateway-service": "gateway-secret",
			"worker-service":  "worker-secret",
		},
		Now: func() time.Time { return now },
	}

	if _, err := VerifyServiceAssertion(
		"/test.Service/Call",
		md,
		verifier,
	); err != nil {
		t.Fatalf("VerifyServiceAssertion() error = %v", err)
	}

	md.Set(ServiceAssertionCallerMetadataKey, "worker-service")

	if _, err := VerifyServiceAssertion(
		"/test.Service/Call",
		md,
		verifier,
	); !errors.Is(err, ErrServiceAssertionInvalid) {
		t.Fatalf(
			"impersonation error = %v, want ErrServiceAssertionInvalid",
			err,
		)
	}
}

func TestServiceAssertionSameMethodReplayWithinWindowIsKnownLimitation(
	t *testing.T,
) {
	now := time.Unix(1_700_000_000, 0)

	cfg := ServiceAssertionConfig{
		ServiceID: "gateway-service",
		Secret:    "test-secret",
		Now:       func() time.Time { return now },
	}

	md := metadata.Pairs(
		ProfileIDMetadataKey, "42",
		RequestIDMetadataKey, "req-123",
	)
	if err := SignServiceAssertion(
		"/test.Service/Call",
		md,
		cfg,
	); err != nil {
		t.Fatal(err)
	}

	for attempt := 0; attempt < 2; attempt++ {
		if _, err := VerifyServiceAssertion(
			"/test.Service/Call",
			md,
			cfg,
		); err != nil {
			t.Fatalf(
				"attempt %d VerifyServiceAssertion() error = %v",
				attempt+1,
				err,
			)
		}
	}
}
