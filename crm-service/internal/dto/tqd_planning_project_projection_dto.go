package dto

import "time"

type TqdPlanningProjectProjection struct {
	ID                uint64                       `json:"id"`
	CanonicalKey      string                       `json:"canonicalKey"`
	Slug              string                       `json:"slug"`
	CanonicalPath     string                       `json:"canonicalPath"`
	Code              string                       `json:"code,omitempty"`
	Name              string                       `json:"name"`
	PlanningType      uint32                       `json:"planningType"`
	PlanningTypeName  string                       `json:"planningTypeName,omitempty"`
	PlanningLevel     uint32                       `json:"planningLevel"`
	PlanningLevelName string                       `json:"planningLevelName,omitempty"`
	TotalArea         float64                      `json:"totalArea,omitempty"`
	LegalStatus       uint32                       `json:"legalStatus"`
	LegalStatusName   string                       `json:"legalStatusName,omitempty"`
	ValidityStatus    string                       `json:"validityStatus,omitempty"`
	Authority         string                       `json:"authority,omitempty"`
	Summary           string                       `json:"summary,omitempty"`
	ResearchScope     string                       `json:"researchScope,omitempty"`
	ApprovalDate      *time.Time                   `json:"approvalDate,omitempty"`
	EffectiveDate     *time.Time                   `json:"effectiveDate,omitempty"`
	ExpiryDate        *time.Time                   `json:"expiryDate,omitempty"`
	CurrentVersion    string                       `json:"currentVersion,omitempty"`
	JurisdictionID    *uint64                      `json:"jurisdictionId,omitempty"`
	JurisdictionName  string                       `json:"jurisdictionName,omitempty"`
	Preview           TqdPlanningPreview           `json:"preview"`
	Indicators        []TqdPlanningIndicator       `json:"indicators"`
	Events            []TqdPlanningEvent           `json:"events"`
	Documents         []TqdPlanningDocument        `json:"documents"`
	RelatedProjects   []TqdPlanningProjectRelation `json:"relatedProjects"`
	Source            TqdPlanningProjectionSource  `json:"source"`
}

type TqdPlanningBounds struct {
	MinLongitude float64 `json:"minLongitude"`
	MinLatitude  float64 `json:"minLatitude"`
	MaxLongitude float64 `json:"maxLongitude"`
	MaxLatitude  float64 `json:"maxLatitude"`
}

type TqdPlanningPoint struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type TqdPlanningPreview struct {
	GeometryType string             `json:"geometryType,omitempty"`
	Centroid     *TqdPlanningPoint  `json:"centroid,omitempty"`
	Bounds       *TqdPlanningBounds `json:"bounds,omitempty"`
	PaddingRatio float64            `json:"paddingRatio"`
	FitMode      string             `json:"fitMode"`
	Completeness string             `json:"completeness"`
}

type TqdPlanningIndicator struct {
	Code        string `json:"code,omitempty"`
	Name        string `json:"name"`
	Value       string `json:"value"`
	Unit        string `json:"unit,omitempty"`
	Description string `json:"description,omitempty"`
}

type TqdPlanningEvent struct {
	ID          uint64     `json:"id"`
	EventType   string     `json:"eventType,omitempty"`
	Name        string     `json:"name"`
	Description string     `json:"description,omitempty"`
	EventDate   *time.Time `json:"eventDate,omitempty"`
	DocumentID  *uint64    `json:"documentId,omitempty"`
	SourceURL   string     `json:"sourceUrl,omitempty"`
}

type TqdPlanningDocument struct {
	ID               uint64     `json:"id"`
	DocumentType     uint32     `json:"documentType"`
	DocumentTypeName string     `json:"documentTypeName,omitempty"`
	Code             string     `json:"code,omitempty"`
	Title            string     `json:"title"`
	Description      string     `json:"description,omitempty"`
	Filepath         string     `json:"filepath,omitempty"`
	Thumbnail        string     `json:"thumbnail,omitempty"`
	VersionNo        string     `json:"versionNo,omitempty"`
	ValidityStatus   string     `json:"validityStatus,omitempty"`
	IssueDate        *time.Time `json:"issueDate,omitempty"`
	EffectiveDate    *time.Time `json:"effectiveDate,omitempty"`
}

type TqdPlanningProjectRelation struct {
	ID         uint64 `json:"id"`
	Slug       string `json:"slug,omitempty"`
	Code       string `json:"code,omitempty"`
	Name       string `json:"name"`
	PublicPath string `json:"publicPath"`
	MapPath    string `json:"mapPath"`
}

type TqdPlanningProjectionSource struct {
	System         string     `json:"system"`
	Dataset        string     `json:"dataset"`
	Authority      string     `json:"authority"`
	AuthorityName  string     `json:"authorityName,omitempty"`
	AuthorityClass string     `json:"authorityClass"`
	DataQuality    string     `json:"dataQuality"`
	UpdatedAt      *time.Time `json:"updatedAt,omitempty"`
	Limitations    []string   `json:"limitations,omitempty"`
}
