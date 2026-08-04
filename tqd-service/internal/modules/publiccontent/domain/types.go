package domain

import "time"

type Bounds struct {
	MinLongitude float64 `json:"minLongitude"`
	MinLatitude  float64 `json:"minLatitude"`
	MaxLongitude float64 `json:"maxLongitude"`
	MaxLatitude  float64 `json:"maxLatitude"`
}

type Point struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type PreviewContext struct {
	GeometryType string  `json:"geometryType,omitempty"`
	Centroid     *Point  `json:"centroid,omitempty"`
	Bounds       *Bounds `json:"bounds,omitempty"`
	PaddingRatio float64 `json:"paddingRatio"`
	FitMode      string  `json:"fitMode"`
	Completeness string  `json:"completeness"`
}

type RelatedPlanningProject struct {
	ID                uint64     `json:"id"`
	Code              string     `json:"code,omitempty"`
	Slug              string     `json:"slug,omitempty"`
	Name              string     `json:"name"`
	PlanningType      uint32     `json:"planningType,omitempty"`
	PlanningTypeName  string     `json:"planningTypeName,omitempty"`
	LegalStatus       uint32     `json:"legalStatus,omitempty"`
	LegalStatusName   string     `json:"legalStatusName,omitempty"`
	ValidityStatus    string     `json:"validityStatus,omitempty"`
	TotalArea         float64    `json:"totalArea,omitempty"`
	IntersectAreaSqm  float64    `json:"intersectAreaSqm,omitempty"`
	IntersectionRatio float64    `json:"intersectionRatio,omitempty"`
	UpdatedAt         *time.Time `json:"updatedAt,omitempty"`
	PublicPath        string     `json:"publicPath"`
	DossierPath       string     `json:"dossierPath"`
	MapPath           string     `json:"mapPath"`
}

type ParcelQuickView struct {
	ParcelID        uint64                   `json:"parcelId"`
	CanonicalKey    string                   `json:"canonicalKey"`
	MapNumber       string                   `json:"mapNumber,omitempty"`
	LandNumber      string                   `json:"landNumber,omitempty"`
	PropertyCode    string                   `json:"propertyCode,omitempty"`
	Address         string                   `json:"address,omitempty"`
	TotalAreaSqm    float64                  `json:"totalAreaSqm,omitempty"`
	ProvinceCode    string                   `json:"provinceCode,omitempty"`
	WardCode        string                   `json:"wardCode,omitempty"`
	IsVerified      bool                     `json:"isVerified"`
	Completeness    string                   `json:"completeness"`
	Warnings        []string                 `json:"warnings,omitempty"`
	Preview         PreviewContext           `json:"preview"`
	RelatedProjects []RelatedPlanningProject `json:"relatedPlanningProjects"`
	UpdatedAt       *time.Time               `json:"updatedAt,omitempty"`
}

type PlanningEvent struct {
	ID               uint64     `json:"id"`
	EventType        string     `json:"eventType,omitempty"`
	Name             string     `json:"name"`
	Description      string     `json:"description,omitempty"`
	EventDate        *time.Time `json:"eventDate,omitempty"`
	DocumentID       *uint64    `json:"documentId,omitempty"`
	DocumentCode     string     `json:"documentCode,omitempty"`
	DocumentTitle    string     `json:"documentTitle,omitempty"`
	IssuingAuthority string     `json:"issuingAuthority,omitempty"`
	SourceURL        string     `json:"sourceUrl,omitempty"`
}

