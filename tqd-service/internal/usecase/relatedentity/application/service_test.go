package application

import (
	"context"
	"errors"
	"sort"
	"sync"
	"testing"

	discoverydomain "tqd/internal/domain/discovery/model"
	relateddomain "tqd/internal/domain/relatedentity/model"
)

type fakeRepository struct {
	primary *discoverydomain.EntityCandidate

	planning       []relateddomain.Candidate
	administrative []relateddomain.Candidate
	parcels        []relateddomain.Candidate
	pois           []relateddomain.Candidate

	planningErr       error
	administrativeErr error
	parcelErr         error
	poiErr            error

	mu      sync.Mutex
	queries []relateddomain.Query
}

func (f *fakeRepository) GetEntity(context.Context, discoverydomain.EntityRef) (*discoverydomain.EntityCandidate, error) {
	return f.primary, nil
}

func (f *fakeRepository) record(query relateddomain.Query) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.queries = append(f.queries, query)
}

func (f *fakeRepository) FindPlanningEntities(_ context.Context, query relateddomain.Query) ([]relateddomain.Candidate, error) {
	f.record(query)
	return f.planning, f.planningErr
}

func (f *fakeRepository) FindAdministrativeUnits(_ context.Context, query relateddomain.Query) ([]relateddomain.Candidate, error) {
	f.record(query)
	return f.administrative, f.administrativeErr
}

func (f *fakeRepository) FindParcels(_ context.Context, query relateddomain.Query) ([]relateddomain.Candidate, error) {
	f.record(query)
	return f.parcels, f.parcelErr
}

func (f *fakeRepository) FindPOIs(_ context.Context, query relateddomain.Query) ([]relateddomain.Candidate, error) {
	f.record(query)
	return f.pois, f.poiErr
}

func TestGetRelatedEntitiesClampsAndOrdersByUserFacingRole(t *testing.T) {
	primary := entity(discoverydomain.EntityKindParcel, "10", 10)
	repo := &fakeRepository{
		primary: &primary,
		planning: []relateddomain.Candidate{
			related(entity(discoverydomain.EntityKindPlanningRegion, "21", 21), relateddomain.RelationshipContainsPrimary, 0, 100),
		},
		administrative: []relateddomain.Candidate{
			related(entity(discoverydomain.EntityKindAdministrativeUnit, "ward:31", 31), relateddomain.RelationshipAdministrativeContext, 20, 95),
		},
		parcels: []relateddomain.Candidate{
			related(primary, relateddomain.RelationshipNearby, 0, 999),
			related(entity(discoverydomain.EntityKindParcel, "11", 11), relateddomain.RelationshipAdjacent, 0, 68),
			related(entity(discoverydomain.EntityKindParcel, "12", 12), relateddomain.RelationshipNearby, 30, 45),
		},
		pois: []relateddomain.Candidate{
			related(entity(discoverydomain.EntityKindPOI, "41", 41), relateddomain.RelationshipNearby, 50, 35),
		},
	}

	response, err := NewService(repo).GetRelatedEntities(context.Background(), relateddomain.Request{
		Key:          "parcel:999",
		RadiusMeters: 9000,
		Limit:        100,
	})
	if err != nil {
		t.Fatalf("GetRelatedEntities() error = %v", err)
	}
	if response.Projection.Coverage != relateddomain.CoverageBoundedComplete {
		t.Fatalf("Coverage = %q, want bounded_complete", response.Projection.Coverage)
	}
	if response.AppliedRadiusMeters != relateddomain.MaxRadiusMeters {
		t.Fatalf("AppliedRadiusMeters = %v, want %v", response.AppliedRadiusMeters, relateddomain.MaxRadiusMeters)
	}
	if response.AppliedLimit != relateddomain.MaxLimit {
		t.Fatalf("AppliedLimit = %v, want %v", response.AppliedLimit, relateddomain.MaxLimit)
	}
	assertContains(t, response.Warnings, "related_radius_clamped_to_max")
	assertContains(t, response.Warnings, "related_limit_clamped_to_max")
	if len(response.Candidates) != 5 {
		t.Fatalf("len(Candidates) = %d, want 5", len(response.Candidates))
	}
	wantKeys := []string{
		"planning_region:21",          // impact
		"administrative_unit:ward:31", // context
		"parcel:11",                   // adjacent
		"parcel:12",                   // higher-score nearby item remains ahead in pure relevance order
	}
	for index, want := range wantKeys {
		if response.Candidates[index].Entity.Entity.Key != want {
			t.Fatalf("Candidates[%d].Key = %q, want %q", index, response.Candidates[index].Entity.Entity.Key, want)
		}
	}
	for _, candidate := range response.Candidates {
		if candidate.Entity.Entity.Key == primary.Entity.Key {
			t.Fatalf("primary entity %q must be excluded", primary.Entity.Key)
		}
		if candidate.Entity.Attributes["relationship"] == nil {
			t.Fatalf("candidate %q has no relationship", candidate.Entity.Entity.Key)
		}
		if candidate.Entity.Attributes["relatedRole"] == nil {
			t.Fatalf("candidate %q has no relatedRole", candidate.Entity.Entity.Key)
		}
	}

	repo.mu.Lock()
	queries := append([]relateddomain.Query(nil), repo.queries...)
	repo.mu.Unlock()
	limits := make([]int, 0, len(queries))
	for _, query := range queries {
		if query.Primary.Key != primary.Entity.Key {
			t.Fatalf("provider primary = %q, want canonical %q", query.Primary.Key, primary.Entity.Key)
		}
		if query.RadiusMeters != relateddomain.MaxRadiusMeters {
			t.Fatalf("provider radius = %v, want %v", query.RadiusMeters, relateddomain.MaxRadiusMeters)
		}
		if query.Limit <= 0 || query.Limit > relateddomain.MaxLimit {
			t.Fatalf("provider limit = %v, want bounded positive value", query.Limit)
		}
		if query.IncludeGeometryPreview {
			t.Fatal("Related list query unexpectedly requested geometry preview")
		}
		limits = append(limits, query.Limit)
	}
	sort.Ints(limits)
	wantLimits := []int{2, 4, 8, 8}
	for index := range wantLimits {
		if limits[index] != wantLimits[index] {
			t.Fatalf("provider limits = %v, want %v", limits, wantLimits)
		}
	}
}

