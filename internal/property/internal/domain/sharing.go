package domain

import (
	"bdspro/internal/enums"
	_models "common/models"
	"strings"
)

type Sharing struct {
	_models.BaseEntity
	TargetID    uint64            `gorm:"column:target_id" json:"targetId"`
	TargetType  enums.ETargetType `gorm:"column:target_type" json:"targetType"`
	OwnerID     uint64            `gorm:"column:owner_id" json:"ownerId"`
	OwnerType   enums.EOwnerOf    `gorm:"column:owner_type" json:"ownerType"`
	ShareID     uint64            `gorm:"column:share_id" json:"shareId"`
	ShareType   string            `gorm:"column:share_type" json:"shareType"`
	Permissions string            `gorm:"column:permissions" json:"permissions"`
}

func (Sharing) TableName() string {
	return "sharings"
}

// Chuyển đổi khi lưu vào DB
func (p *Sharing) SetFields(fields []string) {
	p.Permissions = strings.Join(fields, ",")
}

// Chuyển đổi khi đọc từ DB
func (p *Sharing) GetFields() []string {
	if p.Permissions == "" {
		return []string{}
	}
	return strings.Split(p.Permissions, ",")
}
