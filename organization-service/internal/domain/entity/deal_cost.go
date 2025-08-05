package entity

import (
	_models "common/models"
	"organization/internal/enums"
	"time"
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
