package models

import (
	"time"
)

type ConversationTypeEnum int

const (
	ConversationTypePrivate ConversationTypeEnum = 1
	ConversationTypeGroup   ConversationTypeEnum = 2
	ConversationTypeChannel ConversationTypeEnum = 3
)

type ConversationModel struct {
	ID              uint64               `gorm:"primaryKey;autoIncrement"`
	Type            ConversationTypeEnum `gorm:"not null;index:idx_type_created"`
	Name            *string              `gorm:"size:255"`
	CreatedBy       uint64               `gorm:"not null;column:created_by"`
	IsBroadcast     bool                 `gorm:"default:false;column:is_broadcast"`
	ForbidForward   bool                 `gorm:"default:false;column:forbid_forward"`
	CreatedAt       time.Time            `gorm:"autoCreateTime;column:created_at;index:idx_type_created"`
	UpdatedAt       time.Time            `gorm:"autoUpdateTime;column:updated_at"`
	Participants    []ParticipantModel   `gorm:"foreignKey:ConversationID;references:ID"`
	IsDeleted       bool                 `json:"is_deleted" gorm:"default:false"`
	Avatar          string               `gorm:"size:255"`
	BackgroundImage string               `gorm:"size:255"`
	UnreadCount     int                  `gorm:"-"`
	LatestMessage   MessageModel         `gorm:"-"`
}

func (ConversationModel) TableName() string {
	return "conversations"
}

type ConversationWithMessage struct {
	ConversationModel

	LatestMessageID          uint64    `gorm:"column:latest_message_id"`
	LatestMessageContent     string    `gorm:"column:latest_message_content"`
	LatestMessageContentType string    `gorm:"column:latest_message_content_type"`
	LatestMessageFileURL     *string   `gorm:"column:latest_message_file_url"`
	LatestMessageForwardFrom *uint64   `gorm:"column:latest_message_forward_from"`
	LatestMessageSenderID    uint64    `gorm:"column:latest_message_sender_id"`
	LatestMessageCreatedAt   time.Time `gorm:"column:latest_message_created_at"`
	LatestMessageIsDeleted   bool      `gorm:"column:latest_message_is_deleted"`
	LatestMessagePined       bool      `gorm:"column:latest_message_pined"`
	LatestMessageRecall      bool      `gorm:"column:latest_message_recall"`
	UnreadCount              int       `gorm:"column:unread_count"`
}
