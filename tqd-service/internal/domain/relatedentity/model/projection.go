package domain

import (
	"encoding/json"
	"math"
	"strconv"
	"strings"

	discoverydomain "tqd/internal/domain/discovery/model"
)

// DisplayMode tells clients how a related projection should be presented.
// It is deliberately independent from transport status: a grouped or
// restricted response can still be a successful response.
type DisplayMode string

const (
	DisplayModeItems      DisplayMode = "items"
	DisplayModeGroups     DisplayMode = "groups"
	DisplayModeRestricted DisplayMode = "restricted"
)

type CoverageKind string

const (
	// Bounded coverage describes completeness inside the declared Related query
	// envelope, never completeness across all possible datasets in reality.
	CoverageBoundedComplete CoverageKind = "bounded_complete"
	CoverageBoundedPartial  CoverageKind = "bounded_partial"
	CoverageAggregated      CoverageKind = "aggregated"
	CoverageContextOnly     CoverageKind = "context_only"
)

type QueryStrategy string

const (
	StrategyExactList         QueryStrategy = "exact_list"
	StrategyGroupedSummary    QueryStrategy = "grouped_summary"
	StrategyRestrictedContext QueryStrategy = "restricted_context"
)

type ScopeScale string

const (
	ScopeScaleSmall     ScopeScale = "small"
	ScopeScaleMedium    ScopeScale = "medium"
	ScopeScaleLarge     ScopeScale = "large"
	ScopeScaleVeryLarge ScopeScale = "very_large"
	ScopeScaleUnknown   ScopeScale = "unknown"
)

type ProjectionReason string

const (
	ProjectionReasonNone           ProjectionReason = ""
	ProjectionReasonScopeLarge     ProjectionReason = "scope_too_large_for_items"
	ProjectionReasonViewportNeeded ProjectionReason = "viewport_or_smaller_scope_required"
	ProjectionReasonProjectSummary ProjectionReason = "planning_project_summary"
)

type ScopeSummary struct {
	AreaSquareMeters float64
	AreaSource       string
	Scale            ScopeScale
}

type GroupSummary struct {
	Key           string `json:"key"`
	ReturnedCount int    `json:"returnedCount"`
	Limited       bool   `json:"limited"`
}

// ProjectionPolicy is the business decision made before provider queries run.
// It prevents a large polygon from being treated like a parcel-sized target.
type ProjectionPolicy struct {
	DisplayMode      DisplayMode
	Coverage         CoverageKind
	Strategy         QueryStrategy
	Reason           ProjectionReason
	Scope            ScopeSummary
	AllowedGroups    []string
	SuppressedGroups []string
}

const (
	planningGroupThresholdSquareMeters        = 5_000_000.0   // 5 km²
	planningRestrictedThresholdSquareMeters   = 250_000_000.0 // 250 km²
	administrativeGroupThresholdSquareMeters  = 25_000_000.0  // 25 km²
	administrativeRestrictedThresholdSqMeters = 1_000_000_000.0
)

func ResolveProjectionPolicy(primary *discoverydomain.EntityCandidate) ProjectionPolicy {
	scope := summarizeScope(primary)
	policy := ProjectionPolicy{
		DisplayMode:   DisplayModeItems,
		Coverage:      CoverageBoundedComplete,
		Strategy:      StrategyExactList,
		Scope:         scope,
		AllowedGroups: []string{"planning", "administrative", "parcel", "poi"},
	}
	if primary == nil {
		return policy
	}

	switch primary.Entity.Kind {
	case discoverydomain.EntityKindPlanningProject, discoverydomain.EntityKindPlanningMap:
		policy.DisplayMode = DisplayModeGroups
		policy.Coverage = CoverageAggregated
		policy.Strategy = StrategyGroupedSummary
		policy.Reason = ProjectionReasonProjectSummary
		policy.AllowedGroups = []string{"planning", "administrative"}
		policy.SuppressedGroups = []string{"parcel", "poi"}
		return policy

	case discoverydomain.EntityKindPlanningRegion:
		if scope.AreaSquareMeters >= planningRestrictedThresholdSquareMeters {
			policy.DisplayMode = DisplayModeRestricted
			policy.Coverage = CoverageContextOnly
			policy.Strategy = StrategyRestrictedContext
			policy.Reason = ProjectionReasonViewportNeeded
			policy.AllowedGroups = []string{"planning", "administrative"}
			policy.SuppressedGroups = []string{"parcel", "poi"}
			return policy
		}
		if scope.AreaSquareMeters >= planningGroupThresholdSquareMeters {
			policy.DisplayMode = DisplayModeGroups
			policy.Coverage = CoverageAggregated
			policy.Strategy = StrategyGroupedSummary
			policy.Reason = ProjectionReasonScopeLarge
			policy.AllowedGroups = []string{"planning", "administrative"}
			policy.SuppressedGroups = []string{"parcel", "poi"}
			return policy
		}

	case discoverydomain.EntityKindAdministrativeUnit:
		if scope.AreaSquareMeters >= administrativeRestrictedThresholdSqMeters {
			policy.DisplayMode = DisplayModeRestricted
			policy.Coverage = CoverageContextOnly
			policy.Strategy = StrategyRestrictedContext
			policy.Reason = ProjectionReasonViewportNeeded
			policy.AllowedGroups = []string{"planning", "administrative"}
			policy.SuppressedGroups = []string{"parcel", "poi"}
			return policy
		}
		if scope.AreaSquareMeters >= administrativeGroupThresholdSquareMeters {
			policy.DisplayMode = DisplayModeGroups
			policy.Coverage = CoverageAggregated
			policy.Strategy = StrategyGroupedSummary
			policy.Reason = ProjectionReasonScopeLarge
			policy.AllowedGroups = []string{"planning", "administrative"}
			policy.SuppressedGroups = []string{"parcel", "poi"}
			return policy
		}
	}

	return policy
}

