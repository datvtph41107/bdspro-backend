package qh_domain

import (
	"encoding/json"
	"time"
	"tqd/internal/enums"
)

type QHRegion struct {
	ID          uint64  `gorm:"primaryKey;autoIncrement" json:"id"`
	LayerID     uint64  `gorm:"not null;index" json:"layerId"`
	Name        string  `gorm:"type:varchar(255)" json:"name"`
	DisplayName string  `gorm:"type:varchar(255)" json:"displayName"`
	Description string  `gorm:"type:text" json:"description"`
	LabelID     *uint64 `gorm:"index" json:"labelId,omitempty"`

	Label      *QHLabel `gorm:"foreignKey:LabelID" json:"label"`
	LabelName  string   `gorm:"-:migration;->;column:label_name" json:"labelName"`
	LabelColor string   `gorm:"-:migration;->;column:label_color" json:"labelColor"`

	LandUseID *uint64        `gorm:"column:land_use_id;index" json:"landUseId,omitempty"`
	LandUse   *QHLandUse     `gorm:"foreignKey:LandUseID" json:"landUse,omitempty"`
	LegendID  *uint64        `gorm:"column:legend_id;index" json:"legendId,omitempty"`
	Legend    *QHLayerLegend `gorm:"foreignKey:LegendID" json:"legend,omitempty"`

	LegalDoc     string `gorm:"type:text" json:"legalDoc"`
	PlanningName string `gorm:"type:varchar(255)" json:"planningName"`

	Geometry RawGeometry `gorm:"type:geometry(MultiPolygon,4326);not null" json:"geometry"`

	Centroid   []byte  `gorm:"type:geometry(Point,4326)" json:"-"`
	AreaSqm    float64 `gorm:"default:0" json:"areaSqm"`
	AreaHa     float64 `gorm:"default:0" json:"areaHa"`
	PerimeterM float64 `gorm:"default:0" json:"perimeterM"`

	// Dữ liệu gốc
	OriginalProperties json.RawMessage `gorm:"type:jsonb" json:"originalProperties"`
	SourceFile         string          `gorm:"type:varchar(255)" json:"sourceFile"`
	ImportBatchID      string          `gorm:"type:varchar(100);index" json:"importBatchId"`

	Status           enums.RegionStatus     `gorm:"default:10;index" json:"status"`
	ProcessingStatus enums.ProcessingStatus `gorm:"default:10;index" json:"processingStatus"`

	Version  int  `gorm:"default:1" json:"version"`
	IsLatest bool `gorm:"default:true;index" json:"isLatest"`

	CreatedAt time.Time  `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt *time.Time `gorm:"index" json:"deletedAt"`
}

func (QHRegion) TableName() string {
	return "qh_regions"
}

type RegionStyleConfig struct {
	FillColor   string  `json:"fillColor,omitempty"`
	FillOpacity float64 `json:"fillOpacity,omitempty"`
	StrokeColor string  `json:"strokeColor,omitempty"`
	StrokeWidth int     `json:"strokeWidth,omitempty"`
}

func (r *QHRegion) CalculateColorHex() string {
	return "#ff0000"
}

// func (r *QHRegion) GetStyleConfig() RegionStyleConfig {
// 	if r.StyleConfig == nil {
// 		return RegionStyleConfig{
// 			FillColor:   r.CalculateColorHex(),
// 			FillOpacity: 0.6,
// 			StrokeColor: "#000000",
// 			StrokeWidth: 1,
// 		}
// 	}

// 	var config RegionStyleConfig
// 	data, _ := json.Marshal(r.StyleConfig)
// 	json.Unmarshal(data, &config)

// 	if config.FillColor == "" {
// 		config.FillColor = r.CalculateColorHex()
// 	}
// 	if config.FillOpacity == 0 {
// 		config.FillOpacity = 0.6
// 	}

// 	return config
// }

// THÊM: Kiểm tra region có đang active không
func (r *QHRegion) IsActive() bool {
	return r.Status == enums.RegionStatusActive && r.DeletedAt == nil && r.IsLatest
}

// THÊM: Kiểm tra region đã được xác nhận chưa
func (r *QHRegion) IsVerified() bool {
	return r.ProcessingStatus == enums.ProcessingStatusVerified
}

func (r *QHRegion) GetLabelID() uint64 {
	return *r.LabelID
}

// Gán label cho region
func (r *QHRegion) SetLabelID(labelID uint64) {
	r.LabelID = &labelID
}
