package repo

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	_dto "common/domain/dto"
	"context"
)

type DocTypeRepo interface {
	GetAllItem(c context.Context) ([]domain.DocTypeItem, error)
	SearchItem(c context.Context, dto *dto.DocTypeSearchDTO) ([]domain.DocTypeItem, int64, error)
	InterText(c context.Context, text string, limit int) ([]uint64, error)
	InterTextToItem(c context.Context, text string, limit int) (*_dto.ItemDTO, error)
}
