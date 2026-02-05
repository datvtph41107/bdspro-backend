package domain

import "chat/internal/enums"

type ChatTimeline struct {
	ID             uint64 `gorm:"primaryKey;index:idx_msg_lookup,priority:2"`
	ConversationID uint64 `gorm:"index:idx_msg_lookup,priority:1"`
	Timestamp      int64  `gorm:"column: timestamp;not null;index:idx_timeline,priority:1,sort:desc"`
	Sequence       uint64 `gorm:"not null;index:idx_timeline,priority:2,sort:desc"` // indexKey
	SendFromID     uint64
	Type           enums.MessageType
	Text           string
	FileID         *uint64
	ReplyTo        *uint64 `gorm:"index"`
	ForwardFrom    *uint64 `gorm:"index"`
	IsPinned       bool    `gorm:"default:false;index"`
	IsDeleted      bool    `gorm:"default:false"`
	IsEdited       bool    `gorm:"default:false"`
}

func (ChatTimeline) TableName() string {
	return "chat_timelines"
}
