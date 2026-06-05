package dto

const (
	WorkspaceEntityTypeParcel   uint32 = 1
	WorkspaceEntityTypeLocation uint32 = 2
	WorkspaceEntityTypeReport   uint32 = 3

	ViewHistorySourceSearch    uint32 = 1
	ViewHistorySourceMapClick  uint32 = 2
	ViewHistorySourceReport    uint32 = 3
	ViewHistorySourceShareLink uint32 = 4
	ViewHistorySourceUnknown   uint32 = 99

	GeneratedReportTypePlanning   uint32 = 1
	GeneratedReportTypeSnapshot   uint32 = 2
	GeneratedReportTypeComparison uint32 = 3
	GeneratedReportTypeAnalysis   uint32 = 4

	GeneratedReportStatusReady      uint32 = 10
	GeneratedReportStatusProcessing uint32 = 20
	GeneratedReportStatusExpired    uint32 = 30
	GeneratedReportStatusFailed     uint32 = 40

	PreviewModeGeometry  uint32 = 1
	PreviewModeThumbnail uint32 = 2
	PreviewModeBounds    uint32 = 3
	PreviewModeStatus    uint32 = 4

	FocusTypeBounds   uint32 = 1
	FocusTypePoint    uint32 = 2
	FocusTypeGeometry uint32 = 3
)

type WorkspacePageRequest struct {
	UserID uint64
	Page   uint32
	Size   uint32
	Cursor string
}

type WorkspaceListResponse[T any] struct {
	Items      []T
	NextCursor string
	Total      int64
}

type SpatialPointDTO struct {
	Lat float64
	Lon float64
}

type SpatialBoundsDTO struct {
	MinLon float64
	MinLat float64
	MaxLon float64
	MaxLat float64
}

type GeometryPreviewDTO struct {
	Type     string
	GeoJSON  string
	Centroid SpatialPointDTO
	Bounds   SpatialBoundsDTO
	SRID     uint32
}

type PreviewHintDTO struct {
	Mode            uint32
	Bounds          SpatialBoundsDTO
	Centroid        SpatialPointDTO
	FitPaddingRatio float64
	ThumbnailURL    string
	HighlightLayer  string
}

type SpatialFocusTargetDTO struct {
	Type     uint32
	Bounds   SpatialBoundsDTO
	Centroid SpatialPointDTO
	GeoJSON  string
	SRID     uint32
}

type WorkspaceLocationDTO struct {
	Address      string
	Province     string
	ProvinceCode string
	WardCode     string
}

type WorkspaceParcelMetaDTO struct {
	MapNumber  string
	LandNumber string
	AreaSqm    float64
}

type ParcelWorkspacePreviewRow struct {
	ParcelID uint64 `gorm:"column:parcel_id"`

	Lat float64 `gorm:"column:lat"`
	Lon float64 `gorm:"column:lon"`

	AreaSqm float64 `gorm:"column:area_sqm"`
	Address string  `gorm:"column:address"`

	MapNumber  string `gorm:"column:map_number"`
	LandNumber string `gorm:"column:land_number"`

	PropertyCode string `gorm:"column:property_code"`
	PropertyUUID string `gorm:"column:property_uuid"`

	GeometryGeoJSON string `gorm:"column:geometry_geojson"`
	GeometryType    string `gorm:"column:geometry_type"`

	MinLon float64 `gorm:"column:min_lon"`
	MinLat float64 `gorm:"column:min_lat"`
	MaxLon float64 `gorm:"column:max_lon"`
	MaxLat float64 `gorm:"column:max_lat"`

	CentroidLat float64 `gorm:"column:centroid_lat"`
	CentroidLon float64 `gorm:"column:centroid_lon"`

	Province     string `gorm:"column:province"`
	ProvinceCode string `gorm:"column:province_code"`
	WardCode     string `gorm:"column:ward_code"`
}

type FollowedParcelActionsDTO struct {
	Focus        bool
	Remove       bool
	Compare      bool
	CreateReport bool
}

type FollowedParcelPreviewDTO struct {
	ID       string
	Type     uint32
	FollowID uint64
	UserID   uint64
	ParcelID uint64

	Title    string
	Subtitle string

	Parcel   WorkspaceParcelMetaDTO
	Location WorkspaceLocationDTO

	Lat float64
	Lon float64

	GeometryPreview GeometryPreviewDTO
	Preview         PreviewHintDTO
	FocusTarget     SpatialFocusTargetDTO

	FollowedAt string

	Actions FollowedParcelActionsDTO
}

type ViewHistoryContextDTO struct {
	Source uint32
	Zoom   float64
}

type ViewHistoryActionsDTO struct {
	Focus  bool
	Remove bool
	Follow bool
}

type ViewHistoryPreviewDTO struct {
	ID        string
	Type      uint32
	HistoryID uint64
	UserID    uint64

	EntityType uint32
	EntityID   uint64
	ParcelID   uint64

	Title    string
	Subtitle string

	Parcel   WorkspaceParcelMetaDTO
	Location WorkspaceLocationDTO

	Lat float64
	Lon float64

	GeometryPreview GeometryPreviewDTO
	Preview         PreviewHintDTO
	FocusTarget     SpatialFocusTargetDTO

	ViewedAt    string
	ViewContext ViewHistoryContextDTO

	Actions ViewHistoryActionsDTO
}

type ReportSpatialPreviewDTO struct {
	Bounds   SpatialBoundsDTO
	Centroid SpatialPointDTO
}

type ReportAssetsDTO struct {
	ThumbnailURL string
	ImageURL     string
	PDFURL       string
	ShareURL     string
}

type ReportMetaDTO struct {
	CreatedAt string
	UpdatedAt string
	ExpiresAt string
	FileSize  uint64
	Format    string
}

type ReportComparisonDTO struct {
	FromPlanName string
	ToPlanName   string
	FromYear     uint32
	ToYear       uint32
}

type ReportActionsDTO struct {
	Focus         bool
	DownloadImage bool
	DownloadPDF   bool
	Share         bool
	Regenerate    bool
	Remove        bool
}

type GeneratedReportPreviewDTO struct {
	ID       string
	Type     uint32
	ReportID uint64
	UserID   uint64

	ReportType uint32
	Status     uint32

	Title    string
	Subtitle string

	Location WorkspaceLocationDTO
	Spatial  ReportSpatialPreviewDTO
	Assets   ReportAssetsDTO
	Meta     ReportMetaDTO

	Comparison ReportComparisonDTO
	Actions    ReportActionsDTO
}

type FollowParcelRequestDTO struct {
	UserID   uint64
	ParcelID uint64
	Note     string
}

type RemoveFollowedParcelRequestDTO struct {
	UserID   uint64
	FollowID uint64
	ParcelID uint64
}

type AddViewHistoryRequestDTO struct {
	UserID     uint64
	EntityType uint32
	EntityID   uint64
	ParcelID   uint64
	Source     uint32
	Zoom       float64
}

type RemoveViewHistoryRequestDTO struct {
	UserID    uint64
	HistoryID uint64
}

type RemoveGeneratedReportRequestDTO struct {
	UserID   uint64
	ReportID uint64
}

type ShareReportResultDTO struct {
	Success  bool
	ShareURL string
	Message  string
}
