package domain

import (
	"time"
	"tqd/internal/domain/jsonb"
)

type Feature struct {
	ID          uint64      `json:"id" gorm:"primaryKey;autoIncrement"`
	SessionID   *int        `json:"session_id" gorm:"column:session_id"`
	Layer       string      `json:"layer" gorm:"type:text;not null"`
	Geometry    string      `json:"geometry" gorm:"column:geometry;type:geometry(Geometry,4326)"` // GeoJSON text (from ST_AsGeoJSON)
	Properties  jsonb.JSONB `json:"properties" gorm:"type:jsonb;not null;default:'{}'"`
	Tile        string      `json:"tile" gorm:"type:text"`
	Extent      *int        `json:"extent"`
	FeatureHash string      `json:"feature_hash" gorm:"type:text;uniqueIndex"`
	CreatedAt   time.Time   `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time   `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName overrides the table name
func (Feature) TableName() string {
	return "features"
}
