package domain

import "time"

// ─────────────────────────────────────────────────────────────────────────────
// PropertyLineageHistory — lịch sử phiên bản lineage sau khi merge
//
// Khi VerifyContribute approve:
//   - Bản gốc cũ KHÔNG bị xóa
//   - Được chuyển vào bảng này để lưu vết
//   - OriginLineageID = ID của lineage gốc trước khi thay thế
//   - ContributeID = đóng góp nào trigger sự thay đổi này
// ─────────────────────────────────────────────────────────────────────────────

type PropertyLineageHistory struct {
	ID uint64 `gorm:"primaryKey;autoIncrement" json:"id"`

	// LineageID — ID của PropertyLineage bản gốc bị thay thế
	LineageID uint64 `gorm:"column:lineage_id;not null;index" json:"lineage_id"`

	// ContributeID — đóng góp nào dẫn đến sự thay thế
	ContributeID uint64 `gorm:"column:contribute_id;not null;index" json:"contribute_id"`

	// ReplacedByLineageID — lineage mới thay thế (chính là ProposeLineage được promote)
	ReplacedByLineageID uint64 `gorm:"column:replaced_by_lineage_id;not null" json:"replaced_by_lineage_id"`

	// Snapshot FK tại thời điểm bị thay thế — để reconstruct lịch sử nếu cần
	SnapshotPropertyInfoID *uint64 `gorm:"column:snapshot_property_info_id"  json:"snapshot_property_info_id,omitempty"`
	SnapshotLocationID     *uint64 `gorm:"column:snapshot_location_id"       json:"snapshot_location_id,omitempty"`
	SnapshotLandInfoID     *uint64 `gorm:"column:snapshot_land_info_id"      json:"snapshot_land_info_id,omitempty"`
	SnapshotBuildingInfoID *uint64 `gorm:"column:snapshot_building_info_id"  json:"snapshot_building_info_id,omitempty"`
	SnapshotEdvidenceID    *uint64 `gorm:"column:snapshot_edvidence_id"      json:"snapshot_edvidence_id,omitempty"`

	// MergedByID — admin nào thực hiện merge
	MergedByID uint64    `gorm:"column:merged_by_id;not null"       json:"merged_by_id"`
	MergedAt   time.Time `gorm:"column:merged_at;not null"          json:"merged_at"`
	MergeNote  string    `gorm:"column:merge_note;type:text"        json:"merge_note,omitempty"`

	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
}

func (PropertyLineageHistory) TableName() string { return "property_lineage_history" }
