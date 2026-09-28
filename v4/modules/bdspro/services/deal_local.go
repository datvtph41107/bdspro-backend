package services

import (
	"context"
	"fmt"
	"strconv"
	"time"

	_errors "common/errors"
	_utils "common/utils"
	bdspropb "pb/types/bdspro"
	sharepb "pb/types/shared"

	"userv4/internal/interface/providers"
	"userv4/modules/bdspro/dto"
	"userv4/modules/bdspro/entity"
	"userv4/modules/bdspro/enums"
	"userv4/modules/bdspro/repository"
)

// LocalDealService owns deal user-list against the v4 DSN + SyncUtil.
type LocalDealService struct {
	sync     *_utils.SyncUtil
	repo     repository.DealRepository
	products repository.ProductRepository
	profile  providers.ProfileProvider
}

func NewLocalDealService(
	sync *_utils.SyncUtil,
	repo repository.DealRepository,
	products repository.ProductRepository,
	profile providers.ProfileProvider,
) *LocalDealService {
	return &LocalDealService{sync: sync, repo: repo, products: products, profile: profile}
}

func (s *LocalDealService) GetTimestamps(ctx context.Context, id uint64) (*bdspropb.DealTimestampsResponse, error) {
	if id == 0 {
		return nil, _errors.BadRequestException("deal id is required")
	}
	if s.sync == nil {
		return &bdspropb.DealTimestampsResponse{Timestamps: map[string]int64{}}, nil
	}

	redisKeys := []string{
		s.sync.GetKey(ctx, _utils.SyncKeyDealPermission, id),
		s.sync.GetKey(ctx, _utils.SyncKeyDealProducts, id),
		s.sync.GetKey(ctx, _utils.SyncKeyDealContracts, id),
		s.sync.GetKey(ctx, _utils.SyncKeyDealInvestment, id),
		s.sync.GetKey(ctx, _utils.SyncKeyDealSummary, id),
		s.sync.GetKey(ctx, _utils.SyncKeyDealBank, id),
		s.sync.GetKey(ctx, _utils.SyncKeyDealInvitation, id),
		s.sync.GetKey(ctx, _utils.SyncKeyDealStatistics, id),
		s.sync.GetKey(ctx, _utils.SyncKeyDealCost, id),
		s.sync.GetKey(ctx, _utils.SyncKeyDealDetail, id),
		s.sync.GetKey(ctx, _utils.SyncKeyDealMembers, id),
		s.sync.GetKey(ctx, _utils.SyncKeyDealProcess, id),
		s.sync.GetKey(ctx, _utils.SyncKeyDealCommission, id),
		s.sync.GetKey(ctx, _utils.SyncKeyDealHistory, id),
	}
	names := []string{
		"permission", "products", "transaction", "investment", "summary",
		"bank", "invitation", "statistics", "cost", "detail",
		"members", "process", "commission", "history",
	}

	vals, err := s.sync.MGet(ctx, redisKeys)
	if err != nil {
		return nil, _errors.InternalServerException(err.Error())
	}

	timestamps := make(map[string]int64, len(names))
	for i, v := range vals {
		timestamps[names[i]] = parseSyncTimestamp(v)
	}
	return &bdspropb.DealTimestampsResponse{Timestamps: timestamps}, nil
}

func (s *LocalDealService) GetDetail(ctx context.Context, req dto.DealDetailRequest) (*bdspropb.GroupDeal, int64, error) {
	if req.ID == 0 {
		return nil, 0, _errors.BadRequestException("deal id is required")
	}
	var (
		key string
		ts  int64
	)
	if s.sync != nil {
		key = s.sync.GetKey(ctx, _utils.SyncKeyDealDetail, req.ID)
		ts = s.sync.GetTimestamp(ctx, key)
		if !s.sync.HasUpdated(ctx, key, req.Timestamp) && req.Timestamp > 0 {
			if ts <= 0 {
				ts = req.Timestamp
			}
			return &bdspropb.GroupDeal{Id: uint32(req.ID)}, ts, nil
		}
	}

	deal, err := s.repo.GetByID(ctx, req.ID)
	if err != nil {
		return nil, 0, err
	}
	if deal == nil {
		return nil, 0, _errors.NotFoundException("deal not found")
	}

	if ts <= 0 {
		if s.sync != nil && key != "" {
			ts = s.sync.EnsureTimestamp(ctx, key)
		} else {
			ts = time.Now().UnixMilli()
		}
	}
	return mapDealDetailToPB(deal), ts, nil
}

