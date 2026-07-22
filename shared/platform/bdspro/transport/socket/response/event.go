package socket_resp

import (
	"time"
)

// Event is the standard WebSocket event structure
type Event struct {
	Type      string      `json:"type"` // e.g., "property.created", "error"
	Payload   interface{} `json:"payload,omitempty"`
	Error     *ErrorEvent `json:"error,omitempty"`
	RequestID string      `json:"request_id,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
}

type ErrorEvent struct {
	Code    string       `json:"code"`
	Message string       `json:"message"`
	Fields  []FieldError `json:"fields,omitempty"`
}

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Code    string `json:"code"`
}
