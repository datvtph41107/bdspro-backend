package discovery

import (
	"testing"

	"tqd/internal/domain/discovery/model"
)

func TestPlanningRegionDetailURLUsesPublicApplicationRoute(t *testing.T) {
	t.Parallel()

	row := spatialRow{ID: "7421", Title: "Khu quy hoạch trung tâm"}
	path := detailURL(domain.EntityKindPlanningRegion, row)

	if path != "/vung-quy-hoach/7421" {
		t.Fatalf("unexpected planning region detail URL: %q", path)
	}

	item := candidate(domain.EntityKindPlanningRegion, row, "exact")
	if item.Links.DetailURL != path {
		t.Fatalf("unexpected candidate detail URL: %q", item.Links.DetailURL)
	}
	if item.Links.CanonicalURL != path {
		t.Fatalf("unexpected candidate canonical URL: %q", item.Links.CanonicalURL)
	}
}
