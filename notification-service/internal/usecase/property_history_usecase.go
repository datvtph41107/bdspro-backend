package usecase

import (
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
	data *domain.PropertyHistoryEntity,
) (*domain.PropertyHistoryEntity, error) {
	now := time.Now()
	entity := &domain.PropertyHistoryEntity{
		SubjectID:   data.SubjectID,
		SubjectType: data.SubjectType,
		Action:      data.Action,
		Description: data.Description,
		ActorID:     data.ActorID,
		OccurredAt:  &now,
	}

	if entity.OccurredAt == nil {
		now := time.Now()
		entity.OccurredAt = &now
	}

	return uc.Repo.Create(ctx, entity)
}

func (uc *PropertyHistoryUseCase) Search(
	ctx context.Context,
	q dto.ActivityHistoryQueryDTO,
) ([]domain.PropertyHistoryEntity, error) {

	query := repo.PropertyHistoryQuery{
		SubjectID:   q.SubjectID,
		SubjectType: q.SubjectType,
		// Action:      q.Action,
		ActorID: q.ActorID,
		Keyword: q.Keyword,
		Limit:   q.GetLimit(),
		Offset:  q.GetOffset(),
		Sort:    q.Sort,
	}

	return uc.Repo.Search(ctx, query)
}
