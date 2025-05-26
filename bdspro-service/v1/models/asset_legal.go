package models

import _models "common/models"

type AssetLegal struct {
	_models.BaseEntity
	AssetID             uint64  `json:"assetId" binding:"required"`
	DocumentName        string  `json:"documentName" binding:"required"`
	DocumentURL         string  `json:"documentUrl" binding:"required"`
	IssuedDate          string  `json:"issuedDate"`
	ExpiryDate          string  `json:"expiryDate"`
	Description         string  `json:"description"`
	RelatedSplitMergeID *uint64 `json:"relatedSplitMergeId"`
}

func (AssetLegal) TableName() string {
	return "asset_legals"
}
