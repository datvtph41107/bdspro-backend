package admin

import (
	"context"
	"errors"
	"testing"

	"common/identity"
)

type permissionAuthorizerStub struct {
	allowed bool
	err     error
	actorID uint64
	code    string
}

func (s *permissionAuthorizerStub) HasPermission(
	_ context.Context,
	actorID uint64,
	permissionCode string,
) (bool, error) {
	s.actorID = actorID
	s.code = permissionCode
	return s.allowed, s.err
}

func TestRequirePermissionUsesAssignmentNotRoleName(t *testing.T) {
	allowed := &permissionAuthorizerStub{allowed: true}
	actorID, err := RequirePermission(
		trustedActorContext(t, "ROLE_SUPPORT"),
		allowed,
		PermissionPlanView,
	)
	if err != nil || actorID != 42 {
		t.Fatalf("explicit permission rejected: actor=%d err=%v", actorID, err)
	}
	if allowed.actorID != 42 || allowed.code != PermissionPlanView {
		t.Fatalf("authorizer call = actor:%d code:%q", allowed.actorID, allowed.code)
	}

	denied := &permissionAuthorizerStub{}
	_, err = RequirePermission(
		trustedActorContext(t, "ROLE_SUPER_ADMIN"),
		denied,
		PermissionPlanView,
	)
	if !errors.Is(err, ErrPermissionDenied) {
		t.Fatalf("role-name bypass error = %v", err)
	}
}

func TestRequirePermissionFailsClosed(t *testing.T) {
	ctx := trustedActorContext(t, "ROLE_ADMIN")
	if _, err := RequirePermission(ctx, nil, PermissionPlanView); !errors.Is(err, ErrPermissionAuthorityUnavailable) {
		t.Fatalf("nil authority error = %v", err)
	}
	if _, err := RequirePermission(ctx, &permissionAuthorizerStub{err: errors.New("db unavailable")}, PermissionPlanView); !errors.Is(err, ErrPermissionAuthorityUnavailable) {
		t.Fatalf("authority failure error = %v", err)
	}
}

func TestRequirePermissionStillRequiresTrustedActor(t *testing.T) {
	if _, err := RequirePermission(context.Background(), &permissionAuthorizerStub{allowed: true}, PermissionPlanView); !errors.Is(err, ErrUnauthorized) {
		t.Fatalf("unsigned actor error = %v", err)
	}
}

func trustedActorContext(t *testing.T, role string) context.Context {
	t.Helper()
	ctx, err := identity.BindServiceCaller(context.Background(), identity.ServiceCaller{ServiceID: "gateway-service"})
	if err != nil {
		t.Fatal(err)
	}
	ctx, err = identity.BindCaller(ctx, identity.Caller{Kind: identity.CallerUser})
	if err != nil {
		t.Fatal(err)
	}
	ctx, err = identity.BindActor(ctx, identity.Actor{
		AuthID: 7, ProfileID: 42, SessionID: 8, Role: role, TokenType: "ACCESS",
	})
	if err != nil {
		t.Fatal(err)
	}
	return ctx
}
