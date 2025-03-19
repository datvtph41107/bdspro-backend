package models

import (
	_models "common/models"
	"noti/enums"
)

type NotificationEntity struct {
	_models.BaseEntity
	ID      uint64 `gorm:"primaryKey"`
	UserID  uint64 `gorm:"index"`
	Title   string
	Message string
	Link    string
	Type    enums.TypeEnum
	IsRead  bool `gorm:"default:false"`
}

func (n *NotificationEntity) TableName() string {
	return "notification"
}
