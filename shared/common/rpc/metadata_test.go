package rpc

import (
	"common/identity"
	"common/request"
	"context"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func TestRequestMetadataRoundTrip(t *testing.T) {
	ctx := request.WithRequestID(context.Background(), "req_123")
	var err error
	ctx, err = request.BindOperationID(ctx, "op_123")
	if err != nil {
		t.Fatal(err)
	}
	ctx, err = request.BindIdempotencyKey(ctx, "idem_123")
	if err != nil {
		t.Fatal(err)
	}

	md := AppendRequestFromContext(ctx, metadata.MD{})
	restored, err := RestoreRequestFromMetadata(context.Background(), md)
	if err != nil {
		t.Fatal(err)
	}
	if got, ok := request.RequestIDFromContext(restored); !ok || got != "req_123" {
		t.Fatalf("request ID = %q, %v", got, ok)
	}
	if got, ok := request.OperationIDFromContext(restored); !ok || got != "op_123" {
		t.Fatalf("operation ID = %q, %v", got, ok)
	}
	if got, ok := request.IdempotencyKeyFromContext(restored); !ok || got != "idem_123" {
		t.Fatalf("idempotency key = %q, %v", got, ok)
	}
}

func TestRestoreRequestRejectsDuplicateOperationID(t *testing.T) {
	_, err := RestoreRequestFromMetadata(context.Background(), metadata.MD{
		OperationIDMetadataKey: []string{"op_a", "op_b"},
	})
	if status.Code(err) != codes.InvalidArgument {
		t.Fatalf("code = %v", status.Code(err))
	}
}

func TestIdentityMetadataRoundTrip(t *testing.T) {
	ctx, err := identity.BindActor(context.Background(), identity.Actor{
		AuthID: 10, ProfileID: 42, OriginID: 20, SessionID: 30,
		Role: "USER", TokenType: "access",
	})
	if err != nil {
		t.Fatal(err)
	}
	ctx, err = identity.BindCaller(ctx, identity.Caller{Kind: identity.CallerUser})
	if err != nil {
		t.Fatal(err)
	}

	md := AppendIdentityFromContext(ctx, metadata.MD{})
	restored, err := RestoreIdentityFromMetadata(context.Background(), md, true)
	if err != nil {
		t.Fatal(err)
	}
	actor, ok := identity.ActorFromContext(restored)
	if !ok || actor.ProfileID != 42 {
		t.Fatalf("actor = %#v, %v", actor, ok)
	}
	caller, ok := identity.CallerFromContext(restored)
	if !ok || caller.Kind != identity.CallerUser {
		t.Fatalf("caller = %#v, %v", caller, ok)
	}
}

func TestAppendIdentityDefaultsLocalAndActorCallers(t *testing.T) {
	md := AppendIdentityFromContext(context.Background(), metadata.MD{})
	if got := md.Get(CallerKindMetadataKey); len(got) != 1 || got[0] != string(identity.CallerInternalService) {
		t.Fatalf("local caller = %v", got)
	}

	ctx, err := identity.BindActor(context.Background(), identity.Actor{ProfileID: 42})
	if err != nil {
		t.Fatal(err)
	}
	md = AppendIdentityFromContext(ctx, metadata.MD{})
	if got := md.Get(CallerKindMetadataKey); len(got) != 1 || got[0] != string(identity.CallerUser) {
		t.Fatalf("actor caller = %v", got)
	}
}

func TestRestoreIdentityValidatesWithoutPromoting(t *testing.T) {
	md := metadata.Pairs(
		CallerKindMetadataKey, string(identity.CallerUser),
		ProfileIDMetadataKey, "42",
	)
	ctx, err := RestoreIdentityFromMetadata(context.Background(), md, false)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := identity.ActorFromContext(ctx); ok {
		t.Fatal("unverified actor was promoted")
	}
	if _, ok := identity.CallerFromContext(ctx); ok {
		t.Fatal("unverified caller was promoted")
	}
}

func TestRestoreIdentityRejectsActorWithNonUserCaller(t *testing.T) {
	_, err := RestoreIdentityFromMetadata(context.Background(), metadata.Pairs(
		CallerKindMetadataKey, string(identity.CallerInternalService),
		ProfileIDMetadataKey, "42",
	), true)
	if status.Code(err) != codes.Unauthenticated {
		t.Fatalf("code = %v, err = %v", status.Code(err), err)
	}
}

func TestCallerMetadataCodecRoundTrip(t *testing.T) {
	t.Parallel()

	want := identity.Caller{
		Kind:      identity.CallerUser,
		APIKeyID:  8,
		APIKeyApp: "partner-app",
	}

	md := metadata.MD{}
	AppendCallerMetadata(md, want)

	got, ok, err := CallerFromMetadata(md)
	if err != nil || !ok || got != want {
		t.Fatalf(
			"CallerFromMetadata() = %#v, %v, %v; want %#v, true, nil",
			got,
			ok,
			err,
			want,
		)
	}
}

func TestCallerMetadataCodecRejectsIncompleteAPIKeyEvidence(t *testing.T) {
	t.Parallel()

	md := metadata.Pairs(
		CallerKindMetadataKey,
		string(identity.CallerAnonymous),
		APIKeyIDMetadataKey,
		"42",
	)

	if _, _, err := CallerFromMetadata(md); err == nil {
		t.Fatal("incomplete caller metadata was accepted")
	}
}

func TestCallerMetadataCodecRejectsNonCanonicalText(t *testing.T) {
	t.Parallel()

	md := metadata.Pairs(
		CallerKindMetadataKey,
		string(identity.CallerAPIKey),
		APIKeyIDMetadataKey,
		"42",
		APIKeyAppMetadataKey,
		" partner-app ",
	)

	if _, _, err := CallerFromMetadata(md); err == nil {
		t.Fatal("caller metadata with surrounding whitespace was accepted")
	}
}
