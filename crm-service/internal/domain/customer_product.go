package domain

import (
	_models "common/models"
)

type CustomerProductEntity struct {
	_models.BaseEntity
	CustomerID uint64 `json:"customer_id"`
	ProductID  uint64 `json:"product_id"`
}
