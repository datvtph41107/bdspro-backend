package repo

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	_dto "common/domain/dto"
	"context"
)

type AmenityRepo interface {
	GetAll(c context.Context, dto *dto.AmenitySearchDTO) ([]domain.AmenityItem, error)
	InterText(text string, limit int) ([]uint64, error)
	InterTextToItem(text string, limit int) ([]_dto.ItemDTO, error)
}