func summarizeScope(primary *discoverydomain.EntityCandidate) ScopeSummary {
	if primary == nil {
		return ScopeSummary{Scale: ScopeScaleUnknown}
	}

	if area, ok := numericAttribute(primary.Attributes, "areaSqm"); ok && area > 0 {
		return ScopeSummary{
			AreaSquareMeters: area,
			AreaSource:       "entity_attribute",
			Scale:            scaleForArea(primary.Entity.Kind, area),
		}
	}

	if primary.Spatial != nil && primary.Spatial.Bounds != nil {
		if area := approximateBoundsAreaSquareMeters(primary.Spatial.Bounds); area > 0 {
			return ScopeSummary{
				AreaSquareMeters: area,
				AreaSource:       "bounds_estimate",
				Scale:            scaleForArea(primary.Entity.Kind, area),
			}
		}
	}

	return ScopeSummary{Scale: ScopeScaleUnknown}
}

func scaleForArea(kind discoverydomain.EntityKind, area float64) ScopeScale {
	if kind == discoverydomain.EntityKindAdministrativeUnit {
		switch {
		case area <= 0:
			return ScopeScaleUnknown
		case area < 5_000_000:
			return ScopeScaleSmall
		case area < administrativeGroupThresholdSquareMeters:
			return ScopeScaleMedium
		case area < administrativeRestrictedThresholdSqMeters:
			return ScopeScaleLarge
		default:
			return ScopeScaleVeryLarge
		}
	}

	switch {
	case area <= 0:
		return ScopeScaleUnknown
	case area < 1_000_000:
		return ScopeScaleSmall
	case area < planningGroupThresholdSquareMeters:
		return ScopeScaleMedium
	case area < planningRestrictedThresholdSquareMeters:
		return ScopeScaleLarge
	default:
		return ScopeScaleVeryLarge
	}
}

func approximateBoundsAreaSquareMeters(bounds *discoverydomain.Bounds) float64 {
	if bounds == nil {
		return 0
	}
	latSpan := math.Abs(bounds.MaxLatitude - bounds.MinLatitude)
	lonSpan := math.Abs(bounds.MaxLongitude - bounds.MinLongitude)
	if latSpan <= 0 || lonSpan <= 0 || latSpan > 90 || lonSpan > 180 {
		return 0
	}
	midLat := (bounds.MinLatitude + bounds.MaxLatitude) / 2
	metersPerDegreeLat := 111_320.0
	metersPerDegreeLon := metersPerDegreeLat * math.Cos(midLat*math.Pi/180)
	if metersPerDegreeLon <= 0 {
		return 0
	}
	return latSpan * metersPerDegreeLat * lonSpan * metersPerDegreeLon
}

func numericAttribute(attributes map[string]any, key string) (float64, bool) {
	if attributes == nil {
		return 0, false
	}
	value, exists := attributes[key]
	if !exists {
		return 0, false
	}
	switch typed := value.(type) {
	case float64:
		return finitePositive(typed)
	case float32:
		return finitePositive(float64(typed))
	case int:
		return finitePositive(float64(typed))
	case int32:
		return finitePositive(float64(typed))
	case int64:
		return finitePositive(float64(typed))
	case uint:
		return finitePositive(float64(typed))
	case uint32:
		return finitePositive(float64(typed))
	case uint64:
		return finitePositive(float64(typed))
	case json.Number:
		parsed, err := typed.Float64()
		if err != nil {
			return 0, false
		}
		return finitePositive(parsed)
	case string:
		parsed, err := strconv.ParseFloat(strings.TrimSpace(typed), 64)
		if err != nil {
			return 0, false
		}
		return finitePositive(parsed)
	default:
		return 0, false
	}
}

func finitePositive(value float64) (float64, bool) {
	if value <= 0 || math.IsNaN(value) || math.IsInf(value, 0) {
		return 0, false
	}
	return value, true
}