func TestProviderPolicyForPlanningProjectAvoidsCentroidParcelAndPOIQueries(t *testing.T) {
	policy := providerPolicyForKind(discoverydomain.EntityKindPlanningProject)
	want := []providerKind{providerPlanning, providerAdministrative}
	if len(policy) != len(want) {
		t.Fatalf("policy = %v, want %v", policy, want)
	}
	for index := range want {
		if policy[index] != want[index] {
			t.Fatalf("policy[%d] = %q, want %q", index, policy[index], want[index])
		}
	}
}

func TestPlanningProjectRunsOnlyDirectPlanningAndAdministrativeProviders(t *testing.T) {
	primary := entity(discoverydomain.EntityKindPlanningProject, "77", 10)
	repo := &fakeRepository{
		primary:   &primary,
		parcelErr: errors.New("parcel provider must not run"),
		poiErr:    errors.New("poi provider must not run"),
		planning: []relateddomain.Candidate{
			related(entity(discoverydomain.EntityKindPlanningRegion, "21", 21), relateddomain.RelationshipMemberOfPrimary, 100, 98),
		},
	}

	response, err := NewService(repo).GetRelatedEntities(context.Background(), relateddomain.Request{Key: primary.Entity.Key})
	if err != nil {
		t.Fatalf("GetRelatedEntities() error = %v", err)
	}
	if response.Partial {
		t.Fatalf("Partial = true; parcel/POI providers should not run")
	}
	repo.mu.Lock()
	queryCount := len(repo.queries)
	repo.mu.Unlock()
	if queryCount != 2 {
		t.Fatalf("provider query count = %d, want 2", queryCount)
	}
}

func TestCandidateGroupTreatsProjectAndRegionAsOnePlanningGroup(t *testing.T) {
	if got := candidateGroup(discoverydomain.EntityKindPlanningProject); got != "planning" {
		t.Fatalf("planning project group = %q", got)
	}
	if got := candidateGroup(discoverydomain.EntityKindPlanningRegion); got != "planning" {
		t.Fatalf("planning region group = %q", got)
	}
}

