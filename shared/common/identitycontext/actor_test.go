package identitycontext

import (
	"context"
	"testing"
)

func TestActorContextRoundTrip(t *testing.T) {
	t.Parallel()

	organizationID := uint64(19)
	actor := Actor{
		AuthID:         11,
		ProfileID:      12,
		OriginID:       13,
		OrganizationID: &organizationID,
		SessionID:      14,
		Role:           "user",
		TokenType:      "ACCESS",
	}

	ctx := WithActor(context.Background(), actor)
	got, ok := FromContext(ctx)
	if !ok {
		t.Fatal("actor is missing")
	}
	if got.AuthID != actor.AuthID || got.ProfileID != actor.ProfileID || got.OriginID != actor.OriginID {
		t.Fatalf("actor IDs = %+v, want %+v", got, actor)
	}
	if got.OrganizationID == nil || *got.OrganizationID != organizationID {
		t.Fatalf("organization ID = %v, want %d", got.OrganizationID, organizationID)
	}

	*got.OrganizationID = 99
	again, _ := FromContext(ctx)
	if again.OrganizationID == nil || *again.OrganizationID != organizationID {
		t.Fatalf("context actor was mutated: %+v", again)
	}
}

func TestWithActorIgnoresZeroActor(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	gotCtx := WithActor(ctx, Actor{})
	if gotCtx != ctx {
		t.Fatal("zero actor changed context")
	}
	if _, ok := FromContext(gotCtx); ok {
		t.Fatal("zero actor reached context")
	}
}
