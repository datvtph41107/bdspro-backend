// internal/domain/ward.go
package domain

import (
	"time"
	"tqd/internal/enums"
)

type Ward struct {
	ID           uint64    `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	FullName     string    `json:"fullName" gorm:"column:full_name;size:255;not null;index:idx_wards_name"`
	ShortName    string    `json:"shortName" gorm:"column:short_name;size:100;index:idx_wards_short_name"`
	Lat          float64   `json:"lat" gorm:"column:lat;type:decimal(10,8)"`
	Lng          float64   `json:"lng" gorm:"column:lng;type:decimal(11,8)"`
	Code         string    `json:"code" gorm:"column:code;size:20"`
	ProvinceID   uint64    `json:"provinceId" gorm:"column:province_id;not null;index:idx_wards_province_id"`
	ProvinceCode string    `json:"provinceCode" gorm:"column:province_code;size:20"`
	CreatedAt    time.Time `json:"createdAt" gorm:"column:created_at;default:CURRENT_TIMESTAMP;index:idx_wards_created_at"`
	UpdatedAt    time.Time `json:"updatedAt" gorm:"column:updated_at;default:CURRENT_TIMESTAMP"`

	Geometry RawGeometry `gorm:"column:geometry;type:geometry(MultiPolygon,4326)"`

	// thông tin địa bàn lấy ở đây
	Description string `json:"description" gorm:"column:description;size:255"`

	OfficialName          string     `json:"officialName" gorm:"column:official_name;size:255"`
	DivisionType          string     `json:"divisionType" gorm:"column:division_type;size:100"`
	OfficialOperationDate *time.Time `json:"officialOperationDate,omitempty" gorm:"column:official_operation_date;type:date"`
	NaturalAreaKm2        float64    `json:"naturalAreaKm2" gorm:"column:natural_area_km2;type:decimal(12,4)"`
	Population            uint64     `json:"population" gorm:"column:population"`
	PopulationDensity     float64    `json:"populationDensity" gorm:"column:population_density;type:decimal(12,2)"`

	UnitType        enums.LocationUnitType   `json:"unitType" gorm:"column:unit_type;type:int4;not null;default:0"`
	DataSource      enums.LocationDataSource `json:"dataSource" gorm:"column:data_source;type:int4;not null;default:0"`
	OfficialRole    string                   `json:"officialRole" gorm:"column:official_role;size:500"`
	OfficialAddress string                   `json:"officialAddress" gorm:"column:official_address;size:500"`

	SeoID     *uint64    `json:"seoId,omitempty" gorm:"column:seo_id;index"`
	DeletedAt *time.Time `json:"deletedAt,omitempty" gorm:"column:deleted_at;index"`

	Province *Province `json:"province,omitempty" gorm:"foreignKey:ProvinceID;references:ID;-:migration"`
}

func (Ward) TableName() string {
	return "ward_v2"
}
