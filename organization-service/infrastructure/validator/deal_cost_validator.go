package validator

import (
	"context"
	"errors"
	organizationpb "pb/types/organization"
)

type DealCostValidator struct {
}

func NewDealCostValidator() *DealCostValidator {
	return &DealCostValidator{}
}

func (v *DealCostValidator) ValidateCreateDealCost(ctx context.Context, req *organizationpb.DealCostDTO) error {
	if req.DealId == 0 {
		return errors.New("deal_id is required")
	}
	if req.Amount == 0 {
		return errors.New("amount is required")
	}
	if req.CostTypeId == 0 {
		return errors.New("cost_type_id is required")
	}
	return nil
}
