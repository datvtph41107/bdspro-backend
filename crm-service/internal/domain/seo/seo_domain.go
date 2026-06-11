package seo_domain

import (
	"crm/internal/enums"
	"time"

	"gorm.io/gorm"
)

type SeoDomain struct {
	ID uint64 `gorm:"primaryKey;autoIncrement"`

	OriginURL string `gorm:"not null;index"`

	CanonicalURL string `gorm:"index"`

	Classify enums.ESeoClassify `gorm:"type:smallint;not null;default:50"`

	SiteMapLastedAt *time.Time `gorm:""`

	IsSiteMap bool `gorm:"not null;default:false"`
	IsIndex bool `gorm:"not null;default:false"`
	IsRobot *bool `gorm:"type:bool;default:false"`

	Metadata *string `gorm:"type:jsonb"`
	Note string `gorm:"type:text"`
	InternalLinks []SeoInternalLink `gorm:"foreignKey:ParentSeoID;references:ID"`
	Relatives []SeoRelative `gorm:"foreignKey:ParentSeoID;references:ID"`
	RefType enums.ESEORefType `gorm:"type:smallint;index"`

	RefID *uint64 `gorm:"type:bigint;index"`

	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (SeoDomain) TableName() string {
	return "seo_domain"
}

func (s *SeoDomain) GetSchema() string {
	switch s.RefType {
	case enums.ESEORefTypeAdmUnit:
		return "administrative_unit"
	case enums.ESEORefTypeRegion:
		return "planning_region"
	case enums.ESEORefTypeParcel:
		return "planning_parcel"
	case enums.ESEORefTypeProject:
		return "planning_project"
	case enums.ESEORefTypeDocument:
		return "legal_document"
	case enums.ESEORefTypeMap:
		return "gis_layer"
	case enums.ESEORefTypeReport:
		return "report"
	default:
		return "info"
	}
}
