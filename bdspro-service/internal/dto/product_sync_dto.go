package dto

type SyncProductsRequest struct {
	ProfileID    uint64   `json:"profileId"`
	LastSyncTime int64    `json:"lastSyncTime"` // ms
	PageSize     int      `json:"pageSize"`
	PageToken    string   `json:"pageToken"`
	ProductIDs   []uint64 `json:"productIds"` // optional
}

type SyncProductsResponse struct {
	ProductIDs      []*ChangedProduct `json:"productIds"`
	RemovedIDs      []uint64          `json:"removedIds"`
	CurrentSyncTime int64             `json:"currentSyncTime"`
	NextToken       string            `json:"nextToken"`
	HasMore         bool              `json:"hasMore"`
}

type ChangedProduct struct {
	ID         uint64 `json:"id"`
	UpdatedAt  int64  `json:"updatedAt"`
	ChangeType int32  `json:"changeType"` // 1=grant,2=update,3=revoke
}
