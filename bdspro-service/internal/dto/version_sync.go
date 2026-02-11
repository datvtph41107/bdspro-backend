package dto

type CheckVersionSyncRequest struct {
	Resource string `json:"resource" validate:"required,oneof=product"`
	LastSync int64  `json:"last_sync"` // Unix timestamp lần sync cuối
	Id       uint64 `json:"id"`        // resource ID
}

type CheckVersionSyncResponse struct {
	NeedSync        bool     `json:"need_sync"`
	CurrentVersion  int64    `json:"current_version"`
	ResourceVersion *int64   `json:"product_version,omitempty"`
	ChangeCount     int      `json:"change_count"`
	ChangedIDs      []uint64 `json:"changed_ids"`
}
