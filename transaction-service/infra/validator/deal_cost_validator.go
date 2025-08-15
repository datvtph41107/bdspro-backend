package validator

import (
	"context"
	"errors"
	"tx/internal/domain/entity"
)

type DealCostValidator struct {
}

func NewDealCostValidator() *DealCostValidator {
	return &DealCostValidator{}
}

func (v *DealCostValidator) ValidateCreateDealCost(ctx context.Context, req interface{}) error {
	// TODO: Implement after protobuf generation
	return nil
}

func (v *DealCostValidator) ValidateUpdateDealCost(ctx context.Context, req interface{}) error {
	// TODO: Implement after protobuf generation
	return nil
}

func (v *DealCostValidator) ValidateDealCost(ctx context.Context, dealCost *entity.DealCost) error {
	if dealCost.DealID == 0 {
		return errors.New("deal ID is required")
	}
	if dealCost.Amount <= 0 {
		return errors.New("amount must be greater than 0")
	}
	if dealCost.CostTypeID == 0 {
		return errors.New("cost type ID is required")
	}
	if dealCost.CostName == "" {
		return errors.New("cost name is required")
	}
	return nil
}
