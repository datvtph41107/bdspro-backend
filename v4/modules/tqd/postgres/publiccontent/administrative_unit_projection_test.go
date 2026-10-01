package publiccontent

import "testing"

func TestAdministrativeUnitSlugUsesNameAndAdministrativeCode(t *testing.T) {
	slug, err := administrativeUnitSlug("Phường Ngô Quyền", "00123", "ward")
	if err != nil {
		t.Fatalf("administrativeUnitSlug returned error: %v", err)
	}
	if slug != "001232/phuong-ngo-quyen" {
		t.Fatalf("unexpected slug: %q", slug)
	}
}

func TestAdministrativeUnitSlugFailsWithoutStableCode(t *testing.T) {
	if _, err := administrativeUnitSlug("Phường Ngô Quyền", "", "ward"); err == nil {
		t.Fatal("expected missing administrative code to fail closed")
	}
}

func TestNormalizeAdministrativeIdentityAcceptsCanonicalAndTypedIdentity(t *testing.T) {
	unitType, identity := normalizeAdministrativeIdentity("https://qhpro.vn/tin-tuc/dia-ban/001232/phuong-ngo-quyen?utm=test")
	if unitType != "" || identity != "001232" {
		t.Fatalf("unexpected canonical identity: unitType=%q identity=%q", unitType, identity)
	}
	if got := administrativeCodeCandidate(identity); got != "00123" {
		t.Fatalf("unexpected code candidate from composite url id: %q", got)
	}

	unitType, identity = normalizeAdministrativeIdentity("https://qhpro.vn/tin-tuc/dia-ban/phuong-ngo-quyen-00123?utm=test")
	if unitType != "" || identity != "phuong-ngo-quyen-00123" {
		t.Fatalf("unexpected legacy canonical identity: unitType=%q identity=%q", unitType, identity)
	}

	unitType, identity = normalizeAdministrativeIdentity("ward:42")
	if unitType != "ward" || identity != "42" {
		t.Fatalf("unexpected typed identity: unitType=%q identity=%q", unitType, identity)
	}
}

func TestSelectAdministrativeUnitRowRejectsAmbiguousCode(t *testing.T) {
	rows := []administrativeUnitRow{
		{EntityID: "ward:1", SourceID: "1", UnitType: "ward", FullName: "Phường Một", ShortName: "Một", Code: "001"},
		{EntityID: "province:2", SourceID: "2", UnitType: "province", FullName: "Tỉnh Hai", ShortName: "Hai", Code: "001"},
	}
	if _, err := selectAdministrativeUnitRow(rows, "", "001"); err == nil {
		t.Fatal("expected ambiguous administrative code to fail closed")
	}
}

func TestSelectAdministrativeUnitRowUsesTypedIdentity(t *testing.T) {
	rows := []administrativeUnitRow{
		{EntityID: "ward:1", SourceID: "1", UnitType: "ward", FullName: "Phường Một", ShortName: "Một", Code: "001"},
		{EntityID: "province:2", SourceID: "2", UnitType: "province", FullName: "Tỉnh Hai", ShortName: "Hai", Code: "001"},
	}
	row, err := selectAdministrativeUnitRow(rows, "ward", "1")
	if err != nil {
		t.Fatalf("select typed administrative unit: %v", err)
	}
	if row.EntityID != "ward:1" {
		t.Fatalf("unexpected selected row: %q", row.EntityID)
	}
}
