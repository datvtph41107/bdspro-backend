package mapper

import (
	"hub/internal/dto"
	"hub/internal/enums"
	hubpb "pb/types/hub"
)

type ClickURLMapper struct{}

func NewClickURLMapper() *ClickURLMapper {
	return &ClickURLMapper{}
}

func (m *ClickURLMapper) MapTrackClickURLReq(
	req *hubpb.TrackClickURLRequest,
) *dto.TrackClickURLRequest {
	if req == nil {
		return nil
	}

	return &dto.TrackClickURLRequest{
		EventTime:  req.EventTime,
		DomainID:   req.DomainId,
		DomainType: enums.ClickURLDomainType(req.DomainType),
	}
}
