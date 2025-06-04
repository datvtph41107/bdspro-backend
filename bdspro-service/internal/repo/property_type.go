package repo

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"context"
)

type PropertyTypeRepo interface {
	// common.IBaseRepo[domain.PropertyType, *dto.PropertyTypeSearchDTO]
	GetAllItem(c context.Context, dto *dto.PropertyTypeSearchDTO) ([]domain.PropertyTypeItem, error)
	InterText(text string, limit int) ([]uint64, error)
	InterTextToItem(text string, limit int) ([]domain.PropertyTypeItem, error)
}
