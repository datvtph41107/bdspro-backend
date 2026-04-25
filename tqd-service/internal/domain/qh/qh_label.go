package qh_domain

import _models "common/domain/entity"

type QHLabel struct {
	_models.BaseEntity
	LayerID uint64 `gorm:"type:int" json:"layerId"`
	Name    string `gorm:"type:varchar(255)" json:"name"`
	Color   string ``
}

func (QHLabel) TableName() string {
	return "qh_label"
}
