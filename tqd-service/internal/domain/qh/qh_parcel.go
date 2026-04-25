package qh_domain

import (
	"bytes"
	_models "common/domain/entity"
	"database/sql/driver"
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
	Raw string
}

func (r *RawGeometry) GetType() uint32 {
	return 0 // hoặc type MultiPolygon nếu bạn define enum
}

func (r *RawGeometry) Write(buf *bytes.Buffer) error {
	buf.WriteString(r.Raw)
	return nil
}

// DB scan
func (r *RawGeometry) Scan(value any) error {
	if value == nil {
		return nil
	}
	switch v := value.(type) {
	case []byte:
		r.Raw = string(v)
	case string:
		r.Raw = v
	}
	return nil
}

// DB insert
func (r RawGeometry) Value() (driver.Value, error) {
	return r.Raw, nil
}
