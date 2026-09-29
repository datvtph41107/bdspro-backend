package domain

import _models "common/models"

type AssetIncomeType struct {
	_models.BaseEntity
	TypeName string `json:"typeName"`
	IsCustom bool   `json:"isCustom"`
}
