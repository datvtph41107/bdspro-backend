package identitygrpc

import (
	"common/callercontext"
	"common/identitycontext"
	"context"
	"testing"

	"google.golang.org/grpc/metadata"
)

func TestAppendAndRestoreUserIdentity(t *testing.T) {
	ctx, err := identitycontext.BindActor(context.Background(), identitycontext.Actor{
		AuthID:    10,
		ProfileID: 42,
		OriginID:  20,
		SessionID: 30,
		Role:      "USER",
		TokenType: "access",
	})
	if err != nil {
		t.Fatal(err)
	}
	ctx, err = callercontext.BindCaller(ctx, callercontext.Caller{Kind: callercontext.KindUser})
	if err != nil {
		t.Fatal(err)
	}

	md := AppendFromContext(ctx, metadata.MD{})
	restored, err := RestoreFromMetadata(context.Background(), md)
	if err != nil {
		t.Fatal(err)
	}
	actor, ok := identitycontext.FromContext(restored)
	if !ok || actor.ProfileID != 42 {
		t.Fatalf("actor = %#v, ok=%v", actor, ok)
	}
	caller, ok := callercontext.FromContext(restored)
	if !ok || caller.Kind != callercontext.KindUser {
		t.Fatalf("caller = %#v, ok=%v", caller, ok)
	}
}

func TestRestoreRejectsActorWithoutCaller(t *testing.T) {
	_, err := RestoreFromMetadata(context.Background(), metadata.Pairs(
		identitycontext.ProfileIDMetadataKey, "42",
	))
	if err == nil {
		t.Fatal("expected error")
	}
}
