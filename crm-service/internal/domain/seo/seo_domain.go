package seo_domain

import (
	"crm/internal/enums"
	"time"
)

type SeoClassify uint8

const (
	// chứa thông tin chi tiết
	SeoClassifyA SeoClassify = 10

	// chứa các điều hướng
	SeoClassifyB SeoClassify = 20

	// kết quả của dữ liệu tìm kiếm
	SeoClassifyC SeoClassify = 30

	// link chia sẻ lan truyền ko phục vụ seo, liên kết, bản xem trước
	SeoClassifyD SeoClassify = 40

	// không phục vụ seo: đăng nhập, đăng ký,  admin, 404, 500, ...
	SeoClassifyN SeoClassify = 50
)

type SeoDomain struct {
	ID uint64 `gorm:"primaryKey;autoIncrement"`

	OriginURL string `gorm:"not null;index"`

	CanonicalURL string `gorm:"index"`

	// Phân loại SEO
	Classify SeoClassify `gorm:"not null;default:50"`

	// Thời gian cập nhật sitemap lần cuối
	SiteMapLastedAt *time.Time `gorm:""`

	IsSiteMap bool

	IsIndex bool

	// kiểm tra robot
	IsRobot *bool `gorm:"type:bool;default:false"`

	Metadata *string `gorm:"type:jsonb"`

	Note string `gorm:"type:text"`

	InternalLinks []SeoInternalLink `gorm:"foreignKey:ParentSeoID;references:ID"`

	// schema,đối tượng
	RefType enums.ESEORefType `gorm:"type:smallint"`

	RefID *uint64 `gorm:"type:uint64"`
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
