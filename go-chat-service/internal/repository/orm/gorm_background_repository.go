package orm

import (
	"context"

	"github.com/tuandq3/go-chat-service/models"
)

func (r *gormRepository) GetBackgroundImages(ctx context.Context, page, size *int64) ([]*models.BackgroundImageModel, int64, error) {
	var models []*models.BackgroundImageModel
	var total int64

	queryPage, querySize := 0, 10
	if page != nil {
		queryPage = int(*page)
	}
	if size != nil {
		querySize = int(*size)
	}

	err := r.db.WithContext(ctx).Model(&models).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	err = r.db.WithContext(ctx).Limit(querySize).Offset(queryPage * querySize).Find(&models).Error
	if err != nil {
		return nil, 0, err
	}
	return models, total, nil
}

func (r *gormRepository) CreateBackgroundImage(ctx context.Context, model *models.BackgroundImageModel) (*models.BackgroundImageModel, error) {
	err := r.db.WithContext(ctx).Create(model).Error
	if err != nil {
		return nil, err
	}
	return model, nil
}

func (r *gormRepository) GetBackgroundImageByID(ctx context.Context, id string) (*models.BackgroundImageModel, error) {
	var model models.BackgroundImageModel
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&model).Error
	if err != nil {
		return nil, err
	}
	return &model, nil
}

func (r *gormRepository) UpdateBackgroundImage(ctx context.Context, model *models.BackgroundImageModel) (*models.BackgroundImageModel, error) {
	err := r.db.WithContext(ctx).Model(model).Updates(model).Error
	if err != nil {
		return nil, err
	}
	return model, nil
}

func (r *gormRepository) DeleteBackgroundImage(ctx context.Context, id uint64) error {
	var model models.BackgroundImageModel
	err := r.db.WithContext(ctx).Where("id = ?", id).Delete(&model).Error
	if err != nil {
		return err
	}
	return nil
}
