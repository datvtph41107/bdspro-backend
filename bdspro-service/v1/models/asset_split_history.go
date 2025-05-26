package models

import (
	"bdspro/internal/enums"
	_models "common/models"
)

type AssetHistory struct {
	_models.BaseEntity
	AssetID    *uint64            `json:"assetId"`
	ActionType enums.AssetHistory `json:"actionType"`
	BeforeData *string            `json:"before"`
	AfterData  *string            `json:"after"`
	Target     AssetInfo          `gorm:"foreignKey:AssetID;references:ID" json:"assetInfo"`

	// ParentAssetID *uint64            `json:"parentAssetId"`
	// ChildAssetIDs pq.Int64Array `gorm:"type:bigint[];column:child_ids" json:"childAssetIds"` // JSON string or comma-separated IDs
	// Childs        []AssetInfo   `gorm:"-" json:"childs"`
	// Description   string             `json:"description"`
	// PerformedBy   uint64             `json:"performedBy"`
	// SplitMergeDate string             `json:"splitMergeDate"`
}

type AssetInfo struct {
	ID    uint64  `gorm:"primaryKey;autoIncrement" json:"id"`
	Name  string  `gorm:"size:255" json:"name"`
	Area  float64 `json:"area"`
	Price float64 `json:"price"`
	// Price       *ProductPrice `gorm:"foreignKey:LastPriceID;references:ID" json:"price"`
}

func (AssetHistory) TableName() string {
	return "asset_split_history"
}

func (AssetInfo) TableName() string {
	return "assets"
}
