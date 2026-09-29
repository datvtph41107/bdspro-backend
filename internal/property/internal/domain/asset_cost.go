package domain

import (
	"bdspro/internal/enums"
	_models "common/models"
	"time"
)

type AssetCost struct {
	_models.BaseEntity
	AssetID     uint64          `json:"assetId" binding:"required"`
	OwnerID     uint64          `json:"ownerId"`
	OwnerType   enums.EOwnerOf  `gorm:"default:10" json:"ownerType"`
	Type        enums.ECostType `gorm:"default:10" json:"type" binding:"required"`
	CostTypeID  uint64          `json:"costTypeId"`
	Amount      float64         `json:"amount" binding:"required"`
	Date        *time.Time      `json:"date"` // ngày chi
	Description string          `json:"description"`
	CostType    *AssetCostType  `gorm:"foreignKey:CostTypeID;references:ID" json:"costType"`
}

func (AssetCost) TableName() string {
	return "asset_costs"
}

// type AssetCostQuery struct {
// 	_models.BaseEntity
// 	AssetID      uint64          `json:"assetId"`
// 	Type         enums.ECostType `json:"type"`
// 	CostTypeID   uint64          `json:"costTypeId"`
// 	CostTypeName string          `gorm:"column:type_name" json:"costTypeName"`
// 	Amount       float64         `json:"amount"`
// 	Date         *time.Time      `json:"date"`
// 	Description  string          `json:"description"`
// }

// func (AssetCostQuery) TableName() string {
// 	return "asset_costs"
// }
