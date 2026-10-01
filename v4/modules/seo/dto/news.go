// Package dto defines request/response shapes for modules/seo HTTP API.
// JSON shapes match the CRM public API for client/public compatibility.
package dto

import "time"

// ── Categories ────────────────────────────────────────────────────────────────

type ContentCategory struct {
	ID          uint64 `json:"id"`
	Code        string `json:"code"`
	Slug        string `json:"slug"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Status      string `json:"status,omitempty"`
	SortOrder   int    `json:"sortOrder"`
	ItemCount   int64  `json:"itemCount"`
}

type CreateCategoryRequest struct {
	Code        string `json:"code"`
	Slug        string `json:"slug"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      string `json:"status"`
	SortOrder   int    `json:"sortOrder"`
}

type UpdateCategoryRequest struct {
	ID          uint64  `path:"id"`
	Code        *string `json:"code"`
	Slug        *string `json:"slug"`
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Status      *string `json:"status"`
	SortOrder   *int    `json:"sortOrder"`
}

// ── Planning news ─────────────────────────────────────────────────────────────

type PlanningNewsCategory struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
}

type PlanningNewsItem struct {
	ID              uint64                 `json:"id"`
	SEOPageID       uint64                 `json:"seoDomainId"`
	Slug            string                 `json:"slug"`
	Title           string                 `json:"title"`
	Description     string                 `json:"description,omitempty"`
	Summary         string                 `json:"summary,omitempty"`
	Content         string                 `json:"content,omitempty"`
	RenderedHTML    string                 `json:"renderedHtml,omitempty"`
	CanonicalURL    string                 `json:"canonicalUrl"`
	PublicURL       string                 `json:"publicUrl,omitempty"`
	DeepLink        string                 `json:"deepLink,omitempty"`
	ImageURL        string                 `json:"imageUrl,omitempty"`
	RefSource       string                 `json:"refSource,omitempty"`
	RefLabel        string                 `json:"refLabel,omitempty"`
	SourceStatus    string                 `json:"sourceStatus,omitempty"`
	VisibleOn       uint32                 `json:"visibleOn,omitempty"`
	RefType         uint32                 `json:"refType,omitempty"`
	RefTypeName     string                 `json:"refTypeName,omitempty"`
	RefID           *uint64                `json:"refId,omitempty"`
	PublishedAt     *time.Time             `json:"publishedAt,omitempty"`
	UpdatedAt       *time.Time             `json:"updatedAt,omitempty"`
	SavedAt         *time.Time             `json:"savedAt,omitempty"`
	CategoryID      *uint64                `json:"categoryId,omitempty"`
	CategoryCode    string                 `json:"categoryCode,omitempty"`
	CategoryName    string                 `json:"categoryName,omitempty"`
	Category        *PlanningNewsCategory  `json:"category,omitempty"`
	ProvinceID      *uint64                `json:"provinceId,omitempty"`
	Metadata        map[string]any         `json:"metadata,omitempty"`
	RefSnapshotJSON map[string]any         `json:"refSnapshotJson,omitempty"`
}

type PlanningNewsListFilter struct {
	Page        int    `form:"page"`
	Size        int    `form:"size"`
	Keyword     string `form:"keyword"`
	CategoryID  uint64 `form:"categoryId"`
	ProvinceID  uint64 `form:"provinceId"`
	RefType     uint32 `form:"refType"`
	VisibleOn   uint32 `form:"visibleOn"`
	Sort        string `form:"sort"`
	SortBy      string `form:"sortBy"`
	SortOrder   string `form:"sortOrder"`
	LastDates   uint32 `form:"lastDates"`
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

type GetNewsDetailRequest struct {
	ID uint64 `path:"id"`
}

// ── Click tracking ────────────────────────────────────────────────────────────

// RecordClickRequest mirrors CRM TrackSeoClickRequest:
// prefer id (or seoDomainId), else refType+refId.
type RecordClickRequest struct {
	ID          uint64 `json:"id"`
	SeoDomainID uint64 `json:"seoDomainId"`
	RefType     uint32 `json:"refType"`
	RefID       uint64 `json:"refId"`
	UserID      uint64 `json:"userId,omitempty"`
	IPAddr      string `json:"ipAddr,omitempty"`
}

// RecordClickResult is the Accepted payload for public SEO click tracking.
type RecordClickResult struct {
	Accepted bool   `json:"accepted"`
	ID       uint64 `json:"id,omitempty"`
	RefType  uint32 `json:"refType,omitempty"`
	RefID    uint64 `json:"refId,omitempty"`
}

// ── Saved news ────────────────────────────────────────────────────────────────

type SavedNewsFilter struct {
	Page int `form:"page"`
	Size int `form:"size"`
}

type NewsIDPath struct {
	ID uint64 `path:"id"`
}

// ── Content search ────────────────────────────────────────────────────────────

type ContentSearchRequest struct {
	Query string `form:"query"`
	Kind  string `form:"kind"`
	Limit int    `form:"limit"`
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
		CanOpenDetail bool `json:"canOpenDetail"`
		CanFocusMap   bool `json:"canFocusMap"`
		CanShare      bool `json:"canShare"`
	} `json:"capabilities"`
	Links struct {
		CanonicalURL string `json:"canonicalUrl,omitempty"`
		DetailURL    string `json:"detailUrl,omitempty"`
	} `json:"links"`
	Match struct {
		Type          string   `json:"type,omitempty"`
		MatchedFields []string `json:"matchedFields,omitempty"`
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
	Groups  []SearchGroup `json:"groups"`
	Partial bool          `json:"partial"`
}
