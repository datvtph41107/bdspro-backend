package usecases

import (
	"context"

	"chat/models"
)

type backgroundImageUsecases struct {
	repository Repository
}

func NewBackgroundImageUsecases(repository Repository) *backgroundImageUsecases {
	return &backgroundImageUsecases{
		repository: repository,
	}
}

func (b *backgroundImageUsecases) GetBackgroundImages(ctx context.Context, page, size *int64) ([]*models.BackgroundImageModel, int64, error) {
	models, total, err := b.repository.GetBackgroundImages(ctx, page, size)
	if err != nil {
		return nil, 0, err
	}
	return models, total, nil
}

func (b *backgroundImageUsecases) CreateBackgroundImage(ctx context.Context, model *models.BackgroundImageModel) (*models.BackgroundImageModel, error) {
	savedModel, err := b.repository.CreateBackgroundImage(ctx, model)
	if err != nil {
		return nil, err
	}
	return savedModel, nil
}

func (b *backgroundImageUsecases) GetBackgroundImageByID(ctx context.Context, id string) (*models.BackgroundImageModel, error) {
	model, err := b.repository.GetBackgroundImageByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return model, nil
}

func (b *backgroundImageUsecases) UpdateBackgroundImage(ctx context.Context, model *models.BackgroundImageModel) (*models.BackgroundImageModel, error) {
	savedModel, err := b.repository.UpdateBackgroundImage(ctx, model)
	if err != nil {
		return nil, err
	}
	return savedModel, nil
}

func (b *backgroundImageUsecases) DeleteBackgroundImage(ctx context.Context, id uint64) error {
	err := b.repository.DeleteBackgroundImage(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
