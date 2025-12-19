package domain

import (
	"bdspro/internal/enums"
	_models "common/models"
	"time"
)

type TxDealCost struct {
	_models.BaseEntity
	CostName    string
	DealID      uint64
	Amount      float64
	ConfirmedAt *time.Time
	CostTypeID  uint64
	CostType    *TxCostType `gorm:"foreignKey:CostTypeID"`
	Status      enums.TxApprovedStatus
	Note        string
}

func (TxDealCost) TableName() string {
	return "tx_deal_costs"
}
