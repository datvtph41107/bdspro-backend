package qh_domain

import (
	"encoding/json"
	"time"
)

// ImportErrorStatus trạng thái của một bản ghi lỗi import
type ImportErrorStatus string

const (
	ImportErrorStatusPending    ImportErrorStatus = "pending"    // chưa xử lý
	ImportErrorStatusProcessing ImportErrorStatus = "processing" // đang retry (tránh chạy song song)
	ImportErrorStatusResolved   ImportErrorStatus = "resolved"   // đã xử lý / re-import thành công
	ImportErrorStatusIgnored    ImportErrorStatus = "ignored"    // bỏ qua
)

// QHRegionImportErrorLog lưu lại các batch insert thất bại để kiểm tra và retry sau
type QHRegionImportErrorLog struct {
	ID uint64 `gorm:"primaryKey;autoIncrement" json:"id"`

	// Thông tin định danh batch
	ImportBatchID string `gorm:"type:varchar(100);index" json:"importBatchId"`
	LayerID       uint64 `gorm:"index"                  json:"layerId"`

	// Thông tin lỗi
	ErrorMessage string `gorm:"type:text" json:"errorMessage"`

	// Payload của các region thất bại (serialized JSON)
	RegionsPayload json.RawMessage `gorm:"type:jsonb" json:"regionsPayload"`
	BatchSize      int             `gorm:"default:0"  json:"batchSize"`

	// Quản lý retry
	RetryCount int               `gorm:"default:0"  json:"retryCount"`
	Status     ImportErrorStatus `gorm:"type:varchar(20);default:'pending';index" json:"status"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}

func (QHRegionImportErrorLog) TableName() string {
	return "qh_region_import_error_logs"
}
