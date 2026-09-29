package dto

// ─────────────────────────────────────────────────────────────────────────────
// SubmitReportDTO — command object Handler → Usecase cho SendPropertyReport.
//
// Thiết kế tái sử dụng tối đa:
//
//   contribute_data trong proto là UpdatePropertyRequest → map thẳng sang
//   *UpdatePropertyDTO — không tạo sub-struct mới, không duplicate field.
//
//   Usecase gọi proposeChange(*UpdatePropertyDTO) nguyên vẹn nếu ContributeData != nil.
//   Không sửa bất kỳ chữ nào trong proposeChange().
//
// ─────────────────────────────────────────────────────────────────────────────

type SubmitReportDTO struct {
	// ID lineage bị báo cáo
	LineageID uint64

	// Loại báo cáo: map EReportType (10/20/30/40/50/60)
	ReportType uint32

	// Mô tả thêm
	Note *string

	// Sub-issues được chọn — string constants (VD: "wrong_name", "wrong_coords")
	SubIssues []string

	// BĐS liên quan — bắt buộc khi ReportType=30 (duplicate)
	RelatedLineageID *uint64

	// Ảnh bằng chứng (PropertyMedia IDs)
	AttachmentMediaIDs []uint64

	// Tài liệu đính kèm (File IDs)
	AttachmentFileIDs []uint64

	// ─── Propose data ─────────────────────────────────────────────────────────
	// Tái sử dụng UpdatePropertyDTO nguyên vẹn — zero duplication.
	// nil  = chỉ báo cáo thuần → không tạo propose lineage.
	// !nil = kèm đề xuất sửa → gọi proposeChange(ContributeData).
	ContributeData *UpdatePropertyDTO
}

// ─────────────────────────────────────────────────────────────────────────────
// SubmitReportResult — output từ usecase
// ─────────────────────────────────────────────────────────────────────────────

type SubmitReportResult struct {
	ReportID  uint64
	LineageID uint64

	// nil nếu ContributeData == nil (báo cáo thuần, không propose)
	ProposeLineageID *uint64

	ReportType uint32
	Note       string
	CreatedAt  string
}
