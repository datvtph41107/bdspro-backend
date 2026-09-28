package services

import (
	"context"
	"testing"

	"userv4/modules/bdspro/dto"
	"userv4/modules/bdspro/entity"
)

func TestLocalPublicService_GetEnums(t *testing.T) {
	svc := NewLocalPublicService(nil)
	resp, err := svc.GetEnums(context.Background(), 0)
	if err != nil {
		t.Fatalf("GetEnums: %v", err)
	}
	if resp == nil || !resp.GetLastest() || len(resp.GetEnums()) == 0 {
		t.Fatalf("expected enum catalog, got %+v", resp)
	}
	if resp.GetTimestamp() <= 0 {
		t.Fatalf("expected positive timestamp, got %d", resp.GetTimestamp())
	}
}

func TestLocalPublicService_GetEnumByName(t *testing.T) {
	svc := NewLocalPublicService(nil)
	resp, err := svc.GetEnumByName(context.Background(), "visibility")
	if err != nil {
		t.Fatalf("GetEnumByName: %v", err)
	}
	if resp == nil || len(resp.GetData()) == 0 {
		t.Fatalf("expected visibility enums, got %+v", resp)
	}
}

func TestLocalPublicService_GetEnumByName_NotFound(t *testing.T) {
	svc := NewLocalPublicService(nil)
	_, err := svc.GetEnumByName(context.Background(), "does-not-exist")
	if err == nil {
		t.Fatal("expected not found error")
	}
}

func TestLocalProductService_GetTimestamps_NilSync(t *testing.T) {
	svc := NewLocalProductService(nil, nil)
	resp, err := svc.GetTimestamps(context.Background(), 42)
	if err != nil {
		t.Fatalf("GetTimestamps: %v", err)
	}
	if resp == nil || resp.Timestamps == nil {
		t.Fatalf("expected empty timestamps map, got %+v", resp)
	}
	if len(resp.Timestamps) != 0 {
		t.Fatalf("expected empty map without sync, got %v", resp.Timestamps)
	}
}

func TestLocalProductService_GetTimestamps_RequiresID(t *testing.T) {
	svc := NewLocalProductService(nil, nil)
	_, err := svc.GetTimestamps(context.Background(), 0)
	if err == nil {
		t.Fatal("expected error for id=0")
	}
}

func TestLocalProductService_GetDetail_ReturnsTimestamp(t *testing.T) {
	svc := NewLocalProductService(nil, stubProductRepo{product: &entity.Product{ID: 385, Name: "P385"}})
	resp, ts, err := svc.GetDetail(context.Background(), dto.SyncListRequest{ID: 385})
	if err != nil {
		t.Fatalf("GetDetail: %v", err)
	}
	if resp == nil || resp.GetId() != 385 {
		t.Fatalf("expected product 385, got %+v", resp)
	}
	if ts <= 0 {
		t.Fatalf("expected positive envelope timestamp, got %d", ts)
	}
}

func TestLocalProductService_GetDetail_RequiresID(t *testing.T) {
	svc := NewLocalProductService(nil, stubProductRepo{})
	_, _, err := svc.GetDetail(context.Background(), dto.SyncListRequest{})
	if err == nil {
		t.Fatal("expected error for id=0")
	}
}

func TestMapProductToPB_Nil(t *testing.T) {
	if mapProductToPB(nil) != nil {
		t.Fatal("expected nil")
	}
}

func TestLocalDealService_EmptyRepo(t *testing.T) {
	svc := NewLocalDealService(nil, stubDealRepo{})
	resp, err := svc.GetDealsUser(context.Background(), dto.DealSearchRequest{})
	if err != nil {
		t.Fatalf("GetDealsUser: %v", err)
	}
	if resp == nil {
		t.Fatal("expected response")
	}
	if len(resp.GetData()) != 0 {
		t.Fatalf("expected empty deals, got %d", len(resp.GetData()))
	}
}

func TestLocalPropertyService_EmptyRepo(t *testing.T) {
	svc := NewLocalPropertyService(nil, stubPropertyRepo{})
	resp, err := svc.ListMe(context.Background(), dto.PropertyListRequest{Size: 10})
	if err != nil {
		t.Fatalf("ListMe: %v", err)
	}
	if resp == nil || len(resp.GetData()) != 0 {
		t.Fatalf("expected empty property list, got %+v", resp)
	}
}

type stubDealRepo struct{}

func (stubDealRepo) ListForUser(context.Context, uint64, uint32, uint32, string) ([]*entity.Deal, uint32, error) {
	return nil, 0, nil
}

type stubPropertyRepo struct{}

func (stubPropertyRepo) ListMine(context.Context, uint64, uint32, uint32) ([]*entity.PropertyLineage, int64, error) {
	return nil, 0, nil
}

func (stubPropertyRepo) GetByID(context.Context, uint64) (*entity.PropertyLineage, error) {
	return nil, nil
}

type stubProductRepo struct {
	product *entity.Product
}

func (r stubProductRepo) ListMine(context.Context, uint64, uint32, uint32, *bool) ([]*entity.Product, int64, error) {
	return nil, 0, nil
}

func (r stubProductRepo) GetByID(context.Context, uint64) (*entity.Product, error) {
	return r.product, nil
}

func (stubProductRepo) Create(context.Context, *entity.Product) error { return nil }

func (stubProductRepo) SoftDelete(context.Context, uint64, uint64) error { return nil }
