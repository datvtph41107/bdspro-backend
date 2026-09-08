package mapper

import (
	_models "common/domain/entity"
	"bdspro/internal/domain"
	bdspropb "pb/types/bdspro"
)

type AreaRegionMapper struct{}

func NewAreaRegionMapper() *AreaRegionMapper {
	return &AreaRegionMapper{}
}

func (m *AreaRegionMapper) ToProto(entity *domain.AreaRegion) *bdspropb.AreaRegion {
	if entity == nil {
		return nil
	}
	return &bdspropb.AreaRegion{
		Id:     entity.ID,
		Name:   entity.Name,
		Code:   entity.Code,
		Active: entity.Active,
	}
}

func (m *AreaRegionMapper) ToProtoList(entities []domain.AreaRegion) []*bdspropb.AreaRegion {
	if entities == nil {
		return []*bdspropb.AreaRegion{}
	}
	result := make([]*bdspropb.AreaRegion, 0, len(entities))
	for i := range entities {
		result = append(result, m.ToProto(&entities[i]))
	}
	return result
}

func (m *AreaRegionMapper) FromCreateRequest(req *bdspropb.CreateAreaRegionRequest) *domain.AreaRegion {
	if req == nil {
		return nil
	}
	return &domain.AreaRegion{
		Name:   req.Name,
		Code:   req.Code,
		Active: req.Active,
	}
}

func (m *AreaRegionMapper) FromUpdateRequest(req *bdspropb.UpdateAreaRegionRequest) *domain.AreaRegion {
	if req == nil {
		return nil
	}
	return &domain.AreaRegion{
		BaseEntity: _models.BaseEntity{ID: req.Id},
		Name:       req.Name,
		Code:       req.Code,
		Active:     req.Active,
	}
}
