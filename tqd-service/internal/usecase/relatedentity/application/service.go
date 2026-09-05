package application

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"math"
	"sort"
	"strings"
	"time"

	discoverydomain "tqd/internal/domain/discovery/model"
	relateddomain "tqd/internal/domain/relatedentity/model"
	relatedports "tqd/internal/usecase/relatedentity/ports"
)

type Service struct {
	repository relatedports.Repository
}

func NewService(repository relatedports.Repository) *Service {
	return &Service{repository: repository}
}

type providerKind string

const (
	providerPlanning       providerKind = "planning"
	providerAdministrative providerKind = "administrative"
	providerParcel         providerKind = "parcel"
	providerPOI            providerKind = "poi"
)

// providerPolicyForKind is the explicit business policy for the switcher.
// Direct planning-project membership is useful without querying arbitrary
// parcels/POIs around the project centroid; other spatial entities use the
// complete bounded projection.
func providerPolicyForKind(kind discoverydomain.EntityKind) []providerKind {
	switch kind {
	case discoverydomain.EntityKindPlanningProject, discoverydomain.EntityKindPlanningMap:
		return []providerKind{providerPlanning, providerAdministrative}
	case discoverydomain.EntityKindParcel,
		discoverydomain.EntityKindPlanningRegion,
		discoverydomain.EntityKindAdministrativeUnit,
		discoverydomain.EntityKindPOI:
		return []providerKind{
			providerPlanning,
			providerAdministrative,
			providerParcel,
			providerPOI,
		}
	default:
		return nil
	}
}

func providerPolicyForProjection(
	kind discoverydomain.EntityKind,
	projection relateddomain.ProjectionPolicy,
) []providerKind {
	if len(projection.AllowedGroups) == 0 {
		return providerPolicyForKind(kind)
	}

	providers := make([]providerKind, 0, len(projection.AllowedGroups))
	for _, group := range projection.AllowedGroups {
		switch strings.TrimSpace(group) {
		case "planning":
			providers = append(providers, providerPlanning)
		case "administrative":
			providers = append(providers, providerAdministrative)
		case "parcel":
			providers = append(providers, providerParcel)
		case "poi":
			providers = append(providers, providerPOI)
		}
	}
	return providers
}

func providerBudget(kind providerKind, resultLimit int) int {
	switch kind {
	case providerPlanning:
		return minInt(resultLimit, 4)
	case providerAdministrative:
		return minInt(resultLimit, 2)
	case providerParcel, providerPOI:
		return minInt(resultLimit, 8)
	default:
		return 0
	}
}

// providerTimeout keeps one slow enrichment source from blocking the whole
// Related projection. Selection/overview never uses this path; Related may
// return partial data with explicit warning codes.
func providerTimeout(kind providerKind) time.Duration {
	// Related is explicitly user-triggered and isolated from the selection /
	// overview critical path. These budgets are intentionally long enough for
	// indexed PostGIS queries on production data, while still bounding one slow
	// optional provider. They should be calibrated from provider latency metrics.
	switch kind {
	case providerAdministrative:
		return 800 * time.Millisecond
	case providerPlanning:
		return 1500 * time.Millisecond
	case providerParcel:
		return 1200 * time.Millisecond
	case providerPOI:
		return 900 * time.Millisecond
	default:
		return 1000 * time.Millisecond
	}
}

