package domain

import (
	_models "common/models"
)

// ProductUser đại diện cho bảng liên kết sản phẩm với user
type ProductUser struct {
	_models.BaseEntity
	ProfileID uint64 `gorm:"not null;index" json:"profileId"`
	ProductID uint64 `gorm:"not null;index" json:"productId"`
	IsOwner   bool   `gorm:"not null;default:false" json:"isOwner"`
	RoleID    uint64 `gorm:"not null;default:0" json:"roleId"`
}

func (ProductUser) TableName() string {
	return "product_user"
}
