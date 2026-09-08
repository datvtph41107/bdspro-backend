package qh_domain

import (
	"time"
	"tqd/internal/enums"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type QHUserReported struct {
	ID uint64 `gorm:"primaryKey;autoIncrement" json:"id"`

	UserID uint64 `gorm:"column:user_id;not null;index:idx_generated_reports_user_id;uniqueIndex:uidx_user_report_command_key,priority:1,where:command_key <> '' AND deleted_at IS NULL" json:"userId"`

	ReportType enums.GeneratedReportType   `gorm:"column:report_type;not null;index:idx_generated_reports_report_type" json:"reportType"`
	Status     enums.GeneratedReportStatus `gorm:"column:status;not null;index:idx_generated_reports_status" json:"status"`

	Title    string `gorm:"column:title;type:varchar(255);not null" json:"title"`
	Subtitle string `gorm:"column:subtitle;type:varchar(255);default:''" json:"subtitle"`

	Address      string `gorm:"column:address;type:text;default:''" json:"address"`
	Province     string `gorm:"column:province;type:varchar(100);default:''" json:"province"`
	ProvinceCode string `gorm:"column:province_code;type:varchar(20);default:''" json:"provinceCode"`
	WardCode     string `gorm:"column:ward_code;type:varchar(20);default:''" json:"wardCode"`

	ParcelID *uint64 `gorm:"column:parcel_id;index:idx_generated_reports_parcel_id" json:"parcelId,omitempty"`
	RegionID *uint64 `gorm:"column:region_id;index:idx_generated_reports_region_id" json:"regionId,omitempty"`

	MinLon float64 `gorm:"column:min_lon;default:0" json:"minLon"`
	MinLat float64 `gorm:"column:min_lat;default:0" json:"minLat"`
	MaxLon float64 `gorm:"column:max_lon;default:0" json:"maxLon"`
	MaxLat float64 `gorm:"column:max_lat;default:0" json:"maxLat"`

	CenterLat float64 `gorm:"column:center_lat;default:0" json:"centerLat"`
	CenterLon float64 `gorm:"column:center_lon;default:0" json:"centerLon"`

	ThumbnailURL string `gorm:"column:thumbnail_url;type:text;default:''" json:"thumbnailUrl"`
	ImageURL     string `gorm:"column:image_url;type:text;default:''" json:"imageUrl"`
	PDFURL       string `gorm:"column:pdf_url;type:text;default:''" json:"pdfUrl"`
	ShareURL     string `gorm:"column:share_url;type:text;default:''" json:"shareUrl"`

	FileSize uint64 `gorm:"column:file_size;default:0" json:"fileSize"`
	Format   string `gorm:"column:format;type:varchar(20);default:''" json:"format"`

	Comparison datatypes.JSON `gorm:"column:comparison;type:jsonb;default:'{}'::jsonb" json:"comparison"`
	Metadata   datatypes.JSON `gorm:"column:metadata;type:jsonb;default:'{}'::jsonb" json:"metadata"`

	CommandKey   string `gorm:"column:command_key;type:varchar(128);default:'';uniqueIndex:uidx_user_report_command_key,priority:2" json:"commandKey"`
	RequestHash  string `gorm:"column:request_hash;type:varchar(64);default:''" json:"-"`
	JobID        string `gorm:"column:job_id;type:varchar(100);default:'';index:idx_generated_reports_job_id" json:"jobId"`
	ErrorMessage string `gorm:"column:error_message;type:text;default:''" json:"errorMessage"`

	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime;index:idx_generated_reports_created_at" json:"createdAt"`
	UpdatedAt time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
	ExpiresAt *time.Time     `gorm:"column:expires_at;index:idx_generated_reports_expires_at" json:"expiresAt,omitempty"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deletedAt,omitempty"`
}

func (QHUserReported) TableName() string {
	return "user_reported"
}
