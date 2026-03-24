// internal/domain/entity_snapshot.go
package domain

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

type PropertyImpactVersion struct {
	ID           uint64     `gorm:"primaryKey;autoIncrement"`
	EntityType   string     `gorm:"type:varchar(50);not null;index:idx_entity_type_id"`
	EntityID     uint64     `gorm:"not null;index:idx_entity_type_id"`
	PropertyID   uint64     `gorm:"not null;index:idx_property_id"`
	SnapshotData JSONB      `gorm:"type:jsonb;not null"`
	Version      uint32     `gorm:"default:1"`
	ValidFrom    time.Time  `gorm:"not null"`
	ValidTo      *time.Time `gorm:"index"`
	CreatedAt    time.Time  `gorm:"autoCreateTime"`
	CreatedBy    uint64     `gorm:"not null"`
}

func (PropertyImpactVersion) TableName() string {
	return "property_impact_versions"
}

type JSONB map[string]interface{}

func (j JSONB) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, j)
}

func (j JSONB) Unmarshal(target interface{}) error {
	bytes, err := json.Marshal(j)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, target)
}

// Marshal - mã hóa struct thành JSONB
func MarshalJSONB(v interface{}) (JSONB, error) {
	bytes, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	var result JSONB
	if err := json.Unmarshal(bytes, &result); err != nil {
		return nil, err
	}
	return result, nil
}
