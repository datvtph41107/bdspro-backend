package domain

import (
	"fmt"
	"strings"
)

type EntityKind string

const (
	EntityKindParcel             EntityKind = "parcel"
	EntityKindPlanningRegion     EntityKind = "planning_region"
	EntityKindAdministrativeUnit EntityKind = "administrative_unit"
	EntityKindPOI                EntityKind = "poi"
	EntityKindPlanningProject    EntityKind = "planning_project"
	EntityKindLegalDocument      EntityKind = "legal_document"
	EntityKindPlanningMap        EntityKind = "planning_map"
	EntityKindMapLayer           EntityKind = "map_layer"
	EntityKindAnalysisReport     EntityKind = "analysis_report"
	EntityKindNewsArticle        EntityKind = "news_article"
)

func (k EntityKind) IsValid() bool {
	switch k {
	case EntityKindParcel,
		EntityKindPlanningRegion,
		EntityKindAdministrativeUnit,
		EntityKindPOI,
		EntityKindPlanningProject,
		EntityKindLegalDocument,
		EntityKindPlanningMap,
		EntityKindMapLayer,
		EntityKindAnalysisReport,
		EntityKindNewsArticle:
		return true
	default:
		return false
	}
}

type EntityRef struct {
	Kind EntityKind `json:"kind"`
	ID   string     `json:"id"`
	Key  string     `json:"key"`
}

func NewEntityRef(kind EntityKind, id string) EntityRef {
	id = strings.TrimSpace(id)
	return EntityRef{Kind: kind, ID: id, Key: string(kind) + ":" + id}
}

func ParseEntityKey(key string) (EntityRef, error) {
	kind, id, ok := strings.Cut(strings.TrimSpace(key), ":")
	if !ok || kind == "" || id == "" {
		return EntityRef{}, fmt.Errorf("invalid canonical entity key")
	}
	entityKind := EntityKind(kind)
	if !entityKind.IsValid() {
		return EntityRef{}, fmt.Errorf("unsupported entity kind %q", kind)
	}
	return NewEntityRef(entityKind, id), nil
}

type Point struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type Bounds struct {
	MinLongitude float64 `json:"minLongitude"`
	MinLatitude  float64 `json:"minLatitude"`
	MaxLongitude float64 `json:"maxLongitude"`
	MaxLatitude  float64 `json:"maxLatitude"`
}

type SpatialSummary struct {
	GeometryType    string  `json:"geometryType,omitempty"`
	Centroid        *Point  `json:"centroid,omitempty"`
	Bounds          *Bounds `json:"bounds,omitempty"`
	GeometryRef     string  `json:"geometryRef,omitempty"`
	GeometryPreview any     `json:"geometryPreview,omitempty"`
}

type Presentation struct {
	Title       string   `json:"title"`
	Subtitle    string   `json:"subtitle,omitempty"`
	Description string   `json:"description,omitempty"`
	Status      string   `json:"status,omitempty"`
	Badges      []string `json:"badges,omitempty"`
}

type Location struct {
	Address      string `json:"address,omitempty"`
	Province     string `json:"province,omitempty"`
	ProvinceCode string `json:"provinceCode,omitempty"`
	Ward         string `json:"ward,omitempty"`
	WardCode     string `json:"wardCode,omitempty"`
}

type Source struct {
	System      string `json:"system"`
	Dataset     string `json:"dataset,omitempty"`
	Authority   string `json:"authority,omitempty"`
	UpdatedAt   string `json:"updatedAt,omitempty"`
	DataQuality string `json:"dataQuality,omitempty"`
}

type Capabilities struct {
	CanOpenDetail   bool `json:"canOpenDetail"`
	CanFocusMap     bool `json:"canFocusMap"`
	CanFollow       bool `json:"canFollow"`
	CanCreateReport bool `json:"canCreateReport"`
	CanCompare      bool `json:"canCompare"`
	CanShare        bool `json:"canShare"`
	CanDownload     bool `json:"canDownload"`
}

