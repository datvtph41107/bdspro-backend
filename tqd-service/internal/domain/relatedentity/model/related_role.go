package domain

import discoverydomain "tqd/internal/domain/discovery/model"

// RelatedRole explains why a candidate matters to the current primary entity.
// It is deliberately independent from provider/entity kind so clients can
// filter by user intent instead of exposing a flat technical layer list.
type RelatedRole string

const (
	RelatedRoleImpact    RelatedRole = "impact"
	RelatedRoleContext   RelatedRole = "context"
	RelatedRoleContained RelatedRole = "contained"
	RelatedRoleOverlap   RelatedRole = "overlap"
	RelatedRoleAdjacent  RelatedRole = "adjacent"
	RelatedRoleNearby    RelatedRole = "nearby"
	RelatedRoleUnknown   RelatedRole = "unknown"
)

type RoleSummary struct {
	Key           RelatedRole `json:"key"`
	ReturnedCount int         `json:"returnedCount"`
	Limited       bool        `json:"limited"`
}

func (value RelatedRole) IsValid() bool {
	switch value {
	case RelatedRoleImpact,
		RelatedRoleContext,
		RelatedRoleContained,
		RelatedRoleOverlap,
		RelatedRoleAdjacent,
		RelatedRoleNearby,
		RelatedRoleUnknown:
		return true
	default:
		return false
	}
}

func isPlanningKind(kind discoverydomain.EntityKind) bool {
	return kind == discoverydomain.EntityKindPlanningRegion ||
		kind == discoverydomain.EntityKindPlanningProject ||
		kind == discoverydomain.EntityKindPlanningMap
}

// ResolveRelatedRole converts a spatial/business relationship into a stable
// user-facing role relative to the primary entity.
func ResolveRelatedRole(primaryKind discoverydomain.EntityKind, relationship RelationshipKind) RelatedRole {
	switch relationship {
	case RelationshipAdministrativeContext:
		return RelatedRoleContext
	case RelationshipAdjacent:
		return RelatedRoleAdjacent
	case RelationshipNearby:
		return RelatedRoleNearby
	case RelationshipUnknown:
		return RelatedRoleUnknown
	case RelationshipMemberOfPrimary, RelationshipContainedByPrimary:
		if isPlanningKind(primaryKind) || primaryKind == discoverydomain.EntityKindAdministrativeUnit {
			return RelatedRoleContained
		}
		return RelatedRoleContext
	case RelationshipContainsPrimary:
		if primaryKind == discoverydomain.EntityKindParcel || primaryKind == discoverydomain.EntityKindPOI {
			return RelatedRoleImpact
		}
		return RelatedRoleContext
	case RelationshipIntersectsPrimary:
		if isPlanningKind(primaryKind) || primaryKind == discoverydomain.EntityKindAdministrativeUnit {
			return RelatedRoleOverlap
		}
		return RelatedRoleImpact
	default:
		return RelatedRoleUnknown
	}
}

// RelatedRoleOrder is the decision order used by both ranking and UI metadata.
// It reflects the question users normally ask for each primary entity kind.
func RelatedRoleOrder(primaryKind discoverydomain.EntityKind) []RelatedRole {
	switch {
	case primaryKind == discoverydomain.EntityKindParcel || primaryKind == discoverydomain.EntityKindPOI:
		return []RelatedRole{
			RelatedRoleImpact,
			RelatedRoleContext,
			RelatedRoleAdjacent,
			RelatedRoleNearby,
			RelatedRoleContained,
			RelatedRoleOverlap,
			RelatedRoleUnknown,
		}
	case isPlanningKind(primaryKind):
		return []RelatedRole{
			RelatedRoleContext,
			RelatedRoleContained,
			RelatedRoleOverlap,
			RelatedRoleAdjacent,
			RelatedRoleNearby,
			RelatedRoleImpact,
			RelatedRoleUnknown,
		}
	case primaryKind == discoverydomain.EntityKindAdministrativeUnit:
		return []RelatedRole{
			RelatedRoleContained,
			RelatedRoleContext,
			RelatedRoleOverlap,
			RelatedRoleNearby,
			RelatedRoleAdjacent,
			RelatedRoleImpact,
			RelatedRoleUnknown,
		}
	default:
		return []RelatedRole{
			RelatedRoleContext,
			RelatedRoleImpact,
			RelatedRoleContained,
			RelatedRoleOverlap,
			RelatedRoleAdjacent,
			RelatedRoleNearby,
			RelatedRoleUnknown,
		}
	}
}

func RelatedRolePriority(primaryKind discoverydomain.EntityKind, role RelatedRole) int {
	for index, candidate := range RelatedRoleOrder(primaryKind) {
		if candidate == role {
			return index
		}
	}
	return len(RelatedRoleOrder(primaryKind))
}