func (s *Service) GetRelatedEntities(ctx context.Context, request relateddomain.Request) (*relateddomain.Response, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	ref, err := parsePrimaryRef(request.Key)
	if err != nil {
		return nil, err
	}

	radiusMeters, radiusWarnings, err := normalizeRadius(request.RadiusMeters)
	if err != nil {
		return nil, err
	}
	resultLimit, limitWarnings, err := normalizeLimit(request.Limit)
	if err != nil {
		return nil, err
	}

	primary, err := s.repository.GetEntity(ctx, ref)
	if err != nil {
		return nil, err
	}
	if primary == nil {
		return nil, relateddomain.ErrEntityNotFound
	}
	if !normalizeEntityCandidate(primary) {
		return nil, errors.New("related entity repository returned an invalid primary entity")
	}

	projection := relateddomain.ResolveProjectionPolicy(primary)
	projectionWarnings := projectionWarningCodes(projection)

	response := &relateddomain.Response{
		RequestID:           newRequestID(),
		Primary:             primary,
		Candidates:          []relateddomain.Candidate{},
		Warnings:            append(append(radiusWarnings, limitWarnings...), projectionWarnings...),
		AppliedRadiusMeters: radiusMeters,
		AppliedLimit:        resultLimit,
		Projection:          projection,
	}

	focus, ok := candidateFocusPoint(primary)
	if !ok {
		response.Partial = true
		response.Warnings = append(response.Warnings, "related_entity_focus_unavailable")
		finalizeProjection(response, nil)
		return response, nil
	}

	query := relateddomain.Query{
		Primary:      primary.Entity,
		Focus:        focus,
		RadiusMeters: radiusMeters,
		// Cards only need identity, relation, evidence and source summary.
		// Geometry is hydrated only after the user selects/previews a candidate.
		IncludeGeometryPreview: false,
	}
	if primary.Location != nil {
		query.ProvinceCode = strings.TrimSpace(primary.Location.ProvinceCode)
		query.WardCode = strings.TrimSpace(primary.Location.WardCode)
	}

	providers := s.providers(providerPolicyForProjection(primary.Entity.Kind, projection), resultLimit)
	if len(providers) == 0 {
		response.Partial = true
		response.Warnings = append(response.Warnings, "related_entity_kind_unsupported")
		finalizeProjection(response, nil)
		return response, nil
	}

	type providerResult struct {
		name  providerKind
		limit int
		items []relateddomain.Candidate
		err   error
	}

	results := make(chan providerResult, len(providers))
	for _, provider := range providers {
		go func(provider providerSpec) {
			providerCtx, cancel := context.WithTimeout(ctx, providerTimeout(provider.name))
			defer cancel()
			providerQuery := query
			providerQuery.Limit = provider.limit
			items, loadErr := provider.load(providerCtx, providerQuery)
			select {
			case results <- providerResult{name: provider.name, limit: provider.limit, items: items, err: loadErr}:
			case <-ctx.Done():
			}
		}(provider)
	}

	all := make([]relateddomain.Candidate, 0, resultLimit*2)
	providerStats := make(map[string]providerProjectionStat, len(providers))
	failedProviders := 0
	for range providers {
		var result providerResult
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case result = <-results:
		}

		if result.err != nil {
			if errors.Is(result.err, context.Canceled) && ctx.Err() != nil {
				return nil, ctx.Err()
			}
			failedProviders++
			warningSuffix := "_unavailable"
			if errors.Is(result.err, context.DeadlineExceeded) {
				warningSuffix = "_timeout"
			}
			response.Warnings = append(response.Warnings, "related_"+string(result.name)+warningSuffix)
			continue
		}
		providerStats[string(result.name)] = providerProjectionStat{
			returnedCount: len(result.items),
			limited:       result.limit > 0 && len(result.items) >= result.limit,
		}
		all = append(all, result.items...)
	}

	if failedProviders == len(providers) {
		// Related is enrichment. The primary entity remains valid even when
		// every optional provider is unavailable, so return a localized partial
		// projection instead of turning the whole Quick Overview into an error.
		response.Partial = true
		response.Warnings = append(response.Warnings, "related_all_providers_unavailable")
		finalizeProjection(response, providerStats)
		return response, nil
	}

	response.Partial = failedProviders > 0
	response.Candidates = selectCandidates(primary.Entity.Kind, primary.Entity.Key, all, resultLimit)
	finalizeProjection(response, providerStats)
	return response, nil
}

