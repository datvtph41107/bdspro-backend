package domain

import (
	"chat/internal/enums"
	"time"
)

type Membership struct {
	UserID                 uint64     `gorm:"primaryKey;autoIncrement:false;index:idx_user_conversation,unique"`
	ConversationID         uint64     `gorm:"not null;index:idx_user_conversation,unique;index:idx_conversation_status"`
	Nickname               *string    `gorm:"type:varchar(100)"`
	Role                   enums.Role `gorm:"type:smallint;not null"`
	LastReadSeqID          uint64     `gorm:"not null;default:0"`
	LastReadTs             int64
	IsArchivedConversation bool         `gorm:"not null;default:false"`
	EnableNotification     bool         `gorm:"not null;default:true"`
	NotificationOffUntil   *int64       `gorm:"column:notification_mute_until"`
	Status                 enums.Status `gorm:"type:smallint;not null;default:1;index:idx_conversation_status"`
	JoinedAt               time.Time    `gorm:"column:joined_at"`
	CreatedAt              time.Time    `gorm:"not null;autoCreateTime"`
	UpdatedAt              time.Time    `gorm:"not null;autoUpdateTime"`
}

func (Membership) TableName() string {
	return "conversation_memberships"
}
