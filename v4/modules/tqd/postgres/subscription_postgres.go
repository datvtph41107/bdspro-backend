package postgres

import (
	"context"
	"errors"
	domain "userv4/modules/tqd/entity"
	"userv4/modules/tqd/dto"
	repo "userv4/modules/tqd/repository"

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

func (r *subscriptionRepo) CountFollowersByParcelIDs(ctx context.Context, parcelIDs []uint64, status *uint32) (map[uint64]int64, error) {
	result := make(map[uint64]int64, len(parcelIDs))
	if len(parcelIDs) == 0 {
		return result, nil
	}

	type countRow struct {
		TargetID uint64
		Count    int64
	}
	rows := make([]countRow, 0)
	query := r.db.WithContext(ctx).Model(&domain.UserSubscription{}).
		Select("target_id, COUNT(*) AS count").
		Where("target_type = ? AND deleted_at IS NULL", "parcel").
		Where("target_id IN ?", parcelIDs)
	if status != nil {
		query = query.Where("status = ?", *status)
	}
	if err := query.Group("target_id").Scan(&rows).Error; err != nil {
		return nil, err
	}
	for _, row := range rows {
		result[row.TargetID] = row.Count
	}
	return result, nil
}

func (r *subscriptionRepo) CountSubscribersByLayerIDs(ctx context.Context, layerIDs []uint64, status *uint32) (map[uint64]int64, error) {
	result := make(map[uint64]int64, len(layerIDs))
	if len(layerIDs) == 0 {
		return result, nil
	}

	layerSet := make(map[uint64]struct{}, len(layerIDs))
	userSets := make(map[uint64]map[uint64]struct{}, len(layerIDs))
	for _, layerID := range layerIDs {
		if layerID == 0 {
			continue
		}
		layerSet[layerID] = struct{}{}
		userSets[layerID] = make(map[uint64]struct{})
		result[layerID] = 0
	}
	if len(layerSet) == 0 {
		return result, nil
	}

	var subs []domain.UserSubscription
	query := r.db.WithContext(ctx).Model(&domain.UserSubscription{}).
		Select("user_id", "subscription_scope").
		Where("deleted_at IS NULL")
	if status != nil {
		query = query.Where("status = ?", *status)
	}
	if err := query.Find(&subs).Error; err != nil {
		return nil, err
	}

	for _, sub := range subs {
		if sub.UserID == 0 {
			continue
		}
		scope, err := dto.SubscriptionScopeJSONToProto(sub.SubscriptionScope)
		if err != nil {
			continue
		}
		for layerID := range layerSet {
			if dto.IsLayerInSubscriptionScope(scope, layerID) {
				userSets[layerID][sub.UserID] = struct{}{}
			}
		}
	}

	for layerID, users := range userSets {
		result[layerID] = int64(len(users))
	}
	return result, nil
}
