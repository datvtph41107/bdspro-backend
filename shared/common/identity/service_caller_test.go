package identity

import (
	"context"
	"errors"
	"testing"
)

func TestServiceCallerContextRoundTrip(t *testing.T) {
	t.Parallel()

	want := ServiceCaller{ServiceID: "payment-service"}
	ctx, err := BindServiceCaller(context.Background(), want)
	if err != nil {
		t.Fatalf("BindServiceCaller() error = %v", err)
	}

	got, ok := ServiceCallerFromContext(ctx)
	if !ok || got != want {
		t.Fatalf("ServiceCallerFromContext() = %#v, %v; want %#v, true", got, ok, want)
	}
}

func TestBindServiceCallerIsIdempotentAndRejectsConflict(t *testing.T) {
	t.Parallel()

	first := ServiceCaller{ServiceID: "payment-service"}
	ctx, err := BindServiceCaller(context.Background(), first)
	if err != nil {
		t.Fatalf("first BindServiceCaller() error = %v", err)
	}

	rebound, err := BindServiceCaller(ctx, first)
	if err != nil || rebound != ctx {
		t.Fatalf("same service caller bind = %v, %v", rebound, err)
	}

	unchanged, err := BindServiceCaller(ctx, ServiceCaller{ServiceID: "user-service"})
	if !errors.Is(err, ErrServiceCallerContextConflict) {
		t.Fatalf("conflicting BindServiceCaller() error = %v", err)
	}
	if unchanged != ctx {
		t.Fatal("conflicting service caller changed context")
	}
}

func TestServiceCallerValidation(t *testing.T) {
	t.Parallel()

	for _, serviceID := range []string{
		"gateway-service",
		"user_service",
		"payment.service",
		"service:v1",
		"Service123",
	} {
		if !(ServiceCaller{ServiceID: serviceID}).IsValid() {
			t.Fatalf("service ID should be valid: %q", serviceID)
		}
	}

	for _, serviceID := range []string{
		"",
		" user-service",
		"user/service",
		"user service",
	} {
		if (ServiceCaller{ServiceID: serviceID}).IsValid() {
			t.Fatalf("service ID should be invalid: %q", serviceID)
		}
	}
}
