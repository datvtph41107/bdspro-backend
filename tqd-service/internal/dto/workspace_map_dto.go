package dto

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

type FollowedParcelActionsDTO struct {
	Focus        bool
	Remove       bool
	Compare      bool
	CreateReport bool
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

// workspace_map_TARGET
type ResolveMapTargetRequestDTO struct {
	Latitude  float64
	Longitude float64
	Radius    float64
}

type ResolveMapTargetResponseDTO struct {
	TargetType uint32
	TargetID   uint64

	Parcel *ParcelInfoResponse
	Region *RegionInfoResponse
}

type RegionInfoResponse struct {
	RegionID  uint64
	LayerID   uint64
	LabelID   uint64
	LandUseID uint64
	LegendID  uint64

	Name        string
	DisplayName string
	Description string

	LayerName        string
	LayerDisplayName string

	LandUseCode  string
	LandUseName  string
	LandUseColor string
	CanBuild     bool

	LabelName  string
	LabelColor string

	LegendColor  string
	LegendType   string
	GeometryType string

	LegalDoc     string
	PlanningName string

	AreaSqm    float64
	AreaHa     float64
	PerimeterM float64

	CenterLat float64
	CenterLon float64

	MinLon float64
	MinLat float64
	MaxLon float64
	MaxLat float64

	GeoJSON string

	Province     string
	ProvinceCode string
	WardCode     string
}

// workspace_PREVIEW
type ParcelWorkspacePreviewRow struct {
	ParcelID uint64

	Lat float64
	Lon float64

	AreaSqm float64
	Address string

	MapNumber  string
	LandNumber string

	GeometryGeoJSON string
	GeometryType    string

	MinLon float64
	MinLat float64
	MaxLon float64
	MaxLat float64

	CentroidLat float64
	CentroidLon float64

	Province     string
	ProvinceCode string
	WardCode     string
}

type RegionWorkspacePreviewRow struct {
	RegionID uint64

	LayerID   uint64
	LabelID   uint64
	LandUseID uint64
	LegendID  uint64

	Name        string
	DisplayName string
	Description string

	LayerName        string
	LayerDisplayName string

	LandUseCode  string
	LandUseName  string
	LandUseColor string
	CanBuild     bool

	LabelName  string
	LabelColor string

	LegendColor  string
	LegendType   string
	GeometryType string

	LegalDoc     string
	PlanningName string

	AreaSqm    float64
	AreaHa     float64
	PerimeterM float64

	CenterLat float64
	CenterLon float64

	MinLon float64
	MinLat float64
	MaxLon float64
	MaxLat float64

	GeoJSON string

	Province     string
	ProvinceCode string
	WardCode     string
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

type ViewHistoryPreviewDTO struct {
	ID   string
	Type uint32

	HistoryID uint64
	UserID    uint64

	EntityType uint32
	EntityID   uint64
	ParcelID   uint64
	RegionID   uint64

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

	ViewCount        uint64
	CountedViewCount uint64

	Actions ViewHistoryActionsDTO
}

type GeneratedReportPreviewDTO struct {
	ID   string
	Type uint32

	ReportID uint64
	UserID   uint64

	ReportType uint32
	Status     uint32

	Title    string
	Subtitle string

	Location   WorkspaceLocationDTO
	Spatial    ReportSpatialPreviewDTO
	Assets     ReportAssetsDTO
	Meta       ReportMetaDTO
	Comparison ReportComparisonDTO
	Actions    ReportActionsDTO
	EntityType uint32
	EntityID   uint64
	ParcelID   uint64
	RegionID   uint64
}
