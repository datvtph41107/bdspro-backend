package callercontext

import (
	"context"
	"testing"

	"google.golang.org/grpc/metadata"
)

func TestCallerRoundTrip(t *testing.T) {
	t.Parallel()

	want := Caller{Kind: KindAPIKey, APIKeyID: 42, APIKeyApp: "release-uploader"}
	ctx := WithCaller(context.Background(), want)
	got, ok := FromContext(ctx)
	if !ok || got != want {
		t.Fatalf("FromContext() = %#v, %v; want %#v, true", got, ok, want)
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
