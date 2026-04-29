package qh_domain

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

type QHRegion struct {
	ID          uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string `gorm:"type:varchar(255);not null" json:"name"`
	DisplayName string `gorm:"type:varchar(255);not null" json:"display_name"`
	Description string `gorm:"type:text" json:"description"`
	// Type        enums.RegionType   `gorm:"type:int;default:1;index" json:"type"`
	// Status     enums.RegionStatus `gorm:"type:int;default:10;index" json:"status"`
	Geometry RawGeometry `gorm:"type:geometry(MultiPolygon,4326)"`
	// Centroid   string             `gorm:"type:text" json:"centroid"`
	// BBox       string             `gorm:"type:text" json:"bbox"`
	Area       float64   `gorm:"type:float" json:"area"`
	Length     float64   `gorm:"type:float" json:"length"`
	LayerID    *uint64   `gorm:"index" json:"layer_id,omitempty"`
	Properties JSONMap   `gorm:"type:jsonb" json:"properties"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (QHRegion) TableName() string {
	return "qh_regions"
}

type JSONMap map[string]interface{}

func (j JSONMap) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

func (j *JSONMap) Scan(value interface{}) error {
	if value == nil {
		*j = make(JSONMap)
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, j)
}