func (s *LocalDealService) GetContractProcess(ctx context.Context, dealID uint64) (*bdspropb.ContractProcessResponse, int64, error) {
	if dealID == 0 {
		return nil, 0, _errors.BadRequestException("deal id is required")
	}
	rows, err := s.repo.GetContractProcess(ctx, dealID)
	if err != nil {
		return nil, 0, err
	}

	data := make([]*bdspropb.ContractProcess, 0, len(rows))
	for _, row := range rows {
		action := enums.TxAction(row.Action)
		style := enums.GetTxActionStyle(action)
		item := &bdspropb.ContractProcess{
			StatusName:    enums.TxActionNames[action],
			Timestamp:     _utils.FormatTimeToString(row.Timestamp),
			Code:          fmt.Sprintf("TXT%06d", row.TransactionID),
			TransactionId: row.TransactionID,
			Color:         style.Color,
			BgColor:       style.BgColor,
			BorderColor:   style.BorderColor,
		}
		data = append(data, item)
	}

	var ts int64
	if s.sync != nil {
		key := s.sync.GetKey(ctx, _utils.SyncKeyDealProcess, dealID)
		now := time.Now().UnixMilli()
		_ = s.sync.PutTimestamp(ctx, key, now)
		ts = s.sync.GetTimestamp(ctx, key)
		if ts <= 0 {
			ts = now
		}
	} else {
		ts = time.Now().UnixMilli()
	}
	return &bdspropb.ContractProcessResponse{Data: data}, ts, nil
}

func (s *LocalDealService) GetDealsUser(ctx context.Context, req dto.DealSearchRequest) (*bdspropb.GetDealsResponse, error) {
	var syncKey string
	if s.sync != nil {
		syncKey = s.sync.GetKey(ctx, _utils.SyncKeyDealMe, 0)
		if !s.sync.HasUpdated(ctx, syncKey, req.Timestamp) && req.Page == 0 && req.Timestamp > 0 {
			ts := s.sync.GetTimestamp(ctx, syncKey)
			if ts <= 0 {
				ts = req.Timestamp
			}
			return &bdspropb.GetDealsResponse{Timestamp: ts}, nil
		}
	}

	profileID := _utils.GetProfileIdWithContext(ctx)
	size := req.Size
	if size == 0 {
		size = 20
	}
	list, total, err := s.repo.ListForUser(ctx, profileID, req.Page, size, req)
	if err != nil {
		return nil, err
	}

	data := make([]*bdspropb.GroupDeal, 0, len(list))
	for _, d := range list {
		data = append(data, mapDealToPB(d))
	}

	resp := &bdspropb.GetDealsResponse{
		Data:  data,
		Total: total,
	}
	if s.sync != nil && syncKey != "" && req.Page == 0 {
		now := time.Now().UnixMilli()
		_ = s.sync.PutTimestamp(ctx, syncKey, now)
		ts := s.sync.GetTimestamp(ctx, syncKey)
		if ts <= 0 {
			ts = now
		}
		resp.Timestamp = ts
	} else if req.Page == 0 {
		resp.Timestamp = time.Now().UnixMilli()
	}
	return resp, nil
}

func (s *LocalDealService) GetDealsByProduct(ctx context.Context, productID uint64, req dto.SyncByIDRequest) (*bdspropb.GetOrganizationDealsResponse, error) {
	if productID == 0 {
		return &bdspropb.GetOrganizationDealsResponse{}, nil
	}
	if s.sync != nil {
		key := s.sync.GetKey(ctx, _utils.SyncKeyProductDeals, productID)
		if !s.sync.HasUpdated(ctx, key, req.Timestamp) && req.Page == 0 && req.Timestamp > 0 {
			return &bdspropb.GetOrganizationDealsResponse{}, nil
		}
	}
	size := req.Size
	if size == 0 {
		size = 20
	}
	list, total, err := s.repo.ListByProductID(ctx, productID, req.Page, size)
	if err != nil {
		return nil, err
	}
	data := make([]*bdspropb.GroupDeal, 0, len(list))
	for _, d := range list {
		data = append(data, mapDealToPB(d))
	}
	return &bdspropb.GetOrganizationDealsResponse{Data: data, Total: total}, nil
}

