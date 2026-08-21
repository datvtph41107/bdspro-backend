package services

import (
	"testing"
	"time"

	"auth/internal/permissioncatalog"
)

func testPermissionService() *PermissionService {
	catalog, err := permissioncatalog.Build(
		[]permissioncatalog.RolePermissions{
			{RoleID: 1, PermissionIDs: []uint64{10}},
			{RoleID: 2, PermissionIDs: []uint64{20}},
		},
		[]permissioncatalog.Permission{
			{ID: 10, Code: "PERM_A"},
			{ID: 20, Code: "PERM_B"},
		},
	)
	if err != nil {
		panic(err)
	}
	return &PermissionService{
		snapshot:     permissionSnapshot{catalog: catalog, loadedAt: time.Now()},
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

func TestResolvePermissionIDMustExistInSnapshot(t *testing.T) {
	svc := testPermissionService()
	if _, err := svc.resolvePermissionIDs([]uint32{999}, nil); err == nil {
		t.Fatal("unknown numeric permission ID must fail closed")
	}
}
