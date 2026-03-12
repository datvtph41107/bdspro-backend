package domain

import (
	"bdspro/internal/enums"
	"time"
)

// ─────────────────────────────────────────────────────────────────────────────
// PropertyContribute — bản ghi đóng góp / báo cáo của user
//
// Hai trường hợp:
//   1. Báo cáo thuần (WrongInfo, Spam...): ProposeLineageID = nil
//   2. Đóng góp có dữ liệu (CorrectInfo): ProposeLineageID trỏ đến
//      PropertyLineage mới được tạo bởi proposeChange()
//
// Khi VerifyContribute (approve):
//   - ProposeLineage trở thành bản gốc mới
//   - OriginLineage được lưu lịch sử (IdentifyPropertyID không đổi, IsArchived=true)
//   - ContributeStatus → Approved
// ─────────────────────────────────────────────────────────────────────────────

type PropertyContribute struct {
	ID uint64 `gorm:"primaryKey;autoIncrement" json:"id"`

	// ── Target — BĐS nào đang được đóng góp ─────────────────────────────────
	LineageID uint64 `gorm:"column:lineage_id;not null;index" json:"lineage_id"`

	// ── Contributor ───────────────────────────────────────────────────────────
	ContributorID uint64 `gorm:"column:contributor_id;not null;index" json:"contributor_id"`

	// ── Loại đóng góp ─────────────────────────────────────────────────────────
	ContributeType enums.ContributeType `gorm:"column:contribute_type;not null" json:"contribute_type"`

	// ── Nội dung báo cáo / mô tả đóng góp ────────────────────────────────────
	Note string `gorm:"column:note;type:text" json:"note"`

	// ── Link đến propose lineage (chỉ khi ContributeType = CorrectInfo) ───────
	// nil = báo cáo thuần, không có dữ liệu đề xuất
	// non-nil = đóng góp có dữ liệu, trỏ đến PropertyLineage được tạo bởi proposeChange()
	ProposeLineageID *uint64 `gorm:"column:propose_lineage_id;index" json:"propose_lineage_id,omitempty"`

	// ── Trạng thái xử lý ──────────────────────────────────────────────────────
	Status enums.ContributeStatus `gorm:"column:status;not null;default:10" json:"status"`

	// ── Review ────────────────────────────────────────────────────────────────
	ReviewedByID *uint64    `gorm:"column:reviewed_by_id"          json:"reviewed_by_id,omitempty"`
	ReviewedAt   *time.Time `gorm:"column:reviewed_at"             json:"reviewed_at,omitempty"`
	ReviewNote   string     `gorm:"column:review_note;type:text"   json:"review_note,omitempty"`

	// ── Timestamps ────────────────────────────────────────────────────────────
	CreatedAt time.Time  `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at;index"          json:"deleted_at,omitempty"`
}

func (PropertyContribute) TableName() string { return "property_contribute" }
