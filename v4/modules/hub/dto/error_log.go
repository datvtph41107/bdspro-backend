package dto

// LogErrorRequest is the JSON body for POST /v2/hub/error/log.
type LogErrorRequest struct {
	Message string                 `json:"message"`
	Stack   *string                `json:"stack,omitempty"`
	Screen  *string                `json:"screen,omitempty"`
	UserID  *string                `json:"userId,omitempty"`
	Device  map[string]interface{} `json:"device,omitempty"`
	Extra   map[string]interface{} `json:"extra,omitempty"`
}
