package domain

import _models "common/models"

type AssetHistory struct {
	_models.BaseEntity
	AssetID        uint64  `json:"assetId"`
	ActionType     string  `json:"actionType"`
	ParentAssetID  *uint64 `json:"parentAssetId"`
	ChildAssetIDs  string  `json:"childAssetIds"`
	SplitMergeDate string  `json:"splitMergeDate"`
	Description    string  `json:"description"`
	PerformedBy    uint64  `json:"performedBy"`
}
