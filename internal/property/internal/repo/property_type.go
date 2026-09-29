package repo

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	_dto "common/domain/dto"
	"context"
)

type PropertyTypeRepo interface {
	// common.IBaseRepo[domain.PropertyType, *dto.PropertyTypeSearchDTO]
	GetAllItem(c context.Context, dto *dto.PropertyTypeSearchDTO) ([]domain.PropertyType, error)
	InterText(text string, limit int) ([]uint64, error)
	InterTextToItem(text string, limit int) (*_dto.ItemDTO, error)
}
