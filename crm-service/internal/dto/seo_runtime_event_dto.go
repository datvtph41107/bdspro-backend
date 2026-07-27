package dto

const SeoRuntimeEventSchemaVersion = "qhpro-seo-runtime-event/v1"

type SeoRuntimeEventIngestRequest struct {
	SchemaVersion int                    `json:"schemaVersion"`
	Source        string                 `json:"source"`
	Event         string                 `json:"event"`
	SessionID     string                 `json:"sessionId"`
	CanonicalURL  string                 `json:"canonicalUrl"`
	CanonicalPath string                 `json:"canonicalPath"`
	ModuleID      string                 `json:"moduleId"`
	Environment   string                 `json:"environment"`
	Indexable     *bool                  `json:"indexable"`
	Action        string                 `json:"action"`
	Target        string                 `json:"target"`
	Timestamp     string                 `json:"timestamp"`
	Metrics       map[string]float64     `json:"metrics"`
	Metadata      map[string]interface{} `json:"metadata"`
	UserAgent     string                 `json:"userAgent"`
}

type SeoRuntimeEventIngestResponse struct {
	SchemaVersion string `json:"schemaVersion"`
	Accepted      bool   `json:"accepted"`
	Persisted     bool   `json:"persisted"`
	Idempotent    bool   `json:"idempotent"`
	EventKey      string `json:"eventKey"`
	ReceivedAt    string `json:"receivedAt"`
}
