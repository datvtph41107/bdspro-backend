package dto

type CheckVersionSyncRequest struct {
	Resource string `json:"resource" validate:"required,oneof=product"`
	OwnerID  uint64 `json:"owner_id" validate:"required"`
	LastSync int64  `json:"last_sync"` // Unix timestamp lần sync cuối
}

type CheckVersionSyncResponse struct {
	NeedSync       bool     `json:"need_sync"`
	CurrentVersion int64    `json:"current_version"`
	ChangeCount    int      `json:"change_count"`
	ChangedIDs     []uint64 `json:"changed_ids"`
}
