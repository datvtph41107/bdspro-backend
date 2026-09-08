package domain

import (
	"bdspro/internal/enums"
	_models "common/models"
	"time"
)

type TxContractDeal struct {
	_models.BaseEntity
	DealID        uint64                  `gorm:"not null"`
	Type          enums.TxTransactionType `gorm:"not null"`
	Note          string                  `gorm:"type:text"`
	ProductID     *uint64                 `gorm:"index"`
	TransactionID *uint64                 `gorm:"index"`
	CustomerID    uint64                  `gorm:"index"`
	Transaction   *Tx                     `gorm:"foreignKey:TransactionID"`
	ApprovedBy    *uint64                 `gorm:"index"`
	ApprovedAt    *time.Time
	RejectReason  string  `gorm:"type:text"`
	Amount        float64 `gorm:"-"`
}

func (TxContractDeal) TableName() string {
	return "tx_contract_deals"
}
