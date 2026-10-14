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

// DealService owns deal list/detail/timestamps/process flows for modules/bdspro.
type DealService interface {
	GetDealsUser(ctx context.Context, req dto.DealSearchRequest) (*bdspropb.GetDealsResponse, error)
	// GetDealsByProduct returns deal_products rows for the product, each with nested deal detail.
	GetDealsByProduct(ctx context.Context, productID uint64, req dto.SyncByIDRequest) ([]dto.ProductDealItem, int64, error)
	GetTimestamps(ctx context.Context, id uint64) (*bdspropb.DealTimestampsResponse, error)
	GetDetail(ctx context.Context, req dto.DealDetailRequest) (*bdspropb.GroupDeal, int64, error)
	GetContractProcess(ctx context.Context, dealID uint64) (*bdspropb.ContractProcessResponse, int64, error)
	GetDealProducts(ctx context.Context, req dto.DealListByIDRequest) ([]dto.DealProductItem, int64, error)
	GetContracts(ctx context.Context, req dto.DealListByIDRequest) (*bdspropb.DealContractListResponse, error)
	GetInvestments(ctx context.Context, req dto.DealListByIDRequest) (*bdspropb.ListInvestment, int64, error)
	GetInvitations(ctx context.Context, req dto.DealListByIDRequest) (*bdspropb.GetDealInvitationsByDealResponse, error)
	GetMembers(ctx context.Context, req dto.DealMembersRequest) (*bdspropb.GetDealMembersResponse, int64, error)
	GetDealCosts(ctx context.Context, req dto.DealListByIDRequest) (*bdspropb.DealCostListDTO, error)
	GetSummary(ctx context.Context, req dto.DealSummaryRequest) (*bdspropb.SummaryResponse, int64, error)
	GetBankAccount(ctx context.Context, req dto.DealDetailRequest) (*bdspropb.BankAccountItem, int64, error)
	GetContractStatistics(ctx context.Context, req dto.DealStatisticsRequest) (*bdspropb.ContractDealStatisticsDTO, int64, error)
	CreateGroupDeal(ctx context.Context, req *bdspropb.CreateGroupDealRequest) (uint64, error)
	UpdateGroupDeal(ctx context.Context, req *bdspropb.UpdateGroupDealRequest) (uint64, error)
	SeedDemoRelated(ctx context.Context, dealID, profileID, customerID uint64, dealType uint32) error
	DeleteGroupDeal(ctx context.Context, id uint64) error
}

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

	invitationMemberID := profileID
	if req.InvitedProfileID != nil && *req.InvitedProfileID > 0 {
		invitationMemberID = *req.InvitedProfileID
	}
	invitationRows := make([]*entity.DealMember, 0, len(list))
	for _, d := range list {
		if m := dealMemberForProfile(d.Members, invitationMemberID); m != nil {
			invitationRows = append(invitationRows, m)
		}
	}
	profilesMap := s.profileItemsByMemberIDs(ctx, invitationRows)

	data := make([]*bdspropb.GroupDeal, 0, len(list))
	for _, d := range list {
		pb := mapDealToPB(d)
		if m := dealMemberForProfile(d.Members, invitationMemberID); m != nil {
			pb.Invitation = mapDealMemberToSharePB(m, profilesMap[m.MemberID])
		}
		data = append(data, pb)
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

func (s *LocalDealService) GetDealsByProduct(ctx context.Context, productID uint64, req dto.SyncByIDRequest) ([]dto.ProductDealItem, int64, error) {
	if productID == 0 {
		return []dto.ProductDealItem{}, 0, nil
	}
	if s.sync != nil {
		key := s.sync.GetKey(ctx, _utils.SyncKeyProductDeals, productID)
		if !s.sync.HasUpdated(ctx, key, req.Timestamp) && req.Page == 0 && req.Timestamp > 0 {
			ts := s.sync.GetTimestamp(ctx, key)
			if ts <= 0 {
				ts = req.Timestamp
			}
			return []dto.ProductDealItem{}, ts, nil
		}
	}
	size := req.Size
	if size == 0 {
		size = 20
	}
	rows, err := s.repo.ListDealProductsByProductID(ctx, productID, req.Page, size)
	if err != nil {
		return nil, 0, err
	}

	var ts int64
	if s.sync != nil {
		key := s.sync.GetKey(ctx, _utils.SyncKeyProductDeals, productID)
		ts = s.sync.EnsureTimestamp(ctx, key)
	}
	if ts <= 0 {
		ts = time.Now().UnixMilli()
	}
	if len(rows) == 0 {
		return []dto.ProductDealItem{}, ts, nil
	}

	dealIDs := make([]uint64, 0, len(rows))
	seen := make(map[uint64]struct{}, len(rows))
	for _, row := range rows {
		if row == nil || row.DealID == 0 {
			continue
		}
		if _, ok := seen[row.DealID]; ok {
			continue
		}
		seen[row.DealID] = struct{}{}
		dealIDs = append(dealIDs, row.DealID)
	}

	byID := map[uint64]*entity.Deal{}
	if len(dealIDs) > 0 {
		deals, err := s.repo.ListByIDs(ctx, dealIDs)
		if err != nil {
			return nil, 0, err
		}
		for _, d := range deals {
			if d != nil {
				byID[d.ID] = d
			}
		}
	}

	out := make([]dto.ProductDealItem, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		item := dto.ProductDealItem{
			ID:        row.ID,
			CreatedAt: _utils.FormatTimeToString(row.CreatedAt),
			UpdatedAt: _utils.FormatTimeToString(row.UpdatedAt),
			CreatedBy: row.CreatedBy,
			UpdatedBy: row.UpdatedBy,
			DealID:    row.DealID,
			ProductID: row.ProductID,
		}
		if d := byID[row.DealID]; d != nil {
			item.Deal = mapDealDetailToPB(d)
		} else {
			item.Deal = mapDealDetailToPB(&entity.Deal{ID: row.DealID})
		}
		out = append(out, item)
	}
	return out, ts, nil
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

// UpdateGroupDeal updates core deal fields (parity with bdspro UpdateGroupDeal).
func (s *LocalDealService) UpdateGroupDeal(ctx context.Context, req *bdspropb.UpdateGroupDealRequest) (uint64, error) {
	id := uint64(req.GetId())
	if id == 0 {
		return 0, _errors.BadRequestException("deal id is required")
	}
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return 0, err
	}
	if existing == nil {
		return 0, _errors.NotFoundException("deal not found")
	}

	patch := &entity.Deal{
		ID:              id,
		Name:            req.GetName(),
		Status:          req.GetStatus(),
		ShareVisibility: req.GetShareVisibility(),
		Description:     req.GetDescription(),
		TargetProfit:    req.GetAmount(),
	}
	if pid := _utils.GetProfileIdWithContext(ctx); pid > 0 {
		patch.UpdatedBy = &pid
	}
	if err := s.repo.Update(ctx, patch); err != nil {
		return 0, err
	}

	if s.sync != nil {
		key := s.sync.GetKey(ctx, _utils.SyncKeyDealMe, 0)
		_ = s.sync.PutTimestamp(ctx, key, time.Now().UnixMilli())
	}
	_ = s.syncTS(ctx, _utils.SyncKeyDealDetail, id, true)

	return id, nil
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
		_ = s.sync.PutTimestamp(ctx, key, time.Now().UnixMilli())
	}
	_ = s.syncTS(ctx, _utils.SyncKeyDealDetail, id, true)
	return nil
}

