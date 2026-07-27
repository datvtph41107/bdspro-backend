package dto

const SeoAdministrativeUnitProvisioningSchemaVersion = "qhpro-administrative-unit-provisioning/v2u"

type ProvisionAdministrativeUnitSeoRequest struct {
	Identity string `json:"identity"`
	Publish  *bool  `json:"publish,omitempty"`
	Generate bool   `json:"generate"`
	Reason   string `json:"reason,omitempty"`
}

type ProvisionAdministrativeUnitSeoResponse struct {
	SchemaVersion         string              `json:"schemaVersion"`
	Operation             string              `json:"operation"`
	SeoDomainID           uint64              `json:"seoDomainId"`
	EntityID              string              `json:"entityId"`
	CanonicalSlug         string              `json:"canonicalSlug"`
	CanonicalURL          string              `json:"canonicalUrl"`
	Created               bool                `json:"created"`
	Updated               bool                `json:"updated"`
	Generated             bool                `json:"generated"`
	PageStatus            string              `json:"pageStatus"`
	Published             bool                `json:"published"`
	RenderStatus          string              `json:"renderStatus"`
	NeedGenerate          bool                `json:"needGenerate"`
	SourceStatus          string              `json:"sourceStatus"`
	DocumentVersion       string              `json:"documentVersion"`
	ETag                  string              `json:"etag,omitempty"`
	AnswerReadiness       *SeoAnswerReadiness `json:"answerReadiness"`
	IndexDecision         *SeoIndexDecision   `json:"indexDecision"`
	PolicyEvaluationScope string              `json:"policyEvaluationScope"`
	FinalPublicIndexable  bool                `json:"finalPublicIndexable"`
	Warnings              []string            `json:"warnings,omitempty"`
}
