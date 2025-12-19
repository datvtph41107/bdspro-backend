package validator

import (
	"bdspro/internal/domain"
	"context"
	"errors"
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

func (v *DealCostValidator) ValidateDealCost(ctx context.Context, dealCost *domain.TxDealCost) error {
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
