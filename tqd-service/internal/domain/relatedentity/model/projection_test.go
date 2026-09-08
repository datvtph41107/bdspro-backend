package domain

import (
	"testing"

	discoverydomain "tqd/internal/domain/discovery/model"
)

func TestResolveProjectionPolicyKeepsParcelAsItems(t *testing.T) {
	primary := projectionEntity(discoverydomain.EntityKindParcel, 42_000_000)
	policy := ResolveProjectionPolicy(&primary)
	if policy.DisplayMode != DisplayModeItems {
		t.Fatalf("DisplayMode = %q, want %q", policy.DisplayMode, DisplayModeItems)
	}
	if policy.Strategy != StrategyExactList {
		t.Fatalf("Strategy = %q, want %q", policy.Strategy, StrategyExactList)
	}
}

func TestResolveProjectionPolicyGroupsLargePlanningRegion(t *testing.T) {
	primary := projectionEntity(discoverydomain.EntityKindPlanningRegion, 41_200_000)
	policy := ResolveProjectionPolicy(&primary)
	if policy.DisplayMode != DisplayModeGroups {
		t.Fatalf("DisplayMode = %q, want %q", policy.DisplayMode, DisplayModeGroups)
	}
	if policy.Coverage != CoverageAggregated {
		t.Fatalf("Coverage = %q, want %q", policy.Coverage, CoverageAggregated)
	}
	if policy.Reason != ProjectionReasonScopeLarge {
		t.Fatalf("Reason = %q, want %q", policy.Reason, ProjectionReasonScopeLarge)
	}
	assertProjectionGroups(t, policy.AllowedGroups, []string{"planning", "administrative"})
	assertProjectionGroups(t, policy.SuppressedGroups, []string{"parcel", "poi"})
}

func TestResolveProjectionPolicyRestrictsVeryLargePlanningRegion(t *testing.T) {
	primary := projectionEntity(discoverydomain.EntityKindPlanningRegion, 500_000_000)
	policy := ResolveProjectionPolicy(&primary)
	if policy.DisplayMode != DisplayModeRestricted {
		t.Fatalf("DisplayMode = %q, want %q", policy.DisplayMode, DisplayModeRestricted)
	}
	if policy.Reason != ProjectionReasonViewportNeeded {
		t.Fatalf("Reason = %q, want %q", policy.Reason, ProjectionReasonViewportNeeded)
	}
}

func TestResolveProjectionPolicyGroupsPlanningProjectWithoutArea(t *testing.T) {
	primary := projectionEntity(discoverydomain.EntityKindPlanningProject, 0)
	policy := ResolveProjectionPolicy(&primary)
	if policy.DisplayMode != DisplayModeGroups {
		t.Fatalf("DisplayMode = %q, want %q", policy.DisplayMode, DisplayModeGroups)
	}
	if policy.Reason != ProjectionReasonProjectSummary {
		t.Fatalf("Reason = %q, want %q", policy.Reason, ProjectionReasonProjectSummary)
	}
	assertProjectionGroups(t, policy.SuppressedGroups, []string{"parcel", "poi"})
}

