package domain

import _models "common/models"

type Deposite struct {
	_models.BaseEntity
	ProductID     uint64 `json:"productId"`
	CustomerName  string `json:"customerName"`
	CustomerPhone string `json:"customerPhone"`
	Amount        uint64 `json:"amount"`
	Note          string `json:"note"`
}

func (Deposite) TableName() string {
	return "deposite"
}
