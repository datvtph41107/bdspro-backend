package services

import (
	"context"
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

// InvestmentService owns deal investment flows (migrated from /v2.1/org/investment).
type InvestmentService interface {
	Create(ctx context.Context, req *bdspropb.SubmitInvestRequest) (uint64, error)
	Update(ctx context.Context, req *bdspropb.SubmitInvestRequest) (uint64, error)
	Delete(ctx context.Context, id uint64) error
	Detail(ctx context.Context, id uint64) (*bdspropb.Investment, error)
	List(ctx context.Context, q dto.InvestmentListQuery) ([]*entity.DealInvestment, int64, int64, error)
	Approve(ctx context.Context, id uint64) (uint64, error)
	Reject(ctx context.Context, id uint64) (uint64, error)
	Confirm(ctx context.Context, id uint64, confirmed bool) (uint64, error)
	MemberHistories(ctx context.Context, dealID, memberID uint64) (*bdspropb.GetInvestmentHistoriesResponse, error)
}

type LocalInvestmentService struct {
	sync     *_utils.SyncUtil
	repo     repository.InvestmentRepository
	dealRepo repository.DealRepository
	profile  providers.ProfileProvider
}

func NewLocalInvestmentService(
	sync *_utils.SyncUtil,
	repo repository.InvestmentRepository,
	dealRepo repository.DealRepository,
	profile providers.ProfileProvider,
) *LocalInvestmentService {
	return &LocalInvestmentService{sync: sync, repo: repo, dealRepo: dealRepo, profile: profile}
}

func (s *LocalInvestmentService) Create(ctx context.Context, req *bdspropb.SubmitInvestRequest) (uint64, error) {
	if req == nil || req.GetDealId() == 0 {
		return 0, _errors.BadRequestException("deal id is required")
	}
	userID := _utils.GetProfileIdWithContext(ctx)
	if userID == 0 {
		return 0, _errors.UnauthorizedException("user not authenticated")
	}
	if err := s.ensureMember(ctx, req.GetDealId(), userID); err != nil {
		return 0, err
	}
	if _, err := s.dealRepo.GetByID(ctx, req.GetDealId()); err != nil {
		return 0, err
	}
	row := submitInvestToEntity(req)
	row.ID = 0
	if row.MemberID == 0 {
		row.MemberID = userID
	}
	if row.Status == uint32(enums.TxApprovedStatusRejected) {
		row.Status = uint32(enums.TxApprovedStatusPending)
	}
	if userID > 0 {
		row.CreatedBy = &userID
	}
	if err := s.repo.Create(ctx, row); err != nil {
		return 0, err
	}
	if req.GetDoneInvestment() {
		_ = s.dealRepo.UpdateDoneInvestment(ctx, row.DealID, row.MemberID)
	}
	s.bumpSync(ctx, row.DealID)
	return row.ID, nil
}

func (s *LocalInvestmentService) Update(ctx context.Context, req *bdspropb.SubmitInvestRequest) (uint64, error) {
	if req == nil || req.GetId() == 0 {
		return 0, _errors.BadRequestException("investment id is required")
	}
	userID := _utils.GetProfileIdWithContext(ctx)
	if userID == 0 {
		return 0, _errors.UnauthorizedException("user not authenticated")
	}
	if err := s.ensureMember(ctx, req.GetDealId(), userID); err != nil {
		return 0, err
	}
	current, err := s.repo.GetByID(ctx, req.GetId())
	if err != nil {
		return 0, err
	}
	if current == nil {
		return 0, _errors.NotFoundException("investment not found")
	}
	if !s.canManage(current, userID) {
		return 0, _errors.ForbiddenException("permission denied")
	}
	row := submitInvestToEntity(req)
	row.ID = req.GetId()
	if err := s.repo.Update(ctx, row); err != nil {
		return 0, err
	}
	if req.GetDoneInvestment() {
		_ = s.dealRepo.UpdateDoneInvestment(ctx, row.DealID, row.MemberID)
	}
	s.bumpSync(ctx, row.DealID)
	return row.ID, nil
}

func (s *LocalInvestmentService) Delete(ctx context.Context, id uint64) error {
	if id == 0 {
		return _errors.BadRequestException("investment id is required")
	}
	userID := _utils.GetProfileIdWithContext(ctx)
	current, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if current == nil {
		return nil
	}
	if userID > 0 && !s.canManage(current, userID) {
		return _errors.ForbiddenException("permission denied")
	}
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}
	s.bumpSync(ctx, current.DealID)
	return nil
}

