package callercontext

import (
	"context"
	"errors"
	"testing"

	"google.golang.org/grpc/metadata"
)

func TestCallerRoundTrip(t *testing.T) {
	t.Parallel()

	want := Caller{Kind: KindAPIKey, APIKeyID: 42, APIKeyApp: "release-uploader"}
	ctx, err := BindCaller(context.Background(), want)
	if err != nil {
		t.Fatalf("BindCaller() error = %v", err)
	}
	got, ok := FromContext(ctx)
	if !ok || got != want {
		t.Fatalf("FromContext() = %#v, %v; want %#v, true", got, ok, want)
	}
}

func TestBindCallerRejectsConflictAndAllowsSameValue(t *testing.T) {
	t.Parallel()

	first := Caller{Kind: KindUser}
	ctx, err := BindCaller(context.Background(), first)
	if err != nil {
		t.Fatalf("first BindCaller() error = %v", err)
	}
	if rebound, err := BindCaller(ctx, first); err != nil || rebound != ctx {
		t.Fatalf("same caller bind = %v, %v", rebound, err)
	}

	unchanged, err := BindCaller(ctx, Caller{Kind: KindAnonymous})
	if !errors.Is(err, ErrCallerContextConflict) {
		t.Fatalf("conflicting BindCaller() error = %v", err)
	}
	if unchanged != ctx {
		t.Fatal("conflicting caller changed context")
	}
}

func TestCallerMetadataRoundTrip(t *testing.T) {
	t.Parallel()

	md := metadata.MD{}
	want := Caller{Kind: KindUser}
	AppendMetadata(md, want)
	got, ok, err := FromMetadata(md)
	if err != nil || !ok || got != want {
		t.Fatalf("FromMetadata() = %#v, %v, %v; want %#v, true, nil", got, ok, err, want)
	}
}

func TestAPIKeyCallerRequiresIdentity(t *testing.T) {
	t.Parallel()

	if (Caller{Kind: KindAPIKey}).IsValid() {
		t.Fatal("api key caller without key identity was valid")
	}
}

func TestMetadataRejectsMismatchedAPIKeyFields(t *testing.T) {
	t.Parallel()

	md := metadata.Pairs(
		CallerKindMetadataKey, string(KindAnonymous),
		APIKeyIDMetadataKey, "42",
	)
	if _, _, err := FromMetadata(md); err == nil {
		t.Fatal("mismatched caller metadata was accepted")
	}
}

func TestMetadataRejectsNonCanonicalCallerText(t *testing.T) {
	t.Parallel()

	md := metadata.Pairs(
		CallerKindMetadataKey, string(KindAPIKey),
		APIKeyIDMetadataKey, "42",
		APIKeyAppMetadataKey, " partner-app ",
	)
	if _, _, err := FromMetadata(md); err == nil {
		t.Fatal("caller metadata with surrounding whitespace was accepted")
	}
}

func TestUserCallerMayIncludeVerifiedAPIKeyClient(t *testing.T) {
	t.Parallel()

	caller := Caller{Kind: KindUser, APIKeyID: 8, APIKeyApp: "partner-app"}
	if !caller.IsValid() {
		t.Fatal("user caller with verified API key client was invalid")
	}

	md := metadata.MD{}
	AppendMetadata(md, caller)
	got, ok, err := FromMetadata(md)
	if err != nil || !ok || got != caller {
		t.Fatalf("FromMetadata() = %#v, %v, %v", got, ok, err)
	}
}
