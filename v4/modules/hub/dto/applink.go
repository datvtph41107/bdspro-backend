package dto

// CreateApplinkRequest is the JSON body for POST /v2/hub/applinks.
type CreateApplinkRequest struct {
	RefID  uint64 `json:"refId"`
	Action int32  `json:"action"`
}

// ApplinkResponse is the API response for applink operations.
type ApplinkResponse struct {
	ID     uint64 `json:"id"`
	Code   string `json:"code"`   // base64-encoded uint64 code
	RefID  uint64 `json:"refId"`
	Action int32  `json:"action"`
}

// GetApplinkByCodePath binds the {code} path param.
type GetApplinkByCodePath struct {
	Code string `path:"code"`
}
