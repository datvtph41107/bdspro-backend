package qh_domain

import (
	"encoding/json"
	"fmt"
	"time"
	"tqd/internal/domain/jsonb"
)

type RegionStyleConfig struct {
	FillColor   string  `json:"fillColor,omitempty"`
	FillOpacity float64 `json:"fillOpacity,omitempty"`
	StrokeColor string  `json:"strokeColor,omitempty"`
	StrokeWidth int     `json:"strokeWidth,omitempty"`
}

// QHRegion - Khu vực quy hoạch chi tiết
type QHRegion struct {
	ID      uint64  `gorm:"primaryKey;autoIncrement" json:"id"`
	LayerID *uint64 `gorm:"not null;index;comment:Thuộc layer quy hoạch" json:"layerId"`

	Name        string `gorm:"type:varchar(255)" json:"name,omitempty"`
	DisplayName string `gorm:"type:varchar(255)" json:"displayName,omitempty"`
	Description string `gorm:"type:text" json:"description,omitempty"`

	LandUseCode string `gorm:"type:varchar(20);index;comment:Mã loại đất (ODT, TMD, DGT...)" json:"landUseCode"`
	LandUseName string `gorm:"type:varchar(255);comment:Tên loại đất" json:"landUseName,omitempty"`

	LegalDoc     string `gorm:"type:text;comment:Văn bản pháp lý" json:"legalDoc,omitempty"`
	PlanningName string `gorm:"type:varchar(255);comment:Tên đồ án quy hoạch" json:"planningName,omitempty"`

	Geometry RawGeometry `gorm:"type:geometry(MultiPolygon,4326);not null" json:"-"`

	Centroid   []byte  `gorm:"type:geometry(Point,4326)" json:"-"`
	AreaSqm    float64 `gorm:"default:0;comment:Diện tích mét vuông" json:"areaSqm"`
	AreaHa     float64 `gorm:"default:0;comment:Diện tích hecta" json:"areaHa"`
	PerimeterM float64 `gorm:"default:0;comment:Chu vi mét" json:"perimeterM"`

	ColorRed   uint8 `gorm:"default:200;comment:Màu đỏ 0-255" json:"colorRed"`
	ColorGreen uint8 `gorm:"default:200;comment:Màu xanh lá 0-255" json:"colorGreen"`
	ColorBlue  uint8 `gorm:"default:200;comment:Màu xanh dương 0-255" json:"colorBlue"`

	// Style (có thể override layer default)
	StyleConfig jsonb.JSONB `gorm:"type:jsonb;default:'{}';comment:Style override" json:"styleConfig,omitempty"`

	RawProperties jsonb.JSONB `gorm:"type:jsonb;comment:Dữ liệu gốc từ NDJSON" json:"rawProperties,omitempty"`
	RawGeometry   jsonb.JSONB `gorm:"type:jsonb" json:"rawGeometry,omitempty"`

	Status   int  `gorm:"default:10;index;comment:10:Active, 20:Inactive" json:"status"`
	Version  int  `gorm:"default:1;comment:Phiên bản" json:"version"`
	IsLatest bool `gorm:"default:true;index;comment:Phiên bản mới nhất" json:"isLatest"`

	CreatedAt time.Time  `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt *time.Time `gorm:"index" json:"deletedAt,omitempty"`
}

func (QHRegion) TableName() string {
	return "qh_regions"
}

// CalculateColorHex - Tạo mã màu hex từ RGB
func (r *QHRegion) CalculateColorHex() string {
	return fmt.Sprintf("#%02X%02X%02X", r.ColorRed, r.ColorGreen, r.ColorBlue)
}

// GetStyleConfig - Lấy style config cho region
func (r *QHRegion) GetStyleConfig() RegionStyleConfig {
	if r.StyleConfig == nil {
		return RegionStyleConfig{
			FillColor:   r.CalculateColorHex(),
			FillOpacity: 0.6,
			StrokeColor: "#000000",
			StrokeWidth: 1,
		}
	}

	var config RegionStyleConfig
	data, _ := json.Marshal(r.StyleConfig)
	json.Unmarshal(data, &config)

	if config.FillColor == "" {
		config.FillColor = r.CalculateColorHex()
	}
	if config.FillOpacity == 0 {
		config.FillOpacity = 0.6
	}

	return config
}
