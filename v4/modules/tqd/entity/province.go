package entity

import (
	"time"

	"userv4/modules/tqd/enums"
)

// Province is the TQD administrative province row (table province_v2 on the TQD DSN).
type Province struct {
	ID                    uint64                   `gorm:"column:id;primaryKey"`
	FullName              string                   `gorm:"column:full_name"`
	ShortName             string                   `gorm:"column:short_name"`
	Lat                   float64                  `gorm:"column:lat"`
	Lng                   float64                  `gorm:"column:lng"`
	Code                  string                   `gorm:"column:code"`
	WardCount             int32                    `gorm:"column:ward_count"`
	QHOrder               int32                    `gorm:"column:qh_order"`
	Zoom                  uint32                   `gorm:"column:zoom"`
	Description           string                   `gorm:"column:description"`
	OfficialName          string                   `gorm:"column:official_name"`
	DivisionType          string                   `gorm:"column:division_type"`
	OfficialOperationDate *time.Time               `gorm:"column:official_operation_date"`
	NaturalAreaKm2        float64                  `gorm:"column:natural_area_km2"`
	Population            uint64                   `gorm:"column:population"`
	PopulationDensity     float64                  `gorm:"column:population_density"`
	UnitType              enums.LocationUnitType   `gorm:"column:unit_type"`
	DataSource            enums.LocationDataSource `gorm:"column:data_source"`
	PlanDocID             *uint64                  `gorm:"column:plan_doc_id"`
	OfficialRole          string                   `gorm:"column:official_role"`
	OfficialAddress       string                   `gorm:"column:official_address"`
	SeoID                 *uint64                  `gorm:"column:seo_id"`
	CreatedAt             time.Time                `gorm:"column:created_at"`
	UpdatedAt             time.Time                `gorm:"column:updated_at"`
	DeletedAt             *time.Time               `gorm:"column:deleted_at"`
}

func (Province) TableName() string { return "province_v2" }
