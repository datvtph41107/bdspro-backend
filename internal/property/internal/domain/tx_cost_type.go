package domain

import (
	"bdspro/internal/enums"
	_models "common/models"
)

type TxCostType struct {
	_models.BaseEntity
	TypeName    string
	IsCustom    bool
	CostType    enums.TxCostType
	Description string
}

func (TxCostType) TableName() string {
	return "tx_cost_types"
}
