package types

// CompareRequest — input cho compare engine
type CompareRequest struct {
	ParcelID   uint64
	BaselineAt string // YYYY-MM-DD
	TargetAt   string // YYYY-MM-DD
}

// ChangeType — loại thay đổi
type ChangeType string

const (
	ChangeAdded     ChangeType = "added"
	ChangeRemoved   ChangeType = "removed"
	ChangeModified  ChangeType = "modified"
	ChangeUnchanged ChangeType = "unchanged"
)

// LayerChange — một thay đổi trên 1 layer
type LayerChange struct {
	LayerID      uint64
	LayerName    string
	ChangeType   ChangeType
	ChangeDetail *ChangeDetail
}

// ChangeDetail — chi tiết thay đổi
type ChangeDetail struct {
	FieldChanges  []FieldChange
	OverlapDelta  float64
	ImpactSummary string
}

// FieldChange — thay đổi 1 field
type FieldChange struct {
	Field    string
	OldValue string
	NewValue string
}

// CompareResult — kết quả so sánh
type CompareResult struct {
	ParcelID          uint64
	BaselineAt        string
	TargetAt          string
	TotalLayersBefore int
	TotalLayersAfter  int
	Changes           []*LayerChange
	Summary           *CompareSummary
	Timeline          []*TimelinePoint
}

// CompareSummary — tóm tắt so sánh
type CompareSummary struct {
	AddedCount     int
	RemovedCount   int
	ModifiedCount  int
	UnchangedCount int
	NetChange      string
	RiskTrend      string
}

// TimelinePoint — điểm trên timeline
type TimelinePoint struct {
	Date        string
	EventType   string
	LayerID     uint64
	LayerName   string
	Description string
}

// CompareConfig — cấu hình cho Compare Engine
type CompareConfig struct {
	MaxDiffFields              int
	SignificantChangeThreshold float64
}

// DefaultCompareConfig — cấu hình compare mặc định
func DefaultCompareConfig() *CompareConfig {
	return &CompareConfig{
		MaxDiffFields:              20,
		SignificantChangeThreshold: 5.0,
	}
}
