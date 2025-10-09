package models

import (
	_models "common/models"
)

// FollowEntity đại diện cho bảng follow trong cơ sở dữ liệu
type FollowEntity struct {
	_models.BaseEntity
	ID          uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	FollowingID uint64 `gorm:"column:following_id;not null" json:"followingId"`
	// Following    AccountEntity `gorm:"foreignKey:FollowingID;references:ID" json:"following"`
	// CreatedBy uint64 `gorm:"column:created_by;not null" json:"createdBy"`
	// CreatedUser  AccountEntity `gorm:"foreignKey:CreatedBy;references:ID" json:"createdUser"`
	Status int `gorm:"default:1" json:"status"` // 1: Active, 2: Unfollowed
	// FollowedAt   time.Time  `gorm:"column:followed_at" json:"followedAt"`
	// UnfollowedAt *time.Time `gorm:"column:unfollowed_at" json:"unfollowedAt"`
}

// TableName để định nghĩa tên bảng trong cơ sở dữ liệu
func (FollowEntity) TableName() string {
	return "follow"
}
