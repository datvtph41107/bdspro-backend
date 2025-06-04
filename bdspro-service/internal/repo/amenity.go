package repo

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"context"
)

type AmenityRepo interface {
	GetAll(c context.Context, dto *dto.AmenitySearchDTO) ([]domain.AmenityItem, error)
	InterText(text string, limit int) ([]uint64, error)
	InterTextToEntity(text string, limit int) ([]domain.AmenityItem, error)
}