func TestGetRelatedEntitiesDeduplicatesAcrossProviders(t *testing.T) {
	primary := entity(discoverydomain.EntityKindParcel, "10", 10)
	duplicate := entity(discoverydomain.EntityKindPOI, "41", 41)
	repo := &fakeRepository{
		primary: &primary,
		planning: []relateddomain.Candidate{
			related(duplicate, relateddomain.RelationshipNearby, 15, 80),
		},
		pois: []relateddomain.Candidate{
			related(duplicate, relateddomain.RelationshipNearby, 10, 90),
		},
	}

	response, err := NewService(repo).GetRelatedEntities(context.Background(), relateddomain.Request{Key: primary.Entity.Key})
	if err != nil {
		t.Fatalf("GetRelatedEntities() error = %v", err)
	}
	if len(response.Candidates) != 1 {
		t.Fatalf("len(Candidates) = %d, want 1", len(response.Candidates))
	}
	if response.Candidates[0].Entity.Entity.Key != duplicate.Entity.Key {
		t.Fatalf("candidate key = %q, want %q", response.Candidates[0].Entity.Entity.Key, duplicate.Entity.Key)
	}
}

func TestGetRelatedEntitiesHonorsCanceledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	response, err := NewService(&fakeRepository{}).GetRelatedEntities(ctx, relateddomain.Request{Key: "parcel:10"})
	if response != nil {
		t.Fatalf("response = %#v, want nil", response)
	}
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("error = %v, want context.Canceled", err)
	}
}

func TestGetRelatedEntitiesSkipsInvalidProviderCandidate(t *testing.T) {
	primary := entity(discoverydomain.EntityKindParcel, "10", 10)
	invalid := entity(discoverydomain.EntityKindPOI, "41", 41)
	invalid.Entity.ID = ""
	invalid.Entity.Key = "poi:41"
	repo := &fakeRepository{
		primary: &primary,
		pois: []relateddomain.Candidate{
			related(invalid, relateddomain.RelationshipNearby, 10, 90),
		},
	}

	response, err := NewService(repo).GetRelatedEntities(context.Background(), relateddomain.Request{Key: primary.Entity.Key})
	if err != nil {
		t.Fatalf("GetRelatedEntities() error = %v", err)
	}
	if len(response.Candidates) != 0 {
		t.Fatalf("len(Candidates) = %d, want 0", len(response.Candidates))
	}
}

func TestGetRelatedEntitiesReturnsPartialWhenOneProviderFails(t *testing.T) {
	primary := entity(discoverydomain.EntityKindParcel, "10", 10)
	repo := &fakeRepository{
		primary:     &primary,
		planningErr: errors.New("planning database unavailable"),
		parcels: []relateddomain.Candidate{
			related(entity(discoverydomain.EntityKindParcel, "11", 11), relateddomain.RelationshipAdjacent, 0, 68),
		},
	}

	response, err := NewService(repo).GetRelatedEntities(context.Background(), relateddomain.Request{Key: "parcel:10"})
	if err != nil {
		t.Fatalf("GetRelatedEntities() error = %v", err)
	}
	if !response.Partial {
		t.Fatal("Partial = false, want true")
	}
	if response.Projection.Coverage != relateddomain.CoverageBoundedPartial {
		t.Fatalf("Coverage = %q, want bounded_partial", response.Projection.Coverage)
	}
	assertContains(t, response.Warnings, "related_planning_unavailable")
	if len(response.Candidates) != 1 {
		t.Fatalf("len(Candidates) = %d, want 1", len(response.Candidates))
	}
}

func TestGetRelatedEntitiesKeepsPrimaryWhenAllProvidersFail(t *testing.T) {
	primary := entity(discoverydomain.EntityKindParcel, "10", 10)
	repo := &fakeRepository{
		primary:           &primary,
		planningErr:       errors.New("failed"),
		administrativeErr: errors.New("failed"),
		parcelErr:         errors.New("failed"),
		poiErr:            errors.New("failed"),
	}

	response, err := NewService(repo).GetRelatedEntities(context.Background(), relateddomain.Request{Key: "parcel:10"})
	if err != nil {
		t.Fatalf("GetRelatedEntities() error = %v", err)
	}
	if response == nil || response.Primary == nil {
		t.Fatal("primary overview must survive optional Related provider failures")
	}
	if !response.Partial {
		t.Fatal("Partial = false, want true")
	}
	if len(response.Candidates) != 0 {
		t.Fatalf("len(Candidates) = %d, want 0", len(response.Candidates))
	}
	assertContains(t, response.Warnings, "related_all_providers_unavailable")
}

func TestGetRelatedEntitiesValidatesRequest(t *testing.T) {
	service := NewService(&fakeRepository{})
	cases := []relateddomain.Request{
		{Key: ""},
		{Key: "parcel:10", RadiusMeters: -1},
		{Key: "parcel:10", Limit: -1},
	}
	for _, input := range cases {
		_, err := service.GetRelatedEntities(context.Background(), input)
		var validationError *relateddomain.ValidationError
		if !errors.As(err, &validationError) {
			t.Fatalf("request %#v error = %v, want ValidationError", input, err)
		}
	}
}

