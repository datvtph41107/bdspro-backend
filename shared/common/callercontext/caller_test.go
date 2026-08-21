package callercontext

import (
	"testing"

	"google.golang.org/grpc/metadata"
)

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

func TestUserCallerMetadataMayIncludeVerifiedAPIKeyClient(t *testing.T) {
	t.Parallel()

	caller := Caller{Kind: KindUser, APIKeyID: 8, APIKeyApp: "partner-app"}
	md := metadata.MD{}
	AppendMetadata(md, caller)
	got, ok, err := FromMetadata(md)
	if err != nil || !ok || got != caller {
		t.Fatalf("FromMetadata() = %#v, %v, %v", got, ok, err)
	}
}
