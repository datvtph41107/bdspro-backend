package domain

import (
	"bytes"
	"context"
	"database/sql/driver"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/encoding/ewkb"
	"github.com/paulmach/orb/encoding/wkb"
	"github.com/paulmach/orb/encoding/wkt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PointGeometry struct {
	Raw []byte
}

func (PointGeometry) GormDataType() string {
	return "geometry(Point,4326)"
}

func (p *PointGeometry) Scan(value any) error {
	if value == nil {
		return nil
	}
	switch v := value.(type) {
	case []byte:
		p.Raw = v
	case string:
		p.Raw = []byte(v)
	}
	return nil
}

func (p PointGeometry) Value() (driver.Value, error) {
	return p.Raw, nil
}

func (p PointGeometry) GormValue(_ context.Context, _ *gorm.DB) clause.Expr {
	raw := bytes.TrimSpace(p.Raw)
	if len(raw) == 0 {
		return clause.Expr{
			SQL: "ST_SetSRID(ST_GeomFromText('POINT EMPTY'),4326)::geometry(Point,4326)",
		}
	}
	// Với Point, ST_CollectionExtract(...,1) đảm bảo kết quả là Point
	return clause.Expr{
		SQL:  "ST_CollectionExtract(ST_MakeValid(ST_SetSRID(ST_GeomFromGeoJSON(?),4326)), 1)::geometry(Point,4326)",
		Vars: []interface{}{string(raw)},
	}
}

func (p PointGeometry) MarshalJSON() ([]byte, error) {
	encoded := base64.StdEncoding.EncodeToString(p.Raw)
	return json.Marshal(encoded)
}

func (p *PointGeometry) UnmarshalJSON(data []byte) error {
	var encoded string
	if err := json.Unmarshal(data, &encoded); err != nil {
		return err
	}
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return err
	}
	p.Raw = decoded
	return nil
}

// SetLatLng tạo GeoJSON từ tọa độ (dùng khi insert/update)
func (p *PointGeometry) SetLatLng(lat, lng float64) {
	p.Raw = []byte(fmt.Sprintf(`{"type":"Point","coordinates":[%f,%f]}`, lng, lat))
}

// GetLatLng parse dữ liệu trong Raw (WKB, EWKB, WKT) thành lat/lng
func (p PointGeometry) GetLatLng() (float64, float64, error) {
	if len(p.Raw) == 0 {
		return 0, 0, fmt.Errorf("empty geometry")
	}

	// Chuẩn bị dữ liệu binary: nếu Raw toàn ký tự hex thì decode
	rawBinary := p.Raw
	if isHexBytes(p.Raw) {
		decoded, err := hex.DecodeString(string(p.Raw))
		if err == nil {
			rawBinary = decoded
		}
	}

	var point orb.Point

	// 1. Thử WKB (binary)
	if g, err := wkb.Unmarshal(rawBinary); err == nil {
		point = g.(orb.Point)
		return point.Y(), point.X(), nil
	}

	// 2. Thử EWKB (binary)
	if g, _, err := ewkb.Unmarshal(rawBinary); err == nil {
		point = g.(orb.Point)
		return point.Y(), point.X(), nil
	}

	// 3. Thử WKT (text)
	if g, err := wkt.Unmarshal(string(p.Raw)); err == nil {
		point = g.(orb.Point)
		return point.Y(), point.X(), nil
	}

	return 0, 0, fmt.Errorf("unsupported geometry format: %s", string(p.Raw))
}

// isHexBytes kiểm tra []byte chỉ chứa ký tự hex (0-9, A-F, a-f)
func isHexBytes(b []byte) bool {
	for _, c := range b {
		if !((c >= '0' && c <= '9') || (c >= 'A' && c <= 'F') || (c >= 'a' && c <= 'f')) {
			return false
		}
	}
	return true
}
