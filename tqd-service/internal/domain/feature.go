package domain

import (
	"database/sql/driver"
	"fmt"
	"time"
)

type Feature struct {
	ID          uint64    `json:"id" gorm:"primaryKey;autoIncrement"`
	SessionID   *int      `json:"session_id" gorm:"column:session_id"`
	Layer       string    `json:"layer" gorm:"type:text;not null"`
	Geometry    string    `json:"geometry" gorm:"column:geometry;type:geometry(Geometry,4326)"` // GeoJSON text (from ST_AsGeoJSON)
	Properties  JSONB     `json:"properties" gorm:"type:jsonb;not null;default:'{}'"`
	Tile        string    `json:"tile" gorm:"type:text"`
	Extent      *int      `json:"extent"`
	FeatureHash string    `json:"feature_hash" gorm:"type:text;uniqueIndex"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName overrides the table name
func (Feature) TableName() string {
	return "features"
}

// JSONB is a custom type for JSONB fields (using gorm's datatypes.JSON)
type JSONB []byte

// Scan implements sql.Scanner
func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal JSONB value: %v", value)
	}
	*j = bytes
	return nil
}

// Value implements driver.Valuer
func (j JSONB) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return []byte(j), nil
}
