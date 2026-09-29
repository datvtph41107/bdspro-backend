package validator

import (
	"context"
	"errors"
	"organization/internal/domain/entity"
	"time"

	"github.com/go-playground/validator/v10"
)

var (
	ErrTitleRequired        = errors.New("title is required")
	ErrDealIDRequired      = errors.New("deal id is required")
	ErrExpectedDateRequired = errors.New("expected date is required")
	ErrExpectedDateInvalid  = errors.New("expected date must be in the future")
	ErrIDRequired          = errors.New("id is required")
)

type DealMilestoneValidator struct {
	validate *validator.Validate
}

func NewDealMilestoneValidator(validate *validator.Validate) *DealMilestoneValidator {
	return &DealMilestoneValidator{
		validate: validate,
	}
}

func (v *DealMilestoneValidator) ValidateCreate(ctx context.Context, milestone *entity.DealMilestone) error {
	if milestone.Title == "" {
		return ErrTitleRequired
	}

	if milestone.DealID == 0 {
		return ErrDealIDRequired
	}

	if milestone.ExpectedDate.IsZero() {
		return ErrExpectedDateRequired
	}

	if milestone.ExpectedDate.Before(time.Now()) {
		return ErrExpectedDateInvalid
	}

	return nil
}

func (v *DealMilestoneValidator) ValidateUpdate(ctx context.Context, milestone *entity.DealMilestone) error {
	if milestone.ID == 0 {
		return ErrIDRequired
	}

	if milestone.Title == "" {
		return ErrTitleRequired
	}

	if milestone.ExpectedDate.IsZero() {
		return ErrExpectedDateRequired
	}

	if milestone.ExpectedDate.Before(time.Now()) {
		return ErrExpectedDateInvalid
	}

	return nil
}

func (v *DealMilestoneValidator) ValidateDelete(ctx context.Context, id uint64) error {
	if id == 0 {
		return ErrIDRequired
	}
	return nil
} 