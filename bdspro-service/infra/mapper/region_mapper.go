package mapper

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	_dto "common/domain/dto"
	_utils "common/utils"
	bdspropb "pb/types/bdspro"
)

type RegionMapper struct{}

func NewRegionMapper() *RegionMapper {
	return &RegionMapper{}
}

// PbToRegionListRequest converts protobuf request to DTO
func (m *RegionMapper) PbToRegionListRequest(req *bdspropb.RegionListRequest) *dto.RegionRequest {
	result := &dto.RegionRequest{
		Pagable: _dto.Pagable{
			Page: uint32(req.Page),
			Size: uint32(req.Size),
		},
	}

	if req.Text != nil {
		result.Text = *req.Text
	}

	if req.Level != nil {
		level := int(*req.Level)
		result.Level = &level
	}

	if req.ParentId != nil {
		result.ParentID = req.ParentId
	}

	return result
}

// RegionToProto converts domain Region to protobuf Region
func (m *RegionMapper) RegionToProto(region *domain.Region) *bdspropb.Region {
	result := &bdspropb.Region{
		Id:       region.ID,
		Name:     region.Name,
		Code:     uint32(region.Code),
		CodeName: region.CodeName,
		Level:    uint32(region.Level),
		Unit:     region.Unit,
	}

	if region.ParentID != nil {
		result.ParentId = region.ParentID
	}

	if !region.CreatedAt.IsZero() {
		result.CreatedAt = _utils.FormatTimeToString(region.CreatedAt)
	}

	if !region.UpdatedAt.IsZero() {
		result.UpdatedAt = _utils.FormatTimeToString(region.UpdatedAt)
	}

	return result
}

// RegionListToProto converts list of domain Regions to protobuf Regions
func (m *RegionMapper) RegionListToProto(regions []domain.Region) []*bdspropb.Region {
	result := make([]*bdspropb.Region, 0, len(regions))
	for _, region := range regions {
		result = append(result, m.RegionToProto(&region))
	}
	return result
}
