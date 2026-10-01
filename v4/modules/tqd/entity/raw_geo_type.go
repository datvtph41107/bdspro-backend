package entity

import (
	"bytes"
	"context"
	"database/sql/driver"
	"encoding/base64"
	"encoding/json"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

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
