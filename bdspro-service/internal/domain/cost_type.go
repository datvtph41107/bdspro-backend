package domain

import (
	"bdspro/internal/enums"
	_models "common/models"
)

type AssetCostType struct {
	_models.BaseEntity
	OwnerID     uint64          `gorm:"column:owner_id" json:"ownerId"`
	OwnerType   enums.EOwnerOf  `gorm:"column:owner_type;default:10" json:"ownerType"`
	TypeName    string          `gorm:"column:type_name" json:"typeName" binding:"required"`
	Type        enums.ECostType `gorm:"column:type;default:10" json:"type" binding:"required"`
	IsCustom    bool            `gorm:"column:is_custom" json:"isCustom"`
	Description string          `gorm:"column:description" json:"description"`
}

func (AssetCostType) TableName() string {
	return "asset_cost_types"
}
