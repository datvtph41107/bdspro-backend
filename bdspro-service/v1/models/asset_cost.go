package models

import (
	_models "common/models"
	"time"
)

type AssetCost struct {
	AssetID     uint64     `json:"assetId" binding:"required"`
	Type        uint       `json:"type" binding:"required"`
	CostTypeID  uint64     `json:"costTypeId" binding:"required"`
	Amount      float64    `json:"amount" binding:"required"`
	Date        *time.Time `json:"date"` // ngày chi
	Description string     `json:"description"`
}

func (AssetCost) TableName() string {
	return "asset_costs"
}

type AssetCostQuery struct {
	_models.BaseEntity
	AssetID      uint64     `json:"assetId"`
	Type         uint       `json:"type"`
	CostTypeID   uint64     `json:"costTypeId"`
	CostTypeName string     `gorm:"column:type_name" json:"costTypeName"`
	Amount       float64    `json:"amount"`
	Date         *time.Time `json:"date"`
	Description  string     `json:"description"`
}

func (AssetCostQuery) TableName() string {
	return "asset_costs"
}
