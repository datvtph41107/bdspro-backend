package repo

import (
	_enum "common/domain/enum"
	"context"
	"notification/internal/domain"
)

type PropertyHistoryQuery struct {
	SubjectID   *uint64
	SubjectType *string
	Action      *_enum.PropertyAction
	ActorID     *uint64
	Keyword     *string
	Limit       int
	Offset      int
	Sort        string
}

type PropertyHistoryRepo interface {
	Create(ctx context.Context, entity *domain.PropertyHistoryEntity) (*domain.PropertyHistoryEntity, error)
	Search(
		ctx context.Context,
		query PropertyHistoryQuery,
	) ([]domain.PropertyHistoryEntity, error)
}
