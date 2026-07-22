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
	ID          uint64 `json:"id"`
	DomainID    uint64 `json:"domainId"`
	DomainCode  string `json:"domainCode"`
	Code        string `json:"code"`
	Slug        string `json:"slug"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	SortOrder   int    `json:"sortOrder"`
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
	CategoryCode  string     `json:"categoryCode,omitempty"`
	CategoryName  string     `json:"categoryName,omitempty"`
	PlacementCode string     `json:"placementCode,omitempty"`
	Priority      int        `json:"priority"`
}

type FeedResult struct {
	Feed     string            `json:"feed"`
	Computed bool              `json:"computed"`
	Items    []ContentFeedItem `json:"items"`
	Warnings []string          `json:"warnings,omitempty"`
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
		System      string `json:"system"`
		Dataset     string `json:"dataset,omitempty"`
		Authority   string `json:"authority,omitempty"`
		DataQuality string `json:"dataQuality,omitempty"`
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
		OriginalQuery   string   `json:"originalQuery"`
		NormalizedQuery string   `json:"normalizedQuery"`
		DetectedIntents []string `json:"detectedIntents"`
	} `json:"interpretation"`
	Groups   []SearchGroup `json:"groups"`
	Partial  bool          `json:"partial"`
	Warnings []string      `json:"warnings,omitempty"`
}
