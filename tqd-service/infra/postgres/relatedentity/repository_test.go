package relatedentity

import (
	"math"
	"os"
	"strings"
	"testing"

	discoverydomain "tqd/internal/domain/discovery/model"
	relateddomain "tqd/internal/domain/relatedentity/model"
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
	if got := relationshipKind("unknown"); got != relateddomain.RelationshipUnknown {
		t.Fatalf("relationshipKind(unknown) = %q, want unknown", got)
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

func TestParcelPlanningQueryKeepsCoreRelationIndexFriendly(t *testing.T) {
	contents, err := os.ReadFile("planning.go")
	if err != nil {
		t.Fatalf("ReadFile(planning.go): %v", err)
	}
	source := string(contents)
	start := strings.Index(source, "func (r *Repository) findPlanningRegionsForParcel")
	end := strings.Index(source[start:], "func (r *Repository) findPlanningEntitiesForRegion")
	if start < 0 || end < 0 {
		t.Fatal("cannot locate parcel planning query")
	}
	body := source[start : start+end]
	if strings.Contains(body, "ST_DWithin") || strings.Contains(body, "ST_Buffer") {
		t.Fatalf("parcel core planning query must not pay nearby-distance cost: %s", body)
	}
	if !strings.Contains(body, "ST_Intersects") || !strings.Contains(body, "region.geometry && primary_entity.geometry") {
		t.Fatalf("parcel core planning query must use bbox + direct intersection: %s", body)
	}
}

func TestParcelAdjacencyQueryDoesNotMixNearbyRadiusSearch(t *testing.T) {
	body := sourceFunctionBody(t, "parcel.go", "findParcelsAroundParcel", "findParcelsInsidePlanningRegion")
	if strings.Contains(body, "ST_DWithin") || strings.Contains(body, "ST_Buffer") {
		t.Fatalf("parcel adjacency must not pay nearby-radius cost: %s", body)
	}
	if !strings.Contains(body, "ST_Touches") || !strings.Contains(body, "parcel.geometry && primary_entity.geometry") {
		t.Fatalf("parcel adjacency must use bbox + ST_Touches: %s", body)
	}
}

func TestPlanningContainedParcelQueryUsesWholeRegionGeometry(t *testing.T) {
	body := sourceFunctionBody(t, "parcel.go", "findParcelsInsidePlanningRegion", "findParcelsAroundPoint")
	if strings.Contains(body, "ST_DWithin") || strings.Contains(body, "ST_Buffer") || strings.Contains(body, "focus.point") {
		t.Fatalf("planning containment must not be an undocumented centroid sample: %s", body)
	}
	if !strings.Contains(body, "ST_Covers(primary_entity.geometry, parcel.geometry)") ||
		!strings.Contains(body, "parcel.geometry && primary_entity.geometry") {
		t.Fatalf("planning containment must use bbox + full geometry coverage: %s", body)
	}
}

func sourceFunctionBody(t *testing.T, filename, functionName, nextFunctionName string) string {
	t.Helper()
	contents, err := os.ReadFile(filename)
	if err != nil {
		t.Fatalf("ReadFile(%s): %v", filename, err)
	}
	source := string(contents)
	start := strings.Index(source, "func (r *Repository) "+functionName)
	if start < 0 {
		t.Fatalf("cannot locate %s", functionName)
	}
	endOffset := strings.Index(source[start:], "func (r *Repository) "+nextFunctionName)
	if endOffset < 0 {
		t.Fatalf("cannot locate %s after %s", nextFunctionName, functionName)
	}
	return source[start : start+endOffset]
}
