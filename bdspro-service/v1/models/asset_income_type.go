package models

type AssetIncomeType struct {
	TypeName string `json:"typeName" binding:"required"`
}

func (AssetIncomeType) TableName() string {
	return "asset_income_types"
}
