package domain

import (
	"encoding/json"
	"time"
)

const (
	DocumentStatusNormal    uint32 = 10
	DocumentStatusImportant uint32 = 20
)

type ProjectListFilter struct {
	Page           int
	Size           int
	Search         string
	PlanningType   string
	PlanningLevel  string
	ValidityStatus string
	ProcessStatus  string
	JurisdictionID uint64
	Area           string
	LegalStatus    string
	HasLayer       *bool
	SortBy         string
	SortOrder      string
	Decision       string
	UpdatedFrom    *time.Time
	UpdatedTo      *time.Time
	TagIDs         string
	ProvinceIDs    string
	// AdmID is administrative unit composite ID ({sourceId}{1=tỉnh|2=phường})
	AdmID          uint64
	// ParcelID lọc đồ án liên quan thửa.
	ParcelID       uint64
	// LayerIDs lọc theo nhiều layer (OR) trong phạm vi parcel.
	LayerIDs       []uint64
}

type ProjectOverview struct {
	Description            string          `json:"description,omitempty"`
	TotalArea              float64         `json:"totalArea,omitempty"`
	ResearchScope          string          `json:"researchScope,omitempty"`
	Indicators             json.RawMessage `json:"indicators"`
	DocumentCount          int64           `json:"documentCount"`
	ImportantDocumentCount int64           `json:"importantDocumentCount"`
	MapLayerCount          int64           `json:"mapLayerCount"`
	LegalEventCount        int64           `json:"legalEventCount"`
	LatestUpdatedAt        *time.Time      `json:"latestUpdatedAt,omitempty"`
	HasLayer               bool            `json:"hasLayer"`
}

type Project struct {
	ID                uint64            `json:"id"`
	Slug              string            `json:"slug,omitempty"`
	Code              string            `json:"code,omitempty"`
	Name              string            `json:"name"`
	PlanningType      uint32            `json:"planningType"`
	PlanningTypeName  string            `json:"planningTypeName,omitempty"`
	PlanningLevel     uint32            `json:"planningLevel"`
	PlanningLevelName string            `json:"planningLevelName,omitempty"`
	TotalArea         float64           `json:"totalArea,omitempty"`
	JurisdictionID    *uint64           `json:"jurisdictionId,omitempty"`
	JurisdictionName  string            `json:"jurisdictionName,omitempty"`
	LegalStatus       uint32            `json:"legalStatus"`
	LegalStatusName   string            `json:"legalStatusName,omitempty"`
	ValidityStatus    string            `json:"validityStatus,omitempty"`
	Authority         string            `json:"authority,omitempty"`
	DecisionNumber    string            `json:"decisionNumber,omitempty"`
	Summary           string            `json:"summary,omitempty"`
	ResearchScope     string            `json:"researchScope,omitempty"`
	Indicators        json.RawMessage   `json:"indicators"`
	ApprovalDate      *time.Time        `json:"approvalDate,omitempty"`
	EffectiveDate     *time.Time        `json:"effectiveDate,omitempty"`
	ExpiryDate        *time.Time        `json:"expiryDate,omitempty"`
	CurrentVersion    string            `json:"currentVersion,omitempty"`
	Metadata          json.RawMessage   `json:"metadata"`
	ProcessStatus     uint32            `json:"processStatus"`
	CreatedAt         *time.Time        `json:"createdAt,omitempty"`
	UpdatedAt         *time.Time        `json:"updatedAt,omitempty"`
	HasLayer          bool              `json:"hasLayer"`
	LayerCount        int64             `json:"layerCount"`
	DocumentCount     int64             `json:"documentCount"`
	Overview          *ProjectOverview  `json:"overview,omitempty" gorm:"-"`
	Tags              []ProjectTag      `json:"tags,omitempty" gorm:"-"`
	ProjectTags       []Tag             `json:"projectTags,omitempty" gorm:"-"`
	Provinces         []ProjectProvince `json:"provinces,omitempty" gorm:"-"`
	Wards             []ProjectWard     `json:"wards,omitempty" gorm:"-"`
}

type ProjectProvince struct {
	ID        string          `json:"id"`
	FullName  string          `json:"fullName"`
	ShortName string          `json:"shortName,omitempty"`
	Code      string          `json:"code,omitempty"`
	QHOrder   int32           `json:"qhOrder"`
	Zoom      uint32          `json:"zoom"`
	Lat       float64         `json:"lat,omitempty"`
	Lng       float64         `json:"lng,omitempty"`
	Geometry  json.RawMessage `json:"geometry,omitempty"`
}

type ProjectWard struct {
	ID         string          `json:"id"`
	FullName   string          `json:"fullName"`
	ShortName  string          `json:"shortName,omitempty"`
	Code       string          `json:"code,omitempty"`
	ProvinceID string          `json:"provinceId,omitempty"`
	Lat        float64         `json:"lat,omitempty"`
	Lng        float64         `json:"lng,omitempty"`
	Geometry   json.RawMessage `json:"geometry,omitempty"`
}

type Tag struct {
	ID       uint64 `json:"id"`
	Name     string `json:"name"`
	Code     string `json:"code,omitempty"`
	Type     uint32 `json:"type"`
	TypeName string `json:"typeName,omitempty"`
}

type ProjectTag struct {
	Key   string `json:"key"`
	Label string `json:"label"`
	Value string `json:"value"`
}

type DocumentPlanningProject struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
}

