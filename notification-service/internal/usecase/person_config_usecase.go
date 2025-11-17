package usecase

import (
	"context"

	_errors "common/errors"

	"notification/internal/domain"
	"notification/internal/dto"
	"notification/internal/repo"
)

type PersonConfigUsecase struct {
	repo repo.PersonConfigRepo
}

func NewPersonConfigUsecase(repo repo.PersonConfigRepo) *PersonConfigUsecase {
	return &PersonConfigUsecase{repo: repo}
}

func (u *PersonConfigUsecase) BatchUpsert(ctx context.Context, configs []*dto.PersonConfigDTO) ([]*domain.PersonConfigEntity, error) {
	if len(configs) == 0 {
		return []*domain.PersonConfigEntity{}, nil
	}

	entities := make([]*domain.PersonConfigEntity, len(configs))
	for i, cfg := range configs {
		if err := cfg.Validate(); err != nil {
			return nil, err
		}
		entities[i] = &domain.PersonConfigEntity{
			UserID:    cfg.UserID,
			Key:       cfg.Key,
			Checked:   cfg.Checked,
			IsDefault: cfg.IsDefault,
			Channel:   cfg.Channel,
		}
	}

	return u.repo.BatchUpsert(ctx, entities)
}

func (u *PersonConfigUsecase) ListByUserID(ctx context.Context, userID uint64) ([]*domain.PersonConfigEntity, error) {
	if userID == 0 {
		return nil, _errors.ReturnError(400, "userId không hợp lệ")
	}

	return u.repo.ListByUserID(ctx, userID)
}
