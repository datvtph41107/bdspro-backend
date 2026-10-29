package models

import (
	"time"

	"github.com/lib/pq"
)

type ConversationTypeEnum int

const (
	ConversationTypePrivate ConversationTypeEnum = 1
	ConversationTypeGroup   ConversationTypeEnum = 2
	ConversationTypeChannel ConversationTypeEnum = 3
)

type ConversationModel struct {
	ID         uint64               `gorm:"primaryKey;autoIncrement"`
	Type       ConversationTypeEnum `gorm:"not null;index:idx_type_created"`
	Name       *string              `gorm:"size:255"`
	CreatedBy  uint64               `gorm:"not null;column:created_by"`
	ReceiverId *uint64              `gorm:"column:receiver_id"`

	IsBroadcast       bool                  `gorm:"default:false;column:is_broadcast"`
	ForbidForward     bool                  `gorm:"default:false;column:forbid_forward"`
	CreatedAt         time.Time             `gorm:"autoCreateTime;column:created_at;index:idx_type_created"`
	UpdatedAt         time.Time             `gorm:"autoUpdateTime;column:updated_at"`
	DeletedAt         *time.Time            `gorm:"index" json:"-"`
	Participants      []ParticipantModel    `gorm:"foreignKey:ConversationID;references:ID"`
	Avatar            string                `gorm:"size:255"`
	BackgroundImage   *BackgroundImageModel `gorm:"foreignKey:BackgroundImageID;references:ID"`
	BackgroundImageID *uint64               `gorm:"column:background_image_id"`
	// Color           string             `gorm:"size:255"`
	UnreadCount    int          `gorm:"-"`
	LatestMessage  MessageModel `gorm:"-"`
	ParticipantIDs []uint64     `gorm:"-"`
}

func (ConversationModel) TableName() string {
	return "conversations"
}

type ConversationWithMessage struct {
	ConversationModel

	LatestMessageID          uint64        `gorm:"column:latest_message_id"`
	LatestMessageContent     string        `gorm:"column:latest_message_content"`
	LatestMessageContentType string        `gorm:"column:latest_message_content_type"`
	LatestMessageFileURL     *string       `gorm:"column:latest_message_file_url"`
	LatestMessageForwardFrom *uint64       `gorm:"column:latest_message_forward_from"`
	LatestMessageSenderID    uint64        `gorm:"column:latest_message_sender_id"`
	LatestMessageCreatedAt   time.Time     `gorm:"column:latest_message_created_at"`
	LatestMessageDeletedAt   *time.Time    `gorm:"column:latest_message_deleted_at"`
	LatestMessagePined       bool          `gorm:"column:latest_message_pined"`
	LatestMessageRecall      bool          `gorm:"column:latest_message_recall"`
	UnreadCount              int           `gorm:"column:unread_count"`
	ParticipantIDs           pq.Int64Array `gorm:"type:bigint[];column:participant_ids"`
}
