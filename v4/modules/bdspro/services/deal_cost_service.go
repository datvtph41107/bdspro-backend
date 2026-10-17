package services

import (
	"context"
	"fmt"
	"time"

	_errors "common/errors"
	_utils "common/utils"
	bdspropb "pb/types/bdspro"

	"userv4/modules/bdspro/dto"
	"userv4/modules/bdspro/entity"
	"userv4/modules/bdspro/enums"
	"userv4/modules/bdspro/repository"
	bdsutils "userv4/modules/bdspro/utils"
)

// DealCostService owns tx_deal_costs (migrated from bdspro DealCostService).
type DealCostService interface {
	Search(ctx context.Context, req dto.DealCostSearchRequest) (*bdspropb.DealCostListDTO, error)
	Detail(ctx context.Context, id uint64) (*bdspropb.DealCostDTO, error)
	Create(ctx context.Context, req *bdspropb.DealCostDTO) (*bdspropb.DealCostDTO, error)
	Update(ctx context.Context, req *bdspropb.DealCostDTO) (*bdspropb.DealCostDTO, error)
	Delete(ctx context.Context, id uint64) error
}

type LocalDealCostService struct {
	sync     *_utils.SyncUtil
	repo     repository.DealCostRepository
	costType repository.CostTypeRepository
	history  *bdsutils.DealHistoryLogger
}

func NewLocalDealCostService(
	sync *_utils.SyncUtil,
	repo repository.DealCostRepository,
	costType repository.CostTypeRepository,
	history *bdsutils.DealHistoryLogger,
) *LocalDealCostService {
	return &LocalDealCostService{sync: sync, repo: repo, costType: costType, history: history}
}

func (s *LocalDealCostService) Search(ctx context.Context, req dto.DealCostSearchRequest) (*bdspropb.DealCostListDTO, error) {
	if req.DealID == 0 {
		return nil, _errors.BadRequestException("deal id is required")
	}
	if ts, skip := s.syncUnchanged(ctx, req.DealID, req.Timestamp, req.Page); skip {
		return &bdspropb.DealCostListDTO{Timestamp: ts}, nil
	}
	rows, err := s.repo.ListByDealID(ctx, req)
	if err != nil {
		return nil, err
	}
	out := &bdspropb.DealCostListDTO{
		Data:          make([]*bdspropb.DealCostDTO, 0, len(rows)),
		TotalElements: uint32(len(rows)),
		Timestamp:     s.syncTS(ctx, req.DealID, false),
	}
	for _, row := range rows {
		out.Data = append(out.Data, dealCostEntityToProto(row))
	}
	return out, nil
}

func (s *LocalDealCostService) Detail(ctx context.Context, id uint64) (*bdspropb.DealCostDTO, error) {
	if id == 0 {
		return nil, _errors.BadRequestException("deal cost id is required")
	}
	row, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if row == nil {
		return nil, _errors.NotFoundException("deal cost not found")
	}
	return dealCostEntityToProto(row), nil
}

func (s *LocalDealCostService) Create(ctx context.Context, req *bdspropb.DealCostDTO) (*bdspropb.DealCostDTO, error) {
	if req == nil || req.GetDealId() == 0 {
		return nil, _errors.BadRequestException("deal id is required")
	}
	row := dealCostProtoToEntity(req)
	row.ID = 0
	if ct, err := s.costType.GetByID(ctx, row.CostTypeID); err != nil {
		return nil, err
	} else if ct != nil && ct.CostType == uint32(enums.TxCostTypeExpense) {
		row.Status = uint32(enums.TxApprovedStatusApproved)
	} else {
		row.Status = uint32(enums.TxApprovedStatusPending)
	}
	if err := s.repo.Create(ctx, row); err != nil {
		return nil, err
	}
	s.bumpSync(ctx, row.DealID)
	s.logDealCostHistory(ctx, row, enums.DealHistoryEventAddDealCost,
		fmt.Sprintf("Thêm chi phí ID %d với số tiền %.2f", row.ID, row.Amount))
	return dealCostEntityToProto(row), nil
}

