package postgres

import (
	"context"
	"errors"
	"tqd/internal/domain"
	"tqd/internal/interface/repo"

	"gorm.io/gorm"
)

type subscriptionRepo struct {
	db *gorm.DB
}

func NewSubscriptionRepository(db *gorm.DB) repo.SubscriptionRepository {
	return &subscriptionRepo{db: db}
}

func (r *subscriptionRepo) Create(ctx context.Context, sub *domain.UserSubscription) error {
	return r.db.WithContext(ctx).Create(sub).Error
}

func (r *subscriptionRepo) Update(ctx context.Context, sub *domain.UserSubscription) error {
	return r.db.WithContext(ctx).Save(sub).Error
}

func (r *subscriptionRepo) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&domain.UserSubscription{}, id).Error
}

func (r *subscriptionRepo) GetByID(ctx context.Context, id uint64) (*domain.UserSubscription, error) {
	var sub domain.UserSubscription
	err := r.db.WithContext(ctx).First(&sub, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &sub, err
}

func (r *subscriptionRepo) GetByUserAndTarget(ctx context.Context, userID uint64, targetType string, targetID uint64) (*domain.UserSubscription, error) {
	var sub domain.UserSubscription
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND target_type = ? AND target_id = ? AND deleted_at IS NULL", userID, targetType, targetID).
		First(&sub).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &sub, err
}

func (r *subscriptionRepo) ListByUser(ctx context.Context, userID uint64, targetType *string, status *string, page, limit int) ([]domain.UserSubscription, int64, error) {
	var subs []domain.UserSubscription
	var total int64
	query := r.db.WithContext(ctx).Model(&domain.UserSubscription{}).Where("user_id = ? AND deleted_at IS NULL", userID)
	if targetType != nil {
		query = query.Where("target_type = ?", *targetType)
	}
	if status != nil {
		query = query.Where("status = ?", *status)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * limit
	err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&subs).Error
	return subs, total, err
}

func (r *subscriptionRepo) ListByTarget(ctx context.Context, targetType string, targetID uint64) ([]domain.UserSubscription, error) {
	var subs []domain.UserSubscription
	err := r.db.WithContext(ctx).
		Where("target_type = ? AND target_id = ? AND status = ? AND deleted_at IS NULL", targetType, targetID, "active").
		Find(&subs).Error
	return subs, err
}

func (r *subscriptionRepo) ListByUserWithStatus(ctx context.Context, userID uint64, status *uint32, page, limit int) ([]domain.UserSubscription, int64, error) {
	var subs []domain.UserSubscription
	var total int64

	query := r.db.WithContext(ctx).Model(&domain.UserSubscription{}).
		Where("user_id = ? AND target_type = ? AND deleted_at IS NULL", userID, "parcel")

	if status != nil {
		query = query.Where("status = ?", *status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&subs).Error

	return subs, total, err
}

func (r *subscriptionRepo) AdminList(ctx context.Context, userID *uint64, targetType *string, status *uint32, page, limit int) ([]domain.UserSubscription, int64, error) {
	var subs []domain.UserSubscription
	var total int64
	query := r.db.WithContext(ctx).Model(&domain.UserSubscription{}).Where("deleted_at IS NULL")
	if userID != nil {
		query = query.Where("user_id = ?", *userID)
	}
	if targetType != nil {
		query = query.Where("target_type = ?", *targetType)
	}
	if status != nil {
		query = query.Where("status = ?", *status)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * limit
	err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&subs).Error
	return subs, total, err
}

func (r *subscriptionRepo) AdminDelete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&domain.UserSubscription{}, id).Error
}
