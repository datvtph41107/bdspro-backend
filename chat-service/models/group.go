package models

import (
	_db "common/db"
	"time"
)

// Order đại diện cho thông tin đơn hàng
// @swagger:model
type GroupModel struct {
	ID          uint64    `gorm:"primaryKey" json:"id"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	CreatedBy   string
	Title       string
	Avatar      string
	LastMessage string
	LastUser    string
}

func (GroupModel) TableName() string {
	return "group"
}

func GetGroupsByMemberID(memberId uint64, page, pageSize int) ([]GroupModel, error) {
	var groups []GroupModel

	offset := (page - 1) * pageSize

	err := _db.DB.Table("group_tb").
		Select("group_tb").
		Joins("JOIN group_member ON group_tb.id = group_member.group_id").
		Where("group_member.member_id = ?", memberId).
		Order("group_tb.updated_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&groups).Error

	if err != nil {
		return nil, err
	}

	return groups, nil
}