type providerProjectionStat struct {
	returnedCount int
	limited       bool
}

func projectionWarningCodes(projection relateddomain.ProjectionPolicy) []string {
	warnings := make([]string, 0, 4+len(projection.SuppressedGroups))
	switch projection.DisplayMode {
	case relateddomain.DisplayModeGroups:
		warnings = append(warnings, "related_projection_grouped")
	case relateddomain.DisplayModeRestricted:
		warnings = append(warnings, "related_projection_restricted")
	}
	if projection.Reason != relateddomain.ProjectionReasonNone {
		warnings = append(warnings, "related_"+string(projection.Reason))
	}
	for _, group := range projection.SuppressedGroups {
		group = strings.TrimSpace(group)
		if group != "" {
			warnings = append(warnings, "related_group_"+group+"_suppressed")
		}
	}
	return warnings
}

func finalizeProjection(
	response *relateddomain.Response,
	providerStats map[string]providerProjectionStat,
) {
	if response == nil {
		return
	}

	if response.Partial {
		response.Projection.Coverage = relateddomain.CoverageBoundedPartial
	} else if response.Projection.DisplayMode == relateddomain.DisplayModeItems {
		response.Projection.Coverage = relateddomain.CoverageBoundedComplete
	}

	groups := summarizeProjectionGroups(response.Candidates, providerStats)
	roles := summarizeRelatedRoles(response.Projection, response.Primary, response.Candidates, providerStats)
	response.Groups = groups
	response.Roles = roles
	response.Warnings = uniqueSortedStrings(response.Warnings)
	response.Primary = cloneEntityCandidate(response.Primary)
	if response.Primary == nil {
		return
	}
	if response.Primary.Attributes == nil {
		response.Primary.Attributes = map[string]any{}
	}

	groupValues := make([]map[string]any, 0, len(groups))
	for _, group := range groups {
		groupValues = append(groupValues, map[string]any{
			"key":           group.Key,
			"returnedCount": group.ReturnedCount,
			"limited":       group.Limited,
		})
	}
	roleValues := make([]map[string]any, 0, len(roles))
	for _, role := range roles {
		roleValues = append(roleValues, map[string]any{
			"key":           string(role.Key),
			"returnedCount": role.ReturnedCount,
			"limited":       role.Limited,
		})
	}
	roleOrder := make([]string, 0)
	if response.Primary != nil {
		for _, role := range relateddomain.RelatedRoleOrder(response.Primary.Entity.Kind) {
			roleOrder = append(roleOrder, string(role))
		}
	}

	response.Primary.Attributes["relatedProjection"] = map[string]any{
		"schemaVersion":   2,
		"displayMode":     string(response.Projection.DisplayMode),
		"coverage":        string(response.Projection.Coverage),
		"strategy":        string(response.Projection.Strategy),
		"reasonCode":      string(response.Projection.Reason),
		"scopeAreaSqm":    response.Projection.Scope.AreaSquareMeters,
		"scopeAreaSource": response.Projection.Scope.AreaSource,
		"scopeScale":      string(response.Projection.Scope.Scale),
		"groups":          groupValues,
		"roles":           roleValues,
		"roleOrder":       roleOrder,
		"suppressedGroups": append(
			[]string(nil),
			response.Projection.SuppressedGroups...,
		),
	}
}

