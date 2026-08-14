package domain

import "time"

type PlanningNewsItem struct {
	ID              uint64         `json:"id"`
	SEOPageID       uint64         `json:"seoDomainId"`
	Slug            string         `json:"slug"`
	Title           string         `json:"title"`
	Description     string         `json:"description,omitempty"`
	Summary         string         `json:"summary,omitempty"`
	Content         string         `json:"content,omitempty"`
	RenderedHTML    string         `json:"renderedHtml,omitempty"`
	CanonicalURL    string         `json:"canonicalUrl"`
	PublicURL       string         `json:"publicUrl,omitempty"`
	DeepLink        string         `json:"deepLink,omitempty"`
	RefSource       string         `json:"refSource,omitempty"`
	RefLabel        string         `json:"refLabel,omitempty"`
	SourceStatus    string         `json:"sourceStatus,omitempty"`
	PublishedAt     *time.Time     `json:"publishedAt,omitempty"`
	UpdatedAt       *time.Time     `json:"updatedAt,omitempty"`
	SavedAt         *time.Time     `json:"savedAt,omitempty"`
	Metadata        map[string]any `json:"metadata,omitempty"`
	RefSnapshotJSON map[string]any `json:"refSnapshotJson,omitempty"`
}

type PlanningNewsListFilter struct {
	Page int
	Size int
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
