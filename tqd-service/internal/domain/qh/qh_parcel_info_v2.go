package qh_domain

// import (
// 	"database/sql/driver"
// 	"encoding/json"
// 	"time"
// )

// // PlanningInfo - Thông tin quy hoạch ảnh hưởng từ one_house
// type PlanningInfoV2 struct {
// 	PlanningTypeCode string  `json:"planning_type_code"`
// 	PlanningTypeName string  `json:"planning_type_name"`
// 	PlanningArea     float64 `json:"planning_area"`
// }

// // PlanningInfoList - List planning info
// type PlanningInfoList []PlanningInfoV2

// func (p PlanningInfoList) Value() (driver.Value, error) {
// 	if p == nil {
// 		return nil, nil
// 	}
// 	return json.Marshal(p)
// }

// func (p *PlanningInfoList) Scan(value interface{}) error {
// 	if value == nil {
// 		*p = nil
// 		return nil
// 	}
// 	bytes, ok := value.([]byte)
// 	if !ok {
// 		return nil
// 	}
// 	return json.Unmarshal(bytes, p)
// }

// // QHParcelInfo - Thông tin chi tiết thửa đất từ one_house
// type QHParcelInfoV2 struct {
// 	// Primary identifier
// 	ID       uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
// 	ParcelID uint64 `gorm:"not null;index" json:"parcelId"`

// 	// Property identifiers
// 	PropertyCode string `gorm:"type:varchar(100)" json:"propertyCode"`
// 	PropertyUUID string `gorm:"type:uuid" json:"propertyUuid"`
// 	MapNumber    string `gorm:"type:varchar(50)" json:"mapNumber"`
// 	LandNumber   string `gorm:"type:varchar(50)" json:"landNumber"`

// 	// Physical characteristics
// 	TotalAreaSqm float64 `gorm:"default:0" json:"totalAreaSqm"`
// 	ShapeType    string  `gorm:"type:varchar(50)" json:"shapeType"`
// 	Direction    string  `gorm:"type:varchar(50)" json:"direction"`

// 	// Current land use
// 	ActualLandUseCode string `gorm:"type:varchar(20);index" json:"actualLandUseCode"`

// 	// Administrative info
// 	ProvinceCode string `gorm:"type:varchar(20)" json:"provinceCode"`
// 	DistrictCode string `gorm:"type:varchar(20)" json:"districtCode"`
// 	WardCode     string `gorm:"type:varchar(20)" json:"wardCode"`
// 	AddressText  string `gorm:"type:text" json:"addressText"`

// 	// Planning info (quan trọng cho API Detail)
// 	PlanningInfos        PlanningInfoList `gorm:"type:jsonb" json:"planningInfos"`
// 	PlanningLandType     string           `gorm:"type:text" json:"planningLandType"`
// 	PlanningLandTypeCode string           `gorm:"type:varchar(20)" json:"planningLandTypeCode"`

// 	// Metadata
// 	SourceType int        `gorm:"default:2" json:"sourceType"`
// 	IsVerified bool       `gorm:"default:false" json:"isVerified"`
// 	VerifiedAt *time.Time `json:"verifiedAt"`

// 	// Audit
// 	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
// 	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
// }

// func (QHParcelInfoV2) TableName() string {
// 	return "qh_parcel_info_v2"
// }

// // GetTotalImpactArea - Tính tổng diện tích quy hoạch ảnh hưởng
// func (p *QHParcelInfoV2) GetTotalImpactArea() float64 {
// 	var total float64
// 	for _, info := range p.PlanningInfos {
// 		total += info.PlanningArea
// 	}
// 	return total
// }

// // GetImpactPercentage - Tính % diện tích bị ảnh hưởng
// func (p *QHParcelInfoV2) GetImpactPercentage() float64 {
// 	if p.TotalAreaSqm <= 0 {
// 		return 0
// 	}
// 	return (p.GetTotalImpactArea() / p.TotalAreaSqm) * 100
// }

// // HasPlanningConflict - Kiểm tra có xung đột quy hoạch không
// func (p *QHParcelInfoV2) HasPlanningConflict() bool {
// 	return p.GetImpactPercentage() > 100
// }
