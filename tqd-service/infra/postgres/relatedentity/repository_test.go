package relatedentity

import (
	"math"
	"strings"
	"testing"

	discoverydomain "tqd/internal/modules/discovery/domain"
	relateddomain "tqd/internal/modules/relatedentity/domain"
)

func TestRelationshipKindPreservesSupportedValues(t *testing.T) {
	supported := []relateddomain.RelationshipKind{
		relateddomain.RelationshipContainsPrimary,
		relateddomain.RelationshipContainedByPrimary,
		relateddomain.RelationshipAdministrativeContext,
		relateddomain.RelationshipMemberOfPrimary,
		relateddomain.RelationshipIntersectsPrimary,
		relateddomain.RelationshipAdjacent,
		relateddomain.RelationshipNearby,
	}
	for _, value := range supported {
		if got := relationshipKind("  " + string(value) + "  "); got != value {
			t.Fatalf("relationshipKind(%q) = %q, want %q", value, got, value)
		}
	}
	if got := relationshipKind("unknown"); got != relateddomain.RelationshipNearby {
		t.Fatalf("relationshipKind(unknown) = %q, want nearby", got)
	}
}

func TestNormalizedDistanceRejectsInvalidValues(t *testing.T) {
	for _, value := range []float64{-1, math.NaN(), math.Inf(1)} {
		if got := normalizedDistance(value); got != math.MaxFloat64 {
			t.Fatalf("normalizedDistance(%v) = %v, want MaxFloat64", value, got)
		}
	}
	if got := normalizedDistance(12.5); got != 12.5 {
		t.Fatalf("normalizedDistance(12.5) = %v", got)
	}
}

func TestMapEntityCandidateBuildsCanonicalAdministrativeKey(t *testing.T) {
	candidates := rowsToCandidates(discoverydomain.EntityKindAdministrativeUnit, []spatialRow{{
		ID:           "ward:abc",
		Title:        "Phường mẫu",
		WardCode:     "00001",
		GeometryType: "Point",
		CenterLat:    10.1,
		CenterLon:    106.2,
		Relationship: string(relateddomain.RelationshipAdministrativeContext),
	}})
	if len(candidates) != 1 {
		t.Fatalf("len(candidates) = %d, want 1", len(candidates))
	}
	candidate := candidates[0].Entity
	if candidate.Entity.ID != "ward:abc" {
		t.Fatalf("Entity.ID = %q", candidate.Entity.ID)
	}
	if candidate.Entity.Key != "administrative_unit:ward:abc" {
		t.Fatalf("Entity.Key = %q", candidate.Entity.Key)
	}
	if candidate.Spatial == nil || candidate.Spatial.Centroid == nil {
		t.Fatal("candidate centroid is required")
	}
}

func TestPlanningProjectMembershipIsNotRadiusFiltered(t *testing.T) {
	sql := planningRegionsForProjectSQL(false)
	if strings.Contains(sql, "ST_DWithin") || strings.Contains(sql, "search_envelope") {
		t.Fatalf("direct planning-project membership must not be radius filtered: %s", sql)
	}
	if !strings.Contains(sql, "layer.planning_project_id::text = ?") {
		t.Fatalf("membership query must use planning_project_id: %s", sql)
	}
	if !strings.Contains(sql, "'member_of_primary'") {
		t.Fatalf("membership relationship missing: %s", sql)
	}
}
