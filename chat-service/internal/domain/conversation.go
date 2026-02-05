package domain

import (
	"chat/internal/enums"
	"time"

	"github.com/lib/pq"
)

type Conversation struct {
	ID              uint64                 `gorm:"primaryKey;autoIncrement"`
	Type            enums.ConversationType `gorm:"type:smallint;not null;index"`
	Name            *string                `gorm:"type:varchar(255);index"`
	Status          enums.Status           `gorm:"type:smallint;not null;index;default:1"`
	AvartarURL      *string                `gorm:"type:varchar(255)"`
	Theme           *string                `gorm:"type:varchar(100)"`
	EmotIcon        *string                `gorm:"type:varchar(50)"`  // owner
	ReferralLink    *string                `gorm:"type:varchar(255)"` // link tham gia nhóm
	QrReferralImage *string                `gorm:"type:varchar(255)"`
	EnableApproval  bool                   `gorm:"not null;default:false"` // phê duyệt tham gia
	CreatedBy       uint64                 `gorm:"not null;index"`         // Owner
	CreatedAt       time.Time              `gorm:"not null;autoCreateTime"`
	UpdatedAt       time.Time              `gorm:"not null;autoUpdateTime"`
}

func (Conversation) TableName() string {
	return "conversations"
}

type ConversationWithMessage struct {
	Conversation
	LatestMessageID        uint64
	LatestMessageType      uint32
	LatestMessageContent   *string
	LatestMessageSenderID  *uint64
	LatestMessageCreatedAt time.Time
	UnreadCount            int
	ParticipantIDs         pq.Int64Array `gorm:"type:bigint[]"`
}
