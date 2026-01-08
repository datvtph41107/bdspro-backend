package repo

import (
	"context"
	"notification/internal/domain"
	"notification/internal/enums"
)

type PropertyHistoryQuery struct {
	SubjectID   *uint64
	SubjectType *string
	Action      *enums.PropertyAction
	ActorID     *uint64
	ActorType   *enums.ActorType
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