func (s *LocalInvestmentService) Detail(ctx context.Context, id uint64) (*bdspropb.Investment, error) {
	if id == 0 {
		return nil, _errors.BadRequestException("investment id is required")
	}
	row, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if row == nil {
		return nil, _errors.NotFoundException("investment not found")
	}
	histories, err := s.repo.ListConfirmedByDealMember(ctx, row.DealID, row.MemberID)
	if err != nil {
		return nil, err
	}
	var total float64
	for _, h := range histories {
		if h != nil {
			total += h.Amount
		}
	}
	out := entityToInvestmentPb(row)
	out.Histories = investmentsToHistoryPb(histories)
	out.TotalInvestment = total
	s.attachProfiles(ctx, []*bdspropb.Investment{out})
	return out, nil
}

func (s *LocalInvestmentService) List(ctx context.Context, q dto.InvestmentListQuery) ([]*entity.DealInvestment, int64, int64, error) {
	if q.DealID == 0 {
		return nil, 0, 0, _errors.BadRequestException("dealId is required")
	}
	dealID := q.DealID
	if s.sync != nil && q.Timestamp > 0 {
		key := s.sync.GetKey(ctx, _utils.SyncKeyDealInvestment, dealID)
		if !s.sync.HasUpdated(ctx, key, q.Timestamp) {
			ts := s.sync.GetTimestamp(ctx, key)
			if ts <= 0 {
				ts = q.Timestamp
			}
			return nil, 0, ts, nil
		}
	}
	rows, total, err := s.repo.List(ctx, q)
	if err != nil {
		return nil, 0, 0, err
	}

	// s.attachProfiles(ctx, rows)
	ts := s.investmentSyncTS(ctx, dealID)
	return rows, total, ts, nil
}

func (s *LocalInvestmentService) investmentSyncTS(ctx context.Context, dealID uint64) int64 {
	if s.sync == nil || dealID == 0 {
		return time.Now().UnixMilli()
	}
	key := s.sync.GetKey(ctx, _utils.SyncKeyDealInvestment, dealID)
	ts := s.sync.EnsureTimestamp(ctx, key)
	if ts <= 0 {
		return time.Now().UnixMilli()
	}
	return ts
}

func (s *LocalInvestmentService) Approve(ctx context.Context, id uint64) (uint64, error) {
	row, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return 0, err
	}
	if row == nil {
		return 0, _errors.NotFoundException("investment not found")
	}
	if err := s.repo.ChangeStatus(ctx, id, enums.TxApprovedStatusApproved); err != nil {
		return 0, err
	}
	s.bumpSync(ctx, row.DealID)
	return id, nil
}

func (s *LocalInvestmentService) Reject(ctx context.Context, id uint64) (uint64, error) {
	row, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return 0, err
	}
	if row == nil {
		return 0, _errors.NotFoundException("investment not found")
	}
	if err := s.repo.ChangeStatus(ctx, id, enums.TxApprovedStatusRejected); err != nil {
		return 0, err
	}
	s.bumpSync(ctx, row.DealID)
	return id, nil
}

func (s *LocalInvestmentService) Confirm(ctx context.Context, id uint64, confirmed bool) (uint64, error) {
	row, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return 0, err
	}
	if row == nil {
		return 0, _errors.NotFoundException("investment not found")
	}
	if err := s.repo.ChangeConfirmation(ctx, id, confirmed); err != nil {
		return 0, err
	}
	s.bumpSync(ctx, row.DealID)
	return id, nil
}

func (s *LocalInvestmentService) MemberHistories(ctx context.Context, dealID, memberID uint64) (*bdspropb.GetInvestmentHistoriesResponse, error) {
	if dealID == 0 || memberID == 0 {
		return nil, _errors.BadRequestException("deal id and member id are required")
	}
	rows, err := s.repo.ListConfirmedByDealMember(ctx, dealID, memberID)
	if err != nil {
		return nil, err
	}
	var total float64
	for _, row := range rows {
		if row != nil {
			total += row.Amount
		}
	}
	return &bdspropb.GetInvestmentHistoriesResponse{
		Data:            investmentsToHistoryPb(rows),
		TotalInvestment: total,
	}, nil
}

func (s *LocalInvestmentService) ensureMember(ctx context.Context, dealID, userID uint64) error {
	ok, err := s.dealRepo.IsAcceptedMember(ctx, dealID, userID)
	if err != nil {
		return err
	}
	if !ok {
		return _errors.ForbiddenException("you are not allowed to access this deal")
	}
	return nil
}