func TestGetRelatedEntitiesKeepsPrimaryWhenFocusUnavailable(t *testing.T) {
	primary := entity(discoverydomain.EntityKindAdministrativeUnit, "ward:10", 0)
	primary.Spatial = nil
	repo := &fakeRepository{primary: &primary}

	response, err := NewService(repo).GetRelatedEntities(context.Background(), relateddomain.Request{Key: primary.Entity.Key})
	if err != nil {
		t.Fatalf("GetRelatedEntities() error = %v", err)
	}
	if !response.Partial {
		t.Fatal("Partial = false, want true")
	}
	if len(response.Candidates) != 0 {
		t.Fatalf("len(Candidates) = %d, want 0", len(response.Candidates))
	}
	assertContains(t, response.Warnings, "related_entity_focus_unavailable")
}

func TestSelectCandidatesUsesPrimarySpecificRoleOrderBeforeRawScore(t *testing.T) {
	primary := entity(discoverydomain.EntityKindPlanningRegion, "root", 10)
	parent := related(
		entity(discoverydomain.EntityKindPlanningProject, "parent", 11),
		relateddomain.RelationshipContainsPrimary,
		100,
		30,
	)
	child := related(
		entity(discoverydomain.EntityKindPlanningRegion, "child", 12),
		relateddomain.RelationshipContainedByPrimary,
		10,
		99,
	)
	overlap := related(
		entity(discoverydomain.EntityKindPlanningRegion, "overlap", 13),
		relateddomain.RelationshipIntersectsPrimary,
		5,
		100,
	)

	selected := selectCandidates(
		primary.Entity.Kind,
		primary.Entity.Key,
		[]relateddomain.Candidate{overlap, child, parent},
		3,
	)
	want := []string{"planning_project:parent", "planning_region:child", "planning_region:overlap"}
	for index := range want {
		if selected[index].Entity.Entity.Key != want[index] {
			t.Fatalf("selected[%d] = %q, want %q", index, selected[index].Entity.Entity.Key, want[index])
		}
	}
}

func TestCandidateDecorationKeepsSourceProvenanceFaithful(t *testing.T) {
	candidate := related(
		entity(discoverydomain.EntityKindPlanningRegion, "21", 21),
		relateddomain.RelationshipContainsPrimary,
		0,
		100,
	)
	candidate.Entity.Source = discoverydomain.Source{
		System:      "planning-service",
		Dataset:     "zoning-2026",
		Authority:   "Department A",
		UpdatedAt:   "2026-08-03T00:00:00Z",
		DataQuality: "source_record",
	}
	decorateCandidate(discoverydomain.EntityKindParcel, &candidate)

	attributes := candidate.Entity.Attributes
	if got := attributes["relatedRole"]; got != string(relateddomain.RelatedRoleImpact) {
		t.Fatalf("relatedRole = %#v, want impact", got)
	}
	if got := attributes["sourceDataset"]; got != "zoning-2026" {
		t.Fatalf("sourceDataset = %#v", got)
	}
	if got := attributes["sourceAuthority"]; got != "Department A" {
		t.Fatalf("sourceAuthority = %#v", got)
	}
	if got := attributes["sourceUpdatedAt"]; got != "2026-08-03T00:00:00Z" {
		t.Fatalf("sourceUpdatedAt = %#v", got)
	}
	if _, exists := attributes["official"]; exists {
		t.Fatal("service must not infer official/legal certainty")
	}
}

func entity(kind discoverydomain.EntityKind, id string, coordinate float64) discoverydomain.EntityCandidate {
	return discoverydomain.EntityCandidate{
		Entity: discoverydomain.NewEntityRef(kind, id),
		Spatial: &discoverydomain.SpatialSummary{
			Centroid: &discoverydomain.Point{Latitude: coordinate, Longitude: coordinate},
		},
		Location: &discoverydomain.Location{ProvinceCode: "01", WardCode: "00001"},
		Rank:     discoverydomain.Rank{},
		Match:    discoverydomain.Match{},
	}
}

func related(candidate discoverydomain.EntityCandidate, relationship relateddomain.RelationshipKind, distance, score float64) relateddomain.Candidate {
	return relateddomain.Candidate{
		Entity:         candidate,
		Relationship:   relationship,
		DistanceMeters: distance,
		Score:          score,
	}
}

