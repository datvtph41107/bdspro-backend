package usecase

import (
	"testing"

	"crm/internal/dto"
)

func TestNormalizeAdministrativeUnitProvisioningIdentityAcceptsOpaqueWardSourceID(t *testing.T) {
	t.Parallel()

	got, err := normalizeAdministrativeUnitProvisioningIdentity("ward:D_26740")
	if err != nil {
		t.Fatalf("expected opaque ward typed identity to be accepted: %v", err)
	}
	if got != "ward:D_26740" {
		t.Fatalf("unexpected normalized identity: got %q want %q", got, "ward:D_26740")
	}
}

func TestNormalizeAdministrativeUnitProvisioningIdentityAcceptsOpaqueProvinceSourceID(t *testing.T) {
	t.Parallel()

	got, err := normalizeAdministrativeUnitProvisioningIdentity("province:5e5501caeb80a7245175de0c")
	if err != nil {
		t.Fatalf("expected opaque province typed identity to be accepted: %v", err)
	}
	if got != "province:5e5501caeb80a7245175de0c" {
		t.Fatalf("unexpected normalized identity: got %q", got)
	}
}

func TestAdministrativeUnitReferenceIDUsesAdministrativeCode(t *testing.T) {
	t.Parallel()

	projection := &dto.TqdAdministrativeUnitProjection{
		EntityID:           "ward:D_26740",
		AdministrativeCode: "26740",
	}

	got := administrativeUnitReferenceID(projection)
	if got != 26740 {
		t.Fatalf("expected compatibility reference ID 26740, got %d", got)
	}
}

func TestValidateAdministrativeUnitProvisioningProjectionAcceptsOpaqueTypedIdentity(t *testing.T) {
	t.Parallel()

	projection := &dto.TqdAdministrativeUnitProjection{
		SchemaVersion:      dto.TqdAdministrativeUnitProjectionSchemaVersion,
		EntityID:           "ward:D_26740",
		Slug:               "phuong-sai-gon-26740",
		CanonicalPath:      "/tin-tuc/dia-ban/phuong-sai-gon-26740",
		UnitType:           "ward",
		Name:               "Phường Sài Gòn",
		AdministrativeCode: "26740",
	}

	if err := validateAdministrativeUnitProvisioningProjection(projection); err != nil {
		t.Fatalf("expected opaque typed identity projection to be valid: %v", err)
	}
}

func TestValidateAdministrativeUnitProvisioningProjectionRejectsUnitTypeMismatch(t *testing.T) {
	t.Parallel()

	projection := &dto.TqdAdministrativeUnitProjection{
		SchemaVersion:      dto.TqdAdministrativeUnitProjectionSchemaVersion,
		EntityID:           "ward:D_26740",
		Slug:               "phuong-sai-gon-26740",
		CanonicalPath:      "/tin-tuc/dia-ban/phuong-sai-gon-26740",
		UnitType:           "province",
		Name:               "Phường Sài Gòn",
		AdministrativeCode: "26740",
	}

	if err := validateAdministrativeUnitProvisioningProjection(projection); err == nil {
		t.Fatal("expected projection unit type mismatch to be rejected")
	}
}
