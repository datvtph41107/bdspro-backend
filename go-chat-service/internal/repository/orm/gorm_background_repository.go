package orm

import (
	"context"
	"time"

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

	err := r.db.WithContext(ctx).Model(&models).Where("is_deleted = ?", false).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	err = r.db.WithContext(ctx).Where("is_deleted = ?", false).Limit(querySize).Offset(queryPage * querySize).Find(&models).Error
	if err != nil {
		return nil, 0, err
	}
	return models, total, nil
}

func (r *gormRepository) CreateBackgroundImage(ctx context.Context, model *models.BackgroundImageModel) (*models.BackgroundImageModel, error) {
	model.IsDeleted = false
	err := r.db.WithContext(ctx).Create(model).Error
	if err != nil {
		return nil, err
	}
	return model, nil
}

func (r *gormRepository) GetBackgroundImageByID(ctx context.Context, id string) (*models.BackgroundImageModel, error) {
	var model models.BackgroundImageModel
	err := r.db.WithContext(ctx).Where("id = ? AND is_deleted = ?", id, false).First(&model).Error
	if err != nil {
		return nil, err
	}
	return &model, nil
}

func (r *gormRepository) UpdateBackgroundImage(ctx context.Context, model *models.BackgroundImageModel) (*models.BackgroundImageModel, error) {
	err := r.db.WithContext(ctx).Where("is_deleted = ?", false).Model(model).Updates(model).Error
	if err != nil {
		return nil, err
	}
	return model, nil
}

func (r *gormRepository) DeleteBackgroundImage(ctx context.Context, id uint64) error {
	err := r.db.WithContext(ctx).Model(&models.BackgroundImageModel{}).Where("id = ?", id).Updates(map[string]interface{}{
		"is_deleted":  true,
		"deleted_at": time.Now(),
	}).Error
	if err != nil {
		return err
	}
	return nil
}
