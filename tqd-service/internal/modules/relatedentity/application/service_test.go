package application

import (
	"context"
	"errors"
	"sort"
	"sync"
	"testing"

	discoverydomain "tqd/internal/modules/discovery/domain"
	relateddomain "tqd/internal/modules/relatedentity/domain"
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

func TestGetRelatedEntitiesClampsAndSelectsOnePerKind(t *testing.T) {
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
	wantKinds := []discoverydomain.EntityKind{
		discoverydomain.EntityKindPlanningRegion,
		discoverydomain.EntityKindAdministrativeUnit,
		discoverydomain.EntityKindParcel,
		discoverydomain.EntityKindPOI,
	}
	for index, want := range wantKinds {
		if response.Candidates[index].Entity.Entity.Kind != want {
			t.Fatalf("Candidates[%d].Kind = %q, want %q", index, response.Candidates[index].Entity.Entity.Kind, want)
		}
	}
	for _, candidate := range response.Candidates {
		if candidate.Entity.Entity.Key == primary.Entity.Key {
			t.Fatalf("primary entity %q must be excluded", primary.Entity.Key)
		}
		if candidate.Entity.Attributes["relationship"] == nil {
			t.Fatalf("candidate %q has no relationship", candidate.Entity.Entity.Key)
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
	assertContains(t, response.Warnings, "related_planning_unavailable")
	if len(response.Candidates) != 1 {
		t.Fatalf("len(Candidates) = %d, want 1", len(response.Candidates))
	}
}

func TestGetRelatedEntitiesFailsWhenAllProvidersFail(t *testing.T) {
	primary := entity(discoverydomain.EntityKindParcel, "10", 10)
	repo := &fakeRepository{
		primary:           &primary,
		planningErr:       errors.New("failed"),
		administrativeErr: errors.New("failed"),
		parcelErr:         errors.New("failed"),
		poiErr:            errors.New("failed"),
	}

	response, err := NewService(repo).GetRelatedEntities(context.Background(), relateddomain.Request{Key: "parcel:10"})
	if response != nil {
		t.Fatalf("response = %#v, want nil", response)
	}
	if !errors.Is(err, relateddomain.ErrProvidersFailed) {
		t.Fatalf("error = %v, want ErrProvidersFailed", err)
	}
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
