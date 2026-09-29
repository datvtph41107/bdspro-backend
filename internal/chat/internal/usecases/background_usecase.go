package usecases

import (
	"context"

	"chat/internal/domain"
	_repo "chat/internal/repo"
)

type BackgroundImageUsecases struct {
	backgroundRepo _repo.BackgroundRepo
}

func NewBackgroundImageUsecases(backgroundRepo _repo.BackgroundRepo) *BackgroundImageUsecases {
	return &BackgroundImageUsecases{
		backgroundRepo: backgroundRepo,
	}
}

func (b *BackgroundImageUsecases) GetBackgroundImages(ctx context.Context, page, size *int64) ([]*domain.BackgroundImage, int64, error) {
	images, total, err := b.backgroundRepo.GetBackgroundImages(ctx, page, size)
	if err != nil {
		return nil, 0, err
	}
	return images, total, nil
}

func (b *BackgroundImageUsecases) CreateBackgroundImage(ctx context.Context, model *domain.BackgroundImage) (*domain.BackgroundImage, error) {
	savedModel, err := b.backgroundRepo.CreateBackgroundImage(ctx, model)
	if err != nil {
		return nil, err
	}
	return savedModel, nil
}

func (b *BackgroundImageUsecases) GetBackgroundImageByID(ctx context.Context, id string) (*domain.BackgroundImage, error) {
	model, err := b.backgroundRepo.GetBackgroundImageByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return model, nil
}

func (b *BackgroundImageUsecases) UpdateBackgroundImage(ctx context.Context, model *domain.BackgroundImage) (*domain.BackgroundImage, error) {
	savedModel, err := b.backgroundRepo.UpdateBackgroundImage(ctx, model)
	if err != nil {
		return nil, err
	}
	return savedModel, nil
}

func (b *BackgroundImageUsecases) DeleteBackgroundImage(ctx context.Context, id uint64) error {
	err := b.backgroundRepo.DeleteBackgroundImage(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
