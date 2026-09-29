package main

import (
	"reflect"
	"testing"
)

func TestNormalizeScanExclude(t *testing.T) {
	got := normalizeScanExclude([]string{
		" internal/modules/catalog/ ",
		"internal/modules/catalog",
		"internal/modules/subscription/shadow",
		"",
	})
	want := []string{
		"internal/modules/catalog",
		"internal/modules/subscription/shadow",
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("normalizeScanExclude() = %#v, want %#v", got, want)
	}
}

func TestIsScanExcludedMatchesSubtreeOnly(t *testing.T) {
	excluded := []string{"internal/modules/catalog"}
	for _, path := range []string{
		"internal/modules/catalog",
		"internal/modules/catalog/shadow",
		"internal/modules/catalog/migration/file.go",
	} {
		if !isScanExcluded(path, excluded) {
			t.Fatalf("expected %q to be excluded", path)
		}
	}
	for _, path := range []string{
		"internal/domain/plan",
		"internal/modules/catalogue",
		"internal/modules/subscription",
	} {
		if isScanExcluded(path, excluded) {
			t.Fatalf("did not expect %q to be excluded", path)
		}
	}
}

func TestFilterInjectionPackagesScopesConstructorsByService(t *testing.T) {
	packages := []InjectionPackage{
		{
			Import:       "common/utils",
			Constructors: []string{"NewSyncUtil"},
		},
		{
			Import:       "common/tilesession",
			Constructors: []string{"NewStore"},
			Services:     []string{"user"},
		},
	}

	user := filterInjectionPackages(packages, nil, "user")
	if len(user) != 2 {
		t.Fatalf("user packages = %d, want 2", len(user))
	}

	hub := filterInjectionPackages(packages, nil, "hub")
	if len(hub) != 1 {
		t.Fatalf("hub packages = %d, want 1", len(hub))
	}
	if hub[0].Import != "common/utils" {
		t.Fatalf("hub package = %q, want common/utils", hub[0].Import)
	}
}

func TestFilterInjectionPackagesKeepsExclusionsWithinServiceScope(t *testing.T) {
	packages := []InjectionPackage{
		{
			Import:       "common/tilesession",
			Constructors: []string{"NewStore", "NewReader"},
			Services:     []string{"user"},
		},
	}

	got := filterInjectionPackages(
		packages,
		[]string{"common/tilesession.NewReader"},
		"user",
	)
	if len(got) != 1 {
		t.Fatalf("packages = %d, want 1", len(got))
	}
	if len(got[0].Constructors) != 1 || got[0].Constructors[0] != "NewStore" {
		t.Fatalf("constructors = %#v, want [NewStore]", got[0].Constructors)
	}
}
