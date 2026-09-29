package domain

import (
	"bdspro/internal/enums"
	_models "common/models"
)

type ProductUser struct {
	_models.BaseEntity
	ProfileID uint64 `gorm:"column:profile_id;type:bigint;not null;index:idx_product_user_profile,priority:1"`
	ProductID uint64 `gorm:"column:product_id;type:bigint;not null;index:idx_product_user_profile,priority:2"`

	IsOwner         bool           `gorm:"column:is_owner;type:boolean;not null;default:false"`
	RoleID          uint64         `gorm:"column:role_id;type:bigint;not null;default:0"`
	DistributeID    *uint64        `gorm:"column:distribute_id;type:bigint;index:idx_product_user_distribute"`
	OwnerOf         enums.EOwnerOf `gorm:"column:owner_of;type:smallint;not null;default:10"`
	PriceID         *uint64        `gorm:"column:price_id;type:bigint"`
	OriginProfileID *uint64        `gorm:"column:origin_profile_id;type:bigint;index:idx_product_user_origin"`
}

func (ProductUser) TableName() string { return "product_user" }
