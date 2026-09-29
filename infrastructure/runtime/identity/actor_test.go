package identity

import (
	"context"
	"errors"
	"strings"
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

	ctx, err := BindActor(context.Background(), actor)
	if err != nil {
		t.Fatalf("BindActor() error = %v", err)
	}
	organizationID = 77
	got, ok := ActorFromContext(ctx)
	if !ok {
		t.Fatal("actor is missing")
	}
	if got.AuthID != actor.AuthID || got.ProfileID != actor.ProfileID || got.OriginID != actor.OriginID {
		t.Fatalf("actor IDs = %+v, want %+v", got, actor)
	}
	if got.OrganizationID == nil || *got.OrganizationID != 19 {
		t.Fatalf("organization ID = %v, want 19", got.OrganizationID)
	}

	*got.OrganizationID = 99
	again, _ := ActorFromContext(ctx)
	if again.OrganizationID == nil || *again.OrganizationID != 19 {
		t.Fatalf("context actor was mutated: %+v", again)
	}
}

func TestBindActorRejectsInvalidActors(t *testing.T) {
	t.Parallel()

	zeroOrganization := uint64(0)
	cases := []Actor{
		{},
		{Role: "admin"},
		{ProfileID: 12, OrganizationID: &zeroOrganization},
		{ProfileID: 12, Role: " admin"},
		{ProfileID: 12, Role: "admin\nroot"},
		{ProfileID: 12, TokenType: strings.Repeat("x", MaxActorTextLength+1)},
	}

	for _, actor := range cases {
		ctx := context.Background()
		got, err := BindActor(ctx, actor)
		if !errors.Is(err, ErrInvalidActor) {
			t.Fatalf("BindActor(%+v) error = %v", actor, err)
		}
		if got != ctx {
			t.Fatalf("invalid actor changed context: %+v", actor)
		}
	}
}

func TestBindActorIsIdempotentForSameActor(t *testing.T) {
	t.Parallel()

	actor := Actor{ProfileID: 12, Role: "user"}
	ctx, err := BindActor(context.Background(), actor)
	if err != nil {
		t.Fatalf("first BindActor() error = %v", err)
	}

	rebound, err := BindActor(ctx, actor)
	if err != nil {
		t.Fatalf("same Actor bind error = %v", err)
	}
	if rebound != ctx {
		t.Fatal("same Actor bind created a new context")
	}
}

func TestBindActorRejectsConflictWithoutChangingContext(t *testing.T) {
	t.Parallel()

	ctx, err := BindActor(context.Background(), Actor{ProfileID: 12})
	if err != nil {
		t.Fatalf("first BindActor() error = %v", err)
	}

	unchanged, err := BindActor(ctx, Actor{ProfileID: 99})
	if !errors.Is(err, ErrActorContextConflict) {
		t.Fatalf("conflicting BindActor() error = %v", err)
	}
	if unchanged != ctx {
		t.Fatal("conflicting bind changed context")
	}

	actor, ok := ActorFromContext(unchanged)
	if !ok || actor.ProfileID != 12 {
		t.Fatalf("stored Actor = %+v, %v", actor, ok)
	}
}

func TestBindActorComparesOrganizationIDByValue(t *testing.T) {
	t.Parallel()

	firstOrganizationID := uint64(19)
	secondOrganizationID := uint64(19)
	ctx, err := BindActor(context.Background(), Actor{
		ProfileID:      12,
		OrganizationID: &firstOrganizationID,
	})
	if err != nil {
		t.Fatalf("first BindActor() error = %v", err)
	}

	if _, err := BindActor(ctx, Actor{
		ProfileID:      12,
		OrganizationID: &secondOrganizationID,
	}); err != nil {
		t.Fatalf("equal organization value conflict = %v", err)
	}
}

func TestBindActorRejectsNilContext(t *testing.T) {
	t.Parallel()

	ctx, err := BindActor(nil, Actor{ProfileID: 12})
	if err == nil || ctx != nil {
		t.Fatalf("BindActor(nil) = %v, %v", ctx, err)
	}
}

func TestActorAllowsPartialVerifiedIdentity(t *testing.T) {
	t.Parallel()

	for _, actor := range []Actor{
		{AuthID: 11},
		{ProfileID: 12},
		{OriginID: 13},
		{SessionID: 14},
	} {
		if !actor.IsValid() {
			t.Fatalf("actor should be valid: %+v", actor)
		}
	}
}
