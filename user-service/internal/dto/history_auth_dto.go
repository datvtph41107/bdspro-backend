package dto

// HistoryAuthCreateDTO represents the payload to log auth-service user actions.
type HistoryAuthCreateDTO struct {
	UserID         uint64                 `json:"userId"`
	OrganizationID *uint64                `json:"organizationId,omitempty"`
	ActionType     string                 `json:"actionType"`
	ActionName     string                 `json:"actionName"`
	Description    string                 `json:"description,omitempty"`
	Success        *bool                  `json:"success,omitempty"`
	Reason         string                 `json:"reason,omitempty"`
	IPAddress      string                 `json:"ipAddress,omitempty"`
	UserAgent      string                 `json:"userAgent,omitempty"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
	PerformedBy    *uint64                `json:"performedBy,omitempty"`
	SessionID      string                 `json:"sessionId,omitempty"`
	Channel        string                 `json:"channel,omitempty"`
	DeviceID       string                 `json:"deviceId,omitempty"`
	Location       string                 `json:"location,omitempty"`
	AdditionalNote string                 `json:"additionalNote,omitempty"`
	SourceService  string                 `json:"sourceService,omitempty"`
}
