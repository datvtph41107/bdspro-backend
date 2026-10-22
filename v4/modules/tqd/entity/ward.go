package entity

import (
	"time"

	"mdn4/modules/tqd/enums"
)

// Ward is the TQD administrative ward row (table ward_v2 on the TQD DSN).
type Ward struct {
	ID                    uint64                   `gorm:"column:id;primaryKey"`
	FullName              string                   `gorm:"column:full_name"`
	ShortName             string                   `gorm:"column:short_name"`
	Lat                   float64                  `gorm:"column:lat"`
	Lng                   float64                  `gorm:"column:lng"`
	Code                  string                   `gorm:"column:code"`
	ProvinceID            uint64                   `gorm:"column:province_id"`
	ProvinceCode          string                   `gorm:"column:province_code"`
	Summary               string                   `gorm:"column:summary"`
	Description           string                   `gorm:"column:description"`
	OfficialName          string                   `gorm:"column:official_name"`
	DivisionType          string                   `gorm:"column:division_type"`
	OfficialOperationDate *time.Time               `gorm:"column:official_operation_date"`
	NaturalAreaKm2        float64                  `gorm:"column:natural_area_km2"`
	Population            uint64                   `gorm:"column:population"`
	PopulationDensity     float64                  `gorm:"column:population_density"`
	UnitType              enums.LocationUnitType   `gorm:"column:unit_type"`
	DataSource            enums.LocationDataSource `gorm:"column:data_source"`
	OfficialRole          string                   `gorm:"column:official_role"`
	OfficialAddress       string                   `gorm:"column:official_address"`
	SeoID                 *uint64                  `gorm:"column:seo_id"`
	CreatedAt             time.Time                `gorm:"column:created_at"`
	UpdatedAt             time.Time                `gorm:"column:updated_at"`
	DeletedAt             *time.Time               `gorm:"column:deleted_at"`
}

func (Ward) TableName() string { return "ward_v2" }
