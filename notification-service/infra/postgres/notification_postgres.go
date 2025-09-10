package postgres

import (
	"context"
	"notification/internal/domain"
	"notification/internal/dto"
	"notification/internal/enums"
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

func (r *PostgreNotification) GetQuery(ctx context.Context, profileId uint64, dto *dto.SearchNotiRequest) *gorm.DB {
	query := r.DB.WithContext(ctx).Model(&domain.NotificationEntity{}).
		Where("owner_id = ? and deleted_at is null and (visible_at is null OR visible_at < CURRENT_TIMESTAMP)", profileId)

	if dto.IsRead != nil && *dto.IsRead {
		query = query.Where("is_read = true")
	}

	return query
}

func (r *PostgreNotification) Search(ctx context.Context, profileId uint64, dto *dto.SearchNotiRequest) ([]domain.NotificationEntity, int32, error) {
	var notifications []domain.NotificationEntity
	query := r.GetQuery(ctx, profileId, dto)

	var total int64
	query.Count(&total)

	err := r.GetQuery(ctx, profileId, dto).
		Debug().
		Order("created_at DESC").
		Offset(dto.GetOffset()).
		Limit(dto.GetLimit()).
		Find(&notifications).
		Error

	return notifications, int32(total), err
}

func (r *PostgreNotification) Create(ctx context.Context, notification *domain.NotificationEntity) error {
	return r.DB.WithContext(ctx).Create(notification).Error
}

// Count: Đếm số thông báo chưa đọc của user
func (r *PostgreNotification) Count(ctx context.Context, profileId uint64) (int64, error) {
	var count int64
	err := r.DB.WithContext(ctx).Model(&domain.NotificationEntity{}).
		Where("owner_id = ? AND (is_read IS NULL OR is_read = FALSE) and (visible_at is null OR visible_at < CURRENT_TIMESTAMP)", profileId).
		Count(&count).Error
	return count, err
}

// MarkNotificationAsRead: Đánh dấu một thông báo là đã đọc
func (r *PostgreNotification) MarkAsRead(ctx context.Context, profileId uint64, id uint64) (int64, error) {
	result := r.DB.WithContext(ctx).Model(&domain.NotificationEntity{}).
		Where("id = ? AND owner_id = ? and (visible_at is null OR visible_at < CURRENT_TIMESTAMP)", id, profileId).
		Update("is_read", true)
	return result.RowsAffected, result.Error
}

func (r *PostgreNotification) ReadAll(ctx context.Context, profileId uint64) (int64, error) {
	result := r.DB.WithContext(ctx).
		Model(&domain.NotificationEntity{}).
		Where("owner_id = ? and (visible_at is null OR visible_at < CURRENT_TIMESTAMP) and (is_read = false or is_read is null)", profileId).
		Update("is_read", true)
	return result.RowsAffected, result.Error
}

// ExistsByIdAndUserId: Kiểm tra xem thông báo có tồn tại không
func (r *PostgreNotification) ExistsByIdAndUserId(id uint64, userID uint64) (bool, error) {
	var count int64
	err := r.DB.Model(&domain.NotificationEntity{}).
		Where("id = ? AND owner_id = ? and (visible_at is null OR visible_at < CURRENT_TIMESTAMP)", id, userID).
		Count(&count).Error
	return count > 0, err
}

func (repo *PostgreNotification) Delete(ctx context.Context, profileId, id uint64) error {
	return repo.DB.Model(&domain.NotificationEntity{}).
		Where("owner_id = ? AND id = ? and (visible_at is null OR visible_at < CURRENT_TIMESTAMP)", profileId, id).
		Update("deleted_at", time.Now()).
		Error
}

func (repo *PostgreNotification) ListNotiWithTargetIdAndOwnerType(ctx context.Context, targetId uint64, ownerType int32, hour int) ([]domain.NotificationEntity, error) {
	var notifications []domain.NotificationEntity
	err := repo.DB.WithContext(ctx).
		Where("owner_id = ? AND owner_type = ? and deleted_at is null AND created_at >= NOW() - INTERVAL '? HOURS'", targetId, ownerType, hour).
		Find(&notifications).Error
	return notifications, err
}

func (repo *PostgreNotification) DeleteIds(ctx context.Context, ids []uint64) error {
	return repo.DB.Model(&domain.NotificationEntity{}).
		Where("id IN (?)", ids).
		Update("deleted_at", time.Now()).
		Error
}

func (repo *PostgreNotification) RemoveUnique(
	ctx context.Context,
	ownerId uint64,
	targetId uint64,
	typeNoti enums.TypeEnum,
	ownerType enums.OwnerType,
) (*domain.NotificationEntity, error) {
	var notification *domain.NotificationEntity
	err := repo.DB.WithContext(ctx).
		Debug().
		Model(&domain.NotificationEntity{}).
		Where("owner_id = ? AND target_id = ? AND owner_type = ? AND type = ? and deleted_at is null", ownerId, targetId, ownerType, typeNoti).
		Update("deleted_at", time.Now()).Error
	// Find(&notification).Error
	return notification, err
}
