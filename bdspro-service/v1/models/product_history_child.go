package models

import (
	"bdspro/enums"
	_models "common/models"

	"github.com/lib/pq"
)

// Product đại diện cho bảng sản phẩm

type ProductChildHistory struct {
	_models.BaseEntity
	ID        uint64               `gorm:"primaryKey;autoIncrement" json:"id"`
	Action    enums.ProductHistory `gorm:"int" json:"action"`
	ParentID  *uint64              `json:"-"`
	Parent    *ProductInfo         `gorm:"foreignKey:ParentID;references:ID" json:"parent"`
	ProductID *uint64              `json:"-"`
	Target    *ProductInfo         `gorm:"foreignKey:ProductID;references:ID" json:"product"`
	ChildIds  pq.Int64Array        `gorm:"type:bigint[];column:child_ids" json:"-"`
	Childs    []ProductInfo        `gorm:"-" json:"childs"`
}

type ProductInfo struct {
	ID          uint64  `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string  `gorm:"size:255" json:"name"`
	Code        string  `gorm:"size:15" json:"code"`
	CategoryID  *uint64 `json:"categoryId"`
	Area        float64 `json:"area"`
	LastPriceID *uint64 `json:"lastPriceId"`
	// Price       *ProductPrice `gorm:"foreignKey:LastPriceID;references:ID" json:"price"`
}

func (ProductChildHistory) TableName() string {
	return "product_child_history"
}

func (ProductInfo) TableName() string {
	return "products"
}
