package qh_domain

import "time"

// QHLayerLandUse bảng nối many-to-many giữa qh_layers và qh_land_use.
type QHLayerLandUse struct {
	LayerID   uint64    `gorm:"primaryKey;column:layer_id" json:"layerId"`
	LandUseID uint64    `gorm:"primaryKey;column:land_use_id" json:"landUseId"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
}

func (QHLayerLandUse) TableName() string {
	return "qh_layer_land_use"
}
