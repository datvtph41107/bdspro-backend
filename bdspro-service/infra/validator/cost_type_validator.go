package validator

import (
	_routes "common/routes"
	bdspropb "pb/types/bdspro"
)

type CostTypeValidator struct {
}

func NewCostTypeValidator() *CostTypeValidator {
	return &CostTypeValidator{}
}

func (v *CostTypeValidator) ValidateSaveCostType(pbCostType *bdspropb.AssetCostTypeDTO) error {
	if pbCostType.TypeName == "" || pbCostType.Type == 0 {
		return &_routes.Except{
			Code:    400,
			Message: "name and type are required",
		}
	}
	return nil
}
