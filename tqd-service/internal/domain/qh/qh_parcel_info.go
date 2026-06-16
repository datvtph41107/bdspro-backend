package qh_domain

import (
	_models "common/domain/entity"
	_utils "common/utils"
	"database/sql/driver"
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

// =====================================================
// PLANNING INFO (QUY HOẠCH ẢNH HƯỞNG)
// =====================================================

// QHPlanningUse - Thông tin quy hoạch ảnh hưởng từ one_house
type QHPlanningUse struct {
	_models.BaseEntity
	PlanningTypeCode string `json:"planning_type_code"`
	PlanningTypeName string `json:"planning_type_name"`
	// PlanningArea     float64 `json:"planning_area"`
}

func (QHPlanningUse) TableName() string {
	return "qh_planning_land_use"
}

// PlanningInfoList - List planning info
type PlanningInfoList []QHPlanningUse

func (p PlanningInfoList) Value() (driver.Value, error) {
	if p == nil {
		return nil, nil
	}
	return json.Marshal(p)
}

func (p *PlanningInfoList) Scan(value interface{}) error {
	if value == nil {
		*p = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, p)
}

// =====================================================
// QH PARCEL INFO - THÔNG TIN CHI TIẾT THỬA ĐẤT
// =====================================================

// QHParcelInfo - Thông tin chi tiết thửa đất (cho preview và search)
type QHParcelInfo struct {
	// Primary identifier
	ID       uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	ParcelID uint64 `gorm:"not null;index" json:"parcelId"`

	// Property identifiers
	PropertyCode string `gorm:"type:varchar(100)" json:"propertyCode"`
	PropertyUUID string `gorm:"type:uuid" json:"propertyUuid"`
	MapNumber    string `gorm:"type:varchar(50)" json:"mapNumber"`
	LandNumber   string `gorm:"type:varchar(50)" json:"landNumber"`

	// Physical characteristics
	TotalAreaSqm float64 `gorm:"default:0" json:"totalAreaSqm"`

	// ShapeType    string  `gorm:"type:varchar(50)" json:"shapeType"`
	// Location (THÊM 2 TRƯỜNG NÀY)
	Latitude  float64 `gorm:"column:lat" json:"latitude"`
	Longitude float64 `gorm:"column:lon" json:"longitude"`

	// Administrative info
	ProvinceCode string `gorm:"type:varchar(20)" json:"provinceCode"`
	DistrictCode string `gorm:"type:varchar(20)" json:"districtCode"`
	WardCode     string `gorm:"type:varchar(20)" json:"wardCode"`
	AddressText  string `gorm:"type:text" json:"addressText"`
	AdrSearch    string `gorm:"type:text" json:"adrSearch"`
	IsSeo        bool   `gorm:"column:is_seo;default:false;index" json:"isSeo"`

	// Administrative info (V1 format - legacy)
	ProvinceV1Id *uint64 `gorm:"default:null" json:"provinceV1Id,omitempty"`
	DistrictV1Id *uint64 `gorm:"default:null" json:"districtV1Id,omitempty"`
	WardV1Id     *uint64 `gorm:"default:null" json:"wardV1Id,omitempty"`
	ProvinceV2Id *uint64 `gorm:"default:null" json:"provinceV2Id,omitempty"`
	WardV2Id     *uint64 `gorm:"default:null" json:"wardV2Id,omitempty"`

	// Planning info
	PlanningInfos PlanningInfoList `gorm:"type:jsonb" json:"planningInfos"`
	// PlanningLandType     string           `gorm:"type:text" json:"planningLandType"`
	// PlanningLandTypeCode string           `gorm:"type:varchar(20)" json:"planningLandTypeCode"`
	// // Current land use
	// ActualLandUseCode string `gorm:"type:varchar(20);index" json:"actualLandUseCode"`
	PlanningUseID *uint64        `gorm:"" json:"planning_use_id"`
	PlanningUse   *QHPlanningUse `gorm:"foreignKey:PlanningUseID;references:ID" json:"planning_use"`
	// Shape info
	ShapeId   *uint64  `gorm:"" json:"shapeId"`
	ShapeName string   `gorm:"-:migration;column:shape_name" json:"shapeName"`
	Shape     *QHShape `gorm:"-" json:"shape,omitempty"`

	// DirectionID   *uint64      `gorm:"" json:"directionId"`
	// DirectionName string       `gorm:"->;column:direction_name" json:"directionName"`
	DirectionsRaw []byte        `gorm:"-:migration;column:directions"` // JSON từ DB
	Directions    []QHDirection `gorm:"many2many:qh_parcel_direction;" json:"direction"`

	// Facade
	Facade uint64 `gorm:"default:0" json:"facade"`

	// Land type
	LandTypeId *uint64 `gorm:"default:null" json:"landTypeId,omitempty"`
	// LandType   *QHLandType `gorm:"-" json:"landType,omitempty"`

	// Geometry (GeoJSON bytes)
	Geometry []byte `gorm:"type:geometry(MultiPolygon,4326)" json:"geometry,omitempty"`

	// Metadata
	SourceType int        `gorm:"default:2" json:"sourceType"`
	IsVerified bool       `gorm:"default:false" json:"isVerified"`
	VerifiedAt *time.Time `json:"verifiedAt"`

	// Audit
	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}

func (QHParcelInfo) TableName() string {
	return "qh_parcel_info"
}

// BeforeSave - GORM hook to automatically update AdrSearch
func (p *QHParcelInfo) BeforeSave(tx *gorm.DB) error {
	if p.AddressText != "" {
		p.AdrSearch = _utils.CleanSearchText(p.AddressText)
	}
	return nil
}

type ParcelInfoResponse struct {
	ID         uint64  `json:"id"`
	ParcelID   uint64  `json:"parcelId"`
	AddressStr string  `json:"addressStr"`
	MapNumber  uint32  `json:"mapNumber"`
	LandNumber uint32  `json:"landNumber"`
	Lat        float64 `json:"lat"`
	Lng        float64 `json:"lng"`
	ShapeId    *uint64 `json:"shapeId"`
	Facade     uint32  `json:"facade"`
	TotalArea  float64 `json:"totalArea"`
	LandTypeId uint64  `json:"landTypeId"`
	Geometry   []byte  `json:"geometry,omitempty"`
	IsSeo      bool    `json:"isSeo"`

	Directions []QHDirection `json:"direction"`
	Shape      *QHShape      `json:"shape"`
}

// GetTotalImpactArea - Tính tổng diện tích quy hoạch ảnh hưởng
func (p *QHParcelInfo) GetTotalImpactArea() float64 {
	var total float64
	// for _, info := range p.PlanningInfos {
	// 	total += info.PlanningArea
	// }
	return total
}

// GetImpactPercentage - Tính % diện tích bị ảnh hưởng
func (p *QHParcelInfo) GetImpactPercentage() float64 {
	if p.TotalAreaSqm <= 0 {
		return 0
	}
	return (p.GetTotalImpactArea() / p.TotalAreaSqm) * 100
}

// HasPlanningConflict - Kiểm tra có xung đột quy hoạch không
func (p *QHParcelInfo) HasPlanningConflict() bool {
	return p.GetImpactPercentage() > 100
}
