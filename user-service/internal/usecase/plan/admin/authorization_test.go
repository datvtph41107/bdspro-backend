package admin

import (
	"context"
	"testing"

	"common/identity"
)

func TestActorIDFromContextRequiresTrustedUserAccessIdentityNotRole(t *testing.T) {
	for _, role := range []string{"ROLE_ADMIN", "ROLE_USER", "ROLE_SUPPORT", ""} {
		ctx := trustedIdentityContext(t, role, "ACCESS")
		actorID, err := ActorIDFromContext(ctx)
		if err != nil || actorID != 42 {
			t.Fatalf("role=%q actorID=%d err=%v", role, actorID, err)
		}
	}
	if _, err := ActorIDFromContext(context.Background()); err == nil {
		t.Fatal("unsigned context must not control catalog")
	}
	if _, err := ActorIDFromContext(trustedIdentityContext(t, "ROLE_ADMIN", "REFRESH")); err == nil {
		t.Fatal("refresh token identity must not control catalog")
	}
}

func trustedIdentityContext(t *testing.T, role, tokenType string) context.Context {
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
		AuthID: 7, ProfileID: 42, SessionID: 8, Role: role, TokenType: tokenType,
	})
	if err != nil {
		t.Fatal(err)
	}
	return ctx
}
