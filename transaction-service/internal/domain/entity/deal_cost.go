package entity

import (
	_models "common/models"
	"time"
	"tx/internal/enums"
)

type DealCost struct {
	_models.BaseEntity
	CostName    string
	DealID      uint64
	Amount      float64
	ConfirmedAt *time.Time
	CostTypeID  uint64
	CostType    *CostType `gorm:"foreignKey:CostTypeID"`
	Status      enums.ApprovedStatus
	Note        string
}

func (DealCost) TableName() string {
	return "deal_costs"
}
