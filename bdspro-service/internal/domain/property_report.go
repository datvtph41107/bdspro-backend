package domain

import (
	"bdspro/internal/enums"
	_models "common/domain/entity"
)

type PropertyReport struct {
	_models.BaseEntity

	// FK đến lineage gốc được báo cáo
	LineageID uint64           `gorm:"column:lineage_id;type:bigint;not null;index" json:"lineageId"`
	Lineage   *PropertyLineage `gorm:"foreignKey:LineageID;references:ID" json:"lineage,omitempty"`

	// Người báo cáo (origin_id từ context)
	ReporterOriginID uint64 `gorm:"column:reporter_origin_id;type:bigint;not null;index" json:"reporterOriginId"`

	// Loại báo cáo chính
	ReportType enums.EReportType `gorm:"column:report_type;type:smallint;not null;index" json:"reportType"`

	// Sub-issues được check (bitmask hoặc JSON array string)
	// VD: "wrong_name,wrong_area,wrong_location"
	SubIssues string `gorm:"column:sub_issues;type:text" json:"subIssues,omitempty"`

	// Trạng thái xử lý
	Status enums.EReportStatus `gorm:"column:status;type:smallint;not null;default:10;index" json:"status"`

	// Mô tả chi tiết từ người báo cáo
	Description string `gorm:"column:description;type:text" json:"description,omitempty"`

	// Ảnh đính kèm (JSON array of media IDs hoặc URLs)
	AttachmentMediaIDs string `gorm:"column:attachment_media_ids;type:text" json:"attachmentMediaIds,omitempty"`

	// File đính kèm tài liệu (JSON array)
	AttachmentFileIDs string `gorm:"column:attachment_file_ids;type:text" json:"attachmentFileIds,omitempty"`

	// BDS trùng lặp liên quan
	RelatedLineageID *uint64 `gorm:"column:related_lineage_id;type:bigint;index" json:"relatedLineageId,omitempty"`

	// Ghi chú của admin khi xử lý
	AdminNote string `gorm:"column:admin_note;type:text" json:"adminNote,omitempty"`

	// Admin xử lý
	ReviewerOriginID *uint64 `gorm:"column:reviewer_origin_id;type:bigint;index" json:"reviewerOriginId,omitempty"`
}

func (PropertyReport) TableName() string {
	return "property_report"
}
