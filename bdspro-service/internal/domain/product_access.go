package domain

import (
	"bdspro/internal/enums"
	_models "common/models"
	"strings"
)

type ProductAccess struct {
	// UserId     uint64
	_models.BaseEntity
	Product    *ProductItem      `json:"product"`
	ProductID  uint64            `json:"productId"`
	Fields     string            `json:"fields"`
	TargetType enums.ETargetType `json:"targetType"` // 1: cá nhân, 2: nhóm, 3: tổ chức
	TargetId   uint64            `json:"targetId"`
	Commission float32           `json:"commission"`
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

type ProductAccessQuery struct {
	// UserId     uint64
	// _models.BaseEntity
	// Product    *ProductItem `json:"product"`
	ID         uint64  `gorm:"column:id" json:"id"`
	ProductID  uint64  `json:"productId"`
	Fields     string  `json:"fields"`
	TargetType uint    `json:"targetType"` // 1: cá nhân, 2: nhóm, 3: tổ chức
	TargetId   uint64  `json:"targetId"`
	Commission float32 `json:"commission"`
	TargetName string  `gorm:"column:target_name" json:"targetName"`
	Avatar     string  `gorm:"column:avatar" json:"avatar"`
}
