package dto

type CreateDefaultSeoByRefResult struct {
	Created      bool   `json:"created"`
	SeoID        uint64 `json:"seoId"`
	RefType      uint32 `json:"refType"`
	RefID        uint64 `json:"refId"`
	Slug         string `json:"slug"`
	OriginURL    string `json:"originUrl"`
	CanonicalURL string `json:"canonicalUrl"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	PageStatus   string `json:"pageStatus"`
	RenderStatus string `json:"renderStatus"`
}
