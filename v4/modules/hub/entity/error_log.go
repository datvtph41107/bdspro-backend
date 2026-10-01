package entity

import (
	"time"

	"gorm.io/datatypes"
)

// ErrorLog stores errors sent from mobile/RN apps.
type ErrorLog struct {
	ID        uint64         `gorm:"primaryKey;autoIncrement"`
	Message   string         `gorm:"type:text;not null"`
	Stack     *string        `gorm:"type:text"`
	Screen    *string        `gorm:"type:varchar(200)"`
	UserID    *string        `gorm:"type:varchar(100);index"`
	Device    datatypes.JSON `gorm:"type:jsonb"`
	Extra     datatypes.JSON `gorm:"type:jsonb"`
	CreatedAt time.Time      `gorm:"autoCreateTime;index"`
}

func (ErrorLog) TableName() string { return "error_logs" }
