package domain

import (
	"time"
	"tqd/internal/enums"
)

type UserNotification struct {
	ID               uint64                     `gorm:"primaryKey"`
	UserID           uint64                     `gorm:"column:user_id;not null;index:idx_notif_user_read,priority:1"`
	SubscriptionID   *uint64                    `gorm:"column:subscription_id"`
	Title            string                     `gorm:"column:title;type:varchar(255);not null"`
	Message          []string                   `gorm:"column:message;type:text[];not null"`
	NotificationType enums.NotificationType     `gorm:"column:notification_type;not null"`
	Severity         enums.NotificationSeverity `gorm:"column:severity;default:20"`
	ActionLink       *string                    `gorm:"column:action_link;type:text"`
	IsRead           bool                       `gorm:"column:is_read;default:false;index:idx_notif_user_read,priority:2"`
	ReadAt           *time.Time                 `gorm:"column:read_at"`
	CreatedAt        time.Time                  `gorm:"column:created_at;default:now()"`
	ExpiresAt        *time.Time                 `gorm:"column:expires_at;index:idx_notif_expires"`
}

func (UserNotification) TableName() string { return "user_notifications" }
