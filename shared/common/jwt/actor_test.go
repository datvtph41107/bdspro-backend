package _jwt

import "testing"

func TestActorFromPrincipalExcludesEntitlementAndRoleAssignments(t *testing.T) {
	t.Parallel()

	organizationID := uint64(19)
	planID := uint64(21)
	principal := &Principal{
		AuthID:         11,
		ProfileId:      12,
		OriginId:       13,
		OrganizationId: &organizationID,
		PlanId:         &planID,
		RoleIds:        []uint64{31, 32},
		Session:        14,
		Role:           "user",
		Type:           "ACCESS",
	}

	actor := ActorFromPrincipal(principal)
	if !actor.IsValid() {
		t.Fatalf("actor is invalid: %+v", actor)
	}
	if actor.AuthID != 11 || actor.ProfileID != 12 || actor.OriginID != 13 || actor.SessionID != 14 {
		t.Fatalf("actor IDs = %+v", actor)
	}
	organizationID = 99
	if actor.OrganizationID == nil || *actor.OrganizationID != 19 {
		t.Fatalf("actor organization was not cloned: %v", actor.OrganizationID)
	}
}

func TestActorFromPrincipalRejectsRoleOnlyIdentity(t *testing.T) {
	t.Parallel()

	actor := ActorFromPrincipal(&Principal{Role: "admin", Type: "ACCESS"})
	if actor.IsValid() {
		t.Fatalf("role-only actor is valid: %+v", actor)
	}
}
