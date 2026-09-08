// internal/domain/building.go
package domain

import _models "common/models"

type Building struct {
	_models.BaseEntity
	Name    string `gorm:"size:255;not null" json:"name"`
	BlockID uint64 `gorm:"type:int8;not null;index" json:"blockId"`
	Block   *Block `gorm:"foreignKey:BlockID;references:ID" json:"block,omitempty"`
	// Có thể thêm số tầng, năm xây dựng...
}

func (Building) TableName() string {
	return "buildings"
}
