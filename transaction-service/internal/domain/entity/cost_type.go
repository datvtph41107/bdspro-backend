package entity

import (
	_models "common/models"
	"tx/internal/enums"
)

type CostType struct {
	_models.BaseEntity
	TypeName    string
	IsCustom    bool
	CostType    enums.ECostType
	Description string
}

func (CostType) TableName() string {
	return "cost_types"
}
