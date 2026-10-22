package controller

import "testing"

func TestNormalizeClientLocationCatalogItem_NameAndNoShortName(t *testing.T) {
	item := map[string]any{
		"id":        float64(1),
		"shortName": "HN",
		"fullName":  "Thành phố Hà Nội",
	}
	out := normalizeClientLocationCatalogItem(item).(map[string]any)
	if _, ok := out["shortName"]; ok {
		t.Fatal("expected shortName removed")
	}
	if out["name"] != "Thành phố Hà Nội" {
		t.Fatalf("name = %v, want fullName", out["name"])
	}
}

func TestNormalizeClientLocationCatalogItem_FallsBackToOfficialName(t *testing.T) {
	item := map[string]any{
		"officialName": "Phường Một",
		"shortName":    "Một",
	}
	out := normalizeClientLocationCatalogItem(item).(map[string]any)
	if out["name"] != "Phường Một" {
		t.Fatalf("name = %v", out["name"])
	}
}
