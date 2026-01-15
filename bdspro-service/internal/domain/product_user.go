package domain

import (
	"bdspro/internal/enums"
	_models "common/models"
)

// ProductUser đại diện cho bảng liên kết sản phẩm với user
type ProductUser struct {
	_models.BaseEntity
	ProfileID    uint64         `gorm:"index" json:"profileId"`
	ProductID    uint64         `gorm:"not null;index" json:"productId"`
	IsOwner      bool           `gorm:"not null;default:false" json:"isOwner"`
	RoleID       uint64         `gorm:"default:0" json:"roleId"`
	DistributeID *uint64        `gorm:"column:distribute_id;index" json:"distributeId"`
	OwnerOf      enums.EOwnerOf `gorm:"default:10" json:"ownerOf"`

	OriginProfileID *uint64 `gorm:"index" json:"originProfileId"`
}

func (ProductUser) TableName() string {
	return "product_user"
}
