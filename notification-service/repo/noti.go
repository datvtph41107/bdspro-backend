package repo

import (
	_db "common/db"
	"noti/data"
	"noti/models"

	"gorm.io/gorm"
)

type NotificationRepository struct {
	DB *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) *NotificationRepository {
	return &NotificationRepository{DB: db}
}

func (r *NotificationRepository) Search(dto data.SearchNotiRequest, pageable _db.Page) ([]models.NotificationEntity, error) {
	var notifications []models.NotificationEntity
	query := r.DB.Where("user_id = ?", dto.UserID)

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
		Where("user_id = ? AND (is_read IS NULL OR is_read = FALSE)", userID).
		Count(&count).Error
	return count, err
}

// MarkNotificationAsRead: Đánh dấu một thông báo là đã đọc
func (r *NotificationRepository) MarkNotificationAsRead(id uint64, userID uint64) (int64, error) {
	result := r.DB.Model(&models.NotificationEntity{}).
		Where("id = ? AND user_id = ?", id, userID).
		Update("is_read", true)
	return result.RowsAffected, result.Error
}

// ExistsByIdAndUserId: Kiểm tra xem thông báo có tồn tại không
func (r *NotificationRepository) ExistsByIdAndUserId(id uint64, userID uint64) (bool, error) {
	var count int64
	err := r.DB.Model(&models.NotificationEntity{}).
		Where("id = ? AND user_id = ?", id, userID).
		Count(&count).Error
	return count > 0, err
}
