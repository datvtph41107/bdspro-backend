package validator

import (
	_errors "common/errors"
	bdspropb "pb/types/bdspro"
)

type CostTypeValidator struct {
}

func NewCostTypeValidator() *CostTypeValidator {
	return &CostTypeValidator{}
}

func (v *CostTypeValidator) ValidateSaveCostType(pbCostType *bdspropb.AssetCostTypeDTO) error {
	if pbCostType.TypeName == "" || pbCostType.Type == 0 {
		return _errors.ReturnError(_errors.RequestValidationFailed, _errors.WithPublicMessage("name and type are required"))
	}
	return nil
}
