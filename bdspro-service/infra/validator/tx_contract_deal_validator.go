package validator

import (
	"context"
	"fmt"
	bdspropb "pb/types/bdspro"
)

type ContractDealValidator struct{}

func NewContractDealValidator() *ContractDealValidator {
	return &ContractDealValidator{}
}

func (v *ContractDealValidator) ValidateDealContractPaymentRequest(ctx context.Context, req *bdspropb.DealContractPaymentRequest) error {
	if req == nil {
		return fmt.Errorf("request cannot be nil")
	}

	if req.DealId == 0 {
		return fmt.Errorf("deal ID is required")
	}

	// if req.Type == bdspropb.ContractType_CONTRACT_TYPE_UNKNOWN {
	// 	return fmt.Errorf("transaction type is required")
	// }

	if req.Amount <= 0 {
		return fmt.Errorf("amount must be greater than 0")
	}

	// if req.Description == "" {
	// 	return fmt.Errorf("description is required")
	// }

	// if req.TransactionDate == "" {
	// 	return fmt.Errorf("transaction date is required")
	// }

	// if req.Category == "" {
	// 	return fmt.Errorf("category is required")
	// }

	return nil
}

func (v *ContractDealValidator) ValidateUpdateRequest(req *bdspropb.UpdateDealContractRequest) error {
	if req == nil {
		return fmt.Errorf("request cannot be nil")
	}

	if req.Id == 0 {
		return fmt.Errorf("transaction ID is required")
	}

	if req.Amount <= 0 {
		return fmt.Errorf("amount must be greater than 0")
	}

	if req.Description == "" {
		return fmt.Errorf("description is required")
	}

	// if req.TransactionDate == "" {
	// 	return fmt.Errorf("transaction date is required")
	// }

	if req.Category == "" {
		return fmt.Errorf("category is required")
	}

	return nil
}
