package qh_domain

import (
	_models "common/domain/entity"
	"database/sql/driver"
	"encoding/base64"
	"encoding/json"
)

type Parcel struct {
	_models.BaseEntity
	Geometry RawGeometry `gorm:"type:geometry(MultiPolygon,4326)"`
	// Các thuộc tính khác có thể thêm vào đây
	Lat     float64 `json:"lat" gorm:"column:lat"`
	Lng     float64 `json:"lng" gorm:"column:lon"`
	RefID   uint64  `json:"refId" gorm:"column:ref_id"`
	RefType uint32  `json:"refType" gorm:"column:ref_type"`
	// MinX    float64 `json:"minX" gorm:"column:min_x"`
	// MinY    float64 `json:"minY" gorm:"column:min_y"`
	// MaxX    float64 `json:"maxX" gorm:"column:max_x"`
	// MaxY    float64 `json:"maxY" gorm:"column:max_y"`
}

func (Parcel) TableName() string {
	return "parcels"
}

type RawGeometry struct {
	Raw []byte
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
