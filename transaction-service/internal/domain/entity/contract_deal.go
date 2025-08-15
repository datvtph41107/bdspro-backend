package entity

import (
	_models "common/models"
	"time"
	"tx/internal/domain"
	"tx/internal/enums"
)

type ContractDeal struct {
	_models.BaseEntity
	DealID        uint64                `gorm:"not null"`
	Type          enums.TransactionType `gorm:"not null"`
	Note          string                `gorm:"type:text"`
	ProductID     *uint64               `gorm:"index"`
	TransactionID *uint64               `gorm:"index"`
	CustomerID    uint64                `gorm:"index"`
	Transaction   *domain.Tx            `gorm:"foreignKey:TransactionID"`
	ApprovedBy    *uint64               `gorm:"index"`
	ApprovedAt    *time.Time
	RejectReason  string  `gorm:"type:text"`
	Amount        float64 `gorm:"-"`
}

func (ContractDeal) TableName() string {
	return "contract_deals"
}
