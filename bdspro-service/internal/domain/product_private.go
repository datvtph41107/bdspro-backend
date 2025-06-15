package domain

import (
	_models "common/models"
	"strings"
)

// type PrivateFieldEnum = [""]

type ProductPrivate struct {
	_models.BaseEntity
	ProductId     uint64   `json:"productId,omitempty"`
	ImportPrice   *float64 `json:"importPrice,omitempty"`
	OperatingCost *float64 `gorm:"type:DECIMAL(15,0)" json:"operatingCost,omitempty"`
	InternalNote  string   `json:"internalNote,omitempty"`
	TargetProfit  *float64 `json:"targetProfit,omitempty"`
	PrivateDocs   string   `gorm:"type:text" json:"privateDocs,omitempty"`
}

// Chuyển đổi khi lưu vào DB
func (p *ProductPrivate) SetPrivateDocs(docs []string) {
	p.PrivateDocs = strings.Join(docs, ",")
}

// Chuyển đổi khi đọc từ DB
func (p *ProductPrivate) GetPrivateDocs() []string {
	if p.PrivateDocs == "" {
		return []string{}
	}
	return strings.Split(p.PrivateDocs, ",")
}

func (ProductPrivate) TableName() string {
	return "product_private"
}
