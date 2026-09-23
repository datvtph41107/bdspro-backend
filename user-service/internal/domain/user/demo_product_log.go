package user

import "time"

// DemoProductLog lưu sản phẩm tạo trong chế độ demo để xoá khi tắt.
type DemoProductLog struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement;column:id"`
	ProfileID uint64    `gorm:"column:profile_id;not null;index:idx_demo_product_log_profile_id"`
	ProductID uint64    `gorm:"column:product_id;not null"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (DemoProductLog) TableName() string {
	return "demo_product_log"
}

// DemoContactLog lưu liên hệ CRM tạo trong chế độ demo để xoá khi tắt.
type DemoContactLog struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement;column:id"`
	ProfileID uint64    `gorm:"column:profile_id;not null;index:idx_demo_contact_log_profile_id"`
	ContactID uint64    `gorm:"column:contact_id;not null"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (DemoContactLog) TableName() string {
	return "demo_contact_log"
}
