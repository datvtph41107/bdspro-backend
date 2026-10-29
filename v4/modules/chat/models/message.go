package models

import (
	"time"
)

type ContentTypeEnum string

const (
	ContentTypeText    ContentTypeEnum = "text"
	ContentTypeImage   ContentTypeEnum = "image"
	ContentTypeVideo   ContentTypeEnum = "video"
	ContentTypeAudio   ContentTypeEnum = "audio"
	ContentTypePDF     ContentTypeEnum = "pdf"
	ContentTypeLink    ContentTypeEnum = "link"
	ContentTypeProduct ContentTypeEnum = "product"
	ContentTypeContact ContentTypeEnum = "contact"
	ContentTypeSystem  ContentTypeEnum = "system"

	ContentTypeLocation ContentTypeEnum = "location"
)

type MessageModel struct {
	ID             uint64                  `json:"id" gorm:"primaryKey;autoIncrement"`
	ConversationID uint64                  `json:"conversation_id" gorm:"index:idx_conversation_created"`
	SenderID       uint64                  `json:"sender_id"`
	Content        string                  `json:"content" gorm:"type:text"`
	ContentType    ContentTypeEnum         `json:"content_type" gorm:"type:varchar(50)"`
	FileURL        *string                 `json:"file_url" gorm:"type:varchar(500)"`
	ForwardFrom    *uint64                 `json:"forward_from"`
	DeletedAt      *time.Time              `gorm:"index" json:"-"`
	CreatedAt      time.Time               `json:"created_at" gorm:"index:idx_conversation_created;default:CURRENT_TIMESTAMP"`
	UpdatedAt      time.Time               `json:"updated_at" gorm:"index:autoUpdateTime"`
	Pined          bool                    `json:"pin" gorm:"default:false"`
	PinnedAt       *time.Time              `json:"pinned_at" gorm:"column:pinned_at"`
	Recall         bool                    `json:"recall" gorm:"default:false"`
	ReplyId        *uint64                 `json:"reply_id"`
	ReplyMessage   *MessageModel           `json:"reply_message" gorm:"foreignKey:ReplyId;references:ID"`
	Reactions      []*MessageReactionModel `json:"reactions" gorm:"foreignKey:MessageID;references:ID"`
	ExtraId        *uint64                 `json:"extra_id"`
	ReadAt         *time.Time              `json:"read_at" gorm:"-"`
	SystemMetadata *SystemMessageMetadata  `json:"system_metadata" gorm:"-"`
}

type SystemMessageMetadata struct {
	Action      string      `json:"action"`       // "pin", "add_member", "remove_member", etc.
	ActorID     uint64      `json:"actor_id"`     // Người thực hiện
	TargetIDs   []uint64    `json:"target_ids"`   // IDs của đối tượng bị tác động
	TargetNames []string    `json:"target_names"` // Tên hiển thị
	Timestamp   time.Time   `json:"timestamp"`
	Extra       interface{} `json:"extra,omitempty"` // Thông tin thêm
}

func (MessageModel) TableName() string {
	return "messages"
}
