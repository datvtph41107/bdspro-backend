package dto

// UpdateDataResourcePath binds the {resource} path param.
type UpdateDataResourcePath struct {
	Resource string `path:"resource"`
}

// CheckVersionSyncQuery binds query params for sync-status check.
type CheckVersionSyncQuery struct {
	Resource string `path:"resource"`
	ID       uint64 `form:"id"`
	LastSync int64  `form:"lastSync"`
	Limit    int64  `form:"limit"`
}

// CheckVersionSyncByIDQuery binds path+query for sync/status/{id}.
type CheckVersionSyncByIDQuery struct {
	Resource string `path:"resource"`
	ID       uint64 `path:"id"`
	LastSync int64  `form:"lastSync"`
}

// RefreshSyncIdsQuery binds query params for refresh ids.
type RefreshSyncIdsQuery struct {
	Resource string `path:"resource"`
	LastSync int64  `form:"lastSync"`
	Limit    int64  `form:"limit"`
}

// FlushSyncIdsQuery binds path+JSON body for flush.
type FlushSyncIdsQuery struct {
	Resource string `path:"resource"`
	Limit    int64  `json:"limit"`
}

// CheckVersionSyncResponse is the response for sync status checks.
type CheckVersionSyncResponse struct {
	NeedSync       bool     `json:"needSync"`
	CurrentVersion int64    `json:"currentVersion"`
	ChangeCount    int32    `json:"changeCount"`
	ChangedIDs     []uint64 `json:"changedIds,omitempty"`
}

// RefreshSyncIdsResponse is the response for refresh ids.
type RefreshSyncIdsResponse struct {
	NeedSync       bool     `json:"needSync"`
	CurrentVersion int64    `json:"currentVersion"`
	ChangeCount    int32    `json:"changeCount"`
	ChangedIDs     []uint64 `json:"changedIds,omitempty"`
	Next           *int64   `json:"next,omitempty"`
}
