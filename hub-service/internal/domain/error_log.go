package domain

import (
	"time"

	"gorm.io/datatypes"
)

// ErrorLog lưu lỗi từ app gửi lên (React Native, mobile)
type ErrorLog struct {
	ID        uint64         `gorm:"primaryKey;autoIncrement"`
	Message   string         `gorm:"type:text;not null;comment:Error message"`
	Stack     *string        `gorm:"type:text;comment:Stack trace"`
	Screen    *string        `gorm:"type:varchar(200);comment:Screen/route name"`
	UserID    *string        `gorm:"type:varchar(100);index;comment:User ID from app"`
	Device    datatypes.JSON `gorm:"type:jsonb;comment:Device info (platform, model, osVersion, appVersion)"`
	Extra     datatypes.JSON `gorm:"type:jsonb;comment:Extra data (breadcrumbs, etc)"`
	CreatedAt time.Time      `gorm:"autoCreateTime;index"`
}

func (ErrorLog) TableName() string {
	return "error_logs"
}
