package dto

import "time"

// APIKeyItem is the API representation of an api_keys row.
type APIKeyItem struct {
	ID          uint64     `json:"id"`
	Name        string     `json:"name"`
	AppName     string     `json:"appName"`
	Description string     `json:"description,omitempty"`
	APIKey      string     `json:"apiKey,omitempty"`
	ExpiredAt   *time.Time `json:"expiredAt,omitempty"`
	CreatedAt   *string    `json:"createdAt,omitempty"`
}

// CreateAPIKeyRequest is the JSON body for POST /v2/hub/admin/api-keys.
type CreateAPIKeyRequest struct {
	Name        string     `json:"name"`
	AppName     string     `json:"appName"`
	Description string     `json:"description,omitempty"`
	ExpiredAt   *time.Time `json:"expiredAt,omitempty"`
}

// VerifyAPIKeyRequest is the JSON body for POST /v2/hub/admin/api-keys/verify.
type VerifyAPIKeyRequest struct {
	APIKey string `json:"apiKey"`
}

// APIKeyIDPath binds the {id} path param.
type APIKeyIDPath struct {
	ID uint64 `path:"id"`
}
