package dto

import "encoding/json"

const SeoPublicDocumentSchemaVersion = "qhpro-public-seo-document/v1"

type SeoPublicDocument struct {
	SchemaVersion   string              `json:"schemaVersion"`
	ModuleID        string              `json:"moduleId"`
	ResourceType    string              `json:"resourceType"`
	Identity        SeoPublicIdentity   `json:"identity"`
	Heading         SeoPublicHeading    `json:"heading"`
	Breadcrumbs     []SeoBreadcrumb     `json:"breadcrumbs"`
	Facts           []SeoFact           `json:"facts"`
	Sections        []SeoSection        `json:"sections"`
	RelatedLinks    []SeoRelatedLink    `json:"relatedLinks"`
	RelatedEntities []SeoRelatedEntity  `json:"relatedEntities,omitempty"`
	Media           *SeoMediaBundle     `json:"media,omitempty"`
	FAQ             *SeoFAQ             `json:"faq,omitempty"`
	MapTarget       *SeoMapTarget       `json:"mapTarget,omitempty"`
	AnswerReadiness *SeoAnswerReadiness `json:"answerReadiness,omitempty"`
	SEO             SeoPublicSEO        `json:"seo"`
	Source          SeoPublicSource     `json:"source"`
}

type SeoPublicIdentity struct {
	ID            string `json:"id"`
	Slug          string `json:"slug"`
	CanonicalPath string `json:"canonicalPath"`
}

type SeoPublicHeading struct {
	Eyebrow string `json:"eyebrow,omitempty"`
	Title   string `json:"title"`
	Summary string `json:"summary"`
}

type SeoBreadcrumb struct {
	Label string `json:"label"`
	Path  string `json:"path"`
}

type SeoFact struct {
	ID          string `json:"id"`
	Label       string `json:"label"`
	Value       string `json:"value"`
	Description string `json:"description,omitempty"`
}

type SeoSection struct {
	ID     string            `json:"id"`
	Title  string            `json:"title"`
	Blocks []SeoContentBlock `json:"blocks"`
}

type SeoContentBlock struct {
	Type          string            `json:"type"`
	Text          string            `json:"text,omitempty"`
	Items         []string          `json:"-"`
	Tone          string            `json:"tone,omitempty"`
	Title         string            `json:"title,omitempty"`
	TimelineItems []SeoTimelineItem `json:"-"`
}

func (b *SeoContentBlock) UnmarshalJSON(data []byte) error {
	type blockWire struct {
		Type  string          `json:"type"`
		Text  string          `json:"text,omitempty"`
		Items json.RawMessage `json:"items,omitempty"`
		Tone  string          `json:"tone,omitempty"`
		Title string          `json:"title,omitempty"`
	}
	var wire blockWire
	if err := json.Unmarshal(data, &wire); err != nil {
		return err
	}
	b.Type = wire.Type
	b.Text = wire.Text
	b.Tone = wire.Tone
	b.Title = wire.Title
	b.Items = nil
	b.TimelineItems = nil
	if len(wire.Items) == 0 || string(wire.Items) == "null" {
		return nil
	}
	if wire.Type == "timeline" {
		return json.Unmarshal(wire.Items, &b.TimelineItems)
	}
	return json.Unmarshal(wire.Items, &b.Items)
}

func (b SeoContentBlock) MarshalJSON() ([]byte, error) {
	type blockWire struct {
		Type  string `json:"type"`
		Text  string `json:"text,omitempty"`
		Items any    `json:"items,omitempty"`
		Tone  string `json:"tone,omitempty"`
		Title string `json:"title,omitempty"`
	}
	var items any
	if b.Type == "timeline" {
		if len(b.TimelineItems) > 0 {
			items = b.TimelineItems
		}
	} else if len(b.Items) > 0 {
		items = b.Items
	}
	return json.Marshal(blockWire{Type: b.Type, Text: b.Text, Items: items, Tone: b.Tone, Title: b.Title})
}

type SeoTimelineItem struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Date        string `json:"date,omitempty"`
	Description string `json:"description,omitempty"`
}

type SeoRelatedLink struct {
	ID          string `json:"id"`
	Label       string `json:"label"`
	Path        string `json:"path"`
	Description string `json:"description,omitempty"`
}

type SeoRelatedEntity struct {
	ID           string `json:"id"`
	RelationType string `json:"relationType"`
	EntityType   string `json:"entityType"`
	EntityID     string `json:"entityId"`
	Title        string `json:"title"`
	Description  string `json:"description,omitempty"`
	PublicPath   string `json:"publicPath,omitempty"`
	MapPath      string `json:"mapPath,omitempty"`
}

type SeoMedia struct {
	ID     string     `json:"id"`
	Role   string     `json:"role"`
	URL    string     `json:"url,omitempty"`
	Alt    string     `json:"alt"`
	Width  int        `json:"width,omitempty"`
	Height int        `json:"height,omitempty"`
	Bounds [4]float64 `json:"bounds,omitempty"`
	Order  int        `json:"order"`
}

type SeoMediaBundle struct {
	Primary      *SeoMedia  `json:"primary,omitempty"`
	Items        []SeoMedia `json:"items"`
	Completeness string     `json:"completeness"`
	FitMode      string     `json:"fitMode,omitempty"`
	Note         string     `json:"note,omitempty"`
}

type SeoFAQ struct {
	GroupKey string       `json:"groupKey"`
	Items    []SeoFAQItem `json:"items"`
}

type SeoFAQItem struct {
	ID       string `json:"id"`
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

type SeoMapTarget struct {
	Path  string `json:"path"`
	Label string `json:"label"`
}

type SeoAnswerReadiness struct {
	SpecificationVersion string   `json:"specificationVersion"`
	RequiredFacts        []string `json:"requiredFacts"`
	MissingRequiredFacts []string `json:"missingRequiredFacts"`
	RequiredSections     []string `json:"requiredSections"`
	RenderedSections     []string `json:"renderedSections"`
	Ready                bool     `json:"ready"`
}

type SeoPublicSEO struct {
	Title         string `json:"title"`
	Description   string `json:"description"`
	CanonicalPath string `json:"canonicalPath"`
	IndexClass    string `json:"indexClass"`
	Indexable     bool   `json:"indexable"`
	Follow        bool   `json:"follow"`
	ImageURL      string `json:"imageUrl,omitempty"`
}

type SeoPublicSource struct {
	Environment string   `json:"environment"`
	Authority   string   `json:"authority"`
	SourceName  string   `json:"sourceName,omitempty"`
	SourceURL   string   `json:"sourceUrl,omitempty"`
	Quality     string   `json:"quality"`
	UpdatedAt   string   `json:"updatedAt"`
	PublishedAt string   `json:"publishedAt,omitempty"`
	Limitations []string `json:"limitations,omitempty"`
}
