package domain

import (
	_models "common/models"
	"notification/internal/enums"
	"time"

	"github.com/lib/pq"
)

type NotificationEntity struct {
	_models.BaseEntity
	ID         uint64 `gorm:"primaryKey"`
	TargetID   uint64 `gorm:"index"`
	OwnerType  int32  `gorm:"index"`
	OwnerID    uint64 `gorm:"index"`
	Title      string
	Message    pq.StringArray `gorm:"type:varchar[]"`
	AttachData string
	Link       string
	Type       enums.TypeEnum
	VisibleAt  *time.Time
	IsRead     bool `gorm:"default:false"`
}

func (n *NotificationEntity) TableName() string {
	return "notification"
}