type Links struct {
	CanonicalURL string `json:"canonicalUrl,omitempty"`
	DetailURL    string `json:"detailUrl,omitempty"`
	GeometryURL  string `json:"geometryUrl,omitempty"`
}

type Match struct {
	Type          string   `json:"type,omitempty"`
	MatchedFields []string `json:"matchedFields,omitempty"`
	Highlights    []string `json:"highlights,omitempty"`
	ReasonCodes   []string `json:"reasonCodes,omitempty"`
}

type Rank struct {
	Score       float64  `json:"score"`
	ReasonCodes []string `json:"reasonCodes,omitempty"`
}

type EntityCandidate struct {
	Entity       EntityRef       `json:"entity"`
	Presentation Presentation    `json:"presentation"`
	Spatial      *SpatialSummary `json:"spatial,omitempty"`
	Location     *Location       `json:"location,omitempty"`
	Source       Source          `json:"source"`
	Capabilities Capabilities    `json:"capabilities"`
	Links        Links           `json:"links"`
	Match        Match           `json:"match"`
	Rank         Rank            `json:"rank"`
	Attributes   map[string]any  `json:"attributes,omitempty"`
}

type MapContext struct {
	Zoom           float64      `json:"zoom,omitempty"`
	MapMode        string       `json:"mapMode,omitempty"`
	ActiveLayerIDs []string     `json:"activeLayerIds,omitempty"`
	PreferredKinds []EntityKind `json:"preferredKinds,omitempty"`
	Viewport       *Bounds      `json:"viewport,omitempty"`
}

type IdentifyOptions struct {
	ToleranceMeters        float64 `json:"toleranceMeters,omitempty"`
	Limit                  int     `json:"limit,omitempty"`
	IncludeGeometryPreview bool    `json:"includeGeometryPreview,omitempty"`
}

type IdentifyRequest struct {
	Point   Point           `json:"point"`
	Context MapContext      `json:"context"`
	Options IdentifyOptions `json:"options"`
}

type Resolution struct {
	Status     string  `json:"status"`
	PrimaryKey string  `json:"primaryKey,omitempty"`
	Confidence float64 `json:"confidence"`
	Ambiguous  bool    `json:"ambiguous"`
	Strategy   string  `json:"strategy"`
}

type IdentifyResponse struct {
	RequestID  string            `json:"requestId"`
	Resolution Resolution        `json:"resolution"`
	Primary    *EntityCandidate  `json:"primary,omitempty"`
	Candidates []EntityCandidate `json:"candidates"`
	Partial    bool              `json:"partial"`
	Warnings   []string          `json:"warnings,omitempty"`
}

type SearchRequest struct {
	Query              string       `json:"query"`
	Mode               string       `json:"mode,omitempty"`
	EntityKinds        []EntityKind `json:"entityKinds,omitempty"`
	MapContext         MapContext   `json:"mapContext,omitempty"`
	AdministrativeCode string       `json:"administrativeCode,omitempty"`
	Cursor             string       `json:"cursor,omitempty"`
	Limit              int          `json:"limit,omitempty"`
}

type SearchInterpretation struct {
	OriginalQuery   string            `json:"originalQuery"`
	NormalizedQuery string            `json:"normalizedQuery"`
	DetectedIntents []string          `json:"detectedIntents"`
	StructuredTerms map[string]string `json:"structuredTerms,omitempty"`
}

type SearchGroup struct {
	Kind  EntityKind        `json:"kind"`
	Label string            `json:"label"`
	Total int               `json:"total"`
	Items []EntityCandidate `json:"items"`
}

type SearchResponse struct {
	RequestID      string               `json:"requestId"`
	Interpretation SearchInterpretation `json:"interpretation"`
	Groups         []SearchGroup        `json:"groups"`
	NextCursor     string               `json:"nextCursor,omitempty"`
	Partial        bool                 `json:"partial"`
	Warnings       []string             `json:"warnings,omitempty"`
}
