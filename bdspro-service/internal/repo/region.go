package repo

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"context"
)

type RegionRepo interface {
	GetByID(c context.Context, id uint64) (*domain.Region, error)
	Search(c context.Context, dto *dto.RegionRequest) ([]domain.Region, int64, error)
	InterText(text string, limit int) ([]uint64, error)
	InterTextToItem(text string, limit int) ([]domain.Region, error)
}
