package domain

import (
	"testing"

	discoverydomain "tqd/internal/domain/discovery/model"
)

func TestRelatedRoleOrderKeepsDecisionIntentStable(t *testing.T) {
	cases := []struct {
		name string
		kind discoverydomain.EntityKind
		want []RelatedRole
	}{
		{name: "parcel", kind: discoverydomain.EntityKindParcel, want: []RelatedRole{RelatedRoleImpact, RelatedRoleContext, RelatedRoleAdjacent}},
		{name: "planning", kind: discoverydomain.EntityKindPlanningRegion, want: []RelatedRole{RelatedRoleContext, RelatedRoleContained, RelatedRoleOverlap}},
		{name: "administrative", kind: discoverydomain.EntityKindAdministrativeUnit, want: []RelatedRole{RelatedRoleContained, RelatedRoleContext, RelatedRoleOverlap}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := RelatedRoleOrder(tc.kind)
			if len(got) < len(tc.want) {
				t.Fatalf("order = %v, want prefix %v", got, tc.want)
			}
			for index, role := range tc.want {
				if got[index] != role {
					t.Fatalf("order[%d] = %q, want %q", index, got[index], role)
				}
			}
		})
	}
}

func TestResolveRelatedRoleDoesNotConfuseSpatialExistenceWithUserValue(t *testing.T) {
	cases := []struct {
		name         string
		primary      discoverydomain.EntityKind
		relationship RelationshipKind
		want         RelatedRole
	}{
		{name: "planning covers parcel is impact", primary: discoverydomain.EntityKindParcel, relationship: RelationshipContainsPrimary, want: RelatedRoleImpact},
		{name: "administrative ownership is context", primary: discoverydomain.EntityKindParcel, relationship: RelationshipAdministrativeContext, want: RelatedRoleContext},
		{name: "planning intersection is overlap", primary: discoverydomain.EntityKindPlanningRegion, relationship: RelationshipIntersectsPrimary, want: RelatedRoleOverlap},
		{name: "planning child is contained", primary: discoverydomain.EntityKindPlanningRegion, relationship: RelationshipContainedByPrimary, want: RelatedRoleContained},
		{name: "unknown remains unknown", primary: discoverydomain.EntityKindParcel, relationship: RelationshipUnknown, want: RelatedRoleUnknown},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := ResolveRelatedRole(tc.primary, tc.relationship); got != tc.want {
				t.Fatalf("role = %q, want %q", got, tc.want)
			}
		})
	}
}
