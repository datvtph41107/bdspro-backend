package validator

import (
	"bdspro/internal/domain"
	"errors"
)

type InvestmentValidator struct {
	InvestmentValidator *InvestmentValidator
}

func NewInvestmentValidator() *InvestmentValidator {
	return &InvestmentValidator{}
}

func (v *InvestmentValidator) ValidateCreateInvestment(req *domain.DealInvestment) error {
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

func (v *InvestmentValidator) ValidateUpdateInvestment(req *domain.DealInvestment) error {
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
