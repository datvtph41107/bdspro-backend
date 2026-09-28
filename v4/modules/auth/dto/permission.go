package dto

// PermissionKeysResponse is the modules/auth contract for GET /auth/permission/me.
// Gateway lifts timestamp onto the HTTP envelope root and keeps keys under data.keys.
type PermissionKeysResponse struct {
	Keys      []string `json:"keys"`
	Timestamp int64    `json:"timestamp"`
}
