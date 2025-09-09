package domain

import (
	_models "common/models"
)

// PostOrganization đại diện cho bảng liên kết tin đăng với organization
type PostOrganization struct {
	_models.BaseEntity
	OrganizationID uint64 `gorm:"not null;index" json:"organizationId"`
	PostID         uint64 `gorm:"not null;index" json:"postId"`
	IsOwner        bool   `gorm:"not null;default:false" json:"isOwner"`
	RoleID         uint64 `gorm:"not null;default:0" json:"roleId"`
}

func (PostOrganization) TableName() string {
	return "post_organization"
}
