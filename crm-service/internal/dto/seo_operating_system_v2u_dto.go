package dto

const (
	SeoOperatingSystemRelease         = "seo-operating-system/v2.u"
	SeoPageBootstrapContractVersion   = "qhpro-page-bootstrap/v2u"
	SeoMeasurementPlanSchemaVersion   = "qhpro-seo-measurement-plan/v2u"
	SeoGenerativeAnswerSchemaVersion  = "qhpro-generative-answer/v2u"
	SeoOperatingSystemManifestVersion = "qhpro-seo-operating-system-manifest/v2u"
)

type SeoOperatingModuleCapability struct {
	Number      int      `json:"number"`
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Owner       string   `json:"owner"`
	RuntimeRole string   `json:"runtimeRole"`
	State       string   `json:"state"`
	Evidence    []string `json:"evidence,omitempty"`
}

type SeoOperatingSystemSnapshot struct {
	Release       string                         `json:"release"`
	PageFamily    string                         `json:"pageFamily"`
	DeliveryOwner string                         `json:"deliveryOwner"`
	DeliveryState string                         `json:"deliveryState"`
	Modules       []SeoOperatingModuleCapability `json:"modules"`
}

type SeoPageBootstrapMapFocus struct {
	Kind string `json:"kind"`
	ID   string `json:"id"`
}

type SeoPageBootstrap struct {
	ContractVersion string                    `json:"contractVersion"`
	Presentation    string                    `json:"presentation"`
	Surface         string                    `json:"surface"`
	EntityKind      string                    `json:"entityKind"`
	EntityID        string                    `json:"entityId"`
	CanonicalSlug   string                    `json:"canonicalSlug"`
	CanonicalPath   string                    `json:"canonicalPath"`
	InitialScreen   string                    `json:"initialScreen"`
	MountPolicy     string                    `json:"mountPolicy"`
	DataStrategy    string                    `json:"dataStrategy"`
	DocumentVersion string                    `json:"documentVersion"`
	MapFocus        *SeoPageBootstrapMapFocus `json:"mapFocus,omitempty"`
}

type SeoMeasurementPlan struct {
	SchemaVersion   string   `json:"schemaVersion"`
	RequiredEvents  []string `json:"requiredEvents"`
	CorrelationKeys []string `json:"correlationKeys"`
}

type SeoGenerativeCitation struct {
	ID         string `json:"id"`
	SourceName string `json:"sourceName"`
	SourceURL  string `json:"sourceUrl,omitempty"`
	UpdatedAt  string `json:"updatedAt"`
	Authority  string `json:"authority"`
}

type SeoGenerativeAnswerBlock struct {
	ID          string   `json:"id"`
	Kind        string   `json:"kind"`
	Title       string   `json:"title"`
	Text        string   `json:"text"`
	FactIDs     []string `json:"factIds,omitempty"`
	SectionIDs  []string `json:"sectionIds,omitempty"`
	CitationIDs []string `json:"citationIds"`
}

type SeoGenerativeAnswer struct {
	SchemaVersion string                     `json:"schemaVersion"`
	Ready         bool                       `json:"ready"`
	Blocks        []SeoGenerativeAnswerBlock `json:"blocks"`
	Citations     []SeoGenerativeCitation    `json:"citations"`
	Limitations   []string                   `json:"limitations,omitempty"`
}

type SeoOperatingRuntimeStage struct {
	Order  int    `json:"order"`
	ID     string `json:"id"`
	Name   string `json:"name"`
	Owner  string `json:"owner"`
	Input  string `json:"input"`
	Output string `json:"output"`
}

type SeoOperatingVerificationGate struct {
	ID             string   `json:"id"`
	Name           string   `json:"name"`
	RequiredStatus string   `json:"requiredStatus"`
	CurrentStatus  string   `json:"currentStatus"`
	Evidence       []string `json:"evidence"`
	BlockingReason string   `json:"blockingReason,omitempty"`
	BlocksPilot    bool     `json:"blocksPilot"`
}

type SeoPageFamilyCapability struct {
	ID                 string `json:"id"`
	CanonicalRoute     string `json:"canonicalRoute"`
	ActivationState    string `json:"activationState"`
	DeliveryOwner      string `json:"deliveryOwner"`
	DeliveryState      string `json:"deliveryState"`
	BackendContract    string `json:"backendContract"`
	NextRendering      string `json:"nextRendering"`
	ReactNativeSurface string `json:"reactNativeSurface"`
	AdminGovernance    string `json:"adminGovernance"`
	PilotReadiness     string `json:"pilotReadiness"`
}

type SeoOperatingSystemManifest struct {
	SchemaVersion          string                         `json:"schemaVersion"`
	Release                string                         `json:"release"`
	Milestone              string                         `json:"milestone"`
	ImplementationStatus   string                         `json:"implementationStatus"`
	CurrentPageFamily      string                         `json:"currentPageFamily"`
	MaximumDeliveryState   string                         `json:"maximumDeliveryState"`
	Modules                []SeoOperatingModuleCapability `json:"modules"`
	PageFamilies           []string                       `json:"pageFamilies"`
	PageFamilyCapabilities []SeoPageFamilyCapability      `json:"pageFamilyCapabilities"`
	RuntimeFlow            []SeoOperatingRuntimeStage     `json:"runtimeFlow"`
	PilotGates             []SeoOperatingVerificationGate `json:"pilotGates"`
	Invariants             []string                       `json:"invariants"`
}
