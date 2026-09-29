package domain

import (
	_enum "common/domain/enum"
	_models "common/models"
)

type NotificationHistoryEntity struct {
	_models.BaseEntity
	ID      uint64         `gorm:"primaryKey"`
	OwnerOf _enum.EOwnerOf `gorm:"index"`
	OwnerID uint64         `gorm:"index"`
}

func (n *NotificationHistoryEntity) TableName() string {
	return "notification_history"
}
