package identity

import (
	"context"
	"errors"
	"testing"
)

func TestCallerContextRoundTrip(t *testing.T) {
	t.Parallel()

	want := Caller{Kind: CallerAPIKey, APIKeyID: 42, APIKeyApp: "release-uploader"}
	ctx, err := BindCaller(context.Background(), want)
	if err != nil {
		t.Fatalf("BindCaller() error = %v", err)
	}
	got, ok := CallerFromContext(ctx)
	if !ok || got != want {
		t.Fatalf("CallerFromContext() = %#v, %v; want %#v, true", got, ok, want)
	}
}

func TestBindCallerRejectsConflictAndAllowsSameValue(t *testing.T) {
	t.Parallel()

	first := Caller{Kind: CallerUser}
	ctx, err := BindCaller(context.Background(), first)
	if err != nil {
		t.Fatalf("first BindCaller() error = %v", err)
	}
	if rebound, err := BindCaller(ctx, first); err != nil || rebound != ctx {
		t.Fatalf("same caller bind = %v, %v", rebound, err)
	}

	unchanged, err := BindCaller(ctx, Caller{Kind: CallerAnonymous})
	if !errors.Is(err, ErrCallerContextConflict) {
		t.Fatalf("conflicting BindCaller() error = %v", err)
	}
	if unchanged != ctx {
		t.Fatal("conflicting caller changed context")
	}
}

func TestCallerValidation(t *testing.T) {
	t.Parallel()

	valid := []Caller{
		{Kind: CallerAnonymous},
		{Kind: CallerUser},
		{Kind: CallerUser, APIKeyID: 8, APIKeyApp: "partner-app"},
		{Kind: CallerAPIKey, APIKeyID: 42, APIKeyApp: "release-uploader"},
		{Kind: CallerInternalService},
	}
	for _, caller := range valid {
		if !caller.IsValid() {
			t.Fatalf("caller should be valid: %+v", caller)
		}
	}

	invalid := []Caller{
		{},
		{Kind: CallerAPIKey},
		{Kind: CallerAPIKey, APIKeyID: 42, APIKeyApp: " partner-app "},
		{Kind: CallerAnonymous, APIKeyID: 42, APIKeyApp: "partner-app"},
		{Kind: CallerInternalService, APIKeyID: 42, APIKeyApp: "partner-app"},
	}
	for _, caller := range invalid {
		if caller.IsValid() {
			t.Fatalf("caller should be invalid: %+v", caller)
		}
	}
}

func TestBindCallerRejectsNilContext(t *testing.T) {
	t.Parallel()

	ctx, err := BindCaller(nil, Caller{Kind: CallerUser})
	if err == nil || ctx != nil {
		t.Fatalf("BindCaller(nil) = %v, %v", ctx, err)
	}
}
