package validator

import (
	_errors "common/errors"
	bdspropb "pb/types/bdspro"
)

type AssetCostValidator struct {
}

func NewAssetCostValidator() *AssetCostValidator {
	return &AssetCostValidator{}
}

func (v *AssetCostValidator) ValidateSaveAssetCost(pbAssetCost *bdspropb.AssetCostDTO) error {
	if pbAssetCost.Amount == 0 ||
		pbAssetCost.Date == "" ||
		pbAssetCost.OwnerType == 0 ||
		pbAssetCost.Type == 0 ||
		pbAssetCost.AssetId == 0 ||
		pbAssetCost.CostTypeId == 0 {
		return _errors.ReturnError(_errors.RequestValidationFailed, _errors.WithPublicMessage("amount, date, ownerType, type, assetId and costTypeId are required"))
	}
	return nil
}
