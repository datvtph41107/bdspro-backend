package models

type AssetCostType struct {
	ID       uint64 `json:"id"`
	TypeName string `json:"typeName" binding:"required"`
	Type     uint   `gorm:"default:10" json:"type" binding:"required"`
}

func (AssetCostType) TableName() string {
	return "asset_cost_types"
}
