package useradmin

import (
	"context"
	"errors"
	"testing"

	"common/identity"
)

type authorizerStub struct {
	allowed bool
	err     error
	actorID uint64
	code    string
}

func (s *authorizerStub) HasPermission(_ context.Context, actorID uint64, code string) (bool, error) {
	s.actorID = actorID
	s.code = code
	return s.allowed, s.err
}

func TestRequirePermissionUsesDurableAssignmentNotRoleName(t *testing.T) {
	allowed := &authorizerStub{allowed: true}
	actorID, err := RequirePermission(trustedContext(t, "ROLE_SUPPORT"), allowed, PermissionView)
	if err != nil || actorID != 42 {
		t.Fatalf("explicit permission rejected: actor=%d err=%v", actorID, err)
	}
	if allowed.actorID != 42 || allowed.code != PermissionView {
		t.Fatalf("authorizer call = actor:%d code:%q", allowed.actorID, allowed.code)
	}

	_, err = RequirePermission(trustedContext(t, "ROLE_SUPER_ADMIN"), &authorizerStub{}, PermissionManage)
	if !errors.Is(err, ErrPermissionDenied) {
		t.Fatalf("role-name bypass error = %v", err)
	}
}

func TestRequirePermissionFailsClosed(t *testing.T) {
	ctx := trustedContext(t, "ROLE_ADMIN")
	if _, err := RequirePermission(ctx, nil, PermissionView); !errors.Is(err, ErrPermissionAuthorityUnavailable) {
		t.Fatalf("nil authority error = %v", err)
	}
	if _, err := RequirePermission(ctx, &authorizerStub{err: errors.New("database unavailable")}, PermissionView); !errors.Is(err, ErrPermissionAuthorityUnavailable) {
		t.Fatalf("authority failure error = %v", err)
	}
	if _, err := RequirePermission(context.Background(), &authorizerStub{allowed: true}, PermissionView); !errors.Is(err, ErrActorRequired) {
		t.Fatalf("unsigned actor error = %v", err)
	}
}

func trustedContext(t *testing.T, role string) context.Context {
	t.Helper()
	ctx, err := identity.BindServiceCaller(context.Background(), identity.ServiceCaller{ServiceID: "gateway-service"})
	if err != nil {
		t.Fatal(err)
	}
	ctx, err = identity.BindCaller(ctx, identity.Caller{Kind: identity.CallerUser})
	if err != nil {
		t.Fatal(err)
	}
	ctx, err = identity.BindActor(ctx, identity.Actor{ProfileID: 42, AuthID: 7, TokenType: "ACCESS", Role: role})
	if err != nil {
		t.Fatal(err)
	}
	return ctx
}
