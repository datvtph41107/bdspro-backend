package models

import (
	"time"
)

// Order đại diện cho thông tin đơn hàng
// @swagger:model
type ParticipantModel struct {
	ID        uint64 `gorm:"primaryKey" json:"id"`
	GroupId   uint64
	CreatedAt time.Time `gorm:"column:create_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:update_at;autoUpdateTime" json:"updated_at"`
	CreatedBy uint64
	MemberId  uint64
}

func (ParticipantModel) TableName() string {
	return "participant"
}

// func New(msg *chatpb.ChatMessage) MessageModel {
// 	return MessageModel{
// 		Message: msg.Message,
// 		From:    msg.From,
// 		To:      msg.To,
// 	}
// }
