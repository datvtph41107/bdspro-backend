package transporttrust

import (
	_callercontext "common/callercontext"
	"errors"
	"testing"
	"time"

	"google.golang.org/grpc/metadata"
)

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