type PlanningDocument struct {
	ID               uint64     `json:"id"`
	DocumentType     uint32     `json:"documentType"`
	DocumentTypeName string     `json:"documentTypeName,omitempty"`
	Status           uint32     `json:"status"`
	StatusName       string     `json:"statusName"`
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

type PlanningIndicator struct {
	Code        string `json:"code,omitempty"`
	Name        string `json:"name"`
	Value       string `json:"value"`
	Unit        string `json:"unit,omitempty"`
	Description string `json:"description,omitempty"`
}

type PlanningLayer struct {
	ID                uint64     `json:"id"`
	PlanningProjectID uint64     `json:"planningProjectId"`
	Name              string     `json:"name"`
	DisplayName       string     `json:"displayName"`
	Description       string     `json:"description,omitempty"`
	Type              uint32     `json:"type"`
	Status            uint32     `json:"status"`
	Visible           bool       `json:"visible"`
	DefaultVisible    bool       `json:"defaultVisible"`
	MinZoom           uint32     `json:"minZoom"`
	MaxZoom           uint32     `json:"maxZoom"`
	ImageURL          string     `json:"imageUrl,omitempty"`
	ThumbnailURL      string     `json:"thumbnailUrl,omitempty"`
	EffectiveDate     *time.Time `json:"effectiveDate,omitempty"`
	ExpiryDate        *time.Time `json:"expiryDate,omitempty"`
	UpdatedAt         *time.Time `json:"updatedAt,omitempty"`
}

type PlanningProjectRelation struct {
	ID         uint64 `json:"id"`
	Slug       string `json:"slug,omitempty"`
	Code       string `json:"code,omitempty"`
	Name       string `json:"name"`
	PublicPath string `json:"publicPath"`
	MapPath    string `json:"mapPath"`
}

type PlanningProjectProjection struct {
	ID                uint64                    `json:"id"`
	CanonicalKey      string                    `json:"canonicalKey"`
	Slug              string                    `json:"slug"`
	CanonicalPath     string                    `json:"canonicalPath"`
	Code              string                    `json:"code,omitempty"`
	Name              string                    `json:"name"`
	PlanningType      uint32                    `json:"planningType"`
	PlanningTypeName  string                    `json:"planningTypeName,omitempty"`
	PlanningLevel     uint32                    `json:"planningLevel"`
	PlanningLevelName string                    `json:"planningLevelName,omitempty"`
	TotalArea         float64                   `json:"totalArea,omitempty"`
	LegalStatus       uint32                    `json:"legalStatus"`
	LegalStatusName   string                    `json:"legalStatusName,omitempty"`
	ValidityStatus    string                    `json:"validityStatus,omitempty"`
	Authority         string                    `json:"authority,omitempty"`
	DecisionNumber    string                    `json:"decisionNumber,omitempty"`
	Summary           string                    `json:"summary,omitempty"`
	ResearchScope     string                    `json:"researchScope,omitempty"`
	ApprovalDate      *time.Time                `json:"approvalDate,omitempty"`
	EffectiveDate     *time.Time                `json:"effectiveDate,omitempty"`
	ExpiryDate        *time.Time                `json:"expiryDate,omitempty"`
	CurrentVersion    string                    `json:"currentVersion,omitempty"`
	JurisdictionID    *uint64                   `json:"jurisdictionId,omitempty"`
	JurisdictionName  string                    `json:"jurisdictionName,omitempty"`
	Preview           PreviewContext            `json:"preview"`
	Indicators        []PlanningIndicator       `json:"indicators"`
	Events            []PlanningEvent           `json:"events"`
	Documents         []PlanningDocument        `json:"documents"`
	Layers            []PlanningLayer           `json:"layers"`
	RelatedProjects   []PlanningProjectRelation `json:"relatedProjects"`
	Source            ProjectionSource          `json:"source"`
}

type ProjectionSource struct {
	System         string     `json:"system"`
	Dataset        string     `json:"dataset"`
	Authority      string     `json:"authority"`
	AuthorityName  string     `json:"authorityName,omitempty"`
	AuthorityClass string     `json:"authorityClass"`
	DataQuality    string     `json:"dataQuality"`
	UpdatedAt      *time.Time `json:"updatedAt,omitempty"`
	Limitations    []string   `json:"limitations,omitempty"`
}
