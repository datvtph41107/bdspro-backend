package mapper

import (
	"hub/internal/domain"
	"hub/internal/dto"
	hubpb "pb/types/hub"
)

type InteractiveEventMapper struct{}

func NewInteractiveEventMapper() *InteractiveEventMapper {
	return &InteractiveEventMapper{}
}

func (m *InteractiveEventMapper) MapTrackInteractiveEventReq(
	req *hubpb.TrackInteractiveEventRequest,
) *dto.TrackInteractiveEventRequest {
	if req == nil {
		return nil
	}

	var endTime *int64
	if req.EndTime != 0 {
		endTime = &req.EndTime
	}

	var screen *string
	if req.Screen != "" {
		screen = &req.Screen
	}

	var duration *int64
	if req.Duration != 0 {
		duration = &req.Duration
	}

	return &dto.TrackInteractiveEventRequest{
		StartTime: req.StartTime,
		EndTime:   endTime,
		Event:     req.Event,
		RefID:     req.RefId,
		Screen:    screen,
		Duration:  duration,
	}
}

func (m *InteractiveEventMapper) MapTrackInteractiveEventDTOToDomain(
	dtoReq *dto.TrackInteractiveEventRequest,
	profileID *uint64,
) *domain.InteractiveEvent {
	if dtoReq == nil {
		return nil
	}

	return &domain.InteractiveEvent{
		StartTime: dtoReq.StartTime,
		EndTime:   dtoReq.EndTime,
		Event:     dtoReq.Event,
		RefID:     dtoReq.RefID,
		Screen:    dtoReq.Screen,
		Duration:  dtoReq.Duration,
		ProfileID: profileID,
	}
}
