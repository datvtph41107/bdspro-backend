package domain

import (
	_models "common/domain/entity"
)

type Observation struct {
	_models.BaseEntity
	EntityID   uint64  `gorm:"not null;index;comment:ID của entity được theo dõi"`
	EntityType *string `gorm:"type:varchar(32);not null;index;comment:Loại entity được theo dõi"`
	ProfileID  *uint64 `gorm:"index;comment:ID người dùng nếu đã đăng nhập"`
	Action     string  `gorm:"type:varchar(32);not null;index;comment:Loại hành vi người dùng"`
	Duration   uint32  `gorm:"default:0;comment:Thời gian thực hiện hành vi (giây)"`
	StartTime  uint64  `gorm:"comment:Thời gian bắt đầu event (ms)"`
	EndTime    uint64  `gorm:"comment:Thời gian kết thúc event (ms)"`
	Screen     *string `gorm:"type:varchar(64);comment:Tên screen nơi event xảy ra"`
}

func (Observation) TableName() string {
	return "observations"
}
