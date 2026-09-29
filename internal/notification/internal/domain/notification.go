package domain

import (
	_enum "common/domain/enum"
	_models "common/models"
	"time"

	"github.com/lib/pq"
)

type NotificationEntity struct {
	_models.BaseEntity
	ID               uint64                  `gorm:"primaryKey"`
	OwnerOf          _enum.EOwnerOf          `gorm:"index"`
	OwnerID          uint64                  `gorm:"index"`
	Title            string                  `gorm:"size:255"`
	Message          pq.StringArray          `gorm:"type:varchar[]"`
	AttachData       pq.StringArray          `gorm:"type:varchar[]"`
	TargetID         uint64                  `gorm:"index"`
	NotificationType _enum.ENotificationType `gorm:"index"`
	Avatar           string                  `gorm:"size:255"`
	// Link remains a transport/read-model concern until a durable link contract exists.
	VisibleAt *time.Time `gorm:"index" json:"visibleAt,omitempty"`
	IsRead    bool       `gorm:"default:false"`
}

func (n *NotificationEntity) TableName() string {
	return "notification"
}
