package domain

import "time"

type ContentDomain struct {
	ID          uint64 `json:"id"`
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	SortOrder   int    `json:"sortOrder"`
}

type ContentCategory struct {
	ID          uint64     `json:"id"`
	DomainID    uint64     `json:"domainId"`
	DomainCode  string     `json:"domainCode"`
	Code        string     `json:"code"`
	Slug        string     `json:"slug"`
	Name        string     `json:"name"`
	Description string     `json:"description,omitempty"`
	Status      string     `json:"status,omitempty"`
	SortOrder   int        `json:"sortOrder"`
	ItemCount   int64      `json:"itemCount"`
	FollowedAt  *time.Time `json:"followedAt,omitempty"`
}

type FollowedSeoCategoryListFilter struct {
	Page int
	Size int
}

type FollowedSeoCategoryListResult struct {
	Data  []ContentCategory `json:"data"`
	Total int64             `json:"total"`
	Page  int               `json:"page"`
	Size  int               `json:"size"`
}

type CreateSeoCategoryInput struct {
	DomainID    uint64
	DomainCode  string
	Code        string
	Slug        string
	Name        string
	Description string
	Status      string
	SortOrder   int
}

type UpdateSeoCategoryInput struct {
	DomainID    *uint64
	DomainCode  *string
	Code        *string
	Slug        *string
	Name        *string
	Description *string
	Status      *string
	SortOrder   *int
}

type ContentFeedItem struct {
	SEOPageID     uint64     `json:"seoPageId"`
	Slug          string     `json:"slug"`
	Title         string     `json:"title"`
	Description   string     `json:"description,omitempty"`
	Summary       string     `json:"summary,omitempty"`
	CanonicalURL  string     `json:"canonicalUrl"`
	PublishedAt   *time.Time `json:"publishedAt,omitempty"`
	UpdatedAt     *time.Time `json:"updatedAt,omitempty"`
	ImageURL      string     `json:"imageUrl,omitempty"`
	CategoryCode  string     `json:"categoryCode,omitempty"`
	CategoryName  string     `json:"categoryName,omitempty"`
	PlacementCode string     `json:"placementCode,omitempty"`
	Priority      int        `json:"priority"`
	VisibleOn     uint32     `json:"visibleOn"`
}

type FeedResult struct {
	Feed     string            `json:"feed"`
	Computed bool              `json:"computed"`
	Items    []ContentFeedItem `json:"items"`
	Warnings []string          `json:"warnings,omitempty"`
}

type PlanningNewsCategory struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
}

type PlanningNewsItem struct {
	ID              uint64                `json:"id"`
	SEOPageID       uint64                `json:"seoDomainId"`
	Slug            string                `json:"slug"`
	Title           string                `json:"title"`
	Description     string                `json:"description,omitempty"`
	Summary         string                `json:"summary,omitempty"`
	Content         string                `json:"content,omitempty"`
	RenderedHTML    string                `json:"renderedHtml,omitempty"`
	CanonicalURL    string                `json:"canonicalUrl"`
	PublicURL       string                `json:"publicUrl,omitempty"`
	DeepLink        string                `json:"deepLink,omitempty"`
	ImageURL        string                `json:"imageUrl,omitempty"`
	RefSource       string                `json:"refSource,omitempty"`
	RefLabel        string                `json:"refLabel,omitempty"`
	SourceStatus    string                `json:"sourceStatus,omitempty"`
	QualityScore    uint32                `json:"qualityScore"`
	VisibleOn       uint32                `json:"visibleOn"`
	PublishedAt     *time.Time            `json:"publishedAt,omitempty"`
	UpdatedAt       *time.Time            `json:"updatedAt,omitempty"`
	SavedAt         *time.Time            `json:"savedAt,omitempty"`
	CategoryID      *uint64               `json:"categoryId,omitempty"`
	Category        *PlanningNewsCategory `json:"category,omitempty"`
	CategoryCode    string                `json:"categoryCode,omitempty"`
	CategoryName    string                `json:"categoryName,omitempty"`
	Metadata        map[string]any        `json:"metadata,omitempty"`
	RefSnapshotJSON map[string]any        `json:"refSnapshotJson,omitempty"`
	RefType         uint32                `json:"refType,omitempty"`
	RefTypeName     string                `json:"refTypeName,omitempty"`
	RefID           *uint64               `json:"refId,omitempty"`
}

type PlanningNewsListFilter struct {
	Page       int
	Size       int
	Keyword    string
	CategoryID uint64
	ProvinceID uint64
	Location   string
	Range      string
	VisibleOn  uint32
	VisibleOns []uint32
	Sort       string
	SortBy     string
	SortOrder  string
	RefType    uint32
}

type SavedPlanningNewsListFilter struct {
	Page int
	Size int
}

type PlanningNewsListResult struct {
	Data  []PlanningNewsItem `json:"data"`
	Total int64              `json:"total"`
	Page  int                `json:"page"`
	Size  int                `json:"size"`
}

type PlanningNewsDetailResult struct {
	SEODomain   PlanningNewsItem   `json:"seoDomain"`
	RelatedNews []PlanningNewsItem `json:"relatedNews"`
}

type SearchEntityRef struct {
	Kind string `json:"kind"`
	ID   string `json:"id"`
	Key  string `json:"key"`
}

type SearchCandidate struct {
	Entity       SearchEntityRef `json:"entity"`
	Presentation struct {
		Title       string   `json:"title"`
		Subtitle    string   `json:"subtitle,omitempty"`
		Description string   `json:"description,omitempty"`
		Status      string   `json:"status,omitempty"`
		Badges      []string `json:"badges,omitempty"`
	} `json:"presentation"`
	Source struct {
		System  string `json:"system"`
		Dataset string `json:"dataset,omitempty"`
	} `json:"source"`
	Capabilities struct {
		CanOpenDetail   bool `json:"canOpenDetail"`
		CanFocusMap     bool `json:"canFocusMap"`
		CanFollow       bool `json:"canFollow"`
		CanCreateReport bool `json:"canCreateReport"`
		CanCompare      bool `json:"canCompare"`
		CanShare        bool `json:"canShare"`
		CanDownload     bool `json:"canDownload"`
	} `json:"capabilities"`
	Links struct {
		CanonicalURL string `json:"canonicalUrl,omitempty"`
		DetailURL    string `json:"detailUrl,omitempty"`
	} `json:"links"`
	Match struct {
		Type          string   `json:"type,omitempty"`
		MatchedFields []string `json:"matchedFields,omitempty"`
		ReasonCodes   []string `json:"reasonCodes,omitempty"`
	} `json:"match"`
	Rank struct {
		Score       float64  `json:"score"`
		ReasonCodes []string `json:"reasonCodes,omitempty"`
	} `json:"rank"`
	Attributes map[string]any `json:"attributes,omitempty"`
}

type SearchGroup struct {
	Kind  string            `json:"kind"`
	Label string            `json:"label"`
	Total int               `json:"total"`
	Items []SearchCandidate `json:"items"`
}

type ContentSearchResult struct {
	RequestID      string `json:"requestId"`
	Interpretation struct {
		OriginalQuery   string `json:"originalQuery"`
		NormalizedQuery string `json:"normalizedQuery"`
	} `json:"interpretation"`
	Groups   []SearchGroup `json:"groups"`
	Partial  bool          `json:"partial"`
	Warnings []string      `json:"warnings,omitempty"`
}
