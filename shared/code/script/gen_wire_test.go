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
