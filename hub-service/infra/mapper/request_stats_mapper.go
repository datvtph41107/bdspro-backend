package mapper

import (
	"hub/internal/dto"
	hubpb "pb/types/hub"
)

type RequestStatsMapper struct{}

func NewRequestStatsMapper() *RequestStatsMapper {
	return &RequestStatsMapper{}
}

func (m *RequestStatsMapper) MapTrackRequest(req *hubpb.TrackRequestRequest) *dto.TrackRequestRequest {
	if req == nil {
		return nil
	}
	return &dto.TrackRequestRequest{
		Method:         req.GetMethod(),
		Path:           req.GetPath(),
		Service:        req.GetService(),
		StatusCode:     req.GetStatusCode(),
		DurationMs:     req.GetDurationMs(),
		ProfileID:      req.GetProfileId(),
		OrganizationID: req.GetOrganizationId(),
	}
}
