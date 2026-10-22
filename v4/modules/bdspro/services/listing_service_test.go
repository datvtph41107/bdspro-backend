package services

import (
	"context"
	"testing"

	_enums "common/domain/enum"
	"mdn4/modules/bdspro/dto"
	"mdn4/modules/bdspro/entity"
)

type stubListingRepo struct {
	items []entity.Listing
}

func (s stubListingRepo) CreateForProduct(context.Context, *entity.Listing) error { return nil }

func (s stubListingRepo) List(_ context.Context, _ dto.ListingListRequest, ownerID uint64) ([]entity.Listing, int64, error) {
	if ownerID != 42 {
		return nil, 0, nil
	}
	return s.items, int64(len(s.items)), nil
}

func TestLocalListingService_ListMine(t *testing.T) {
	svc := NewLocalListingService(stubListingRepo{
		items: []entity.Listing{{ID: 7, Title: "A", OwnerID: 42}},
	}, nil)
	ctx := context.WithValue(context.Background(), _enums.ProfileIDKey, uint64(42))
	out, err := svc.ListMine(ctx, dto.ListingListRequest{})
	if err != nil {
		t.Fatal(err)
	}
	resp := out.List
	if resp.GetTotal() != 1 || len(resp.GetData()) != 1 || resp.GetData()[0].GetTitle() != "A" {
		t.Fatalf("unexpected: %+v", resp)
	}
}
