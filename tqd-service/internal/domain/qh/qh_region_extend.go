package qh_domain

import (
	"bytes"
	"context"
	"database/sql/driver"
	"encoding/json"
	"time"
	"tqd/internal/enums"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// QHRegionExtend lưu các feature có geometry không phải Polygon/MultiPolygon (Point, LineString, v.v.)
// được lọc ra trong quá trình import region.
type QHRegionExtend struct {
	ID          uint64  `gorm:"primaryKey;autoIncrement" json:"id"`
	LayerID     uint64  `gorm:"not null;index" json:"layerId"`
	Name        string  `gorm:"type:varchar(255)" json:"name"`
	DisplayName string  `gorm:"type:varchar(255)" json:"displayName"`
	Description string  `gorm:"type:text" json:"description"`
	LabelID     *uint64 `gorm:"index" json:"labelId,omitempty"`

	GeomType string `gorm:"type:varchar(50);index" json:"geomType"` // "Point", "LineString", "MultiLineString", v.v.

	Geometry RawExtendGeometry `gorm:"type:geometry(Geometry,4326);not null" json:"geometry"`

	LegalDoc     string `gorm:"type:text" json:"legalDoc"`
	PlanningName string `gorm:"type:varchar(255)" json:"planningName"`

	OriginalProperties json.RawMessage `gorm:"type:jsonb" json:"originalProperties"`
	SourceFile         string          `gorm:"type:varchar(255)" json:"sourceFile"`
	ImportBatchID      string          `gorm:"type:varchar(100);index" json:"importBatchId"`

	Status   enums.RegionStatus `gorm:"default:10;index" json:"status"`
	Version  int                `gorm:"default:1" json:"version"`
	IsLatest bool               `gorm:"default:true;index" json:"isLatest"`

	CreatedAt time.Time  `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt *time.Time `gorm:"index" json:"deletedAt"`
}

func (QHRegionExtend) TableName() string {
	return "qh_region_extends"
}

// RawExtendGeometry lưu GeoJSON nguyên bản cho Point/LineString/MultiLineString.
// Không dùng RawGeometry vì RawGeometry luôn ép về MultiPolygon cho bảng qh_regions.
type RawExtendGeometry struct {
	Raw []byte
}

func (RawExtendGeometry) GormDataType() string {
	return "geometry(Geometry,4326)"
}

func (r *RawExtendGeometry) Scan(value any) error {
	if value == nil {
		return nil
	}
	switch v := value.(type) {
	case []byte:
		r.Raw = v
	case string:
		r.Raw = []byte(v)
	}
	return nil
}

func (r RawExtendGeometry) Value() (driver.Value, error) {
	return r.Raw, nil
}

func (r RawExtendGeometry) GormValue(_ context.Context, _ *gorm.DB) clause.Expr {
	raw := bytes.TrimSpace(r.Raw)
	if len(raw) == 0 {
		return clause.Expr{
			SQL: "ST_SetSRID(ST_GeomFromText('GEOMETRYCOLLECTION EMPTY'),4326)::geometry(Geometry,4326)",
		}
	}
	return clause.Expr{
		SQL:  "ST_MakeValid(ST_SetSRID(ST_GeomFromGeoJSON(?),4326))::geometry(Geometry,4326)",
		Vars: []interface{}{string(raw)},
	}
}
