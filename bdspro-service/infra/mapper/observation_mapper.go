package mapper

import (
	"bdspro/internal/dto"
	bdspropb "pb/types/bdspro"
)

type ObservationMapper struct{}

func NewObservationMapper() *ObservationMapper {
	return &ObservationMapper{}
}

func (m *ObservationMapper) MapTrackObservationReq(
	req *bdspropb.TrackObservationRequest,
) *dto.TrackObservationRequest {

	if req == nil {
		return nil
	}

	return &dto.TrackObservationRequest{
		Events: req.Events,
	}
}

func (m *ObservationMapper) MapProductStatsToPB(
	ps *dto.ProductStats,
) *bdspropb.ProductStatsResponse {
	if ps == nil {
		return nil
	}

	return &bdspropb.ProductStatsResponse{
		ProductId:     ps.ProductID,
		ViewCount:     ps.ViewCount,
		TotalDuration: ps.TotalDuration,
	}
}

func (m *ObservationMapper) MapProductStatsRowToDTO(
	productID uint64,
	viewCount uint64,
	totalDuration uint64,
) *dto.ProductStats {

	return &dto.ProductStats{
		ProductID:     productID,
		ViewCount:     viewCount,
		TotalDuration: totalDuration,
	}
}

func (m *ObservationMapper) MapProductStatsListToPB(
	stats []*dto.ProductStats,
) []*bdspropb.ProductStatsResponse {

	result := make([]*bdspropb.ProductStatsResponse, 0)
	for _, s := range stats {
		result = append(result, m.MapProductStatsToPB(s))
	}
	return result
}
