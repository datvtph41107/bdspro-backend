package services

import (
	"context"
	"time"

	_errors "common/errors"
	_utils "common/utils"
	bdspropb "pb/types/bdspro"

	"mdn4/modules/bdspro/dto"
	"mdn4/modules/bdspro/entity"
	"mdn4/modules/bdspro/repository"
)

// CostTypeListOutcome is the cost-type catalog response + sync metadata.
type CostTypeListOutcome struct {
	List      *bdspropb.CostTypeListDTO
	Timestamp int64
	Unchanged bool
}

// CostTypeService owns tx_cost_types CRUD (migrated from bdspro CostTypeService).
type CostTypeService interface {
	Create(ctx context.Context, req *bdspropb.CostTypeDTO) (*bdspropb.CostTypeDTO, error)
	Update(ctx context.Context, req *bdspropb.CostTypeDTO) (*bdspropb.CostTypeDTO, error)
	Delete(ctx context.Context, id uint64) error
	List(ctx context.Context, req dto.CostTypeSearchRequest) (*CostTypeListOutcome, error)
}

// @bind: mdn4/modules/bdspro/services.CostTypeService
type LocalCostTypeService struct {
	repo repository.CostTypeRepository
	sync *_utils.SyncUtil
}

func NewLocalCostTypeService(repo repository.CostTypeRepository, sync *_utils.SyncUtil) *LocalCostTypeService {
	return &LocalCostTypeService{repo: repo, sync: sync}
}

func (s *LocalCostTypeService) Create(ctx context.Context, req *bdspropb.CostTypeDTO) (*bdspropb.CostTypeDTO, error) {
	if req == nil {
		return nil, _errors.BadRequestException("cost type body is required")
	}
	model := costTypeProtoToEntity(req)
	if err := s.repo.Create(ctx, model); err != nil {
		return nil, err
	}
	s.bumpCostTypesSync(ctx)
	return costTypeEntityToProto(model), nil
}

func (s *LocalCostTypeService) Update(ctx context.Context, req *bdspropb.CostTypeDTO) (*bdspropb.CostTypeDTO, error) {
	if req == nil {
		return nil, _errors.BadRequestException("cost type body is required")
	}
	if req.GetId() == 0 {
		return nil, _errors.BadRequestException("cost type id is required")
	}
	existing, err := s.repo.GetByID(ctx, req.GetId())
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, _errors.NotFoundException("cost type not found")
	}
	model := costTypeProtoToEntity(req)
	if err := s.repo.Update(ctx, model); err != nil {
		return nil, err
	}
	s.bumpCostTypesSync(ctx)
	return costTypeEntityToProto(model), nil
}

func (s *LocalCostTypeService) Delete(ctx context.Context, id uint64) error {
	if id == 0 {
		return _errors.BadRequestException("cost type id is required")
	}
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}
	s.bumpCostTypesSync(ctx)
	return nil
}

func (s *LocalCostTypeService) List(ctx context.Context, req dto.CostTypeSearchRequest) (*CostTypeListOutcome, error) {
	if s.sync != nil {
		catalogTS := s.sync.EnsureTimestamp(ctx, _utils.SyncKeyTxCostTypes)
		if req.Timestamp > 0 && !s.sync.HasUpdated(ctx, _utils.SyncKeyTxCostTypes, req.Timestamp) {
			return &CostTypeListOutcome{Timestamp: catalogTS, Unchanged: true}, nil
		}
	}

	rows, err := s.repo.List(ctx, req)
	if err != nil {
		return nil, err
	}
	list := &bdspropb.CostTypeListDTO{
		Data:  make([]*bdspropb.CostTypeDTO, 0, len(rows)),
		Total: uint32(len(rows)),
	}
	for _, row := range rows {
		list.Data = append(list.Data, costTypeEntityToProto(row))
	}

	ts := time.Now().UnixMilli()
	if s.sync != nil {
		ts = s.sync.EnsureTimestamp(ctx, _utils.SyncKeyTxCostTypes)
	}
	return &CostTypeListOutcome{List: list, Timestamp: ts}, nil
}

func (s *LocalCostTypeService) bumpCostTypesSync(ctx context.Context) {
	if s.sync == nil {
		return
	}
	_ = s.sync.PutTimestamp(ctx, _utils.SyncKeyTxCostTypes, time.Now().UnixMilli())
}

func costTypeEntityToProto(model *entity.TxCostType) *bdspropb.CostTypeDTO {
	if model == nil {
		return nil
	}
	return &bdspropb.CostTypeDTO{
		Id:          model.ID,
		TypeName:    model.TypeName,
		IsCustom:    model.IsCustom,
		CostType:    model.CostType,
		Description: model.Description,
	}
}

func costTypeProtoToEntity(proto *bdspropb.CostTypeDTO) *entity.TxCostType {
	return &entity.TxCostType{
		ID:          proto.GetId(),
		TypeName:    proto.GetTypeName(),
		IsCustom:    proto.GetIsCustom(),
		CostType:    proto.GetCostType(),
		Description: proto.GetDescription(),
	}
}
