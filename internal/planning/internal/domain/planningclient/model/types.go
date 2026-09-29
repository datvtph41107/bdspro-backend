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
	ID                uint64           `json:"id"`
	Slug              string           `json:"slug,omitempty"`
	Code              string           `json:"code,omitempty"`
	Name              string           `json:"name"`
	PlanningType      uint32           `json:"planningType"`
	PlanningTypeName  string           `json:"planningTypeName,omitempty"`
	PlanningLevel     uint32           `json:"planningLevel"`
	PlanningLevelName string           `json:"planningLevelName,omitempty"`
	TotalArea         float64          `json:"totalArea,omitempty"`
	JurisdictionID    *uint64          `json:"jurisdictionId,omitempty"`
	JurisdictionName  string           `json:"jurisdictionName,omitempty"`
	LegalStatus       uint32           `json:"legalStatus"`
	LegalStatusName   string           `json:"legalStatusName,omitempty"`
	ValidityStatus    string           `json:"validityStatus,omitempty"`
	Authority         string           `json:"authority,omitempty"`
	DecisionNumber    string           `json:"decisionNumber,omitempty"`
	Summary           string           `json:"summary,omitempty"`
	ResearchScope     string           `json:"researchScope,omitempty"`
	Indicators        json.RawMessage  `json:"indicators"`
	ApprovalDate      *time.Time       `json:"approvalDate,omitempty"`
	EffectiveDate     *time.Time       `json:"effectiveDate,omitempty"`
	ExpiryDate        *time.Time       `json:"expiryDate,omitempty"`
	CurrentVersion    string           `json:"currentVersion,omitempty"`
	Metadata          json.RawMessage  `json:"metadata"`
	ProcessStatus     uint32           `json:"processStatus"`
	CreatedAt         *time.Time       `json:"createdAt,omitempty"`
	UpdatedAt         *time.Time       `json:"updatedAt,omitempty"`
	HasLayer          bool             `json:"hasLayer"`
	LayerCount        int64            `json:"layerCount"`
	Overview          *ProjectOverview `json:"overview,omitempty" gorm:"-"`
	Tags              []ProjectTag     `json:"tags,omitempty" gorm:"-"`
}

type ProjectTag struct {
	Key   string `json:"key"`
	Label string `json:"label"`
	Value string `json:"value"`
}

type Document struct {
	ID                uint64          `json:"id"`
	PlanningProjectID uint64          `json:"planningProjectId"`
	DocumentType      uint32          `json:"documentType"`
	DocumentTypeName  string          `json:"documentTypeName,omitempty"`
	Status            uint32          `json:"status"`
	StatusName        string          `json:"statusName"`
	Code              string          `json:"code,omitempty"`
	Title             string          `json:"title"`
	Description       string          `json:"description,omitempty"`
	Filepath          string          `json:"filepath,omitempty"`
	Thumbnail         string          `json:"thumbnail,omitempty"`
	VersionNo         string          `json:"versionNo,omitempty"`
	ValidityStatus    string          `json:"validityStatus,omitempty"`
	IssueDate         *time.Time      `json:"issueDate,omitempty"`
	EffectiveDate     *time.Time      `json:"effectiveDate,omitempty"`
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
