package usecase

import (
	"context"
	"notification/internal/domain"
	"notification/internal/dto"
	"notification/internal/repo"
	"time"

	"gorm.io/datatypes"
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
		Metadata:    datatypes.JSON(data.Metadata),
		Description: data.Description,
		ActorID:     data.ActorID,
		ActorType:   data.ActorType,
		OccurredAt:  &now,
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
		Action:      q.Action,
		ActorID:     q.ActorID,
		ActorType:   q.ActorType,
		Keyword:     q.Keyword,
		Limit:       q.GetLimit(),
		Offset:      q.GetOffset(),
		Sort:        q.Sort,
	}

	return uc.Repo.Search(ctx, query)
}
