package domain

import (
	_models "common/models"
)

// AssetOrganization đại diện cho bảng liên kết tài sản với organization
type AssetOrganization struct {
	_models.BaseEntity
	OrganizationID uint64 `gorm:"not null;index" json:"organizationId"`
	AssetID        uint64 `gorm:"not null;index" json:"assetId"`
	IsOwner        bool   `gorm:"not null;default:false" json:"isOwner"`
	RoleID         uint64 `gorm:"not null;default:0" json:"roleId"`
}

func (AssetOrganization) TableName() string {
	return "asset_organization"
}
