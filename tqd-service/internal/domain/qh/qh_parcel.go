package qh_domain

import (
	"bytes"
	_models "common/domain/entity"
	"context"
	"database/sql/driver"
	"encoding/base64"
	"encoding/json"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Parcel struct {
	_models.BaseEntity
	Geometry RawGeometry `gorm:"column:geometry;type:geometry(MultiPolygon,4326)"`
	// Các thuộc tính khác có thể thêm vào đây
	Lat     float64 `json:"lat" gorm:"column:lat"`
	Lng     float64 `json:"lng" gorm:"column:lon"`
	RefID   uint64  `json:"refId" gorm:"column:ref_id"`
	RefType uint32  `json:"refType" gorm:"column:ref_type"`
	SeoID   *uint64 `json:"seoId" gorm:"column:seo_id"`
	// MinX    float64 `json:"minX" gorm:"column:min_x"`
	// MinY    float64 `json:"minY" gorm:"column:min_y"`
	// MaxX    float64 `json:"maxX" gorm:"column:max_x"`
	// MaxY    float64 `json:"maxY" gorm:"column:max_y"`

	// Geometry []byte `gorm:"type:geometry(MultiPolygon,4326)" json:"geometry,omitempty"`
}

func (Parcel) TableName() string {
	return "parcels"
}

type RawGeometry struct {
	Raw []byte
}

// GormDataType — GORM dùng đúng geometry(MultiPolygon,4326), không suy ra GEOMETRY(Geometry,4326).
func (RawGeometry) GormDataType() string {
	return "geometry(MultiPolygon,4326)"
}

func (r *RawGeometry) Scan(value any) error {
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

func (r RawGeometry) Value() (driver.Value, error) {
	return r.Raw, nil
}

// GormValue binds GeoJSON bytes via PostGIS so inserts/updates are valid for geometry columns.
// Raw GeoJSON text is not accepted as implicit geometry cast (unlike WKT/EWKB).
// ST_MakeValid có thể trả về GeometryCollection; ST_CollectionExtract(..., 3) gom phần polygon thành MultiPolygon.
func (r RawGeometry) GormValue(_ context.Context, _ *gorm.DB) clause.Expr {
	raw := bytes.TrimSpace(r.Raw)
	if len(raw) == 0 {
		return clause.Expr{
			SQL: "ST_SetSRID(ST_GeomFromText('MULTIPOLYGON EMPTY'),4326)::geometry(MultiPolygon,4326)",
		}
	}
	return clause.Expr{
		SQL:  "ST_Multi(ST_CollectionExtract(ST_MakeValid(ST_SetSRID(ST_GeomFromGeoJSON(?),4326)), 3))::geometry(MultiPolygon,4326)",
		Vars: []interface{}{string(raw)},
	}
}

func (r RawGeometry) MarshalJSON() ([]byte, error) {
	encoded := base64.StdEncoding.EncodeToString(r.Raw)
	return json.Marshal(encoded)
}

func (r *RawGeometry) UnmarshalJSON(data []byte) error {
	var encoded string
	if err := json.Unmarshal(data, &encoded); err != nil {
		return err
	}
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return err
	}
	r.Raw = decoded
	return nil
}
