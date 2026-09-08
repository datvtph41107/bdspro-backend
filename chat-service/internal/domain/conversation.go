package domain

import (
	"chat/internal/enums"
	_models "common/domain/entity"
	"github.com/lib/pq"
	"time"
)

type Conversation struct {
	_models.BaseEntity
	Type              enums.ConversationType `gorm:"type:smallint;not null;index"`
	Name              *string                `gorm:"type:varchar(255);index"`
	CreatedBy         uint64                 `gorm:"not null;index"`
	ReceiverId        *uint64                `gorm:"column:receiver_id"`
	IsBroadcast       bool                   `gorm:"default:false;column:is_broadcast"`
	ForbidForward     bool                   `gorm:"default:false;column:forbid_forward"`
	Avatar            string                 `gorm:"size:255"`
	BackgroundImage   *BackgroundImage       `gorm:"foreignKey:BackgroundImageID;references:ID"`
	BackgroundImageID *uint64                `gorm:"column:background_image_id"`
	Participants      []Participant          `gorm:"foreignKey:ConversationID;references:ID"`
	UnreadCount       int                    `gorm:"-"`
	LatestMessage     Message                `gorm:"-"`
	ParticipantIDs    []uint64               `gorm:"-"`
}

func (Conversation) TableName() string {
	return "conversations"
}

type ConversationWithMessage struct {
	Conversation
	LatestMessageID          uint64
	LatestMessageContent     string
	LatestMessageContentType string
	LatestMessageFileURL     *string
	LatestMessageForwardFrom *uint64
	LatestMessageSenderID    uint64
	LatestMessageCreatedAt   time.Time
	LatestMessageDeletedAt   *time.Time
	LatestMessagePined       bool
	LatestMessageRecall      bool
	UnreadCount              int
	ParticipantIDs           pq.Int64Array `gorm:"type:bigint[]"`
}

type ConversationTypeEnum = enums.ConversationType

const (
	ConversationTypePrivate = enums.CONVERSATION_TYPE_PRIVATE
	ConversationTypeGroup   = enums.CONVERSATION_TYPE_GROUP
)
