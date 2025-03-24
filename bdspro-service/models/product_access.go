package models

import (
	"time"
)

type ProductAccess struct {
	UserId    uint64
	ProductId uint64
	DeletedAt *time.Time
}

// GORM table name override
func (ProductAccess) TableName() string {
	return "product_access"
}
