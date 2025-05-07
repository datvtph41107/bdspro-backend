package repo

import (
	_utils "common/utils"
	"noti/data"
	"noti/models"
	"time"

	"gorm.io/gorm"
)

type NotificationRepository struct {
	DB *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) *NotificationRepository {
	return &NotificationRepository{DB: db}
}

func (r *NotificationRepository) Search(dto data.SearchNotiRequest, pageable _utils.Page) ([]models.NotificationEntity, error) {
	var notifications []models.NotificationEntity
	query := r.DB.Where("user_id = ? and deleted_at is null and (visible_at is null OR visible_at < CURRENT_TIMESTAMP)", dto.UserID)

	// Phân trang
	if pageable.Limit > 0 {
		query = query.Limit(pageable.Limit).Offset(pageable.Skip)
	}

	// Sắp xếp
	if pageable.Sort != "" {
		query = query.Order(pageable.Sort)
	}

	err := query.Find(&notifications).Error
	return notifications, err
}

func (r *NotificationRepository) Create(notification *models.NotificationEntity) error {
	return r.DB.Create(notification).Error
}

// Count: Đếm số thông báo chưa đọc của user
func (r *NotificationRepository) Count(userID uint64) (int64, error) {
	var count int64
	err := r.DB.Model(&models.NotificationEntity{}).
		Where("user_id = ? AND (is_read IS NULL OR is_read = FALSE) and (visible_at is null OR visible_at < CURRENT_TIMESTAMP)", userID).
		Count(&count).Error
	return count, err
}

// MarkNotificationAsRead: Đánh dấu một thông báo là đã đọc
func (r *NotificationRepository) MarkNotificationAsRead(id uint64, userID uint64) (int64, error) {
	result := r.DB.Model(&models.NotificationEntity{}).
		Where("id = ? AND user_id = ? and (visible_at is null OR visible_at < CURRENT_TIMESTAMP)", id, userID).
		Update("is_read", true)
	return result.RowsAffected, result.Error
}

// ExistsByIdAndUserId: Kiểm tra xem thông báo có tồn tại không
func (r *NotificationRepository) ExistsByIdAndUserId(id uint64, userID uint64) (bool, error) {
	var count int64
	err := r.DB.Model(&models.NotificationEntity{}).
		Where("id = ? AND user_id = ? and (visible_at is null OR visible_at < CURRENT_TIMESTAMP)", id, userID).
		Count(&count).Error
	return count > 0, err
}

func (repo *NotificationRepository) DeleteByID(profileId, id uint64) error {
	return repo.DB.Model(&models.NotificationEntity{}).
		Where("created_by = ? AND id = ? and (visible_at is null OR visible_at < CURRENT_TIMESTAMP)", profileId, id).
		Update("deleted_at", time.Now()).
		Error
}
