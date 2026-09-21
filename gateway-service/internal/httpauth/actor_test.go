package httpauth

import (
	"common/jwtverify"
	"testing"
)

func TestActorFromPrincipalTreatsZeroOrganizationAsAbsent(t *testing.T) {
	t.Parallel()
	zero := uint64(0)
	actor := ActorFromPrincipal(&jwtverify.Principal{
		AuthID:         11,
		ProfileID:      42,
		OrganizationID: &zero,
		Type:           "ACCESS",
	})
	if !actor.IsValid() {
		t.Fatal("personal token with organization=0 should project a valid actor")
	}
	if actor.OrganizationID != nil {
		t.Fatalf("organization=0 should become nil, got %v", *actor.OrganizationID)
	}
}

func TestActorFromPrincipalKeepsPositiveOrganization(t *testing.T) {
	t.Parallel()
	org := uint64(19)
	actor := ActorFromPrincipal(&jwtverify.Principal{
		AuthID:         11,
		ProfileID:      42,
		OrganizationID: &org,
		Type:           "ACCESS",
	})
	if actor.OrganizationID == nil || *actor.OrganizationID != 19 {
		t.Fatalf("organization was not cloned: %v", actor.OrganizationID)
	}
}