func TestResolveProjectionPolicyUsesAdministrativeThresholds(t *testing.T) {
	itemized := projectionEntity(discoverydomain.EntityKindAdministrativeUnit, 10_000_000)
	itemizedPolicy := ResolveProjectionPolicy(&itemized)
	if itemizedPolicy.DisplayMode != DisplayModeItems {
		t.Fatalf("10 km² administrative DisplayMode = %q, want %q", itemizedPolicy.DisplayMode, DisplayModeItems)
	}
	if itemizedPolicy.Scope.Scale != ScopeScaleMedium {
		t.Fatalf("10 km² administrative Scope.Scale = %q, want %q", itemizedPolicy.Scope.Scale, ScopeScaleMedium)
	}

	grouped := projectionEntity(discoverydomain.EntityKindAdministrativeUnit, 30_000_000)
	if policy := ResolveProjectionPolicy(&grouped); policy.DisplayMode != DisplayModeGroups {
		t.Fatalf("30 km² administrative DisplayMode = %q, want %q", policy.DisplayMode, DisplayModeGroups)
	} else if policy.Scope.Scale != ScopeScaleLarge {
		t.Fatalf("30 km² administrative Scope.Scale = %q, want %q", policy.Scope.Scale, ScopeScaleLarge)
	}

	restricted := projectionEntity(discoverydomain.EntityKindAdministrativeUnit, 1_200_000_000)
	if policy := ResolveProjectionPolicy(&restricted); policy.DisplayMode != DisplayModeRestricted {
		t.Fatalf("1200 km² administrative DisplayMode = %q, want %q", policy.DisplayMode, DisplayModeRestricted)
	} else if policy.Scope.Scale != ScopeScaleVeryLarge {
		t.Fatalf("1200 km² administrative Scope.Scale = %q, want %q", policy.Scope.Scale, ScopeScaleVeryLarge)
	}
}

func TestResolveProjectionPolicyUsesBoundsWhenAreaAttributeMissing(t *testing.T) {
	primary := projectionEntity(discoverydomain.EntityKindPlanningRegion, 0)
	primary.Attributes = nil
	primary.Spatial.Bounds = &discoverydomain.Bounds{
		MinLongitude: 105.80,
		MinLatitude:  20.98,
		MaxLongitude: 105.90,
		MaxLatitude:  21.08,
	}
	policy := ResolveProjectionPolicy(&primary)
	if policy.Scope.AreaSquareMeters <= 0 {
		t.Fatal("bounds estimate must produce a positive area")
	}
	if policy.Scope.AreaSource != "bounds_estimate" {
		t.Fatalf("AreaSource = %q, want bounds_estimate", policy.Scope.AreaSource)
	}
}

func projectionEntity(kind discoverydomain.EntityKind, area float64) discoverydomain.EntityCandidate {
	attributes := map[string]any{}
	if area > 0 {
		attributes["areaSqm"] = area
	}
	return discoverydomain.EntityCandidate{
		Entity:     discoverydomain.NewEntityRef(kind, "1"),
		Attributes: attributes,
		Spatial: &discoverydomain.SpatialSummary{
			Centroid: &discoverydomain.Point{Latitude: 21, Longitude: 105.8},
		},
	}
}

func assertProjectionGroups(t *testing.T, got, want []string) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("groups = %v, want %v", got, want)
	}
	for index := range want {
		if got[index] != want[index] {
			t.Fatalf("groups[%d] = %q, want %q", index, got[index], want[index])
		}
	}
}

func TestRelatedRoleOrderMatchesPrimaryDecisionIntent(t *testing.T) {
	parcel := RelatedRoleOrder(discoverydomain.EntityKindParcel)
	if parcel[0] != RelatedRoleImpact || parcel[1] != RelatedRoleContext {
		t.Fatalf("parcel role order = %v", parcel)
	}
	planning := RelatedRoleOrder(discoverydomain.EntityKindPlanningRegion)
	if planning[0] != RelatedRoleContext || planning[1] != RelatedRoleContained || planning[2] != RelatedRoleOverlap {
		t.Fatalf("planning role order = %v", planning)
	}
}

func TestResolveRelatedRoleSeparatesImpactContextAndOverlap(t *testing.T) {
	if got := ResolveRelatedRole(discoverydomain.EntityKindParcel, RelationshipContainsPrimary); got != RelatedRoleImpact {
		t.Fatalf("parcel contains-primary role = %q", got)
	}
	if got := ResolveRelatedRole(discoverydomain.EntityKindParcel, RelationshipAdministrativeContext); got != RelatedRoleContext {
		t.Fatalf("parcel administrative role = %q", got)
	}
	if got := ResolveRelatedRole(discoverydomain.EntityKindPlanningRegion, RelationshipIntersectsPrimary); got != RelatedRoleOverlap {
		t.Fatalf("planning intersects role = %q", got)
	}
}
