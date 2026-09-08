package postgres

import (
	"context"
	"tqd/internal/domain"
	"tqd/internal/interface/repo"

	"gorm.io/gorm"
)

type notificationRepo struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) repo.NotificationRepository {
	return &notificationRepo{db: db}
}

func (r *notificationRepo) Create(ctx context.Context, notif *domain.UserNotification) error {
	return r.db.WithContext(ctx).Create(notif).Error
}

func (r *notificationRepo) MarkAsRead(ctx context.Context, userID uint64, notificationID uint64) error {
	return r.db.WithContext(ctx).Model(&domain.UserNotification{}).
		Where("id = ? AND user_id = ?", notificationID, userID).
		Updates(map[string]interface{}{"is_read": true, "read_at": gorm.Expr("NOW()")}).Error
}

func (r *notificationRepo) MarkAllAsRead(ctx context.Context, userID uint64) error {
	return r.db.WithContext(ctx).Model(&domain.UserNotification{}).
		Where("user_id = ? AND is_read = false", userID).
		Updates(map[string]interface{}{"is_read": true, "read_at": gorm.Expr("NOW()")}).Error
}

func (r *notificationRepo) ListByUser(ctx context.Context, userID uint64, notifType *uint32, read *bool, page, limit int) ([]domain.UserNotification, int64, error) {
	var notifs []domain.UserNotification
	var total int64
	query := r.db.WithContext(ctx).Model(&domain.UserNotification{}).Where("user_id = ?", userID)
	if notifType != nil {
		query = query.Where("notification_type = ?", *notifType)
	}
	if read != nil {
		query = query.Where("is_read = ?", *read)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * limit
	err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&notifs).Error
	return notifs, total, err
}

func (r *notificationRepo) CountUnread(ctx context.Context, userID uint64) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&domain.UserNotification{}).
		Where("user_id = ? AND is_read = false", userID).
		Count(&count).Error
	return count, err
}

func (r *notificationRepo) DeleteExpired(ctx context.Context) error {
	return r.db.WithContext(ctx).Where("expires_at IS NOT NULL AND expires_at < NOW()").Delete(&domain.UserNotification{}).Error
}
