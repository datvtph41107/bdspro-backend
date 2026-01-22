package mapper

import (
	"bdspro/internal/domain"
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
		EntityType: req.EntityType,
		EntityID:   req.EntityId,
		Action:     req.Action,
		Duration:   req.Duration,
	}
}

func (m *ObservationMapper) MapTrackObservationDTOToDomain(
	dtoReq *dto.TrackObservationRequest,
	profileID *uint64,
	clientID *string,
	sessionID *string,
	role *string,
	orgID *uint64,
	pageURL *string,
	referrer *string,
	device *string,
	userAgent *string,
	geo *string,
) *domain.Observation {

	if dtoReq == nil {
		return nil
	}

	entityType := dtoReq.EntityType

	return &domain.Observation{
		EntityType: &entityType,
		EntityID:   dtoReq.EntityID,
		ProfileID:  profileID,
		ClientID:   clientID,
		SessionID:  sessionID,
		Role:       role,
		OrgID:      orgID,
		Action:     dtoReq.Action,
		Duration:   dtoReq.Duration,
		PageURL:    pageURL,
		Referrer:   referrer,
		Device:     device,
		UserAgent:  userAgent,
		Geo:        geo,
	}
}

func (m *ObservationMapper) MapProductStatsToPB(
	stats *dto.ProductStats,
) *bdspropb.ProductStatsResponse {

	if stats == nil {
		return nil
	}

	return &bdspropb.ProductStatsResponse{
		ProductId:     stats.ProductID,
		ViewCount:     stats.ViewCount,
		TotalDuration: stats.TotalDuration,
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