func (s *LocalDealCostService) Update(ctx context.Context, req *bdspropb.DealCostDTO) (*bdspropb.DealCostDTO, error) {
	if req == nil || req.GetId() == 0 {
		return nil, _errors.BadRequestException("deal cost id is required")
	}
	row := dealCostProtoToEntity(req)
	if err := s.repo.Update(ctx, row); err != nil {
		return nil, err
	}
	s.bumpSync(ctx, row.DealID)
	s.logDealCostHistory(ctx, row, enums.DealHistoryEventUpdateDealCost,
		fmt.Sprintf("Cập nhật chi phí ID %d với số tiền %.2f", row.ID, row.Amount))
	return dealCostEntityToProto(row), nil
}

func (s *LocalDealCostService) Delete(ctx context.Context, id uint64) error {
	if id == 0 {
		return _errors.BadRequestException("deal cost id is required")
	}
	row, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if row == nil {
		return nil
	}
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}
	s.bumpSync(ctx, row.DealID)
	s.logDealCostHistory(ctx, row, enums.DealHistoryEventRemoveDealCost,
		fmt.Sprintf("Xóa chi phí ID %d với số tiền %.2f", row.ID, row.Amount))
	return nil
}

func (s *LocalDealCostService) logDealCostHistory(
	ctx context.Context,
	row *entity.TxDealCost,
	event enums.DealHistoryEventType,
	content string,
) {
	if s.history == nil || row == nil {
		return
	}
	s.history.LogDealHistory(ctx, row.DealID, event, content, map[string]interface{}{
		"deal_id":      row.DealID,
		"deal_cost_id": row.ID,
		"amount":       row.Amount,
		"cost_type_id": row.CostTypeID,
		"status":       row.Status,
	})
}

func (s *LocalDealCostService) bumpSync(ctx context.Context, dealID uint64) {
	if s.sync == nil || dealID == 0 {
		return
	}
	now := time.Now().UnixMilli()
	_ = s.sync.PutTimestamp(ctx, s.sync.GetKey(ctx, _utils.SyncKeyDealCost, dealID), now)
	_ = s.sync.PutTimestamp(ctx, s.sync.GetKey(ctx, _utils.SyncKeyDealStatistics, dealID), now)
}

func (s *LocalDealCostService) syncTS(ctx context.Context, dealID uint64, put bool) int64 {
	if s.sync == nil || dealID == 0 {
		return 0
	}
	key := s.sync.GetKey(ctx, _utils.SyncKeyDealCost, dealID)
	if put {
		_ = s.sync.PutTimestamp(ctx, key, time.Now().UnixMilli())
	}
	return s.sync.EnsureTimestamp(ctx, key)
}

func (s *LocalDealCostService) syncUnchanged(ctx context.Context, dealID uint64, timestamp int64, page uint32) (int64, bool) {
	if s.sync == nil || dealID == 0 || timestamp == 0 || page != 0 {
		return 0, false
	}
	key := s.sync.GetKey(ctx, _utils.SyncKeyDealCost, dealID)
	if s.sync.HasUpdated(ctx, key, timestamp) {
		return 0, false
	}
	return s.sync.GetTimestamp(ctx, key), true
}

func dealCostEntityToProto(row *entity.TxDealCost) *bdspropb.DealCostDTO {
	if row == nil {
		return nil
	}
	out := &bdspropb.DealCostDTO{
		Id:          row.ID,
		DealId:      row.DealID,
		Amount:      row.Amount,
		ConfirmedAt: _utils.FormatTimeToString(row.ConfirmedAt),
		Note:        row.Note,
		CostTypeId:  row.CostTypeID,
		Status:      row.Status,
		CostName:    row.CostName,
	}
	if row.CostType != nil {
		out.CostTypeName = row.CostType.TypeName
	}
	return out
}

func dealCostProtoToEntity(p *bdspropb.DealCostDTO) *entity.TxDealCost {
	row := &entity.TxDealCost{
		DealID:     p.GetDealId(),
		Amount:     p.GetAmount(),
		Note:       p.GetNote(),
		CostTypeID: p.GetCostTypeId(),
		CostName:   p.GetCostName(),
		Status:     p.GetStatus(),
	}
	if p.GetConfirmedAt() != "" {
		row.ConfirmedAt = _utils.ParseStringToTime(p.GetConfirmedAt())
	}
	if p.GetId() != 0 {
		row.ID = p.GetId()
	}
	return row
}