func dealMemberForProfile(members []entity.DealMember, profileID uint64) *entity.DealMember {
	if profileID == 0 {
		return nil
	}
	for i := range members {
		if members[i].MemberID == profileID {
			return &members[i]
		}
	}
	return nil
}

func mapDealMemberToSharePB(m *entity.DealMember, profile *sharepb.ProfileItem) *sharepb.DealMember {
	if m == nil {
		return nil
	}
	item := &sharepb.DealMember{
		Id:              m.ID,
		DealId:          m.DealID,
		InviterId:       m.InviterID,
		MemberId:        m.MemberID,
		Status:          m.Status,
		Message:         m.Message,
		RoleKey:         m.RoleKey,
		IsOwner:         m.IsOwner,
		AmountCommit:    m.AmountCommit,
		CommissionValue: m.CommissionValue,
		CommissionType:  m.CommissionType,
		Note:            m.Note,
		DoneInvestment:  m.DoneInvestment,
		RoleId:          &m.RoleID,
		Profile:         profile,
	}
	if m.InvitedAt != nil {
		item.InvitedAt = _utils.FormatTimeToString(m.InvitedAt)
	}
	if m.RespondedAt != nil {
		item.RespondedAt = _utils.FormatTimeToString(m.RespondedAt)
	}
	if m.WithdrawnAt != nil {
		item.WithdrawnAt = _utils.FormatTimeToString(m.WithdrawnAt)
	}
	if m.RoleKey != 0 {
		item.Role = dealMemberRoleFromKey(m.RoleKey)
	}
	return item
}

