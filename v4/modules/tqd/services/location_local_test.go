package services

import (
	"context"
	"testing"

	"userv4/modules/tqd/dto"
	"userv4/modules/tqd/entity"
	"userv4/modules/tqd/enums"
)

type stubLocationRepo struct {
	provinces []*entity.Province
	wards     []*entity.Ward
	wardTotal int64
}

func (s stubLocationRepo) ListAllProvinces(context.Context) ([]*entity.Province, error) {
	return s.provinces, nil
}

func (s stubLocationRepo) ListAllWards(context.Context, int, int) ([]*entity.Ward, int64, error) {
	return s.wards, s.wardTotal, nil
}

func TestLocalLocationService_ListProvinces(t *testing.T) {
	svc := NewLocalLocationService(stubLocationRepo{
		provinces: []*entity.Province{{
			ID:       13,
			Code:     "01",
			FullName: "Hà Nội",
			ShortName: "HN",
			UnitType: enums.LocationUnitTypeCity,
		}},
	}, nil)

	resp, err := svc.ListProvinces(context.Background(), dto.ListProvincesRequest{})
	if err != nil {
		t.Fatalf("ListProvinces: %v", err)
	}
	if resp == nil || len(resp.GetProvinces()) != 1 {
		t.Fatalf("expected 1 province, got %+v", resp)
	}
	p := resp.GetProvinces()[0]
	if p.GetId() != 13 || p.GetFullName() != "Hà Nội" || p.GetCompositeId() != 131 {
		t.Fatalf("unexpected province mapping: %+v", p)
	}
	if resp.GetTimestamp() <= 0 {
		t.Fatalf("expected positive timestamp, got %d", resp.GetTimestamp())
	}
}

func TestLocalLocationService_ListProvinces_UnchangedCatalog(t *testing.T) {
	svc := NewLocalLocationService(stubLocationRepo{
		provinces: []*entity.Province{{ID: 1, FullName: "X"}},
	}, nil)

	first, err := svc.ListProvinces(context.Background(), dto.ListProvincesRequest{})
	if err != nil {
		t.Fatalf("first ListProvinces: %v", err)
	}
	// Without SyncUtil, every call refetches (timestamp always "now").
	second, err := svc.ListProvinces(context.Background(), dto.ListProvincesRequest{Timestamp: first.GetTimestamp()})
	if err != nil {
		t.Fatalf("second ListProvinces: %v", err)
	}
	if len(second.GetProvinces()) != 1 {
		t.Fatalf("nil sync always refetches; got %+v", second)
	}
}

func TestLocalLocationService_ListWards(t *testing.T) {
	svc := NewLocalLocationService(stubLocationRepo{
		wards: []*entity.Ward{{
			ID:         142,
			ShortName:  "P.1",
			ProvinceID: 13,
			UnitType:   enums.LocationUnitTypeCity,
		}},
		wardTotal: 1,
	}, nil)

	page := int32(1)
	size := int32(50)
	resp, err := svc.ListWards(context.Background(), dto.ListWardsRequest{Page: &page, PageSize: &size})
	if err != nil {
		t.Fatalf("ListWards: %v", err)
	}
	if resp == nil || len(resp.GetData()) != 1 {
		t.Fatalf("expected 1 ward, got %+v", resp)
	}
	w := resp.GetData()[0]
	if w.GetId() != 142 || w.GetCompositeId() != 1422 || w.GetProvinceId() != 13 {
		t.Fatalf("unexpected ward mapping: %+v", w)
	}
	if resp.GetPagination() == nil || resp.GetPagination().GetTotal() != 1 {
		t.Fatalf("expected pagination total=1, got %+v", resp.GetPagination())
	}
}