func summarizeRelatedRoles(
	projection relateddomain.ProjectionPolicy,
	primary *discoverydomain.EntityCandidate,
	candidates []relateddomain.Candidate,
	providerStats map[string]providerProjectionStat,
) []relateddomain.RoleSummary {
	if primary == nil {
		return []relateddomain.RoleSummary{}
	}
	counts := make(map[relateddomain.RelatedRole]int)
	limited := make(map[relateddomain.RelatedRole]bool)
	for _, candidate := range candidates {
		role := candidate.Role
		if !role.IsValid() {
			role = relateddomain.ResolveRelatedRole(primary.Entity.Kind, candidate.Relationship)
		}
		counts[role]++
		group := candidateGroup(candidate.Entity.Entity.Kind)
		if stat := providerStats[group]; stat.limited {
			limited[role] = true
		}
	}

	roles := make([]relateddomain.RoleSummary, 0, len(counts))
	for _, role := range relateddomain.RelatedRoleOrder(primary.Entity.Kind) {
		count := counts[role]
		if count <= 0 {
			continue
		}
		roles = append(roles, relateddomain.RoleSummary{
			Key:           role,
			ReturnedCount: count,
			Limited:       limited[role] || projection.DisplayMode != relateddomain.DisplayModeItems,
		})
	}
	return roles
}

func summarizeProjectionGroups(
	candidates []relateddomain.Candidate,
	providerStats map[string]providerProjectionStat,
) []relateddomain.GroupSummary {
	counts := map[string]int{}
	for _, candidate := range candidates {
		counts[candidateGroup(candidate.Entity.Entity.Kind)]++
	}

	order := []string{"planning", "administrative", "parcel", "poi", "other"}
	groups := make([]relateddomain.GroupSummary, 0, len(counts))
	for _, key := range order {
		count := counts[key]
		if count <= 0 {
			continue
		}
		stat := providerStats[key]
		groups = append(groups, relateddomain.GroupSummary{
			Key:           key,
			ReturnedCount: count,
			Limited:       stat.limited || stat.returnedCount > count,
		})
	}
	return groups
}

func cloneEntityCandidate(
	candidate *discoverydomain.EntityCandidate,
) *discoverydomain.EntityCandidate {
	if candidate == nil {
		return nil
	}
	cloned := *candidate
	cloned.Attributes = make(map[string]any, len(candidate.Attributes)+1)
	for key, value := range candidate.Attributes {
		cloned.Attributes[key] = value
	}
	cloned.Match.ReasonCodes = append([]string(nil), candidate.Match.ReasonCodes...)
	cloned.Rank.ReasonCodes = append([]string(nil), candidate.Rank.ReasonCodes...)
	return &cloned
}

type providerSpec struct {
	name  providerKind
	limit int
	load  func(context.Context, relateddomain.Query) ([]relateddomain.Candidate, error)
}

func (s *Service) providers(policy []providerKind, resultLimit int) []providerSpec {
	providers := make([]providerSpec, 0, len(policy))
	for _, kind := range policy {
		limit := providerBudget(kind, resultLimit)
		if limit <= 0 {
			continue
		}
		var load func(context.Context, relateddomain.Query) ([]relateddomain.Candidate, error)
		switch kind {
		case providerPlanning:
			load = s.repository.FindPlanningEntities
		case providerAdministrative:
			load = s.repository.FindAdministrativeUnits
		case providerParcel:
			load = s.repository.FindParcels
		case providerPOI:
			load = s.repository.FindPOIs
		}
		if load != nil {
			providers = append(providers, providerSpec{name: kind, limit: limit, load: load})
		}
	}
	return providers
}

func parsePrimaryRef(key string) (discoverydomain.EntityRef, error) {
	key = strings.TrimSpace(key)
	if key == "" {
		return discoverydomain.EntityRef{}, &relateddomain.ValidationError{
			Field:   "key",
			Message: "canonical entity key is required",
		}
	}
	ref, err := discoverydomain.ParseEntityKey(key)
	if err != nil {
		return discoverydomain.EntityRef{}, &relateddomain.ValidationError{
			Field:   "key",
			Message: err.Error(),
		}
	}
	return ref, nil
}

func normalizeRadius(value float64) (float64, []string, error) {
	if math.IsNaN(value) || math.IsInf(value, 0) || value < 0 {
		return 0, nil, &relateddomain.ValidationError{
			Field:   "radiusMeters",
			Message: "must be a finite number greater than or equal to zero",
		}
	}
	if value == 0 {
		return relateddomain.DefaultRadiusMeters, nil, nil
	}
	if value > relateddomain.MaxRadiusMeters {
		return relateddomain.MaxRadiusMeters, []string{"related_radius_clamped_to_max"}, nil
	}
	return value, nil, nil
}

