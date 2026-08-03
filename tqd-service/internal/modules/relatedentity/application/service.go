package application

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"math"
	"sort"
	"strings"
	"sync"

	discoverydomain "tqd/internal/modules/discovery/domain"
	relateddomain "tqd/internal/modules/relatedentity/domain"
	relatedports "tqd/internal/modules/relatedentity/ports"
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
	case discoverydomain.EntityKindPlanningProject:
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

	response := &relateddomain.Response{
		RequestID:           newRequestID(),
		Primary:             primary,
		Candidates:          []relateddomain.Candidate{},
		Warnings:            append(radiusWarnings, limitWarnings...),
		AppliedRadiusMeters: radiusMeters,
		AppliedLimit:        resultLimit,
	}

	focus, ok := candidateFocusPoint(primary)
	if !ok {
		response.Partial = true
		response.Warnings = append(response.Warnings, "related_entity_focus_unavailable")
		response.Warnings = uniqueSortedStrings(response.Warnings)
		return response, nil
	}

	query := relateddomain.Query{
		Primary:                primary.Entity,
		Focus:                  focus,
		RadiusMeters:           radiusMeters,
		IncludeGeometryPreview: true,
	}
	if primary.Location != nil {
		query.ProvinceCode = strings.TrimSpace(primary.Location.ProvinceCode)
		query.WardCode = strings.TrimSpace(primary.Location.WardCode)
	}

	providers := s.providers(providerPolicyForKind(primary.Entity.Kind), resultLimit)
	if len(providers) == 0 {
		response.Partial = true
		response.Warnings = append(response.Warnings, "related_entity_kind_unsupported")
		response.Warnings = uniqueSortedStrings(response.Warnings)
		return response, nil
	}

	type providerResult struct {
		name  providerKind
		items []relateddomain.Candidate
		err   error
	}

	results := make(chan providerResult, len(providers))
	var wait sync.WaitGroup
	for _, provider := range providers {
		wait.Add(1)
		go func(provider providerSpec) {
			defer wait.Done()
			providerQuery := query
			providerQuery.Limit = provider.limit
			items, loadErr := provider.load(ctx, providerQuery)
			results <- providerResult{name: provider.name, items: items, err: loadErr}
		}(provider)
	}
	wait.Wait()
	close(results)

	all := make([]relateddomain.Candidate, 0, resultLimit*2)
	failedProviders := 0
	for result := range results {
		if result.err != nil {
			if errors.Is(result.err, context.Canceled) || errors.Is(result.err, context.DeadlineExceeded) {
				return nil, result.err
			}
			failedProviders++
			response.Warnings = append(response.Warnings, "related_"+string(result.name)+"_unavailable")
			continue
		}
		all = append(all, result.items...)
	}

	if failedProviders == len(providers) {
		return nil, relateddomain.ErrProvidersFailed
	}

	response.Partial = failedProviders > 0
	response.Candidates = selectCandidates(primary.Entity.Key, all, resultLimit)
	response.Warnings = uniqueSortedStrings(response.Warnings)
	return response, nil
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

func selectCandidates(primaryKey string, candidates []relateddomain.Candidate, resultLimit int) []relateddomain.Candidate {
	if resultLimit <= 0 {
		return []relateddomain.Candidate{}
	}

	normalized := make([]relateddomain.Candidate, 0, len(candidates))
	for index := range candidates {
		if !normalizeEntityCandidate(&candidates[index].Entity) {
			continue
		}
		decorateCandidate(&candidates[index])
		normalized = append(normalized, candidates[index])
	}
	candidates = normalized
	sort.SliceStable(candidates, func(i, j int) bool {
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

	// One representative per business group keeps the switcher useful in dense
	// parcel/POI areas. Remaining slots are then filled by global score.
	priorityGroups := []string{"planning", "administrative", "parcel", "poi"}
	for _, group := range priorityGroups {
		for _, candidate := range candidates {
			if candidateGroup(candidate.Entity.Entity.Kind) == group && appendCandidate(candidate) {
				break
			}
		}
	}
	for _, candidate := range candidates {
		appendCandidate(candidate)
	}
	return selected
}

func candidateGroup(kind discoverydomain.EntityKind) string {
	switch kind {
	case discoverydomain.EntityKindPlanningRegion, discoverydomain.EntityKindPlanningProject:
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

func decorateCandidate(candidate *relateddomain.Candidate) {
	if candidate == nil {
		return
	}
	if !candidate.Relationship.IsValid() {
		candidate.Relationship = relateddomain.RelationshipNearby
	}
	if candidate.Entity.Attributes == nil {
		candidate.Entity.Attributes = map[string]any{}
	}
	candidate.Entity.Attributes["relationship"] = string(candidate.Relationship)
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
