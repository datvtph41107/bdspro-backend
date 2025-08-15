package domain

import (
	_models "common/models"
	"time"
	"tx/internal/enums"
)

type TxAction struct {
	_models.BaseEntity
	TxId      uint64
	FromId    uint64
	FromOf    enums.EOwnerType
	ToId      uint64
	ToOf      enums.EOwnerType
	Action    enums.EAction
	Value     float64
	Note      string
	Message   string
	Timestamp *time.Time
}

func (TxAction) TableName() string {
	return "tb_transaction_action"
}
