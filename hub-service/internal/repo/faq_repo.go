package repo

import (
	"context"

	_crud "common/domain/crud"
	_dto "common/domain/dto"
	"hub/internal/domain"
)

type IFAQRepo interface {
	_crud.ICrudRepo[domain.FAQEntity]
	GetListWithFilter(ctx context.Context, question string, groupKey string, pagable _dto.IPagable) ([]*domain.FAQEntity, int64, error)
	GetSimpleListWithText(ctx context.Context, text string, groupKey string, pagable _dto.IPagable) ([]*domain.FAQEntity, int64, error)
}
