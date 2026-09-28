package services

import (
	"context"
	"fmt"
	"time"

	_errors "common/errors"
	_utils "common/utils"
	bdspropb "pb/types/bdspro"
	sharepb "pb/types/shared"

	"userv4/modules/bdspro/dto"
	"userv4/modules/bdspro/enums"
)

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

func (s *LocalDealService) GetDealProducts(ctx context.Context, req dto.DealListByIDRequest) (*bdspropb.GetDealProductsResponse, int64, error) {
	if req.DealID == 0 {
		return nil, 0, _errors.BadRequestException("deal id is required")
	}
	if ts, skip := s.syncUnchanged(ctx, _utils.SyncKeyDealProducts, req.DealID, req.Timestamp, req.Page); skip {
		return &bdspropb.GetDealProductsResponse{}, ts, nil
	}

	ids, err := s.repo.ListDealProductIDs(ctx, req.DealID, req.Page, req.Size)
	if err != nil {
		return nil, 0, err
	}
	ts := s.syncTS(ctx, _utils.SyncKeyDealProducts, req.DealID, false)
	if len(ids) == 0 || s.products == nil {
		return &bdspropb.GetDealProductsResponse{Data: []*bdspropb.ProductAttachment{}}, ts, nil
	}

	products, err := s.products.ListByIDs(ctx, ids)
	if err != nil {
		return nil, 0, err
	}
	byID := make(map[uint64]*bdspropb.ProductAttachment, len(products))
	for _, p := range products {
		if p == nil {
			continue
		}
		att := &bdspropb.ProductAttachment{
			Id:             p.ID,
			Name:           p.Name,
			Code:           p.Code,
			Area:           p.Area,
			RentStatus:     p.RentStatus,
			SaleStatus:     p.SaleStatus,
			SourceType:     int32(p.SourceType),
			SourceTypeName: enums.SourceTypeMap[enums.ESourceType(p.SourceType)],
		}
		if p.AssetID != nil {
			att.AssetId = p.AssetID
		}
		if p.Price != nil && p.Price.SalePrice != nil {
			att.Price = *p.Price.SalePrice
		}
		byID[p.ID] = att
	}
	data := make([]*bdspropb.ProductAttachment, 0, len(ids))
	for _, id := range ids {
		if att, ok := byID[id]; ok {
			data = append(data, att)
		}
	}
	return &bdspropb.GetDealProductsResponse{Data: data}, ts, nil
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
	data := make([]*sharepb.DealMember, 0, len(list))
	for _, m := range list {
		if m == nil {
			continue
		}
		item := &sharepb.DealMember{
			Id:           m.ID,
			DealId:       m.DealID,
			InviterId:    m.InviterID,
			MemberId:     m.MemberID,
			Status:       m.Status,
			Message:      m.Message,
			RoleKey:      m.RoleKey,
			IsOwner:      m.IsOwner,
			AmountCommit: m.AmountCommit,
			RoleId:       &m.RoleID,
		}
		if m.InvitedAt != nil {
			item.InvitedAt = _utils.FormatTimeToString(m.InvitedAt)
		}
		if m.RespondedAt != nil {
			item.RespondedAt = _utils.FormatTimeToString(m.RespondedAt)
		}
		data = append(data, item)
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

	profileIDs := make([]uint64, 0, len(list))
	seen := map[uint64]struct{}{}
	for _, m := range list {
		if m == nil || m.MemberID == 0 {
			continue
		}
		if _, ok := seen[m.MemberID]; ok {
			continue
		}
		seen[m.MemberID] = struct{}{}
		profileIDs = append(profileIDs, m.MemberID)
	}

	profilesMap := map[uint64]*sharepb.ProfileItem{}
	if s.profile != nil && len(profileIDs) > 0 {
		if m, err := s.profile.GetMapProfileByIds(ctx, profileIDs); err == nil {
			for id, p := range m {
				if p == nil {
					continue
				}
				profilesMap[id] = &sharepb.ProfileItem{
					Id:       id,
					FullName: p.FullName,
					Avatar:   p.Avatar,
				}
			}
		}
	}

	data := make([]*sharepb.DealMember, 0, len(list))
	for _, m := range list {
		if m == nil {
			continue
		}
		item := &sharepb.DealMember{
			Id:              m.ID,
			DealId:          m.DealID,
			MemberId:        m.MemberID,
			AmountCommit:    m.AmountCommit,
			CommissionValue: m.CommissionValue,
			CommissionType:  m.CommissionType,
			Note:            m.Note,
			RoleId:          &m.RoleID,
			IsOwner:         m.IsOwner,
			DoneInvestment:  m.DoneInvestment,
			RoleKey:         m.RoleKey,
			Status:          m.Status,
			Profile:         profilesMap[m.MemberID],
		}
		data = append(data, item)
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
