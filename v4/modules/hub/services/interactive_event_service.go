package services

import (
	"context"

	_errors "common/errors"
	"mdn4/modules/hub/dto"
	"mdn4/modules/hub/entity"
	"mdn4/modules/hub/repository"
)

// InteractiveEventService tracks user interaction events.
type InteractiveEventService interface {
	TrackBatch(ctx context.Context, req dto.TrackInteractiveEventBatchRequest, profileID *uint64) error
}

// @bind: mdn4/modules/hub/services.InteractiveEventService
type LocalInteractiveEventService struct {
	repo repository.InteractiveEventRepository
}

func NewLocalInteractiveEventService(repo repository.InteractiveEventRepository) *LocalInteractiveEventService {
	return &LocalInteractiveEventService{repo: repo}
}

func (s *LocalInteractiveEventService) TrackBatch(ctx context.Context, req dto.TrackInteractiveEventBatchRequest, profileID *uint64) error {
	if len(req.Events) == 0 {
		return _errors.BadRequestException("events list is empty")
	}
	evs := make([]*entity.InteractiveEvent, 0, len(req.Events))
	for _, e := range req.Events {
		evs = append(evs, &entity.InteractiveEvent{
			StartTime: e.StartTime, EndTime: e.EndTime, Event: e.Event,
			RefID: e.RefID, Screen: e.Screen, Duration: e.Duration, ProfileID: profileID,
		})
	}
	return s.repo.InsertBatch(ctx, evs)
}
