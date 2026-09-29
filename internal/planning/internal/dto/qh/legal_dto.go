package qh_dto

// ============================================================
// LEGAL DTO — Dữ liệu trao đổi giữa Service và Handler
// ============================================================
// LOGICAL:
// - Surface: Cần DTO cho legal documents response
// - Root: Tách biệt domain model khỏi API contract
// - Mechanism: DTO chuyển đổi domain → proto → JSON

// LegalDocumentDTO dùng cho API response
type LegalDocumentDTO struct {
	ID               string              `json:"id"`
	Type             string              `json:"type"`
	TypeName         string              `json:"typeName"`
	TypeIcon         string              `json:"typeIcon"`
	Name             string              `json:"name"`
	FullName         string              `json:"fullName"`
	IssuedBy         string              `json:"issuedBy"`
	IssuedAt         string              `json:"issuedAt"`
	EffectiveDate    string              `json:"effectiveDate"`
	ExpiryDate       *string             `json:"expiryDate"`
	Status           string              `json:"status"`
	StatusLabel      string              `json:"statusLabel"`
	StatusColor      string              `json:"statusColor"`
	SupersededBy     *string             `json:"supersededBy"`
	Supersedes       []string            `json:"supersedes"`
	RelevantArticles []string            `json:"relevantArticles"`
	File             *LegalFileDTO       `json:"file"`
	IsPublic         bool                `json:"isPublic"`
	RequiresAuth     bool                `json:"requiresAuth"`
	AffectsEntities  *AffectsEntitiesDTO `json:"affectsEntities"`
}

// LegalFileDTO thông tin file đính kèm
type LegalFileDTO struct {
	Type        string  `json:"type"`
	SizeKb      int     `json:"sizeKb"`
	Pages       int     `json:"pages"`
	URL         *string `json:"url"`
	PreviewURL  *string `json:"previewUrl"`
	DownloadURL *string `json:"downloadUrl"`
}

// AffectsEntitiesDTO document áp dụng cho entity nào
type AffectsEntitiesDTO struct {
	LayerIDs []uint64 `json:"layerIds"`
	ZoneIDs  []uint64 `json:"zoneIds"`
}

// LegalDocumentsSummaryDTO tổng hợp documents
type LegalDocumentsSummaryDTO struct {
	Items             []LegalDocumentDTO `json:"items"`
	Total             int                `json:"total"`
	ByType            map[string]int     `json:"byType"`
	ByStatus          map[string]int     `json:"byStatus"`
	HasSupersededDocs bool               `json:"hasSupersededDocs"`
}

// LegalDocumentsResponseDTO response cho GET /legal-documents
type LegalDocumentsResponseDTO struct {
	Data    []LegalDocumentDTO       `json:"data"`
	Summary LegalDocumentsSummaryDTO `json:"summary"`
	Filters *LegalFiltersDTO         `json:"filters,omitempty"`
	Meta    *ResponseMetaDTO         `json:"_meta"`
	Links   *LinksDTO                `json:"_links"`
}

// LegalFiltersDTO filter options
type LegalFiltersDTO struct {
	Applied           *LegalFiltersAppliedDTO `json:"applied"`
	AvailableTypes    []string                `json:"availableTypes"`
	AvailableStatuses []string                `json:"availableStatuses"`
}

// LegalFiltersAppliedDTO filter đang áp dụng
type LegalFiltersAppliedDTO struct {
	Type   string `json:"type"`
	Status string `json:"status"`
}

// ResponseMetaDTO metadata cho response
type ResponseMetaDTO struct {
	DataVersion    string `json:"dataVersion"`
	DataQuality    string `json:"dataQuality"`
	QueriedAt      string `json:"queriedAt"`
	ResponseTimeMs int64  `json:"responseTimeMs"`
	PartialData    bool   `json:"partialData,omitempty"`
}

// LinksDTO HATEOAS links
type LinksDTO struct {
	Self           string `json:"self"`
	Layers         string `json:"layers,omitempty"`
	LegalDocuments string `json:"legalDocuments,omitempty"`
	ShareUrl       string `json:"shareUrl,omitempty"`
	Next           string `json:"next,omitempty"`
	Prev           string `json:"prev,omitempty"`
}
