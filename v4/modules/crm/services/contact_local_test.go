package services

import (
	"context"
	"testing"
	"time"

	base_enum "base/enum"
	_enums "common/domain/enum"
	"userv4/modules/crm/dto"
	"userv4/modules/crm/entity"
)

type stubContactRepo struct {
	items   []entity.ContactListItem
	total   int64
	byID    map[uint64]*entity.Contact
	friends map[[2]uint64]*entity.FriendRelation
	products map[uint64][]uint64
}

func (s stubContactRepo) SearchByOwner(context.Context, uint64, base_enum.EOwnerOf, uint32, uint32, string) ([]entity.ContactListItem, int64, error) {
	return s.items, s.total, nil
}

func (s stubContactRepo) GetByID(_ context.Context, id uint64) (*entity.Contact, error) {
	if s.byID == nil {
		return nil, nil
	}
	return s.byID[id], nil
}

func (s stubContactRepo) GetFriendBetween(_ context.Context, a, b uint64) (*entity.FriendRelation, error) {
	if s.friends == nil {
		return nil, nil
	}
	if rel, ok := s.friends[[2]uint64{a, b}]; ok {
		return rel, nil
	}
	if rel, ok := s.friends[[2]uint64{b, a}]; ok {
		return rel, nil
	}
	return nil, nil
}

func (s stubContactRepo) ListProductIDs(_ context.Context, contactID uint64) ([]uint64, error) {
	if s.products == nil {
		return nil, nil
	}
	return s.products[contactID], nil
}

func (s stubContactRepo) Create(_ context.Context, _ *entity.Contact) error { return nil }

func (s stubContactRepo) SoftDelete(_ context.Context, _ uint64) error { return nil }

func (s stubContactRepo) FindIDByPhone(_ context.Context, _ uint64, _ base_enum.EOwnerOf, _ string) (uint64, error) {
	return 0, nil
}

func (s stubContactRepo) GetOriginIDByProfileID(_ context.Context, _ uint64) (uint64, error) {
	return 0, nil
}

func TestLocalContactService_ListUser(t *testing.T) {
	now := time.Now()
	profileID := uint64(42)
	contactProfileID := uint64(99)
	svc := NewLocalContactService(stubContactRepo{
		items: []entity.ContactListItem{{
			ID: 1, FullName: "Alice", Phone: "090", OwnerID: profileID,
			OwnerOf: base_enum.EOwnerOfMember, ProfileID: &contactProfileID, Following: true,
		}},
		total: 1,
		byID: map[uint64]*entity.Contact{
			1: {
				ID: 1, FullName: "Alice", Phone: "090", OwnerID: profileID,
				OwnerOf: base_enum.EOwnerOfMember, ProfileID: &contactProfileID,
				CreatedAt: &now, UpdatedAt: &now,
			},
		},
		friends: map[[2]uint64]*entity.FriendRelation{
			{profileID, contactProfileID}: {ID: 7, CreatedBy: profileID, ReceiverID: contactProfileID, Status: 20},
		},
		products: map[uint64][]uint64{1: {11, 22}},
	}, nil, nil)

	ctx := context.WithValue(context.Background(), _enums.ProfileIDKey, profileID)
	resp, err := svc.ListUser(ctx, 0, dto.ContactSearchRequest{Size: 50})
	if err != nil {
		t.Fatalf("ListUser: %v", err)
	}
	if resp == nil || len(resp.GetData()) != 1 || resp.GetTotal() != 1 {
		t.Fatalf("unexpected response: %+v", resp)
	}
	item := resp.GetData()[0]
	if item.GetFullName() != "Alice" || item.GetPhone() != "090" {
		t.Fatalf("unexpected contact mapping: %+v", item)
	}
	if !item.GetFollowing() {
		t.Fatal("expected following=true")
	}
	if item.GetFriendStatus() == nil || item.GetFriendStatus().GetId() != 7 {
		t.Fatalf("expected friend status, got %+v", item.GetFriendStatus())
	}
	if len(item.GetProductIds()) != 2 {
		t.Fatalf("expected product ids, got %+v", item.GetProductIds())
	}
}

func TestLocalContactService_ListUser_Unauthorized(t *testing.T) {
	svc := NewLocalContactService(stubContactRepo{}, nil, nil)
	_, err := svc.ListUser(context.Background(), 0, dto.ContactSearchRequest{})
	if err == nil {
		t.Fatal("expected unauthorized error")
	}
}
