package domain

import (
	_models "common/models"
	"time"
	"tx/internal/enums"
)

type Tx struct {
	_models.BaseEntity
	TransactionName string
	FromID          uint64
	FromOf          enums.EOwnerType
	ToID            uint64
	ToOf            enums.EOwnerType
	Method          enums.EMethod
	Value           float64
	Status          enums.TxStatus
	Timestamp       time.Time
	Note            string `gorm:"-"`
}

func (Tx) TableName() string {
	return "tx_transaction"
}
