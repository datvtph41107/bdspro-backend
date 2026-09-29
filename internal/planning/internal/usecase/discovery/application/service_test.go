package application

import (
	"context"
	"testing"
	"tqd/internal/domain/discovery/model"
)

type fakeRepo struct {
	parcels []domain.EntityCandidate
	regions []domain.EntityCandidate
}

func (f fakeRepo) IdentifyParcels(context.Context, domain.IdentifyRequest) ([]domain.EntityCandidate, error) {
	return f.parcels, nil
}
func (f fakeRepo) IdentifyRegions(context.Context, domain.IdentifyRequest) ([]domain.EntityCandidate, error) {
	return f.regions, nil
}
func (f fakeRepo) IdentifyAdministrativeUnits(context.Context, domain.IdentifyRequest) ([]domain.EntityCandidate, error) {
	return nil, nil
}
func (f fakeRepo) IdentifyPOIs(context.Context, domain.IdentifyRequest) ([]domain.EntityCandidate, error) {
	return nil, nil
}
func (f fakeRepo) SearchParcels(context.Context, domain.SearchRequest) ([]domain.EntityCandidate, error) {
	return f.parcels, nil
}
func (f fakeRepo) SearchRegions(context.Context, domain.SearchRequest) ([]domain.EntityCandidate, error) {
	return f.regions, nil
}
func (f fakeRepo) SearchPlanningProjects(context.Context, domain.SearchRequest) ([]domain.EntityCandidate, error) {
	return nil, nil
}
func (f fakeRepo) SearchAdministrativeUnits(context.Context, domain.SearchRequest) ([]domain.EntityCandidate, error) {
	return nil, nil
}
func (f fakeRepo) SearchPOIs(context.Context, domain.SearchRequest) ([]domain.EntityCandidate, error) {
	return nil, nil
}
func (f fakeRepo) GetEntity(context.Context, domain.EntityRef) (*domain.EntityCandidate, error) {
	return nil, nil
}

func item(kind domain.EntityKind, id string) domain.EntityCandidate {
	return domain.EntityCandidate{Entity: domain.NewEntityRef(kind, id), Presentation: domain.Presentation{Title: id}, Rank: domain.Rank{Score: 1}}
}

func TestIdentifyPlanningModePrefersRegion(t *testing.T) {
	svc := NewService(fakeRepo{parcels: []domain.EntityCandidate{item(domain.EntityKindParcel, "1")}, regions: []domain.EntityCandidate{item(domain.EntityKindPlanningRegion, "2")}})
	res, err := svc.Identify(context.Background(), domain.IdentifyRequest{Point: domain.Point{Latitude: 10, Longitude: 106}, Context: domain.MapContext{MapMode: "planning"}})
	if err != nil {
		t.Fatal(err)
	}
	if res.Primary == nil || res.Primary.Entity.Kind != domain.EntityKindPlanningRegion {
		t.Fatalf("expected planning region primary, got %#v", res.Primary)
	}
}

func TestSearchParcelIntentRanksParcel(t *testing.T) {
	svc := NewService(fakeRepo{parcels: []domain.EntityCandidate{item(domain.EntityKindParcel, "1")}, regions: []domain.EntityCandidate{item(domain.EntityKindPlanningRegion, "2")}})
	res, err := svc.Search(context.Background(), domain.SearchRequest{Query: "thửa 12 tờ 3"})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Groups) == 0 {
		t.Fatal("expected search groups")
	}
	var first domain.EntityCandidate
	for _, g := range res.Groups {
		for _, v := range g.Items {
			if first.Entity.Key == "" || v.Rank.Score > first.Rank.Score {
				first = v
			}
		}
	}
	if first.Entity.Kind != domain.EntityKindParcel {
		t.Fatalf("expected parcel, got %s", first.Entity.Kind)
	}
}

func TestSearchFiltersEntityKinds(t *testing.T) {
	svc := NewService(fakeRepo{
		parcels: []domain.EntityCandidate{item(domain.EntityKindParcel, "1")},
		regions: []domain.EntityCandidate{item(domain.EntityKindPlanningRegion, "2")},
	})
	res, err := svc.Search(context.Background(), domain.SearchRequest{
		Query:       "quy hoạch",
		EntityKinds: []domain.EntityKind{domain.EntityKindPlanningRegion},
	})
	if err != nil {
		t.Fatal(err)
	}
	for _, group := range res.Groups {
		if group.Kind != domain.EntityKindPlanningRegion {
			t.Fatalf("unexpected group kind %s", group.Kind)
		}
	}
}

func TestIdentifyHighDetailZoomPrefersParcelOverUnboundedRegionProviderScore(t *testing.T) {
	parcel := item(domain.EntityKindParcel, "1")
	parcel.Rank.Score = 30
	region := item(domain.EntityKindPlanningRegion, "2")
	region.Rank.Score = 5520

	svc := NewService(fakeRepo{
		parcels: []domain.EntityCandidate{parcel},
		regions: []domain.EntityCandidate{region},
	})
	res, err := svc.Identify(context.Background(), domain.IdentifyRequest{
		Point: domain.Point{Latitude: 21.035748849459495, Longitude: 106.0167378848052},
		Context: domain.MapContext{
			Zoom:    18,
			MapMode: "overview",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if res.Primary == nil || res.Primary.Entity.Kind != domain.EntityKindParcel {
		t.Fatalf("expected parcel primary at high-detail zoom, got %#v", res.Primary)
	}
	if !contains(res.Primary.Rank.ReasonCodes, "high_detail_zoom") {
		t.Fatalf("expected high_detail_zoom reason, got %#v", res.Primary.Rank.ReasonCodes)
	}
}

func TestIdentifyLowDetailZoomPrefersPlanningRegion(t *testing.T) {
	parcel := item(domain.EntityKindParcel, "1")
	parcel.Rank.Score = 30
	region := item(domain.EntityKindPlanningRegion, "2")
	region.Rank.Score = 5520

	svc := NewService(fakeRepo{
		parcels: []domain.EntityCandidate{parcel},
		regions: []domain.EntityCandidate{region},
	})
	res, err := svc.Identify(context.Background(), domain.IdentifyRequest{
		Point: domain.Point{Latitude: 21.035748849459495, Longitude: 106.0167378848052},
		Context: domain.MapContext{
			Zoom:    12,
			MapMode: "overview",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if res.Primary == nil || res.Primary.Entity.Kind != domain.EntityKindPlanningRegion {
		t.Fatalf("expected planning region primary at low-detail zoom, got %#v", res.Primary)
	}
}

func TestIdentifyPlanningModeOverridesHighDetailParcelPreference(t *testing.T) {
	parcel := item(domain.EntityKindParcel, "1")
	parcel.Rank.Score = 30
	region := item(domain.EntityKindPlanningRegion, "2")
	region.Rank.Score = 5520

	svc := NewService(fakeRepo{
		parcels: []domain.EntityCandidate{parcel},
		regions: []domain.EntityCandidate{region},
	})
	res, err := svc.Identify(context.Background(), domain.IdentifyRequest{
		Point: domain.Point{Latitude: 21.035748849459495, Longitude: 106.0167378848052},
		Context: domain.MapContext{
			Zoom:    18,
			MapMode: "planning",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if res.Primary == nil || res.Primary.Entity.Kind != domain.EntityKindPlanningRegion {
		t.Fatalf("expected planning region primary in planning mode, got %#v", res.Primary)
	}
}
