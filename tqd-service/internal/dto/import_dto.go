package dto

// ImportRequest - Request import file GeoJSON
type ImportRequest struct {
	GeoJSON          string
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

// ImportResult - Kết quả import
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
