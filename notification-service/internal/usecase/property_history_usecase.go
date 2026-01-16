package usecase

import (
	_errors "common/errors"
	"context"
	"notification/internal/domain"
	"notification/internal/dto"
	"notification/internal/repo"
	"time"
)

type PropertyHistoryUseCase struct {
	Repo repo.PropertyHistoryRepo
}

func NewPropertyHistoryUseCase(repo repo.PropertyHistoryRepo) *PropertyHistoryUseCase {
	return &PropertyHistoryUseCase{
		Repo: repo,
	}
}

func (uc *PropertyHistoryUseCase) Create(
	ctx context.Context,
	dto *dto.CreatePropertyHistoryDTO,
) (*domain.PropertyHistoryEntity, error) {
	now := time.Now()
	entity := &domain.PropertyHistoryEntity{
		SubjectID: dto.SubjectID,
		// SubjectType: dto.SubjectType,
		Action:      dto.Action,
		Description: dto.Description,
		Metadata:    dto.Metadata,
		ActorID:     dto.ActorID,
		OccurredAt:  &now,
	}

	return uc.Repo.Create(ctx, entity)
}

func (uc *PropertyHistoryUseCase) GetDetail(
	ctx context.Context,
	id uint64,
) (*domain.PropertyHistoryEntity, error) {

	history, err := uc.Repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if history == nil {
		return nil, _errors.NotFoundException("History not found")
	}

	return history, nil
}

func (uc *PropertyHistoryUseCase) SoftDelete(
	ctx context.Context,
	id uint64,
	actorID uint64,
) error {

	history, err := uc.Repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if history == nil {
		return _errors.NotFoundException("History not found")
	}

	return uc.Repo.Delete(ctx, id)
}

func (uc *PropertyHistoryUseCase) Search(
	ctx context.Context,
	q dto.PropertyHistoryQueryDTO,
) ([]domain.PropertyHistoryEntity, int64, error) {
	query := dto.PropertyHistoryQueryDTO{
		SubjectID:   q.SubjectID,
		SubjectType: q.SubjectType,
		Action:      q.Action,
		ActorID:     q.ActorID,
		Keyword:     q.Keyword,
		Limit:       q.GetLimit(),
		Offset:      q.GetOffset(),
	}

	return uc.Repo.Search(ctx, query)
}
