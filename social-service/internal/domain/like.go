package domain

import (
	"social/internal/enums"
	"time"
)

type Like struct {
	TargetID   uint64           `gorm:"primaryKey" json:"targetId"`
	TargetType enums.TargetType `gorm:"primaryKey" json:"targetType"`
	UserID     uint64           `gorm:"primaryKey" json:"userId"`
	DisLike    bool             `gorm:"default:false" json:"dislike"`
	CreatedAt  *time.Time       `json:"createdAt"`
	DeletedAt  *time.Time       `json:"deletedAt"`
}

func (Like) TableName() string {
	return "tb_like"
}
