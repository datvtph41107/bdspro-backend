package mapper

import (
	"bdspro/internal/dto"
	_dto "common/dto"
	pb_bdspro "pb/types/bdspro"
)

type PublicMapper struct {
}

func NewPublicMapper() *PublicMapper {
	return &PublicMapper{}
}

func (m *PublicMapper) PbToRegionRequest(region *pb_bdspro.SearchQueryRequest) *dto.RegionRequest {
	result := &dto.RegionRequest{
		Pagable: _dto.Pagable{
			Page: int(region.Page),
			Size: int(region.Size),
		},
		Text: region.Text,
	}

	if region.ParentId != nil {
		result.ParentID = region.ParentId
	}

	return result
}
