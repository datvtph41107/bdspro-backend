package _jwt

import (
	"context"
	"testing"
)

func TestPrincipalRequestContextUsesDefensiveCopies(t *testing.T) {
	t.Parallel()

	organizationID := uint64(19)
	planID := uint64(21)
	principal := &Principal{
		AuthID:         11,
		ProfileId:      12,
		OrganizationId: &organizationID,
		PlanId:         &planID,
		RoleIds:        []uint64{31, 32},
	}
	ctx := WithPrincipal(context.Background(), principal)

	principal.ProfileId = 99
	organizationID = 77
	planID = 88
	principal.RoleIds[0] = 100

	got, ok := PrincipalFromRequestContext(ctx)
	if !ok {
		t.Fatal("principal is missing")
	}
	if got == principal {
		t.Fatal("context returned the original principal pointer")
	}
	if got.ProfileId != 12 || got.OrganizationId == nil || *got.OrganizationId != 19 {
		t.Fatalf("principal mutated through source: %+v", got)
	}
	if got.PlanId == nil || *got.PlanId != 21 || len(got.RoleIds) != 2 || got.RoleIds[0] != 31 {
		t.Fatalf("principal nested data mutated: %+v", got)
	}

	got.ProfileId = 55
	*got.OrganizationId = 66
	got.RoleIds[0] = 77
	again, _ := PrincipalFromRequestContext(ctx)
	if again.ProfileId != 12 || *again.OrganizationId != 19 || again.RoleIds[0] != 31 {
		t.Fatalf("stored principal mutated through returned copy: %+v", again)
	}
}

func TestWithPrincipalIgnoresNil(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	if got := WithPrincipal(ctx, nil); got != ctx {
		t.Fatal("nil principal changed context")
	}
}
