package validator

import (
	"errors"
	"organization/internal/domain/entity"
)

type InvestmentValidator struct {
	InvestmentValidator *InvestmentValidator
}

func NewInvestmentValidator() *InvestmentValidator {
	return &InvestmentValidator{}
}

func (v *InvestmentValidator) ValidateCreateInvestment(req *entity.Investment) error {
	if req.DealID == 0 {
		return errors.New("deal_id is required")
	}
	// if req.MemberID == 0 {
	// 	return errors.New("member_id is required")
	// }
	if req.Amount <= 0 {
		return errors.New("amount is required")
	}
	return nil
}

func (v *InvestmentValidator) ValidateUpdateInvestment(req *entity.Investment) error {
	if req.ID == 0 {
		return errors.New("id is required")
	}
	if req.MemberID == 0 {
		return errors.New("member_id is required")
	}
	if req.Amount <= 0 {
		return errors.New("amount is required")
	}
	return nil
}
