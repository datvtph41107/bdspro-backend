package admin

import (
	"context"
	"testing"
	"time"

	"common/callercontext"
	"common/identitycontext"
	"common/transporttrust"
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
	ctx, err := transporttrust.BindCaller(context.Background(), transporttrust.Caller{
		ServiceID: "gateway-service",
		Level:     transporttrust.LevelVerified,
		IssuedAt:  time.Now().UTC(),
	})
	if err != nil {
		t.Fatal(err)
	}
	ctx, err = callercontext.BindCaller(ctx, callercontext.Caller{Kind: callercontext.KindUser})
	if err != nil {
		t.Fatal(err)
	}
	ctx, err = identitycontext.BindActor(ctx, identitycontext.Actor{
		AuthID: 7, ProfileID: 42, SessionID: 8, Role: role, TokenType: tokenType,
	})
	if err != nil {
		t.Fatal(err)
	}
	return ctx
}
