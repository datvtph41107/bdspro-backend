package mapper

import (
	"bdspro/internal/dto"
	_dto "common/domain/dto"
	bdspropb "pb/types/bdspro"
)

type PublicMapper struct {
}

func NewPublicMapper() *PublicMapper {
	return &PublicMapper{}
}

func (m *PublicMapper) PbToRegionRequest(region *bdspropb.SearchQueryRequest) *dto.RegionRequest {
	result := &dto.RegionRequest{
		Pagable: _dto.Pagable{
			Page: region.Page,
			Size: region.Size,
		},
		Text: region.Text,
	}

	if region.ParentId != nil {
		result.ParentID = region.ParentId
	}

	return result
}
