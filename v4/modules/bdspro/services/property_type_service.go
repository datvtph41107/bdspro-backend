package services

import (
	"context"
	"time"

	_utils "common/utils"
	"mdn4/modules/bdspro/dto"
	"mdn4/modules/bdspro/entity"
	"mdn4/modules/bdspro/repository"
	bdspropb "pb/types/bdspro"
)

type PropertyTypeService interface {
	List(ctx context.Context, req dto.PropertyTypeListRequest) (*PropertyTypeListOutcome, error)
}

type PropertyTypeListOutcome struct {
	List      *bdspropb.PropertyTypeListResponse
	Timestamp int64
	Unchanged bool
}

// LocalPropertyTypeService owns property type catalog reads and sync metadata.
// @bind: mdn4/modules/bdspro/services.PropertyTypeService
type LocalPropertyTypeService struct {
	repo repository.PropertyTypeRepository
	sync *_utils.SyncUtil
}

func NewLocalPropertyTypeService(repo repository.PropertyTypeRepository, sync *_utils.SyncUtil) *LocalPropertyTypeService {
	return &LocalPropertyTypeService{repo: repo, sync: sync}
}

func (s *LocalPropertyTypeService) List(ctx context.Context, req dto.PropertyTypeListRequest) (*PropertyTypeListOutcome, error) {
	if s.sync != nil {
		catalogTS := s.sync.EnsureTimestamp(ctx, _utils.SyncKeyPropertyTypes)
		if req.Timestamp > 0 && !s.sync.HasUpdated(ctx, _utils.SyncKeyPropertyTypes, req.Timestamp) {
			return &PropertyTypeListOutcome{Timestamp: catalogTS, Unchanged: true}, nil
		}
	}

	rows, err := s.repo.List(ctx, req.Text)
	if err != nil {
		return nil, err
	}
	items := make([]*bdspropb.PropertyType, 0, len(rows))
	for i := range rows {
		items = append(items, propertyTypeEntityToProto(&rows[i]))
	}
	ts := time.Now().UnixMilli()
	if s.sync != nil {
		ts = s.sync.EnsureTimestamp(ctx, _utils.SyncKeyPropertyTypes)
	}
	return &PropertyTypeListOutcome{
		List:      &bdspropb.PropertyTypeListResponse{Data: items, TotalElements: int64(len(items))},
		Timestamp: ts,
	}, nil
}

func (s *LocalPropertyTypeService) BumpSync(ctx context.Context) {
	if s.sync != nil {
		_ = s.sync.PutTimestamp(ctx, _utils.SyncKeyPropertyTypes, time.Now().UnixMilli())
	}
}

func propertyTypeEntityToProto(model *entity.PropertyType) *bdspropb.PropertyType {
	if model == nil {
		return nil
	}
	out := &bdspropb.PropertyType{
		Id:               model.ID,
		Name:             model.Name,
		Active:           model.Active,
		ClassifyProperty: model.ClassifyProperty,
		FieldStrs:        model.FieldStrs,
	}
	if model.CreatedAt != nil {
		out.CreatedAt = _utils.FormatTimeToString(model.CreatedAt)
	}
	if model.UpdatedAt != nil {
		out.UpdatedAt = _utils.FormatTimeToString(model.UpdatedAt)
	}
	return out
}
