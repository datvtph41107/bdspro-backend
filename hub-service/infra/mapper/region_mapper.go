package mapper

import (
	_models "common/domain/entity"
	"hub/internal/domain"
	hubpb "pb/types/hub"
)

// RegionMapper translates between domain and protobuf for region
// simple binary mapping.
type RegionMapper struct{}

func NewRegionMapper() *RegionMapper {
	return &RegionMapper{}
}

func (m *RegionMapper) ToProto(entity *domain.Region) *hubpb.Region {
	if entity == nil {
		return nil
	}
	return &hubpb.Region{
		Id:     entity.ID,
		Name:   entity.Name,
		Code:   entity.Code,
		Active: entity.Active,
	}
}

func (m *RegionMapper) ToProtoList(entities []domain.Region) []*hubpb.Region {
	if entities == nil {
		return []*hubpb.Region{}
	}
	result := make([]*hubpb.Region, 0, len(entities))
	for _, e := range entities {
		result = append(result, m.ToProto(&e))
	}
	return result
}

func (m *RegionMapper) FromCreateRequest(req *hubpb.CreateRegionRequest) *domain.Region {
	if req == nil {
		return nil
	}
	return &domain.Region{
		Name:   req.Name,
		Code:   req.Code,
		Active: req.Active,
	}
}

func (m *RegionMapper) FromUpdateRequest(req *hubpb.UpdateRegionRequest) *domain.Region {
	if req == nil {
		return nil
	}
	return &domain.Region{
		BaseEntity: _models.BaseEntity{
			ID: req.Id,
		},
		Name:   req.Name,
		Code:   req.Code,
		Active: req.Active,
	}
}
