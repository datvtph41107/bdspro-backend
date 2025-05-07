package domain

import _models "common/models"

type AssetSplitMergeHistory struct {
	_models.BaseEntity
	AssetID        uint64  `json:"assetId"`
	ActionType     string  `json:"actionType"`
	ParentAssetID  *uint64 `json:"parentAssetId"`
	ChildAssetIDs  string  `json:"childAssetIds"` // JSON string or comma-separated IDs
	SplitMergeDate string  `json:"splitMergeDate"`
	Description    string  `json:"description"`
	PerformedBy    uint64  `json:"performedBy"`
}

func (AssetSplitMergeHistory) TableName() string {
	return "asset_split_merge_history"
}
