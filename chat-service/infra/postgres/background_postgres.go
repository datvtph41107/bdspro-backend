package postgres

import (
	"context"
	"strconv"
	"time"

	"chat/internal/domain"
	_repo "chat/internal/repo"
	_db "common/db"

	"gorm.io/gorm"
)

type BackgroundRepo struct {
	*_db.TransactionRepo
}

func NewBackgroundRepo(db *_db.TransactionRepo) _repo.BackgroundRepo {
	return &BackgroundRepo{db}
}

func (r *BackgroundRepo) GetBackgroundImages(ctx context.Context, page, size *int64) ([]*domain.BackgroundImage, int64, error) {
	var images []*domain.BackgroundImage
	var total int64

	queryPage, querySize := int64(0), int64(10)
	if page != nil {
		queryPage = *page
	}
	if size != nil {
		querySize = *size
	}

	err := r.GetDB(ctx).Model(&domain.BackgroundImage{}).Where("deleted_at IS NULL").Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.GetDB(ctx).
		Where("deleted_at IS NULL").
		Limit(int(querySize)).
		Offset(int(queryPage * querySize)).
		Find(&images).Error
	if err != nil {
		return nil, 0, err
	}

	return images, total, nil
}

func (r *BackgroundRepo) GetBackgroundImageByID(ctx context.Context, id string) (*domain.BackgroundImage, error) {
	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return nil, err
	}

	var image domain.BackgroundImage
	err = r.GetDB(ctx).
		Where("id = ? AND deleted_at IS NULL", idUint).
		First(&image).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &image, nil
}

func (r *BackgroundRepo) CreateBackgroundImage(ctx context.Context, model *domain.BackgroundImage) (*domain.BackgroundImage, error) {
	err := r.GetDB(ctx).Create(model).Error
	if err != nil {
		return nil, err
	}
	return model, nil
}

func (r *BackgroundRepo) UpdateBackgroundImage(ctx context.Context, model *domain.BackgroundImage) (*domain.BackgroundImage, error) {
	err := r.GetDB(ctx).Where("deleted_at IS NULL").Model(model).Updates(model).Error
	if err != nil {
		return nil, err
	}
	return model, nil
}

func (r *BackgroundRepo) DeleteBackgroundImage(ctx context.Context, id uint64) error {
	now := time.Now()
	err := r.GetDB(ctx).Model(&domain.BackgroundImage{}).Where("id = ?", id).Update("deleted_at", now).Error
	return err
}
