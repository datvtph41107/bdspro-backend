package repo

import (
	"chat/internal/domain"
	"context"
)

type BackgroundRepo interface {
	GetBackgroundImages(ctx context.Context, page, size *int64) ([]*domain.BackgroundImage, int64, error)
	CreateBackgroundImage(ctx context.Context, model *domain.BackgroundImage) (*domain.BackgroundImage, error)
	GetBackgroundImageByID(ctx context.Context, id string) (*domain.BackgroundImage, error)
	UpdateBackgroundImage(ctx context.Context, model *domain.BackgroundImage) (*domain.BackgroundImage, error)
	DeleteBackgroundImage(ctx context.Context, id uint64) error
}
