package identitycontext

import (
	"common/identity"
	"context"
	"testing"
)

func TestLegacyActorContextSharesCanonicalAuthority(t *testing.T) {
	t.Parallel()

	actor := Actor{ProfileID: 42, Role: "user"}
	legacyCtx, err := BindActor(context.Background(), actor)
	if err != nil {
		t.Fatalf("BindActor() error = %v", err)
	}
	if got, ok := identity.ActorFromContext(legacyCtx); !ok || got.ProfileID != actor.ProfileID {
		t.Fatalf("canonical read after legacy bind = %+v, %v", got, ok)
	}

	canonicalCtx, err := identity.BindActor(context.Background(), identity.Actor(actor))
	if err != nil {
		t.Fatalf("identity.BindActor() error = %v", err)
	}
	if got, ok := FromContext(canonicalCtx); !ok || got.ProfileID != actor.ProfileID {
		t.Fatalf("legacy read after canonical bind = %+v, %v", got, ok)
	}
}
