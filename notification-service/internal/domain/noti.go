package domain

import (
	_models "common/models"
	"noti/internal/enums"
	"time"
)

type NotificationEntity struct {
	_models.BaseEntity
	ID        uint64 `gorm:"primaryKey"`
	TargetID  uint64 `gorm:"index"`
	OwnerType int32  `gorm:"index"`
	Title     string
	Message   string
	Link      string
	Type      enums.TypeEnum
	VisibleAt *time.Time
	IsRead    bool `gorm:"default:false"`
}

func (n *NotificationEntity) TableName() string {
	return "notification"
}
