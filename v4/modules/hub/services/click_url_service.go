package services

import (
	"context"

	_errors "common/errors"
	"userv4/modules/hub/dto"
	"userv4/modules/hub/entity"
	"userv4/modules/hub/repository"
)

// ClickURLService tracks link-click events.
type ClickURLService interface {
	TrackBatch(ctx context.Context, req dto.TrackClickURLBatchRequest, profileID *uint64) error
}

type LocalClickURLService struct {
	repo repository.ClickURLRepository
}

func NewLocalClickURLService(repo repository.ClickURLRepository) *LocalClickURLService {
	return &LocalClickURLService{repo: repo}
}

func (s *LocalClickURLService) TrackBatch(ctx context.Context, req dto.TrackClickURLBatchRequest, profileID *uint64) error {
	if len(req.Events) == 0 {
		return _errors.BadRequestException("events list is empty")
	}
	evs := make([]*entity.ClickURL, 0, len(req.Events))
	for _, e := range req.Events {
		evs = append(evs, &entity.ClickURL{
			EventTime:  e.EventTime,
			DomainID:   e.DomainID,
			DomainType: e.DomainType,
			ProfileID:  profileID,
		})
	}
	return s.repo.InsertBatch(ctx, evs)
}
