package domain

import (
	_models "common/models"
)

// ProductOrganization đại diện cho bảng liên kết sản phẩm với organization
type ProductOrganization struct {
	_models.BaseEntity
	OrganizationID uint64 `gorm:"not null;index" json:"organizationId"`
	ProductID      uint64 `gorm:"not null;index" json:"productId"`
	IsOwner        bool   `gorm:"not null;default:false" json:"isOwner"`
	RoleID         uint64 `gorm:"not null;default:0" json:"roleId"`
}

func (ProductOrganization) TableName() string {
	return "product_organization"
}
