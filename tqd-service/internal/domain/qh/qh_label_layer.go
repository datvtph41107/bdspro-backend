package qh_domain

import "time"

// QHLabelLayer bảng nối many-to-many giữa qh_labels và qh_layers.
// Cột layer_id trên qh_labels được giữ tạm để tương thích logic cũ.
type QHLabelLayer struct {
	LabelID   uint64    `gorm:"primaryKey;column:label_id" json:"labelId"`
	LayerID   uint64    `gorm:"primaryKey;column:layer_id" json:"layerId"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
}

func (QHLabelLayer) TableName() string {
	return "qh_label_layers"
}
