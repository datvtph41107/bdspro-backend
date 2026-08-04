package services

import (
	"testing"
	"time"

	_algorithm "common/algorithm"
	_enum "common/domain/enum"
)

func testPermissionService() *PermissionService {
	roleOne := _algorithm.NewBitCache(_enum.PermissionIDMax)
	roleOne.Set(10, true)
	roleTwo := _algorithm.NewBitCache(_enum.PermissionIDMax)
	roleTwo.Set(20, true)
	return &PermissionService{
		snapshot: permissionSnapshot{
			rolePermissions:  map[uint64]*_algorithm.BitCache{1: roleOne, 2: roleTwo},
			permissionByCode: map[string]uint32{"PERM_A": 10, "PERM_B": 20},
			version:          "v1",
			loadedAt:         time.Now(),
		},
		maxStaleness: time.Minute,
	}
}

func TestPermissionMatchingIsExplicit(t *testing.T) {
	svc := testPermissionService()
	roles := []uint64{1, 2}
	if !svc.hasAllPermissions(roles, []uint32{10, 20}) {
		t.Fatal("ALL should combine permissions across assigned roles")
	}
	if svc.hasAllPermissions([]uint64{1}, []uint32{10, 20}) {
		t.Fatal("ALL must reject a missing permission")
	}
	if !svc.hasAnyPermission([]uint64{1}, []uint32{10, 20}) {
		t.Fatal("ANY should accept one matching permission")
	}
}

func TestResolvePermissionCodes(t *testing.T) {
	svc := testPermissionService()
	ids, err := svc.resolvePermissionIDs([]uint32{20}, []string{"PERM_A"})
	if err != nil {
		t.Fatalf("resolvePermissionIDs() error = %v", err)
	}
	if len(ids) != 2 || ids[0] != 10 || ids[1] != 20 {
		t.Fatalf("ids = %v", ids)
	}
}

func TestSnapshotStatusMarksStale(t *testing.T) {
	svc := testPermissionService()
	svc.snapshot.loadedAt = time.Now().Add(-2 * time.Minute)
	status := svc.Status(time.Now())
	if !status.Stale {
		t.Fatal("stale snapshot should be reported")
	}
}

func TestNormalizeRoleIDs(t *testing.T) {
	got := normalizeRoleIDs([]uint64{3, 0, 2, 3})
	if len(got) != 2 || got[0] != 2 || got[1] != 3 {
		t.Fatalf("roles = %v", got)
	}
}
