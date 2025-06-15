package repo

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"context"
)

type DocTypeRepo interface {
	GetAllItem(c context.Context) ([]domain.DocTypeItem, error)
	SearchItem(c context.Context, dto *dto.DocTypeSearchDTO) ([]domain.DocTypeItem, int64, error)
	InterText(text string, limit int) ([]uint64, error)
	InterTextToItem(text string, limit int) ([]domain.DocTypeItem, error)
}
