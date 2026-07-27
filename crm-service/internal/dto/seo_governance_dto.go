package dto

const SeoGovernanceDocumentSchemaVersion = "qhpro-seo-governance-document/v1"

type SeoGovernanceRequirementStatus struct {
	ID      string `json:"id"`
	Present bool   `json:"present"`
	Value   string `json:"value,omitempty"`
}

type SeoGovernancePublication struct {
	PageStatus      string `json:"pageStatus"`
	Published       bool   `json:"published"`
	PublishedAt     string `json:"publishedAt,omitempty"`
	Scope           string `json:"scope"`
	IndexConfigured bool   `json:"indexConfigured"`
	SitemapEnabled  bool   `json:"sitemapEnabled"`
	RenderStatus    string `json:"renderStatus"`
	NeedGenerate    bool   `json:"needGenerate"`
	LastRenderError string `json:"lastRenderError,omitempty"`
}

type SeoGovernanceEvaluation struct {
	LastEvaluatedAt   string `json:"lastEvaluatedAt"`
	EvaluationBasis   string `json:"evaluationBasis"`
	PublicDocumentVer string `json:"publicDocumentVersion"`
	TemplateVersion   string `json:"templateVersion,omitempty"`
	ETag              string `json:"etag"`
	RenderedHTMLHash  string `json:"renderedHtmlHash,omitempty"`
	GeneratedAt       string `json:"generatedAt,omitempty"`
	UpdatedAt         string `json:"updatedAt"`
}

type SeoGovernanceReview struct {
	ReviewState    string `json:"reviewState"`
	AuditState     string `json:"auditState"`
	OperatorNote   string `json:"operatorNote,omitempty"`
	AuditAvailable bool   `json:"auditAvailable"`
}

type SeoGovernanceDocument struct {
	SchemaVersion              string                           `json:"schemaVersion"`
	ModuleID                   string                           `json:"moduleId"`
	PageFamily                 string                           `json:"pageFamily"`
	SeoDomainID                uint64                           `json:"seoDomainId"`
	EntityID                   string                           `json:"entityId"`
	CanonicalSlug              string                           `json:"canonicalSlug"`
	CanonicalURL               string                           `json:"canonicalUrl"`
	AnswerSpecificationVersion string                           `json:"answerSpecificationVersion"`
	AnswerReadiness            *SeoAnswerReadiness              `json:"answerReadiness"`
	IndexDecision              *SeoIndexDecision                `json:"indexDecision"`
	RequiredFacts              []SeoGovernanceRequirementStatus `json:"requiredFacts"`
	RequiredSections           []SeoGovernanceRequirementStatus `json:"requiredSections"`
	Source                     SeoPublicSource                  `json:"source"`
	Publication                SeoGovernancePublication         `json:"publication"`
	Evaluation                 SeoGovernanceEvaluation          `json:"evaluation"`
	Review                     SeoGovernanceReview              `json:"review"`
	OperatingSystem            *SeoOperatingSystemSnapshot      `json:"operatingSystem,omitempty"`
	ClientHandoff              *SeoPageBootstrap                `json:"clientHandoff,omitempty"`
	Measurement                *SeoMeasurementPlan              `json:"measurement,omitempty"`
	GenerativeAI               *SeoGenerativeAnswer             `json:"generativeAi,omitempty"`
	Warnings                   []string                         `json:"warnings,omitempty"`
}

type SeoGovernanceDocumentResponse struct {
	Document   *SeoGovernanceDocument `json:"document"`
	ETag       string                 `json:"etag"`
	StatusCode int                    `json:"statusCode"`
}