type Document struct {
	ID                uint64                   `json:"id"`
	PlanningProjectID uint64                   `json:"planningProjectId"`
	PlanningProject   *DocumentPlanningProject `json:"planningProject,omitempty" gorm:"-"`
	DocumentType      uint32          `json:"documentType"`
	DocumentTypeName  string          `json:"documentTypeName,omitempty"`
	Status            uint32          `json:"status"`
	StatusName        string          `json:"statusName"`
	Code              string          `json:"code,omitempty"`
	Title             string          `json:"title"`
	Description       string          `json:"description,omitempty"`
	Filepath          string          `json:"filepath,omitempty"`
	RelativePath      string          `json:"relativePath,omitempty"`
	Ext               string          `json:"ext,omitempty"`
	Thumbnail         string          `json:"thumbnail,omitempty"`
	VersionNo         string          `json:"versionNo,omitempty"`
	ValidityStatus    string          `json:"validityStatus,omitempty"`
	IssueDate         *time.Time      `json:"issueDate,omitempty"`
	EffectiveDate     *time.Time      `json:"effectiveDate,omitempty"`
	Indicators        json.RawMessage `json:"indicators"`
	Authority         string          `json:"authority,omitempty"`
	Sector            string          `json:"sector,omitempty"`
	Field             string          `json:"field,omitempty"`
	Signer            string          `json:"signer,omitempty"`
	Metadata          json.RawMessage `json:"metadata"`
	CreatedAt         *time.Time      `json:"createdAt,omitempty"`
	UpdatedAt         *time.Time      `json:"updatedAt,omitempty"`
}

type Event struct {
	ID                uint64          `json:"id"`
	PlanningProjectID uint64          `json:"planningProjectId"`
	EventType         string          `json:"eventType,omitempty"`
	EventName         string          `json:"eventName"`
	Description       string          `json:"description,omitempty"`
	EventDate         *time.Time      `json:"eventDate,omitempty"`
	DocumentID        *uint64         `json:"documentId,omitempty"`
	SourceURL         string          `json:"sourceUrl,omitempty"`
	Metadata          json.RawMessage `json:"metadata"`
	Document          *Document       `json:"document,omitempty"`
	CreatedAt         *time.Time      `json:"createdAt,omitempty"`
	UpdatedAt         *time.Time      `json:"updatedAt,omitempty"`
}

type Follow struct {
	ID                uint64     `json:"followId"`
	UserID            uint64     `json:"-"`
	PlanningProjectID uint64     `json:"planningProjectId"`
	Note              string     `json:"note,omitempty"`
	FollowedAt        *time.Time `json:"followedAt,omitempty"`
}

type FollowedProject struct {
	FollowID   uint64     `json:"followId"`
	FollowedAt *time.Time `json:"followedAt,omitempty"`
	Note       string     `json:"note,omitempty"`
	Project    Project    `json:"project"`
}

type Page[T any] struct {
	Data  []T   `json:"data"`
	Total int64 `json:"total"`
	Page  int   `json:"page"`
	Size  int   `json:"size"`
}

type ProjectLayer struct {
	ID                uint64          `json:"id"`
	PlanningProjectID uint64          `json:"planningProjectId"`
	Name              string          `json:"name"`
	DisplayName       string          `json:"displayName"`
	Description       string          `json:"description,omitempty"`
	Type              uint32          `json:"type"`
	Status            uint32          `json:"status"`
	Visible           bool            `json:"visible"`
	DefaultVisible    bool            `json:"defaultVisible"`
	DisplayOrder      int             `json:"displayOrder"`
	MinZoom           uint32          `json:"minZoom"`
	MaxZoom           uint32          `json:"maxZoom"`
	SourceCode        string          `json:"sourceCode,omitempty"`
	SourceType        string          `json:"sourceType,omitempty"`
	LayerURL          string          `json:"layerUrl,omitempty"`
	StyleConfig       json.RawMessage `json:"styleConfig"`
	ImageURL          string          `json:"imageUrl,omitempty"`
	ThumbnailURL      string          `json:"thumbnailUrl,omitempty"`
	EffectiveDate     *time.Time      `json:"effectiveDate,omitempty"`
	ExpiryDate        *time.Time      `json:"expiryDate,omitempty"`
	UpdatedAt         *time.Time      `json:"updatedAt,omitempty"`
}

type ProvinceListFilter struct {
	Scope  uint32
	Search string
}

const ProvinceScopeQH uint32 = 10

type ProvinceItem struct {
	ID             string  `json:"id" gorm:"column:id"`
	FullName       string  `json:"fullName" gorm:"column:full_name"`
	ShortName      string  `json:"shortName,omitempty" gorm:"column:short_name"`
	Code           string  `json:"code,omitempty" gorm:"column:code"`
	Lat            float64 `json:"lat,omitempty" gorm:"column:lat"`
	Lng            float64 `json:"lng,omitempty" gorm:"column:lng"`
	WardCount      int32   `json:"wardCount" gorm:"column:ward_count"`
	QHOrder        int32   `json:"qhOrder" gorm:"column:qh_order"`
	Zoom           uint32  `json:"zoom" gorm:"column:zoom"`
	Description    string  `json:"description,omitempty" gorm:"column:description"`
	UnitType       uint32  `json:"unitType" gorm:"column:unit_type"`
	UnitTypeName   string  `json:"unitTypeName,omitempty" gorm:"-"`
	DataSource     uint32  `json:"dataSource" gorm:"column:data_source"`
	DataSourceName string  `json:"dataSourceName,omitempty" gorm:"-"`
	PlanDocID      *uint64 `json:"planDocId,omitempty" gorm:"column:plan_doc_id"`
}