// CreateGroupDeal inserts a deal with its owner member, customers and products.
// Owner member is Accepted (status=20) so it appears in ListForUser / GET /v4/deals/me,
// matching bdspro CreateDeal membership rules.
func (s *LocalDealService) CreateGroupDeal(ctx context.Context, req *bdspropb.CreateGroupDealRequest) (uint64, error) {
	profileID := _utils.GetProfileIdWithContext(ctx)

	ownerID := req.GetOwnerId()
	if ownerID == 0 {
		ownerID = profileID
	}

	deal := &entity.Deal{
		Name:         req.GetName(),
		TargetProfit: req.GetTargetProfit(),
		Status:       req.GetStatus(),
		DealType:     req.GetDealType(),
		OwnerID:      ownerID,
		OwnerType:    req.GetOwnerType(),
		Note:         req.GetNote(),
		Description:  req.GetDescription(),
		IsUnilateral: req.GetIsUnilateral(),
		IsManual:     req.GetIsManual(),
	}
	if pid := profileID; pid > 0 {
		deal.CreatedBy = &pid
	}
	if req.GetChargePersonId() > 0 {
		cid := req.GetChargePersonId()
		deal.ChargePersonID = &cid
	}
	if bid := req.GetBankAccountId(); bid > 0 {
		deal.BankAccountID = &bid
	}

	var members []entity.DealMember
	// Creator/owner: Accepted so deals/me lists this deal (same as bdspro CreateDeal).
	members = append(members, entity.DealMember{
		MemberID:   ownerID,
		MemberType: uint32(enums.DealMemberTypeMember),
		IsOwner:    true,
		RoleKey:    uint32(enums.RoleKeyDealOwner),
		Status:     uint32(enums.DealMemberStatusAccepted),
	})
	if cid := req.GetChargePersonId(); cid > 0 && cid != ownerID {
		members = append(members, entity.DealMember{
			MemberID:   cid,
			MemberType: uint32(enums.DealMemberTypeMember),
			RoleKey:    uint32(enums.RoleKeyDealAdmin),
			Status:     uint32(enums.DealMemberStatusInvited),
		})
	}
	for _, cid := range req.GetCustomerIds() {
		if cid == 0 || cid == ownerID {
			continue
		}
		members = append(members, entity.DealMember{
			MemberID:   cid,
			MemberType: uint32(enums.DealMemberTypeCustomer),
			RoleKey:    uint32(enums.RoleKeyDealCustomer),
			Status:     uint32(enums.DealMemberStatusInvited),
		})
	}
	for _, mid := range req.GetMemberIds() {
		if mid == 0 || mid == ownerID {
			continue
		}
		members = append(members, entity.DealMember{
			MemberID:   mid,
			MemberType: uint32(enums.DealMemberTypeMember),
			RoleKey:    uint32(enums.RoleKeyDealMember),
			Status:     uint32(enums.DealMemberStatusInvited),
		})
	}
	for _, pid := range req.GetPartnerIds() {
		if pid == 0 || pid == ownerID {
			continue
		}
		members = append(members, entity.DealMember{
			MemberID:   pid,
			MemberType: uint32(enums.DealMemberTypePartner),
			RoleKey:    uint32(enums.RoleKeyDealPartner),
			Status:     uint32(enums.DealMemberStatusInvited),
		})
	}

	if err := s.repo.Create(ctx, deal, members, req.GetProductIds()); err != nil {
		return 0, err
	}

	if s.sync != nil {
		key := s.sync.GetKey(ctx, _utils.SyncKeyDealMe, 0)
		_ = s.sync.PutTimestamp(ctx, key, time.Now().UnixMilli())
	}

	return deal.ID, nil
}

