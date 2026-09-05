package client

import (
	"common/identity"
	"context"
	"errors"
	"testing"

	"tqd/internal/enums"
)

type permissionCheckerProbe struct {
	calls [][]string
	err   error
}

func (p *permissionCheckerProbe) HasPermissions(_ context.Context, codes []string) error {
	p.calls = append(p.calls, append([]string(nil), codes...))
	return p.err
}

func actorContext(t *testing.T, profileID uint64) context.Context {
	t.Helper()
	ctx, err := identity.BindActor(context.Background(), identity.Actor{ProfileID: profileID})
	if err != nil {
		t.Fatalf("BindActor() error = %v", err)
	}
	return ctx
}

func TestPermissionClientDelegatesDecisionToAuthByCode(t *testing.T) {
	probe := &permissionCheckerProbe{}
	client := &PermissionClient{auth: probe}

	allowed, err := client.CheckPermission(actorContext(t, 42), 42, enums.PermissionDirectorySupplierRead)
	if err != nil {
		t.Fatalf("CheckPermission() error = %v", err)
	}
	if !allowed {
		t.Fatal("CheckPermission() = false, want true after Auth allows")
	}
	if len(probe.calls) != 1 || len(probe.calls[0]) != 1 || probe.calls[0][0] != "directory_supplier_read" {
		t.Fatalf("Auth permission calls = %#v", probe.calls)
	}
}

func TestPermissionClientNeverConvertsAuthFailureToAllow(t *testing.T) {
	want := errors.New("auth denied")
	client := &PermissionClient{auth: &permissionCheckerProbe{err: want}}

	allowed, err := client.CheckPermission(actorContext(t, 42), 42, enums.PermissionPOIRead)
	if allowed {
		t.Fatal("CheckPermission() allowed an Auth failure")
	}
	if !errors.Is(err, want) {
		t.Fatalf("CheckPermission() error = %v, want Auth error", err)
	}
}

func TestPermissionClientRejectsSubjectMismatchBeforeAuth(t *testing.T) {
	probe := &permissionCheckerProbe{}
	client := &PermissionClient{auth: probe}

	allowed, err := client.CheckPermission(actorContext(t, 42), 99, enums.PermissionPOIRead)
	if allowed || !errors.Is(err, ErrAuthorizationSubjectMismatch) {
		t.Fatalf("CheckPermission() = %v, %v", allowed, err)
	}
	if len(probe.calls) != 0 {
		t.Fatalf("Auth was called on subject mismatch: %#v", probe.calls)
	}
}

func TestPermissionClientUnsupportedHistoricalOperationsFailClosed(t *testing.T) {
	client := &PermissionClient{auth: &permissionCheckerProbe{}}
	allowed, err := client.UserCanAccessTarget(actorContext(t, 42), 42, 7, "poi")
	if allowed || !errors.Is(err, ErrAuthorizationUnsupported) {
		t.Fatalf("UserCanAccessTarget() = %v, %v", allowed, err)
	}
}