func dealMemberRoleFromKey(roleKey uint32) *sharepb.Role {
	name := enums.RoleKeyMap[enums.RoleKey(roleKey)]
	return &sharepb.Role{
		RoleKey:  roleKey,
		Key:      strconv.FormatUint(uint64(roleKey), 10),
		RoleName: name,
	}
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
func (s *LocalDealService) syncTS(ctx context.Context, prefix string, dealID uint64, put bool) int64 {
	if s.sync == nil || dealID == 0 {
		return time.Now().UnixMilli()
	}
	key := s.sync.GetKey(ctx, prefix, dealID)
	if put {
		now := time.Now().UnixMilli()
		_ = s.sync.PutTimestamp(ctx, key, now)
	}
	ts := s.sync.EnsureTimestamp(ctx, key)
	if ts <= 0 {
		ts = time.Now().UnixMilli()
	}
	return ts
}

func (s *LocalDealService) syncUnchanged(ctx context.Context, prefix string, dealID uint64, timestamp int64, page uint32) (ts int64, skip bool) {
	if s.sync == nil || dealID == 0 || timestamp <= 0 {
		return 0, false
	}
	key := s.sync.GetKey(ctx, prefix, dealID)
	if s.sync.HasUpdated(ctx, key, timestamp) {
		return 0, false
	}
	if page != 0 {
		return 0, false
	}
	ts = s.sync.GetTimestamp(ctx, key)
	if ts <= 0 {
		ts = timestamp
	}
	return ts, true
}

func (s *LocalDealService) GetDealProducts(ctx context.Context, req dto.DealListByIDRequest) ([]dto.DealProductItem, int64, error) {
	if req.DealID == 0 {
		return nil, 0, _errors.BadRequestException("deal id is required")
	}
	if ts, skip := s.syncUnchanged(ctx, _utils.SyncKeyDealProducts, req.DealID, req.Timestamp, req.Page); skip {
		return []dto.DealProductItem{}, ts, nil
	}

	rows, err := s.repo.ListDealProducts(ctx, req.DealID, req.Page, req.Size)
	if err != nil {
		return nil, 0, err
	}
	ts := s.syncTS(ctx, _utils.SyncKeyDealProducts, req.DealID, false)
	if len(rows) == 0 {
		return []dto.DealProductItem{}, ts, nil
	}

	productIDs := make([]uint64, 0, len(rows))
	seen := make(map[uint64]struct{}, len(rows))
	for _, row := range rows {
		if row == nil || row.ProductID == 0 {
			continue
		}
		if _, ok := seen[row.ProductID]; ok {
			continue
		}
		seen[row.ProductID] = struct{}{}
		productIDs = append(productIDs, row.ProductID)
	}

	byID := map[uint64]*entity.Product{}
	if s.products != nil && len(productIDs) > 0 {
		products, err := s.products.ListByIDs(ctx, productIDs)
		if err != nil {
			return nil, 0, err
		}
		for _, p := range products {
			if p != nil {
				byID[p.ID] = p
			}
		}
	}

	out := make([]dto.DealProductItem, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		item := dto.DealProductItem{
			ID:        row.ID,
			CreatedAt: _utils.FormatTimeToString(row.CreatedAt),
			UpdatedAt: _utils.FormatTimeToString(row.UpdatedAt),
			CreatedBy: row.CreatedBy,
			UpdatedBy: row.UpdatedBy,
			DealID:    row.DealID,
			ProductID: row.ProductID,
		}
		if p := byID[row.ProductID]; p != nil {
			item.Product = mapProductToPB(p)
		} else {
			item.Product = mapProductToPB(&entity.Product{ID: row.ProductID})
		}
		out = append(out, item)
	}

	// Client already received product detail here — stamp prd:id so later
	// GET product/{id}/detail can skip when unchanged.
	if s.sync != nil && len(productIDs) > 0 {
		now := time.Now().UnixMilli()
		for _, pid := range productIDs {
			_ = s.sync.PutTimestamp(ctx, s.sync.GetKey(ctx, _utils.SyncKeyProductDetail, pid), now)
		}
	}
	return out, ts, nil
}

func (s *LocalDealService) GetContracts(ctx context.Context, req dto.DealListByIDRequest) (*bdspropb.DealContractListResponse, error) {
	if req.DealID == 0 {
		return nil, _errors.BadRequestException("deal id is required")
	}
	if ts, skip := s.syncUnchanged(ctx, _utils.SyncKeyDealContracts, req.DealID, req.Timestamp, req.Page); skip {
		return &bdspropb.DealContractListResponse{Timestamp: ts}, nil
	}

	rows, total, err := s.repo.ListContractsByDealID(ctx, req.DealID)
	if err != nil {
		return nil, err
	}
	data := make([]*bdspropb.DealContractItem, 0, len(rows))
	for _, row := range rows {
		item := &bdspropb.DealContractItem{
			Id:         row.ContractID,
			CustomerId: row.CustomerID,
			CreatedAt:  _utils.FormatTimeToString(row.Timestamp),
			ContractId: row.ContractID,
			Amount:     row.Amount,
			Status:     uint64(row.Status),
			Code:       fmt.Sprintf("TXT%06d", row.ContractID),
		}
		if row.TransactionName == "" {
			item.TransactionName = enums.TxMethodNames[enums.TxMethod(row.Method)]
			item.StatusName = enums.TxApprovedStatusMap[enums.TxApprovedStatus(row.Status)]
		} else {
			item.TransactionName = row.TransactionName
			item.StatusName = enums.TxApprovedStatusMap[enums.TxApprovedStatus(row.Status)]
		}
		data = append(data, item)
	}
	ts := s.syncTS(ctx, _utils.SyncKeyDealContracts, req.DealID, false)
	return &bdspropb.DealContractListResponse{
		Data:      data,
		Total:     int32(total),
		Timestamp: ts,
	}, nil
}

func (s *LocalDealService) GetInvestments(ctx context.Context, req dto.DealListByIDRequest) (*bdspropb.ListInvestment, int64, error) {
	if req.DealID == 0 {
		return nil, 0, _errors.BadRequestException("deal id is required")
	}
	// Investments short-circuit ignores page (matches bdspro handler).
	if s.sync != nil && req.Timestamp > 0 {
		key := s.sync.GetKey(ctx, _utils.SyncKeyDealInvestment, req.DealID)
		if !s.sync.HasUpdated(ctx, key, req.Timestamp) {
			ts := s.sync.GetTimestamp(ctx, key)
			if ts <= 0 {
				ts = req.Timestamp
			}
			return &bdspropb.ListInvestment{}, ts, nil
		}
	}

	var memberID uint64
	if req.MemberID != nil {
		memberID = *req.MemberID
	}
	list, total, err := s.repo.ListInvestmentsByDealID(ctx, req.DealID, memberID)
	if err != nil {
		return nil, 0, err
	}
	data := make([]*bdspropb.Investment, 0, len(list))
	for _, inv := range list {
		if inv == nil {
			continue
		}
		item := &bdspropb.Investment{
			Id:                 inv.ID,
			DealId:             inv.DealID,
			MemberId:           inv.MemberID,
			Amount:             inv.Amount,
			TransferTime:       _utils.FormatTimeToString(inv.TransferTime),
			Note:               inv.Note,
			TransferProofImage: inv.TransferProofImage,
			InvestType:         inv.InvestType,
			Status:             inv.Status,
			StatusName:         enums.TxApprovedStatusMap[enums.TxApprovedStatus(inv.Status)],
			ChangedAt:          _utils.FormatTimeToString(inv.TransferTime),
			Confirmed:          inv.Confirmed,
		}
		if inv.CreatedBy != nil {
			item.CreatedBy = inv.CreatedBy
		}
		if inv.ProxyUserID != nil {
			item.ProxyUserId = inv.ProxyUserID
		}
		data = append(data, item)
	}
	ts := s.syncTS(ctx, _utils.SyncKeyDealInvestment, req.DealID, false)
	return &bdspropb.ListInvestment{Data: data, Total: int32(total)}, ts, nil
}

func (s *LocalDealService) GetInvitations(ctx context.Context, req dto.DealListByIDRequest) (*bdspropb.GetDealInvitationsByDealResponse, error) {
	if req.DealID == 0 {
		return nil, _errors.BadRequestException("deal id is required")
	}
	if ts, skip := s.syncUnchanged(ctx, _utils.SyncKeyDealInvitation, req.DealID, req.Timestamp, req.Page); skip {
		return &bdspropb.GetDealInvitationsByDealResponse{Timestamp: ts}, nil
	}

	page := int(req.Page)
	size := int(req.Size)
	if size == 0 {
		size = 10
	}
	list, total, err := s.repo.ListMembersByDealID(ctx, req.DealID, page, size)
	if err != nil {
		return nil, err
	}
	profilesMap := s.profileItemsByMemberIDs(ctx, list)
	data := make([]*sharepb.DealMember, 0, len(list))
	for _, m := range list {
		if m == nil {
			continue
		}
		data = append(data, mapDealMemberToSharePB(m, profilesMap[m.MemberID]))
	}
	ts := s.syncTS(ctx, _utils.SyncKeyDealInvitation, req.DealID, false)
	return &bdspropb.GetDealInvitationsByDealResponse{
		Data:      data,
		Total:     int64(total),
		Timestamp: ts,
	}, nil
}

func (s *LocalDealService) GetMembers(ctx context.Context, req dto.DealMembersRequest) (*bdspropb.GetDealMembersResponse, int64, error) {
	if req.DealID == 0 {
		return nil, 0, _errors.BadRequestException("deal id is required")
	}
	if ts, skip := s.syncUnchanged(ctx, _utils.SyncKeyDealMembers, req.DealID, req.Timestamp, req.Page); skip {
		return &bdspropb.GetDealMembersResponse{}, ts, nil
	}

	list, err := s.repo.ListAcceptedMembers(ctx, req)
	if err != nil {
		return nil, 0, err
	}

	profilesMap := s.profileItemsByMemberIDs(ctx, list)
	data := make([]*sharepb.DealMember, 0, len(list))
	for _, m := range list {
		if m == nil {
			continue
		}
		data = append(data, mapDealMemberToSharePB(m, profilesMap[m.MemberID]))
	}

	ts := s.syncTS(ctx, _utils.SyncKeyDealMembers, req.DealID, false)
	return &bdspropb.GetDealMembersResponse{
		Data:  data,
		Total: uint32(len(data)),
	}, ts, nil
}

func (s *LocalDealService) GetDealCosts(ctx context.Context, req dto.DealListByIDRequest) (*bdspropb.DealCostListDTO, error) {
	if req.DealID == 0 {
		return nil, _errors.BadRequestException("deal id is required")
	}
	if ts, skip := s.syncUnchanged(ctx, _utils.SyncKeyDealCost, req.DealID, req.Timestamp, req.Page); skip {
		return &bdspropb.DealCostListDTO{Timestamp: ts}, nil
	}

	list, err := s.repo.ListCostsByDealID(ctx, req.DealID, req.Page, req.Size)
	if err != nil {
		return nil, err
	}
	data := make([]*bdspropb.DealCostDTO, 0, len(list))
	for _, c := range list {
		if c == nil {
			continue
		}
		item := &bdspropb.DealCostDTO{
			Id:          c.ID,
			DealId:      c.DealID,
			Amount:      c.Amount,
			ConfirmedAt: _utils.FormatTimeToString(c.ConfirmedAt),
			Note:        c.Note,
			CostTypeId:  c.CostTypeID,
			Status:      c.Status,
			CostName:    c.CostName,
		}
		if c.CostType != nil {
			item.CostTypeName = c.CostType.TypeName
		}
		data = append(data, item)
	}
	ts := s.syncTS(ctx, _utils.SyncKeyDealCost, req.DealID, false)
	return &bdspropb.DealCostListDTO{
		Data:          data,
		TotalElements: uint32(len(data)),
		Timestamp:     ts,
	}, nil
}

func (s *LocalDealService) GetSummary(ctx context.Context, req dto.DealSummaryRequest) (*bdspropb.SummaryResponse, int64, error) {
	if req.DealID == 0 {
		return nil, 0, _errors.BadRequestException("deal id is required")
	}
	if ts, skip := s.syncUnchanged(ctx, _utils.SyncKeyDealSummary, req.DealID, req.Timestamp, 0); skip {
		return &bdspropb.SummaryResponse{}, ts, nil
	}

	deal, err := s.repo.GetByID(ctx, req.DealID)
	if err != nil {
		return nil, 0, err
	}
	if deal == nil {
		return nil, 0, _errors.NotFoundException("deal not found")
	}

	info, err := s.repo.CountDashboard(ctx, req)
	if err != nil {
		return nil, 0, err
	}
	if info == nil {
		info = &dto.DealSummaryRow{}
	}

	percentDone := 0.0
	if deal.TargetProfit > 0 && (req.AmountApproved || req.All) {
		percentDone = info.AmountApproved / deal.TargetProfit
	}
	amountRest := 0.0
	if deal.TargetProfit > info.AmountApproved && (req.AmountApproved || req.All) {
		amountRest = deal.TargetProfit - info.AmountApproved
	}

	ts := s.syncTS(ctx, _utils.SyncKeyDealSummary, req.DealID, false)
	return &bdspropb.SummaryResponse{
		AmountPending:   info.AmountPending,
		AmountApproved:  info.AmountApproved,
		AmountRejected:  info.AmountRejected,
		AmountRest:      amountRest,
		NumPending:      info.NumPending,
		NumApproved:     info.NumApproved,
		NumRejected:     info.NumRejected,
		NumInvestment:   info.NumInvestment,
		NumMemberSubmit: info.NumMemberSubmit,
		AmountTarget:    deal.TargetProfit,
		PercentDone:     percentDone,
		AmountCost:      info.AmountCost,
		AmountProfit:    info.AmountApproved,
		AmountTotal:     10000000000,
	}, ts, nil
}

func (s *LocalDealService) GetBankAccount(ctx context.Context, req dto.DealDetailRequest) (*bdspropb.BankAccountItem, int64, error) {
	if req.ID == 0 {
		return nil, 0, _errors.BadRequestException("deal id is required")
	}
	if ts, skip := s.syncUnchanged(ctx, _utils.SyncKeyDealBank, req.ID, req.Timestamp, 0); skip {
		return &bdspropb.BankAccountItem{Timestamp: ts}, ts, nil
	}

	deal, err := s.repo.GetByID(ctx, req.ID)
	if err != nil {
		return nil, 0, err
	}
	if deal == nil {
		return nil, 0, _errors.NotFoundException("deal not found")
	}

	ts := s.syncTS(ctx, _utils.SyncKeyDealBank, req.ID, false)
	if deal.BankAccountID == nil || *deal.BankAccountID == 0 {
		return &bdspropb.BankAccountItem{Timestamp: ts}, ts, nil
	}

	account, err := s.repo.GetBankAccountByID(ctx, *deal.BankAccountID)
	if err != nil {
		return nil, 0, err
	}
	if account == nil {
		return &bdspropb.BankAccountItem{Timestamp: ts}, ts, nil
	}

	return &bdspropb.BankAccountItem{
		Name:          account.BankName,
		BankName:      account.BankName,
		AccountNumber: account.AccountNumber,
		AccountName:   account.AccountName,
		BankId:        account.BankID,
		Timestamp:     ts,
	}, ts, nil
}

func statsToPB(s *dto.TxStats) *bdspropb.StatsDTO {
	if s == nil {
		return nil
	}
	return &bdspropb.StatsDTO{
		Approved: s.Approved,
		Reject:   s.Rejected,
		Pending:  s.Pending,
		Total:    s.Total,
	}
}

func (s *LocalDealService) GetContractStatistics(ctx context.Context, req dto.DealStatisticsRequest) (*bdspropb.ContractDealStatisticsDTO, int64, error) {
	if req.DealID == 0 {
		return nil, 0, _errors.BadRequestException("deal id is required")
	}
	if ts, skip := s.syncUnchanged(ctx, _utils.SyncKeyDealStatistics, req.DealID, req.Timestamp, 0); skip {
		return &bdspropb.ContractDealStatisticsDTO{}, ts, nil
	}

	result := &bdspropb.ContractDealStatisticsDTO{DealId: req.DealID}

	if req.Cost {
		costStats, err := s.repo.GetCostStatisticsByDealID(ctx, req.DealID)
		if err != nil {
			return nil, 0, err
		}
		if costStats != nil {
			result.AmountCost = statsToPB(&costStats.AmountCost)
			result.AmountRevenue = statsToPB(&costStats.AmountRevenue)
			result.NumCost = statsToPB(&costStats.NumCost)
			result.NumRevenue = statsToPB(&costStats.NumRevenue)
		}
	}

	if req.Payment {
		payStats, err := s.repo.GetPaymentStatisticsByDealID(ctx, req.DealID)
		if err != nil {
			return nil, 0, err
		}
		if payStats != nil {
			result.Payment = statsToPB(&payStats.Payment)
			result.NumPayment = statsToPB(&payStats.NumPayment)
		}
	}

	if req.Payment && req.Cost {
		var paymentApproved, revenueApproved, costApproved float64
		if result.Payment != nil {
			paymentApproved = result.Payment.Approved
		}
		if result.AmountRevenue != nil {
			revenueApproved = result.AmountRevenue.Approved
		}
		if result.AmountCost != nil {
			costApproved = result.AmountCost.Approved
		}
		result.TotalAmount = paymentApproved + revenueApproved
		result.TotalCost = costApproved
		result.TotalProfit = result.TotalAmount - result.TotalCost
	}

	ts := s.syncTS(ctx, _utils.SyncKeyDealStatistics, req.DealID, false)
	return result, ts, nil
}

func (s *LocalDealService) profileItemsByMemberIDs(ctx context.Context, members []*entity.DealMember) map[uint64]*sharepb.ProfileItem {
	out := map[uint64]*sharepb.ProfileItem{}
	if s.profile == nil || len(members) == 0 {
		return out
	}
	ids := make([]uint64, 0, len(members))
	seen := map[uint64]struct{}{}
	for _, m := range members {
		if m == nil || m.MemberID == 0 {
			continue
		}
		if _, ok := seen[m.MemberID]; ok {
			continue
		}
		seen[m.MemberID] = struct{}{}
		ids = append(ids, m.MemberID)
	}
	if len(ids) == 0 {
		return out
	}
	profiles, err := s.profile.GetMapProfileByIds(ctx, ids)
	if err != nil {
		return out
	}
	for id, p := range profiles {
		if p == nil {
			continue
		}
		out[id] = &sharepb.ProfileItem{
			Id:       id,
			FullName: p.FullName,
			Avatar:   p.Avatar,
			Phone:    p.Phone,
		}
	}
	return out
}