func normalizeLimit(value int) (int, []string, error) {
	if value < 0 {
		return 0, nil, &relateddomain.ValidationError{
			Field:   "limit",
			Message: "must be greater than or equal to zero",
		}
	}
	if value == 0 {
		return relateddomain.DefaultLimit, nil, nil
	}
	if value > relateddomain.MaxLimit {
		return relateddomain.MaxLimit, []string{"related_limit_clamped_to_max"}, nil
	}
	return value, nil, nil
}

func minInt(left, right int) int {
	if left < right {
		return left
	}
	return right
}

func candidateFocusPoint(candidate *discoverydomain.EntityCandidate) (discoverydomain.Point, bool) {
	if candidate == nil || candidate.Spatial == nil {
		return discoverydomain.Point{}, false
	}
	if candidate.Spatial.Centroid != nil && validPoint(*candidate.Spatial.Centroid) {
		return *candidate.Spatial.Centroid, true
	}
	if candidate.Spatial.Bounds != nil {
		bounds := candidate.Spatial.Bounds
		point := discoverydomain.Point{
			Latitude:  (bounds.MinLatitude + bounds.MaxLatitude) / 2,
			Longitude: (bounds.MinLongitude + bounds.MaxLongitude) / 2,
		}
		if validPoint(point) {
			return point, true
		}
	}
	return discoverydomain.Point{}, false
}

func validPoint(point discoverydomain.Point) bool {
	return point.Latitude >= -90 && point.Latitude <= 90 &&
		point.Longitude >= -180 && point.Longitude <= 180 &&
		!math.IsNaN(point.Latitude) && !math.IsNaN(point.Longitude) &&
		!math.IsInf(point.Latitude, 0) && !math.IsInf(point.Longitude, 0)
}

func selectCandidates(
	primaryKind discoverydomain.EntityKind,
	primaryKey string,
	candidates []relateddomain.Candidate,
	resultLimit int,
) []relateddomain.Candidate {
	if resultLimit <= 0 {
		return []relateddomain.Candidate{}
	}

	normalized := make([]relateddomain.Candidate, 0, len(candidates))
	for index := range candidates {
		if !normalizeEntityCandidate(&candidates[index].Entity) {
			continue
		}
		decorateCandidate(primaryKind, &candidates[index])
		normalized = append(normalized, candidates[index])
	}
	candidates = normalized
	sort.SliceStable(candidates, func(i, j int) bool {
		leftPriority := relateddomain.RelatedRolePriority(primaryKind, candidates[i].Role)
		rightPriority := relateddomain.RelatedRolePriority(primaryKind, candidates[j].Role)
		if leftPriority != rightPriority {
			return leftPriority < rightPriority
		}
		if candidates[i].Score != candidates[j].Score {
			return candidates[i].Score > candidates[j].Score
		}
		if candidates[i].DistanceMeters != candidates[j].DistanceMeters {
			return candidates[i].DistanceMeters < candidates[j].DistanceMeters
		}
		return candidates[i].Entity.Entity.Key < candidates[j].Entity.Entity.Key
	})

	seen := map[string]struct{}{strings.TrimSpace(primaryKey): {}}
	selected := make([]relateddomain.Candidate, 0, resultLimit)
	appendCandidate := func(candidate relateddomain.Candidate) bool {
		key := strings.TrimSpace(candidate.Entity.Entity.Key)
		if key == "" || len(selected) >= resultLimit {
			return false
		}
		if _, exists := seen[key]; exists {
			return false
		}
		seen[key] = struct{}{}
		selected = append(selected, candidate)
		return true
	}

	// Keep a pure, explainable relevance order. Diversity belongs to role
	// lanes/filters in the projection UI, not to a hidden interleaving step.
	for _, candidate := range candidates {
		appendCandidate(candidate)
	}
	return selected
}

