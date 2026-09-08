package permissioncatalog

import (
	"errors"
	"testing"
)

func TestBuildAllowsDurableIDsBeyondLegacyCompileTimeCeiling(t *testing.T) {
	snapshot, err := Build(
		[]RolePermissions{{RoleID: 1, PermissionIDs: []uint64{5000}}},
		[]Permission{{ID: 5000, Code: "CATALOG_PLAN_VIEW"}},
	)
	if err != nil {
		t.Fatal(err)
	}
	id, ok := snapshot.ResolveCode("CATALOG_PLAN_VIEW")
	if !ok || id != 5000 {
		t.Fatalf("resolved id = %d, ok=%v", id, ok)
	}
	if !snapshot.RoleHas(1, 5000) {
		t.Fatal("role should own dynamically catalogued permission")
	}
}

func TestBuildRejectsRoleAssignmentOutsideDurableCatalog(t *testing.T) {
	_, err := Build(
		[]RolePermissions{{RoleID: 1, PermissionIDs: []uint64{10, 11}}},
		[]Permission{{ID: 10, Code: "KNOWN"}},
	)
	if !errors.Is(err, ErrUnknownRolePermission) {
		t.Fatalf("err = %v", err)
	}
}

func TestBuildRejectsConflictingCodeOrID(t *testing.T) {
	_, err := Build(nil, []Permission{{ID: 10, Code: "A"}, {ID: 11, Code: "A"}})
	if !errors.Is(err, ErrDuplicatePermissionCode) {
		t.Fatalf("duplicate code err = %v", err)
	}
	_, err = Build(nil, []Permission{{ID: 10, Code: "A"}, {ID: 10, Code: "B"}})
	if !errors.Is(err, ErrDuplicatePermissionID) {
		t.Fatalf("duplicate id err = %v", err)
	}
}

func TestBuildRejectsUnreasonableCatalogID(t *testing.T) {
	_, err := Build(nil, []Permission{{ID: MaxReasonablePermissionID + 1, Code: "BROKEN"}})
	if !errors.Is(err, ErrInvalidPermission) {
		t.Fatalf("err = %v", err)
	}
}

func TestVersionChangesWithAuthoritySnapshot(t *testing.T) {
	a, err := Build([]RolePermissions{{RoleID: 1, PermissionIDs: []uint64{10}}}, []Permission{{ID: 10, Code: "A"}})
	if err != nil {
		t.Fatal(err)
	}
	b, err := Build([]RolePermissions{{RoleID: 1, PermissionIDs: []uint64{10, 20}}}, []Permission{{ID: 10, Code: "A"}, {ID: 20, Code: "B"}})
	if err != nil {
		t.Fatal(err)
	}
	if a.Version() == b.Version() {
		t.Fatal("snapshot version must change when permission authority changes")
	}
}