func assertContains(t *testing.T, values []string, want string) {
	t.Helper()
	for _, value := range values {
		if value == want {
			return
		}
	}
	t.Fatalf("%q not found in %v", want, values)
}

func TestLargePlanningRegionUsesGroupedProjectionAndSkipsParcelPOI(t *testing.T) {
	primary := entity(discoverydomain.EntityKindPlanningRegion, "88", 21)
	primary.Attributes = map[string]any{"areaSqm": 41_200_000.0}
	repo := &fakeRepository{
		primary:   &primary,
		parcelErr: errors.New("parcel provider must not run for large region"),
		poiErr:    errors.New("poi provider must not run for large region"),
		planning: []relateddomain.Candidate{
			related(entity(discoverydomain.EntityKindPlanningProject, "7", 21), relateddomain.RelationshipContainsPrimary, 0, 100),
		},
		administrative: []relateddomain.Candidate{
			related(entity(discoverydomain.EntityKindAdministrativeUnit, "ward:9", 21), relateddomain.RelationshipAdministrativeContext, 0, 95),
		},
	}

	response, err := NewService(repo).GetRelatedEntities(context.Background(), relateddomain.Request{Key: primary.Entity.Key})
	if err != nil {
		t.Fatalf("GetRelatedEntities() error = %v", err)
	}
	if response.Projection.DisplayMode != relateddomain.DisplayModeGroups {
		t.Fatalf("DisplayMode = %q, want %q", response.Projection.DisplayMode, relateddomain.DisplayModeGroups)
	}
	repo.mu.Lock()
	queryCount := len(repo.queries)
	repo.mu.Unlock()
	if queryCount != 2 {
		t.Fatalf("provider query count = %d, want 2", queryCount)
	}
	projection, ok := response.Primary.Attributes["relatedProjection"].(map[string]any)
	if !ok {
		t.Fatalf("relatedProjection metadata = %#v", response.Primary.Attributes["relatedProjection"])
	}
	if _, mutated := primary.Attributes["relatedProjection"]; mutated {
		t.Fatal("service must not mutate the repository-owned primary attributes")
	}
	if projection["displayMode"] != string(relateddomain.DisplayModeGroups) {
		t.Fatalf("projection displayMode = %#v", projection["displayMode"])
	}
	if projection["schemaVersion"] != 2 {
		t.Fatalf("projection schemaVersion = %#v, want 2", projection["schemaVersion"])
	}
	roleOrder, ok := projection["roleOrder"].([]string)
	if !ok || len(roleOrder) < 3 {
		t.Fatalf("projection roleOrder = %#v", projection["roleOrder"])
	}
	if roleOrder[0] != string(relateddomain.RelatedRoleContext) ||
		roleOrder[1] != string(relateddomain.RelatedRoleContained) ||
		roleOrder[2] != string(relateddomain.RelatedRoleOverlap) {
		t.Fatalf("projection roleOrder = %v", roleOrder)
	}
	if len(response.Roles) != 1 || response.Roles[0].Key != relateddomain.RelatedRoleContext {
		t.Fatalf("response Roles = %#v, want one context role", response.Roles)
	}
	assertContains(t, response.Warnings, "related_projection_grouped")
	assertContains(t, response.Warnings, "related_group_parcel_suppressed")
	assertContains(t, response.Warnings, "related_group_poi_suppressed")
}

func TestRestrictedProjectionCanAlsoBePartial(t *testing.T) {
	primary := entity(discoverydomain.EntityKindPlanningRegion, "99", 21)
	primary.Attributes = map[string]any{"areaSqm": 500_000_000.0}
	repo := &fakeRepository{
		primary:     &primary,
		planningErr: errors.New("planning provider unavailable"),
		administrative: []relateddomain.Candidate{
			related(entity(discoverydomain.EntityKindAdministrativeUnit, "province:1", 21), relateddomain.RelationshipAdministrativeContext, 0, 95),
		},
	}

	response, err := NewService(repo).GetRelatedEntities(context.Background(), relateddomain.Request{Key: primary.Entity.Key})
	if err != nil {
		t.Fatalf("GetRelatedEntities() error = %v", err)
	}
	if response.Projection.DisplayMode != relateddomain.DisplayModeRestricted {
		t.Fatalf("DisplayMode = %q, want restricted", response.Projection.DisplayMode)
	}
	if !response.Partial {
		t.Fatal("Partial = false, want true when one allowed provider fails")
	}
	assertContains(t, response.Warnings, "related_projection_restricted")
	assertContains(t, response.Warnings, "related_planning_unavailable")
}
