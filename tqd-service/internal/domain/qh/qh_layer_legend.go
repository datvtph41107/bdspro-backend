package qh_domain

import _models "common/domain/entity"

// QHLayerLegend — chú giải: gắn label với layer kèm ghi chú.
type QHLayerLegend struct {
	_models.BaseEntity
	LayerID uint64 `gorm:"column:layer_id;not null;index" json:"layerId"`
	LabelID uint64 `gorm:"column:label_id;not null;index" json:"labelId"`
	Note    string `gorm:"type:text" json:"note,omitempty"`
}

func (QHLayerLegend) TableName() string {
	return "qh_layer_legends"
}