func candidateGroup(kind discoverydomain.EntityKind) string {
	switch kind {
	case discoverydomain.EntityKindPlanningRegion,
		discoverydomain.EntityKindPlanningProject,
		discoverydomain.EntityKindPlanningMap:
		return "planning"
	case discoverydomain.EntityKindAdministrativeUnit:
		return "administrative"
	case discoverydomain.EntityKindParcel:
		return "parcel"
	case discoverydomain.EntityKindPOI:
		return "poi"
	default:
		return "other"
	}
}

func normalizeEntityCandidate(candidate *discoverydomain.EntityCandidate) bool {
	if candidate == nil || !candidate.Entity.Kind.IsValid() {
		return false
	}
	id := strings.TrimSpace(candidate.Entity.ID)
	if id == "" {
		return false
	}
	candidate.Entity = discoverydomain.NewEntityRef(candidate.Entity.Kind, id)
	return true
}

func decorateCandidate(primaryKind discoverydomain.EntityKind, candidate *relateddomain.Candidate) {
	if candidate == nil {
		return
	}
	if !candidate.Relationship.IsValid() {
		candidate.Relationship = relateddomain.RelationshipUnknown
		candidate.ReasonCodes = append(candidate.ReasonCodes, "relationship_unclassified")
	}
	candidate.Role = relateddomain.ResolveRelatedRole(primaryKind, candidate.Relationship)
	if candidate.Entity.Attributes == nil {
		candidate.Entity.Attributes = map[string]any{}
	}
	candidate.Entity.Attributes["relationship"] = string(candidate.Relationship)
	candidate.Entity.Attributes["relatedRole"] = string(candidate.Role)
	copyCandidateSourceMetadata(&candidate.Entity)
	if candidate.DistanceMeters >= 0 && candidate.DistanceMeters < math.MaxFloat64 &&
		!math.IsInf(candidate.DistanceMeters, 0) && !math.IsNaN(candidate.DistanceMeters) {
		candidate.Entity.Attributes["distanceMeters"] = candidate.DistanceMeters
	}
	candidate.Entity.Rank.Score = candidate.Score
	candidate.Entity.Rank.ReasonCodes = uniqueSortedStrings(append(
		candidate.Entity.Rank.ReasonCodes,
		candidate.ReasonCodes...,
	))
	candidate.Entity.Match.Type = "related_entity"
	candidate.Entity.Match.ReasonCodes = uniqueSortedStrings(append(
		candidate.Entity.Match.ReasonCodes,
		string(candidate.Relationship),
	))
}

func copyCandidateSourceMetadata(candidate *discoverydomain.EntityCandidate) {
	if candidate == nil || candidate.Attributes == nil {
		return
	}
	set := func(key, value string) {
		value = strings.TrimSpace(value)
		if value != "" {
			candidate.Attributes[key] = value
		}
	}
	set("sourceSystem", candidate.Source.System)
	set("sourceDataset", candidate.Source.Dataset)
	set("sourceAuthority", candidate.Source.Authority)
	set("sourceUpdatedAt", candidate.Source.UpdatedAt)
	set("sourceDataQuality", candidate.Source.DataQuality)
}

func uniqueSortedStrings(values []string) []string {
	if len(values) == 0 {
		return []string{}
	}
	seen := make(map[string]struct{}, len(values))
	out := make([]string, 0, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}
		if _, exists := seen[value]; exists {
			continue
		}
		seen[value] = struct{}{}
		out = append(out, value)
	}
	sort.Strings(out)
	return out
}

func newRequestID() string {
	var raw [16]byte
	if _, err := rand.Read(raw[:]); err != nil {
		return "related-entity-request"
	}
	return hex.EncodeToString(raw[:])
}
