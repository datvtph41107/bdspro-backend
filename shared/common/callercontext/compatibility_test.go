package callercontext

import (
	"common/identity"
	"context"
	"testing"
)

func TestLegacyCallerContextSharesCanonicalAuthority(t *testing.T) {
	t.Parallel()

	caller := Caller{Kind: KindAPIKey, APIKeyID: 42, APIKeyApp: "release-uploader"}
	legacyCtx, err := BindCaller(context.Background(), caller)
	if err != nil {
		t.Fatalf("BindCaller() error = %v", err)
	}
	if got, ok := identity.CallerFromContext(legacyCtx); !ok || got != identity.Caller(caller) {
		t.Fatalf("canonical read after legacy bind = %+v, %v", got, ok)
	}

	canonicalCtx, err := identity.BindCaller(context.Background(), identity.Caller(caller))
	if err != nil {
		t.Fatalf("identity.BindCaller() error = %v", err)
	}
	if got, ok := FromContext(canonicalCtx); !ok || got != caller {
		t.Fatalf("legacy read after canonical bind = %+v, %v", got, ok)
	}
}
