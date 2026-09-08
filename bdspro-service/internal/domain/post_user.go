package domain

import (
	_models "common/models"
)

// PostUser đại diện cho bảng liên kết tin đăng với user
type PostUser struct {
	_models.BaseEntity
	ProfileID uint64 `gorm:"not null;index" json:"profileId"`
	PostID    uint64 `gorm:"not null;index" json:"postId"`
	IsOwner   bool   `gorm:"not null;default:false" json:"isOwner"`
	RoleID    uint64 `gorm:"not null;default:0" json:"roleId"`
}

func (PostUser) TableName() string {
	return "post_user"
}
