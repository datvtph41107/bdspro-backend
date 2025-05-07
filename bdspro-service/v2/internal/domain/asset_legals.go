package domain

import _models "common/models"

type AssetLegal struct {
	_models.BaseEntity
	AssetID             uint64  `json:"assetId"`
	DocumentName        string  `json:"documentName"`
	DocumentURL         string  `json:"documentUrl"`
	IssuedDate          string  `json:"issuedDate"`
	ExpiryDate          string  `json:"expiryDate"`
	Description         string  `json:"description"`
	RelatedSplitMergeID *uint64 `json:"relatedSplitMergeId"`
}

func (AssetLegal) TableName() string {
	return "asset_legals"
}
