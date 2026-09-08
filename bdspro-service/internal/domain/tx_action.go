package domain

import (
	"bdspro/internal/enums"
	_models "common/models"
	"time"
)

type TxAction struct {
	_models.BaseEntity
	TxId      uint64 `gorm:"column:tx_id"`
	FromId    uint64
	FromOf    enums.TxOwnerType
	ToId      uint64
	ToOf      enums.TxOwnerType
	Action    enums.TxAction
	Value     float64
	Note      string
	Message   string
	Timestamp *time.Time
}

func (TxAction) TableName() string {
	return "tx_transaction_action"
}
