package publiccontent

import (
	"testing"
	"time"

	"tqd/internal/domain/publiccontent/model"
)

func TestBuildRelatedPlanningProjectKeepsUpdateAndDossierPath(t *testing.T) {
	t.Parallel()

	updatedAt := time.Date(2026, time.July, 30, 8, 0, 0, 0, time.UTC)
	project := buildRelatedPlanningProject(relatedProjectRow{
		ID:        9,
		Name:      "Đồ án A",
		UpdatedAt: &updatedAt,
	}, "do-an-a", 42.5)

	if project.DossierPath != "/ban-do/do-an/9" {
		t.Fatalf("unexpected dossier path: %q", project.DossierPath)
	}
	if project.PublicPath != "/do-an-quy-hoach/do-an-a" {
		t.Fatalf("unexpected public SEO path: %q", project.PublicPath)
	}
	if project.UpdatedAt == nil || !project.UpdatedAt.Equal(updatedAt) {
		t.Fatalf("unexpected updatedAt: %#v", project.UpdatedAt)
	}
}

func TestBuildParcelQuickViewKeepsExistingParcelWhenMetadataIsMissing(t *testing.T) {
	t.Parallel()

	view := buildParcelQuickView(parcelRow{
		ParcelID:      123,
		HasParcelInfo: false,
		MapNumber:     "3",
		LandNumber:    "12",
		Address:       "Phường Sơn Tây",
		TotalAreaSqm:  128.5,
		GeometryType:  "POLYGON",
		CenterLat:     21.03,
		CenterLon:     105.84,
		MinLon:        105.83,
		MinLat:        21.02,
		MaxLon:        105.85,
		MaxLat:        21.04,
	}, nil)

	if view.ParcelID != 123 {
		t.Fatalf("unexpected parcel id: %d", view.ParcelID)
	}
	if view.CanonicalKey != "parcel:123" {
		t.Fatalf("unexpected canonical key: %q", view.CanonicalKey)
	}
	if view.Completeness != "partial" {
		t.Fatalf("expected partial completeness, got %q", view.Completeness)
	}
	if !containsWarning(view.Warnings, "parcel_info_missing") {
		t.Fatalf("expected parcel_info_missing warning: %#v", view.Warnings)
	}
	if !containsWarning(view.Warnings, "planning_context_empty") {
		t.Fatalf("expected planning_context_empty warning: %#v", view.Warnings)
	}
	if view.TotalAreaSqm != 128.5 {
		t.Fatalf("unexpected total area: %v", view.TotalAreaSqm)
	}
	if view.Preview.Completeness != "complete" {
		t.Fatalf("expected complete spatial preview, got %q", view.Preview.Completeness)
	}
}

func TestBuildParcelQuickViewIsCompleteWhenMetadataExists(t *testing.T) {
	t.Parallel()

	view := buildParcelQuickView(
		parcelRow{ParcelID: 456, HasParcelInfo: true},
		[]domain.RelatedPlanningProject{{ID: 9, Name: "Đồ án A"}},
	)

	if view.Completeness != "complete" {
		t.Fatalf("expected complete metadata, got %q", view.Completeness)
	}
	if containsWarning(view.Warnings, "parcel_info_missing") {
		t.Fatalf("unexpected parcel_info_missing warning: %#v", view.Warnings)
	}
	if containsWarning(view.Warnings, "planning_context_empty") {
		t.Fatalf("unexpected planning_context_empty warning: %#v", view.Warnings)
	}
}

func containsWarning(values []string, expected string) bool {
	for _, value := range values {
		if value == expected {
			return true
		}
	}
	return false
}
