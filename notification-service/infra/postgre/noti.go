package postgre

import (
	"context"
	"noti/internal/domain"
	"noti/internal/dto"
	"time"

	"gorm.io/gorm"
)

// @bind: noti/internal/repo.NotificationRepo
type PostgreNotification struct {
	DB *gorm.DB
}

func NewPostgreNotification(db *gorm.DB) *PostgreNotification {
	return &PostgreNotification{DB: db}
}

func (r *PostgreNotification) Search(ctx context.Context, profileId uint64, dto *dto.SearchNotiRequest) ([]domain.NotificationEntity, int32, error) {
	var notifications []domain.NotificationEntity
	query := r.DB.
		WithContext(ctx).
		Where(`target_id = ? and deleted_at is null and 
			(visible_at is null OR visible_at < CURRENT_TIMESTAMP)
		`, profileId)

	if dto.IsRead != nil && *dto.IsRead {
		query = query.Where("is_read = true")
	}

	err := query.Find(&notifications).
		Order("created_at DESC").
		Offset(dto.GetOffset()).
		Limit(dto.GetLimit()).
		Error

	var total int64
	query.Count(&total)
	return notifications, int32(total), err
}

func (r *PostgreNotification) Create(ctx context.Context, notification *domain.NotificationEntity) error {
	return r.DB.WithContext(ctx).Create(notification).Error
}

// Count: Đếm số thông báo chưa đọc của user
func (r *PostgreNotification) Count(ctx context.Context, profileId uint64) (int64, error) {
	var count int64
	err := r.DB.WithContext(ctx).Model(&domain.NotificationEntity{}).
		Where("target_id = ? AND (is_read IS NULL OR is_read = FALSE) and (visible_at is null OR visible_at < CURRENT_TIMESTAMP)", profileId).
		Count(&count).Error
	return count, err
}

// MarkNotificationAsRead: Đánh dấu một thông báo là đã đọc
func (r *PostgreNotification) MarkAsRead(ctx context.Context, profileId uint64, id uint64) (int64, error) {
	result := r.DB.WithContext(ctx).Model(&domain.NotificationEntity{}).
		Where("id = ? AND target_id = ? and (visible_at is null OR visible_at < CURRENT_TIMESTAMP)", id, profileId).
		Update("is_read", true)
	return result.RowsAffected, result.Error
}

func (r *PostgreNotification) ReadAll(ctx context.Context, profileId uint64) (int64, error) {
	result := r.DB.WithContext(ctx).
		Model(&domain.NotificationEntity{}).
		Where("target_id = ? and (visible_at is null OR visible_at < CURRENT_TIMESTAMP) and (is_read = false or is_read is null)", profileId).
		Update("is_read", true)
	return result.RowsAffected, result.Error
}

// ExistsByIdAndUserId: Kiểm tra xem thông báo có tồn tại không
func (r *PostgreNotification) ExistsByIdAndUserId(id uint64, userID uint64) (bool, error) {
	var count int64
	err := r.DB.Model(&domain.NotificationEntity{}).
		Where("id = ? AND target_id = ? and (visible_at is null OR visible_at < CURRENT_TIMESTAMP)", id, userID).
		Count(&count).Error
	return count > 0, err
}

func (repo *PostgreNotification) Delete(ctx context.Context, profileId, id uint64) error {
	return repo.DB.Model(&domain.NotificationEntity{}).
		Where("target_id = ? AND id = ? and (visible_at is null OR visible_at < CURRENT_TIMESTAMP)", profileId, id).
		Update("deleted_at", time.Now()).
		Error
}

func (repo *PostgreNotification) ListNotiWithTargetIdAndOwnerType(ctx context.Context, targetId uint64, ownerType int32, hour int) ([]domain.NotificationEntity, error) {
	var notifications []domain.NotificationEntity
	err := repo.DB.WithContext(ctx).
		Where("target_id = ? AND owner_type = ? and deleted_at is null AND created_at >= NOW() - INTERVAL '? HOURS'", targetId, ownerType, hour).
		Find(&notifications).Error
	return notifications, err
}
