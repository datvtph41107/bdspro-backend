package usecase

import (
	"context"
	"organization/internal/domain/entity"
	"organization/internal/domain/repository"
)

type ColorUsecase interface {
	CreateColor(ctx context.Context, color *entity.Color) error
	UpdateColor(ctx context.Context, color *entity.Color) error
	DeleteColor(ctx context.Context, id uint32) error
	GetColorByID(ctx context.Context, id uint32) (*entity.Color, error)
	GetColors(ctx context.Context, page, size int, isActive *bool) ([]*entity.Color, int64, error)
}

type colorUsecase struct {
	colorRepo repository.ColorRepository
}

func NewColorUsecase(colorRepo repository.ColorRepository) ColorUsecase {
	return &colorUsecase{
		colorRepo: colorRepo,
	}
}

func (u *colorUsecase) CreateColor(ctx context.Context, color *entity.Color) error {
	return u.colorRepo.Create(ctx, color)
}

func (u *colorUsecase) UpdateColor(ctx context.Context, color *entity.Color) error {
	return u.colorRepo.Update(ctx, color)
}

func (u *colorUsecase) DeleteColor(ctx context.Context, id uint32) error {
	return u.colorRepo.Delete(ctx, id)
}

func (u *colorUsecase) GetColorByID(ctx context.Context, id uint32) (*entity.Color, error) {
	return u.colorRepo.GetByID(ctx, id)
}

func (u *colorUsecase) GetColors(ctx context.Context, page, size int, isActive *bool) ([]*entity.Color, int64, error) {
	return u.colorRepo.GetAll(ctx, page, size, isActive)
} 