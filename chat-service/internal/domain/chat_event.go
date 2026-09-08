package domain

import "chat/internal/enums"

type ChatEvent struct {
	ID             uint64              `gorm:"primaryKey"`
	ConversationID uint64              `gorm:"not null;index:idx_event,priority:1"`
	Sequence       uint64              `gorm:"not null;index:idx_event,priority:2"`
	Timestamp      int64               `gorm:"not null;index"`
	EventType      enums.ChatEventType `gorm:"not null;index"` // message_create,edit,revoke,pin,system...
	ActorID        uint64              `gorm:"not null;index"`
	MessageID      *uint64             `gorm:"index"` // nếu event liên quan message
	FileID         *uint64
	ReplyTo        *uint64
	ForwardFrom    *uint64
	EventText      *string
}

func (ChatEvent) TableName() string {
	return "chat_events"
}
