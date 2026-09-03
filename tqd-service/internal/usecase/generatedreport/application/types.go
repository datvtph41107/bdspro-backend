package application

import (
	"common/operation"
	"time"
	"tqd/internal/enums"
)

/**
 * Outcome cho biết create report kết thúc theo nhánh nào.
 */
type Outcome string

const (
	OutcomeAccepted Outcome = "accepted"
	OutcomeExisting Outcome = "existing"
)

/**
 * Input là dữ liệu business cần để tạo generated report.
 *
 * User identity không nằm trong đây vì service lấy từ Actor đã được bind vào context.
 */
type Input struct {
	ReportType uint32
	EntityType uint32
	EntityID   uint64
	ParcelID   uint64
	RegionID   uint64

	Title    string
	Subtitle string

	Location Location
	Spatial  Spatial

	Comparison   Comparison
	MetadataJSON string
}

type Location struct {
	Address      string
	Province     string
	ProvinceCode string
	WardCode     string
}

type Bounds struct {
	MinLon float64
	MinLat float64
	MaxLon float64
	MaxLat float64
}

type Point struct {
	Lat float64
	Lon float64
}

type Spatial struct {
	Bounds   Bounds
	Centroid Point
}

type Comparison struct {
	FromPlanName string `json:"fromPlanName"`
	ToPlanName   string `json:"toPlanName"`
	FromYear     uint32 `json:"fromYear"`
	ToYear       uint32 `json:"toYear"`
}

/**
 * Report là generated report ở mức business.
 *
 * Struct này không chứa GORM tags để internal/report không phụ thuộc cách PostgreSQL lưu dữ liệu.
 */
type Report struct {
	ID     uint64
	UserID uint64

	ReportType enums.GeneratedReportType
	Status     enums.GeneratedReportStatus

	Title    string
	Subtitle string

	Location Location
	Spatial  Spatial

	ParcelID *uint64
	RegionID *uint64

	Format string

	ComparisonJSON []byte
	MetadataJSON   []byte

	Operation   operation.Code
	OperationID string
	CommandKey  string
	RequestHash string
	JobID       string

	ThumbnailURL string
	ImageURL     string
	PDFURL       string
	ShareURL     string
	FileSize     uint64
	ErrorMessage string

	CreatedAt time.Time
	UpdatedAt time.Time
	ExpiresAt *time.Time
}

/**
 * Result trả report cùng outcome để caller biết đây là lần tạo mới hay request lặp lại.
 */
type Result struct {
	Report  Report
	Outcome Outcome
}

/**
 * Parcel chứa phần dữ liệu cần thiết để dựng report từ một thửa đất.
 */
type Parcel struct {
	ID uint64

	MapNumber  string
	LandNumber string
	AreaSqm    float64

	Location Location
	Spatial  Spatial
}

/**
 * Region chứa phần dữ liệu cần thiết để dựng report từ một vùng quy hoạch.
 */
type Region struct {
	ID uint64

	Name             string
	DisplayName      string
	PlanningName     string
	LayerName        string
	LayerDisplayName string
	LandUseName      string
	AreaSqm          float64

	Location Location
	Spatial  Spatial
}
