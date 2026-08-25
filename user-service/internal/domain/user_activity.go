package domain

import "time"

// UserActivityLog lưu hoạt động login/refresh theo khung giờ (Asia/Ho_Chi_Minh).
// Unique (profile_id, slot_at): 1 user chỉ 1 bản ghi / giờ; sang giờ mới mới tạo mới.
// Thống kê theo ngày = COUNT DISTINCT profile_id trong ngày.
type UserActivityLog struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement"`
	ProfileID uint64    `gorm:"not null;uniqueIndex:idx_user_activity_profile_slot;column:profile_id"`
	SlotAt    time.Time `gorm:"not null;uniqueIndex:idx_user_activity_profile_slot;index;column:slot_at"`
	Source    string    `gorm:"size:32;not null;column:source"` // login | refresh
	CreatedAt time.Time `gorm:"not null;column:created_at"`
}

func (UserActivityLog) TableName() string {
	return "user_activity_log"
}
