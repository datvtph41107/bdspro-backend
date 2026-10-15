package services

import (
	"context"
	"testing"

	bdspropb "pb/types/bdspro"

	"userv4/modules/bdspro/dto"
	"userv4/modules/bdspro/entity"
	"userv4/modules/bdspro/enums"
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
	svc := NewLocalDealService(nil, &stubDealRepo{}, stubProductRepo{}, nil)
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

func TestLocalDealService_GetDealsUser_IncludesInvitationForInvitedProfile(t *testing.T) {
	pid := uint64(40)
	repo := &stubDealRepoList{
		deals: []*entity.Deal{{
			ID:   7,
			Name: "Deal 7",
			Members: []entity.DealMember{{
				ID:       501,
				DealID:   7,
				MemberID: pid,
				Status:   uint32(enums.DealMemberStatusAccepted),
				RoleKey:  uint32(enums.RoleKeyDealOwner),
				IsOwner:  true,
			}},
		}},
		total: 1,
	}
	svc := NewLocalDealService(nil, repo, stubProductRepo{}, nil)
	invited := pid
	resp, err := svc.GetDealsUser(context.Background(), dto.DealSearchRequest{
		InvitedProfileID: &invited,
	})
	if err != nil {
		t.Fatalf("GetDealsUser: %v", err)
	}
	if len(resp.GetData()) != 1 {
		t.Fatalf("expected 1 deal, got %d", len(resp.GetData()))
	}
	inv := resp.GetData()[0].GetInvitation()
	if inv == nil || inv.GetId() != 501 || inv.GetMemberId() != pid {
		t.Fatalf("expected invitation for profile 40, got %+v", inv)
	}
	role := inv.GetRole()
	if role == nil || role.GetRoleKey() != uint32(enums.RoleKeyDealOwner) || role.GetRoleName() != enums.RoleKeyMap[enums.RoleKeyDealOwner] {
		t.Fatalf("expected owner role on invitation, got %+v", role)
	}
}

type stubDealRepoList struct {
	stubDealRepo
	deals []*entity.Deal
	total uint32
}

func (r *stubDealRepoList) ListForUser(context.Context, uint64, uint32, uint32, dto.DealSearchRequest) ([]*entity.Deal, uint32, error) {
	return r.deals, r.total, nil
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

func TestLocalPropertyService_GetTimestamps_NilSync(t *testing.T) {
	svc := NewLocalPropertyService(nil, nil)
	resp, err := svc.GetTimestamps(context.Background(), 30)
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

func TestLocalPropertyService_GetTimestamps_RequiresID(t *testing.T) {
	svc := NewLocalPropertyService(nil, nil)
	_, err := svc.GetTimestamps(context.Background(), 0)
	if err == nil {
		t.Fatal("expected error for id=0")
	}
}

type stubDealRepo struct {
	createdDeal    *entity.Deal
	createdMembers []entity.DealMember
	createdProduct []uint64
}

func (r *stubDealRepo) ListForUser(context.Context, uint64, uint32, uint32, dto.DealSearchRequest) ([]*entity.Deal, uint32, error) {
	return nil, 0, nil
}

func (r *stubDealRepo) ListIDsByProductID(context.Context, uint64, uint32, uint32) ([]uint64, error) {
	return nil, nil
}

func (r *stubDealRepo) ListDealProductsByProductID(context.Context, uint64, uint32, uint32) ([]*entity.DealProduct, error) {
	return nil, nil
}

func (r *stubDealRepo) GetByID(context.Context, uint64) (*entity.Deal, error) {
	return nil, nil
}

func (r *stubDealRepo) ListByIDs(context.Context, []uint64) ([]*entity.Deal, error) {
	return nil, nil
}

func (r *stubDealRepo) GetContractProcess(context.Context, uint64) ([]dto.TxProcessRow, error) {
	return nil, nil
}

func (r *stubDealRepo) ListDealProducts(context.Context, uint64, uint32, uint32) ([]*entity.DealProduct, error) {
	return nil, nil
}

func (r *stubDealRepo) ListContractsByDealID(context.Context, uint64) ([]dto.TxContractDealListRow, int64, error) {
	return nil, 0, nil
}

func (r *stubDealRepo) ListMembersByDealID(context.Context, uint64, int, int) ([]*entity.DealMember, uint32, error) {
	return nil, 0, nil
}

func (r *stubDealRepo) ListAcceptedMembers(context.Context, dto.DealMembersRequest) ([]*entity.DealMember, error) {
	return nil, nil
}

func (r *stubDealRepo) ListInvestmentsByDealID(context.Context, uint64, uint64) ([]*entity.DealInvestment, int64, error) {
	return nil, 0, nil
}

func (r *stubDealRepo) ListCostsByDealID(context.Context, uint64, uint32, uint32) ([]*entity.TxDealCost, error) {
	return nil, nil
}

func (r *stubDealRepo) CountDashboard(context.Context, dto.DealSummaryRequest) (*dto.DealSummaryRow, error) {
	return &dto.DealSummaryRow{}, nil
}

func (r *stubDealRepo) GetBankAccountByID(context.Context, uint64) (*entity.BankAccount, error) {
	return nil, nil
}

func (r *stubDealRepo) GetCostStatisticsByDealID(context.Context, uint64) (*dto.DealCostStatsRow, error) {
	return &dto.DealCostStatsRow{}, nil
}

func (r *stubDealRepo) GetPaymentStatisticsByDealID(context.Context, uint64) (*dto.DealPaymentStatsRow, error) {
	return &dto.DealPaymentStatsRow{}, nil
}

func (r *stubDealRepo) SeedDemoRelated(context.Context, uint64, uint64, uint64, uint32) error {
	return nil
}

func (r *stubDealRepo) Update(context.Context, *entity.Deal) error { return nil }

func (r *stubDealRepo) UpdateSetting(context.Context, uint64, *dto.DealSetting) error { return nil }

func (r *stubDealRepo) Create(_ context.Context, deal *entity.Deal, members []entity.DealMember, productIDs []uint64) error {
	if deal != nil {
		deal.ID = 99
		r.createdDeal = deal
	}
	r.createdMembers = append([]entity.DealMember(nil), members...)
	r.createdProduct = append([]uint64(nil), productIDs...)
	return nil
}

func (r *stubDealRepo) SoftDelete(context.Context, uint64) error { return nil }

func (r *stubDealRepo) IsAcceptedMember(context.Context, uint64, uint64) (bool, error) {
	return true, nil
}

func (r *stubDealRepo) UpdateDoneInvestment(context.Context, uint64, uint64) error { return nil }

func TestLocalDealService_GetTimestamps_RequiresID(t *testing.T) {
	svc := NewLocalDealService(nil, &stubDealRepo{}, stubProductRepo{}, nil)
	_, err := svc.GetTimestamps(context.Background(), 0)
	if err == nil {
		t.Fatal("expected error for id=0")
	}
}

func TestLocalDealService_GetTimestamps_NilSync(t *testing.T) {
	svc := NewLocalDealService(nil, &stubDealRepo{}, stubProductRepo{}, nil)
	resp, err := svc.GetTimestamps(context.Background(), 14)
	if err != nil {
		t.Fatalf("GetTimestamps: %v", err)
	}
	if resp == nil || len(resp.Timestamps) != 0 {
		t.Fatalf("expected empty timestamps, got %+v", resp)
	}
}

func TestLocalDealService_GetContractProcess_Empty(t *testing.T) {
	svc := NewLocalDealService(nil, &stubDealRepo{}, stubProductRepo{}, nil)
	resp, ts, err := svc.GetContractProcess(context.Background(), 14)
	if err != nil {
		t.Fatalf("GetContractProcess: %v", err)
	}
	if resp == nil || len(resp.GetData()) != 0 {
		t.Fatalf("expected empty process, got %+v", resp)
	}
	if ts <= 0 {
		t.Fatalf("expected positive envelope timestamp, got %d", ts)
	}
}

func TestLocalDealService_GetDetail_NotFound(t *testing.T) {
	svc := NewLocalDealService(nil, &stubDealRepo{}, stubProductRepo{}, nil)
	_, _, err := svc.GetDetail(context.Background(), dto.DealDetailRequest{ID: 14})
	if err == nil {
		t.Fatal("expected not found")
	}
}

func TestLocalDealService_GetDetail_ReturnsTimestamp(t *testing.T) {
	repo := &stubDealRepo{}
	repo.createdDeal = &entity.Deal{ID: 14, Name: "D14"}
	// reuse Create capture style — override GetByID via dedicated stub
	svc := NewLocalDealService(nil, &stubDealRepoDetail{deal: &entity.Deal{ID: 14, Name: "D14"}}, stubProductRepo{}, nil)
	resp, ts, err := svc.GetDetail(context.Background(), dto.DealDetailRequest{ID: 14})
	if err != nil {
		t.Fatalf("GetDetail: %v", err)
	}
	if resp == nil || resp.GetId() != 14 {
		t.Fatalf("expected deal 14, got %+v", resp)
	}
	if ts <= 0 {
		t.Fatalf("expected positive envelope timestamp, got %d", ts)
	}
}

type stubDealRepoDetail struct {
	stubDealRepo
	deal *entity.Deal
}

func (r *stubDealRepoDetail) GetByID(context.Context, uint64) (*entity.Deal, error) {
	return r.deal, nil
}

func TestLocalDealService_CreateGroupDeal_OwnerAccepted(t *testing.T) {
	repo := &stubDealRepo{}
	svc := NewLocalDealService(nil, repo, stubProductRepo{}, nil)
	ctx := context.WithValue(context.Background(), struct{}{}, nil)
	// Profile id via common utils key if available — CreateGroupDeal falls back to req.OwnerId.
	id, err := svc.CreateGroupDeal(ctx, &bdspropb.CreateGroupDealRequest{
		Name:           "Thương vụ demo bán căn hộ",
		Status:         10,
		DealType:       40,
		OwnerType:      10,
		OwnerId:        7,
		ChargePersonId: 7,
		CustomerIds:    []uint64{11},
		ProductIds:     []uint64{22},
		Description:    "demo",
	})
	if err != nil {
		t.Fatalf("CreateGroupDeal: %v", err)
	}
	if id != 99 {
		t.Fatalf("expected id 99, got %d", id)
	}
	if repo.createdDeal == nil || repo.createdDeal.Description != "demo" {
		t.Fatalf("expected description on deal, got %+v", repo.createdDeal)
	}
	if len(repo.createdMembers) < 1 {
		t.Fatal("expected owner member")
	}
	owner := repo.createdMembers[0]
	if owner.MemberID != 7 || owner.Status != uint32(enums.DealMemberStatusAccepted) {
		t.Fatalf("owner must be Accepted, got %+v", owner)
	}
	if owner.RoleKey != uint32(enums.RoleKeyDealOwner) || !owner.IsOwner {
		t.Fatalf("owner role/flag mismatch: %+v", owner)
	}
	if len(repo.createdProduct) != 1 || repo.createdProduct[0] != 22 {
		t.Fatalf("expected product link, got %v", repo.createdProduct)
	}
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

func (stubProductRepo) ListByIDs(context.Context, []uint64) ([]*entity.Product, error) {
	return nil, nil
}

func (stubProductRepo) Create(context.Context, *entity.Product) error { return nil }

func (stubProductRepo) SoftDelete(context.Context, uint64, uint64) error { return nil }
