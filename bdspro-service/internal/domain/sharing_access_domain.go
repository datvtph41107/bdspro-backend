package domain

import (
	"bdspro/internal/enums"
	_models "common/models"
	"strings"
)

type SharingAccess struct {
	_models.BaseEntity
	// Product     *ProductItem        `gorm:"foreignKey:ProductID" json:"product"`
	DomainID    uint64              `json:"productId"`
	Domain      enums.EDomainAccess `gorm:"column:domain;default:10" json:"domain"`
	FromType    enums.EOwnerOf      `gorm:"column:from_type;default:10" json:"fromType"`
	FromId      uint64              `json:"fromId"`
	ToType      enums.EOwnerOf      `json:"toType"`
	ToId        uint64              `json:"toId"`
	Fields      string              `json:"fields"`
	Commission  float32             `json:"commission"`
	Permissions string              `json:"permissions"`
}

// GORM table name override
func (SharingAccess) TableName() string {
	return "sharing_access"
}

// Chuyển đổi khi lưu vào DB
func (p *SharingAccess) SetFields(fields []string) {
	p.Fields = strings.Join(fields, ",")
}

// Chuyển đổi khi đọc từ DB
func (p *SharingAccess) GetFields() []string {
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
	DomainID   uint64  `json:"productId"`
	Fields     string  `json:"fields"`
	Commission float32 `json:"commission"`
	ToId       uint64  `json:"targetId"`
	ToType     uint    `json:"targetType"` // 1: cá nhân, 2: nhóm, 3: tổ chức
	ToName     string  `gorm:"column:target_name" json:"targetName"`
	Avatar     string  `gorm:"column:avatar" json:"avatar"`
}
