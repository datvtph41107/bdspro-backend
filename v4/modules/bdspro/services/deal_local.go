package services

import (
	"context"
	"strconv"
	"time"

	_utils "common/utils"
	bdspropb "pb/types/bdspro"
	sharepb "pb/types/shared"

	"userv4/modules/bdspro/dto"
	"userv4/modules/bdspro/entity"
	"userv4/modules/bdspro/enums"
	"userv4/modules/bdspro/repository"
)

// LocalDealService owns deal user-list against the v4 DSN + SyncUtil.
type LocalDealService struct {
	sync *_utils.SyncUtil
	repo repository.DealRepository
}

func NewLocalDealService(sync *_utils.SyncUtil, repo repository.DealRepository) *LocalDealService {
	return &LocalDealService{sync: sync, repo: repo}
}

func (s *LocalDealService) GetDealsUser(ctx context.Context, req dto.DealSearchRequest) (*bdspropb.GetDealsResponse, error) {
	var syncKey string
	if s.sync != nil {
		syncKey = s.sync.GetKey(ctx, _utils.SyncKeyDealMe, 0)
		if !s.sync.HasUpdated(ctx, syncKey, req.Timestamp) && req.Page == 0 && req.Timestamp > 0 {
			return &bdspropb.GetDealsResponse{
				Timestamp: s.sync.GetTimestamp(ctx, syncKey),
			}, nil
		}
	}

	profileID := _utils.GetProfileIdWithContext(ctx)
	size := req.Size
	if size == 0 {
		size = 20
	}
	list, total, err := s.repo.ListForUser(ctx, profileID, req.Page, size, req.Text)
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
		ts := time.Now().UnixMilli()
		_ = s.sync.PutTimestamp(ctx, syncKey, ts)
		resp.Timestamp = ts
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
