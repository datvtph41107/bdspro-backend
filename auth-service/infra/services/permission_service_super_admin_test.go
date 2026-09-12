package services

import (
	"testing"
	"time"

	"common/identity"
)

func TestRequiredPermissionsSuperAdminBypass(t *testing.T) {
	svc := testPermissionService()
	svc.roleCache = map[uint64]roleCacheEntry{
		99: {
			roleIDs:    nil,
			superAdmin: 1,
			expiresAt:  time.Now().Add(time.Minute),
		},
	}

	ctx, err := identity.BindActor(t.Context(), identity.Actor{
		ProfileID: 99,
		Role:      "ROLE_ADMIN",
		TokenType: "ACCESS",
	})
	if err != nil {
		t.Fatalf("BindActor() error = %v", err)
	}

	if err := svc.RequiredPermissions(ctx, nil, []string{"PERM_A"}, MatchAny); err != nil {
		t.Fatalf("super_admin=1 must pass RequiredPermissions: %v", err)
	}
}

func TestRequiredPermissionsStillRequiresGrantsWithoutSuperAdmin(t *testing.T) {
	svc := testPermissionService()
	svc.roleCache = map[uint64]roleCacheEntry{
		88: {
			roleIDs:    nil,
			superAdmin: 0,
			expiresAt:  time.Now().Add(time.Minute),
		},
	}

	ctx, err := identity.BindActor(t.Context(), identity.Actor{
		ProfileID: 88,
		Role:      "ROLE_ADMIN",
		TokenType: "ACCESS",
	})
	if err != nil {
		t.Fatalf("BindActor() error = %v", err)
	}

	if err := svc.RequiredPermissions(ctx, nil, []string{"PERM_A"}, MatchAny); err == nil {
		t.Fatal("non-super admin without roles must be denied")
	}
}
