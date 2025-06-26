package models

import "time"

type ContentTypeEnum string

const (
	ContentTypeText  ContentTypeEnum = "text"
	ContentTypeImage ContentTypeEnum = "image"
	ContentTypeVideo ContentTypeEnum = "video"
	ContentTypeAudio ContentTypeEnum = "audio"
	ContentTypePDF   ContentTypeEnum = "pdf"
	ContentTypeLink  ContentTypeEnum = "link"
)

type MessageModel struct {
	ID             uint64          `json:"id" gorm:"primaryKey;autoIncrement"`
	ConversationID uint64          `json:"conversation_id" gorm:"index:idx_conversation_created"`
	SenderID       uint64          `json:"sender_id"`
	Content        string          `json:"content" gorm:"type:text"`
	ContentType    ContentTypeEnum `json:"content_type" gorm:"type:varchar(50)"`
	FileURL        *string         `json:"file_url" gorm:"type:varchar(500)"`
	ForwardFrom    *uint64         `json:"forward_from"`
	IsDeleted      bool            `json:"is_deleted" gorm:"default:false"`
	CreatedAt      time.Time       `json:"created_at" gorm:"index:idx_conversation_created;default:CURRENT_TIMESTAMP"`
	Pined          bool            `json:"pin" gorm:"default:false"`
	Recall         bool            `json:"recall" gorm:"default:false"`
	ReplyId        *uint64         `json:"reply_id"`
	ReplyMessage   *MessageModel   `json:"reply_message" gorm:"foreignKey:ReplyId;references:ID"`
	Reactions      []*MessageReactionModel `json:"reactions" gorm:"foreignKey:MessageID;references:ID"`
}

func (MessageModel) TableName() string {
	return "messages"
}