// SeedDemoRelated ensures bank/investment/cost/contract rows exist for a demo deal.
func (s *LocalDealService) SeedDemoRelated(ctx context.Context, dealID, profileID, customerID uint64, dealType uint32) error {
	if dealID == 0 {
		return _errors.BadRequestException("deal id is required")
	}
	if err := s.repo.SeedDemoRelated(ctx, dealID, profileID, customerID, dealType); err != nil {
		return err
	}
	if s.sync != nil {
		now := time.Now().UnixMilli()
		for _, prefix := range []string{
			_utils.SyncKeyDealBank,
			_utils.SyncKeyDealInvestment,
			_utils.SyncKeyDealCost,
			_utils.SyncKeyDealContracts,
			_utils.SyncKeyDealSummary,
			_utils.SyncKeyDealStatistics,
			_utils.SyncKeyDealInvitation,
			_utils.SyncKeyDealProducts,
		} {
			_ = s.sync.PutTimestamp(ctx, s.sync.GetKey(ctx, prefix, dealID), now)
		}
	}
	return nil
}

// DeleteGroupDeal soft-deletes a deal by id.
func (s *LocalDealService) DeleteGroupDeal(ctx context.Context, id uint64) error {
	if err := s.repo.SoftDelete(ctx, id); err != nil {
		return err
	}
	if s.sync != nil {
		key := s.sync.GetKey(ctx, _utils.SyncKeyDealMe, 0)
		_ = s.sync.PutTimestamp(ctx, key, 0)
	}
	return nil
}

func mapDealToPB(d *entity.Deal) *bdspropb.GroupDeal {
	if d == nil {
		return nil
	}
	out := &bdspropb.GroupDeal{
		Id:                      uint32(d.ID),
		Name:                    d.Name,
		TargetProfit:            d.TargetProfit,
		Status:                  d.Status,
		DealType:                d.DealType,
		OwnerId:                 uint32(d.OwnerID),
		OwnerType:               d.OwnerType,
		Note:                    d.Note,
		IsUnilateral:            d.IsUnilateral,
		IsManual:                d.IsManual,
		CancelReason:            d.CancelReason,
		AllowSharing:            d.AllowSharing,
		MemberCanAddTransaction: d.MemberCanAddTransaction,
		OnlyOwnerGetCommission:  d.OnlyOwnerGetCommission,
		InternalNote:            d.InternalNote,
		ShareVisibility:         d.ShareVisibility,
		Description:             d.Description,
		StatusName:              enums.DealStatusMap[enums.DealStatus(d.Status)],
		DealTypeName:            enums.DealTypeNames[enums.DealType(d.DealType)],
		NumMember:               uint32(len(d.Members)),
		Code:                    fmt.Sprintf("TV%06d", d.ID),
	}
	if d.ChargePersonID != nil {
		out.ChargePersonId = *d.ChargePersonID
	}
	if d.CreatedAt != nil {
		out.CreatedAt = strconv.FormatInt(d.CreatedAt.UnixMilli(), 10)
	}
	if d.UpdatedAt != nil {
		out.UpdatedAt = strconv.FormatInt(d.UpdatedAt.UnixMilli(), 10)
	}
	if d.FromDate != nil {
		out.FromDate = d.FromDate.Format(time.RFC3339)
	}
	if d.ToDate != nil {
		out.ToDate = d.ToDate.Format(time.RFC3339)
	}

	members := make([]*sharepb.ProfileItem, 0)
	customers := make([]*sharepb.ProfileItem, 0)
	partners := make([]*sharepb.ProfileItem, 0)
	for _, m := range d.Members {
		item := &sharepb.ProfileItem{Id: m.MemberID}
		switch enums.DealMemberType(m.MemberType) {
		case enums.DealMemberTypeCustomer:
			customers = append(customers, item)
		case enums.DealMemberTypePartner:
			partners = append(partners, item)
		default:
			members = append(members, item)
		}
	}
	out.Members = members
	out.Customers = customers
	out.Partners = partners
	return out
}

func mapDealDetailToPB(d *entity.Deal) *bdspropb.GroupDeal {
	out := mapDealToPB(d)
	if out == nil {
		return nil
	}
	// Detail response uses string times like bdspro EntityToDealDetailResponse.
	if d.CreatedAt != nil {
		out.CreatedAt = _utils.FormatTimeToString(d.CreatedAt)
	}
	if d.FromDate != nil {
		out.FromDate = _utils.FormatTimeToString(d.FromDate)
	}
	if d.ToDate != nil {
		out.ToDate = _utils.FormatTimeToString(d.ToDate)
	}
	return out
}
