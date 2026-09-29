package domain

import (
	_models "common/models"
)

// AssetUser đại diện cho bảng liên kết tài sản với user
type AssetUser struct {
	_models.BaseEntity
	ProfileID uint64 `gorm:"not null;index" json:"profileId"`
	AssetID   uint64 `gorm:"not null;index" json:"assetId"`
	IsOwner   bool   `gorm:"not null;default:false" json:"isOwner"`
	RoleID    uint64 `gorm:"not null;default:0" json:"roleId"`
}

func (AssetUser) TableName() string {
	return "asset_user"
}
