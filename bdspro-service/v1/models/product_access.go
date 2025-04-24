package models

import (
	_models "common/models"
	"strings"
)

type ProductAccess struct {
	// UserId     uint64
	_models.BaseEntity
	Product    *ProductItem `json:"product"`
	ProductID  uint64       `json:"productId"`
	Fields     string       `json:"fields"`
	TargetType uint         `json:"targetType"`
	TargetId   uint64       `json:"targetId"`
}

// GORM table name override
func (ProductAccess) TableName() string {
	return "product_access"
}

// Chuyển đổi khi lưu vào DB
func (p *ProductAccess) SetFields(fields []string) {
	p.Fields = strings.Join(fields, ",")
}

// Chuyển đổi khi đọc từ DB
func (p *ProductAccess) GetFields() []string {
	if p.Fields == "" {
		return []string{}
	}
	return strings.Split(p.Fields, ",")
}
