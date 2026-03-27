package dto

import "encoding/json"

type LogErrorRequest struct {
	Message string                 `json:"message"`
	Stack   *string                `json:"stack,omitempty"`
	Screen  *string                `json:"screen,omitempty"`
	UserID  *string                `json:"userId,omitempty"`
	Device  *ErrorLogDevice        `json:"device,omitempty"`
	Extra   map[string]interface{} `json:"extra,omitempty"`
}

type ErrorLogDevice struct {
	Platform   *string `json:"platform,omitempty"`
	Model      *string `json:"model,omitempty"`
	OSVersion  *string `json:"osVersion,omitempty"`
	AppVersion *string `json:"appVersion,omitempty"`
}

func DeviceToJSON(d *ErrorLogDevice) (json.RawMessage, error) {
	if d == nil {
		return nil, nil
	}
	return json.Marshal(d)
}

func ExtraToJSON(e map[string]interface{}) (json.RawMessage, error) {
	if e == nil || len(e) == 0 {
		return nil, nil
	}
	return json.Marshal(e)
}
