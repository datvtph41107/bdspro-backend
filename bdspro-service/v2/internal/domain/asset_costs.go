package domain

import (
	_models "common/models"
	"time"
)

type AssetCost struct {
	_models.BaseEntity
	AssetID     uint64     `json:"assetId"`
	CostType    uint64     `json:"costType"` // Tham chiếu asset_cost_types
	Amount      float64    `json:"amount"`
	Date        *time.Time `json:"date"` // ngày chi
	Description string     `json:"description"`
}

func (AssetCost) TableName() string {
	return "asset_costs"
}
