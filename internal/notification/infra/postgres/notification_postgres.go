package postgres

import (
	_enum "common/domain/enum"
	"context"
	"fmt"
	"notification/internal/domain"
	"notification/internal/dto"
	usecase "notification/internal/usecase"

	"gorm.io/gorm"
)

type PostgreNotification struct {
	DB *gorm.DB
}

func NewPostgreNotification(db *gorm.DB) usecase.NotificationStore {
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

func (r *PostgreNotification) CreateBatch(ctx context.Context, notifications []domain.NotificationEntity) error {
	if len(notifications) == 0 {
		return nil
	}
	return r.DB.WithContext(ctx).Create(&notifications).Error
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

// MarkAsReadMany: Đánh dấu nhiều thông báo là đã đọc
func (r *PostgreNotification) MarkAsReadMany(ctx context.Context, profileId uint64, ids []uint64) (int64, error) {
	if len(ids) == 0 {
		return 0, nil
	}
	result := r.DB.WithContext(ctx).Model(&domain.NotificationEntity{}).
		Where("id IN ? AND owner_id = ? and (visible_at is null OR visible_at < CURRENT_TIMESTAMP)", ids, profileId).
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
	return repo.DB.WithContext(ctx).
		Unscoped().
		Where("owner_id = ? AND id = ? and (visible_at is null OR visible_at < CURRENT_TIMESTAMP)", profileId, id).
		Delete(&domain.NotificationEntity{}).
		Error
}

func (repo *PostgreNotification) ListNotiWithTargetIdAndOwnerType(ctx context.Context, targetId uint64, ownerType int32, hour int) ([]domain.NotificationEntity, error) {
	var notifications []domain.NotificationEntity
	err := repo.DB.WithContext(ctx).
		Where("owner_id = ? AND owner_of = ? and deleted_at is null AND created_at >= NOW() - INTERVAL '? HOURS'", targetId, ownerType, hour).
		Find(&notifications).Error
	return notifications, err
}

func (repo *PostgreNotification) DeleteIds(ctx context.Context, ids []uint64) error {
	return repo.DB.WithContext(ctx).
		Unscoped().
		Where("id IN (?)", ids).
		Delete(&domain.NotificationEntity{}).
		Error
}

func (repo *PostgreNotification) RemoveUnique(ctx context.Context, ownerId, targetId uint64, notificationType _enum.ENotificationType, ownerOf _enum.EOwnerOf) (*domain.NotificationEntity, error) {
	var notification *domain.NotificationEntity
	err := repo.DB.WithContext(ctx).
		Debug().
		Unscoped().
		Where("owner_id = ? AND target_id = ? AND owner_of = ? AND notification_type = ? and deleted_at is null", ownerId, targetId, ownerOf, notificationType).
		Delete(&domain.NotificationEntity{}).Error
	return notification, err
}

func (repo *PostgreNotification) RemoveUnique2(ctx context.Context, ownerId uint64, targetId *uint64, notificationType _enum.ENotificationType, ownerOf _enum.EOwnerOf, attachData []string) (*domain.NotificationEntity, error) {
	var notification *domain.NotificationEntity
	query := repo.DB.WithContext(ctx).
		Debug().
		Unscoped().
		Where("owner_id = ? AND owner_of = ? AND notification_type = ? AND deleted_at is null", ownerId, ownerOf, notificationType)
	if targetId != nil {
		query = query.Where("target_id = ?", targetId)
	}

	// Lọc theo các phần tử có giá trị trong mảng attachData
	// PostgreSQL array index bắt đầu từ 1
	for i, value := range attachData {
		if value != "" {
			// Index trong PostgreSQL array bắt đầu từ 1, không phải 0
			arrayIndex := i + 1
			// Dùng raw SQL với fmt.Sprintf để build query string động cho array index
			query = query.Where(fmt.Sprintf("attach_data[%d] = ?", arrayIndex), value)
		}
	}

	err := query.Delete(&domain.NotificationEntity{}).Error
	return notification, err
}
