package dto

// ImportFileRequest — upload file + tùy chọn import (xử lý nền)
type ImportFileRequest struct {
	FileContent      []byte
	FileFormat       string
	LayerID          uint64
	LabelField       string
	LabelMappings    map[string]uint64
	DefaultLabelID   *uint64
	SourceFileName   string
	UserID           uint64
	ValidateGeometry bool
	AutoFixGeometry  bool
	SkipInvalid      bool
}

// ImportEnqueueResult — phản hồi ngay sau khi ghi file và gắn job
type ImportEnqueueResult struct {
	Code    uint32
	BatchID string
	Status  string
	Message string
}

// ImportRequest — tùy chọn xử lý feature (dùng nội bộ job)
type ImportRequest struct {
	LayerID          uint64
	LabelField       string
	LabelMappings    map[string]uint64
	DefaultLabelID   *uint64
	SourceFileName   string
	UserID           uint64
	ValidateGeometry bool
	AutoFixGeometry  bool
	SkipInvalid      bool
}

// ImportResult — tổng hợp sau khi job chạy xong (log / debug)
type ImportResult struct {
	BatchID             string
	TotalFeatures       int
	SuccessCount        int
	FailedCount         int
	Failures            []ImportFailure
	LabelStatistics     map[uint64]int32
	TotalAreaHa         float64
	TotalRegionsCreated int
}

type ImportFailure struct {
	FeatureIndex int
	Reason       string
}

type PreviewResult struct {
	TotalFeatures   int
	AvailableFields []string
	FieldStats      []FieldStat
	Samples         []FeatureSample
}

type FieldStat struct {
	FieldName      string
	Coverage       int
	DistinctValues []string
}

type FeatureSample struct {
	Index        int
	GeometryType string
	Properties   map[string]interface{}
}

type BatchStatus struct {
	BatchID       string
	Status        string
	TotalFeatures int
	SuccessCount  int
	FailedCount   int
	CreatedAt     string
	CompletedAt   string
	ErrorMessages []string
}

// RetryImportErrorResult — kết quả retry một bản ghi qh_region_import_error_logs
type RetryImportErrorResult struct {
	Success      bool
	Message      string
	TotalItems   int
	SuccessCount int
	FailedCount  int
}
