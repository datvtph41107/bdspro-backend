package domain

import (
	"errors"
	"fmt"

	discoverydomain "tqd/internal/domain/discovery/model"
)

const (
	DefaultRadiusMeters = 250.0
	MaxRadiusMeters     = 5000.0
	DefaultLimit        = 12
	MaxLimit            = 50
)

type RelationshipKind string

const (
	RelationshipContainsPrimary       RelationshipKind = "contains_primary"
	RelationshipContainedByPrimary    RelationshipKind = "contained_by_primary"
	RelationshipAdministrativeContext RelationshipKind = "administrative_context"
	RelationshipMemberOfPrimary       RelationshipKind = "member_of_primary"
	RelationshipIntersectsPrimary     RelationshipKind = "intersects_primary"
	RelationshipAdjacent              RelationshipKind = "adjacent"
	RelationshipNearby                RelationshipKind = "nearby"
	RelationshipUnknown               RelationshipKind = "unknown"
)

func (value RelationshipKind) IsValid() bool {
	switch value {
	case RelationshipContainsPrimary,
		RelationshipContainedByPrimary,
		RelationshipAdministrativeContext,
		RelationshipMemberOfPrimary,
		RelationshipIntersectsPrimary,
		RelationshipAdjacent,
		RelationshipNearby,
		RelationshipUnknown:
		return true
	default:
		return false
	}
}

var (
	ErrEntityNotFound  = errors.New("related entity primary not found")
	ErrProvidersFailed = errors.New("all related entity providers failed")
)

type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	if e == nil {
		return "invalid related entity request"
	}
	if e.Field == "" {
		return e.Message
	}
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

type Request struct {
	Key          string
	RadiusMeters float64
	Limit        int
}

type Query struct {
	Primary                discoverydomain.EntityRef
	Focus                  discoverydomain.Point
	ProvinceCode           string
	WardCode               string
	RadiusMeters           float64
	Limit                  int
	IncludeGeometryPreview bool
}

type Candidate struct {
	Entity         discoverydomain.EntityCandidate
	Relationship   RelationshipKind
	Role           RelatedRole
	DistanceMeters float64
	Score          float64
	ReasonCodes    []string
}

type Response struct {
	RequestID           string                           `json:"requestId"`
	Primary             *discoverydomain.EntityCandidate `json:"primary,omitempty"`
	Candidates          []Candidate                      `json:"candidates"`
	Partial             bool                             `json:"partial"`
	Warnings            []string                         `json:"warnings,omitempty"`
	AppliedRadiusMeters float64                          `json:"appliedRadiusMeters"`
	AppliedLimit        int                              `json:"appliedLimit"`
	Projection          ProjectionPolicy                 `json:"projection"`
	Groups              []GroupSummary                   `json:"groups,omitempty"`
	Roles               []RoleSummary                    `json:"roles,omitempty"`
}