func (s *LocalInvestmentService) canManage(row *entity.DealInvestment, userID uint64) bool {
	if row == nil {
		return false
	}
	if row.Status == uint32(enums.TxApprovedStatusApproved) {
		return false
	}
	return row.MemberID == userID
}

func (s *LocalInvestmentService) bumpSync(ctx context.Context, dealID uint64) {
	if s.sync == nil || dealID == 0 {
		return
	}
	now := time.Now().UnixMilli()
	_ = s.sync.PutTimestamp(ctx, s.sync.GetKey(ctx, _utils.SyncKeyDealInvestment, dealID), now)
	_ = s.sync.PutTimestamp(ctx, s.sync.GetKey(ctx, _utils.SyncKeyDealSummary, dealID), now)
}

func (s *LocalInvestmentService) attachProfiles(ctx context.Context, items []*bdspropb.Investment) {
	if s.profile == nil || len(items) == 0 {
		return
	}
	ids := make([]uint64, 0, len(items)*2)
	for _, item := range items {
		if item == nil {
			continue
		}
		if item.CreatedBy != nil {
			ids = append(ids, *item.CreatedBy)
		}
		if item.ProxyUserId != nil {
			ids = append(ids, *item.ProxyUserId)
		}
	}
	if len(ids) == 0 {
		return
	}
	profiles, err := s.profile.GetMapProfileByIds(ctx, ids)
	if err != nil || len(profiles) == 0 {
		return
	}
	for _, item := range items {
		if item == nil {
			continue
		}
		if item.CreatedBy != nil {
			if p, ok := profiles[*item.CreatedBy]; ok && p != nil {
				item.CreatedUser = &sharepb.ProfileItem{
					Id:       *item.CreatedBy,
					FullName: p.FullName,
					Avatar:   p.Avatar,
				}
			}
		}
		if item.ProxyUserId != nil {
			if p, ok := profiles[*item.ProxyUserId]; ok && p != nil {
				item.ProxyUser = &sharepb.ProfileItem{
					Id:       *item.ProxyUserId,
					FullName: p.FullName,
					Avatar:   p.Avatar,
				}
			}
		}
	}
}

func submitInvestToEntity(req *bdspropb.SubmitInvestRequest) *entity.DealInvestment {
	row := &entity.DealInvestment{
		DealID:             req.GetDealId(),
		MemberID:           req.GetMemberId(),
		InvestType:         req.GetInvestType(),
		Amount:             req.GetAmount(),
		Note:               req.GetNote(),
		TransferProofImage: req.GetTransferProofImage(),
		Status:             req.GetStatus(),
	}
	if req.GetTransferTime() != "" {
		row.TransferTime = _utils.ParseStringToTime(req.GetTransferTime())
	}
	if req.ProxyUserId != nil {
		row.ProxyUserID = req.ProxyUserId
	}
	return row
}

func entityToInvestmentPb(row *entity.DealInvestment) *bdspropb.Investment {
	if row == nil {
		return nil
	}
	changed := row.ChangedAt
	if changed == nil {
		changed = row.TransferTime
	}
	out := &bdspropb.Investment{
		Id:                 row.ID,
		DealId:             row.DealID,
		MemberId:           row.MemberID,
		InvestType:         row.InvestType,
		Amount:             row.Amount,
		TransferTime:       _utils.FormatTimeToString(row.TransferTime),
		Note:               row.Note,
		TransferProofImage: row.TransferProofImage,
		Status:             row.Status,
		StatusName:         enums.TxApprovedStatusMap[enums.TxApprovedStatus(row.Status)],
		ChangedAt:          _utils.FormatTimeToString(changed),
		Confirmed:          row.Confirmed,
	}
	if row.CreatedBy != nil {
		out.CreatedBy = row.CreatedBy
	}
	if row.ProxyUserID != nil {
		out.ProxyUserId = row.ProxyUserID
	}
	return out
}

func investmentsToHistoryPb(rows []*entity.DealInvestment) []*bdspropb.ApprovedInvestmentItem {
	out := make([]*bdspropb.ApprovedInvestmentItem, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		out = append(out, &bdspropb.ApprovedInvestmentItem{
			Id:           row.ID,
			Amount:       row.Amount,
			TransferTime: _utils.FormatTimeToString(row.TransferTime),
			Note:         row.Note,
			StatusName:   enums.TxApprovedStatusMap[enums.TxApprovedStatus(row.Status)],
			Confirmed:    row.Confirmed,
		})
	}
	return out
}
