package validator

import (
	"organization/internal/domain/entity"
	organizationpb "pb/types/organization"

	"github.com/go-playground/validator/v10"
)

type ColorValidator struct {
	validate *validator.Validate
}

func NewColorValidator() *ColorValidator {
	return &ColorValidator{
		validate: validator.New(),
	}
}

func (v *ColorValidator) ValidateCreateColorRequest(req *organizationpb.CreateColorRequest) error {
	if req == nil {
		return validator.New().Var(nil, "required")
	}

	return v.validate.Struct(req)
}

func (v *ColorValidator) ValidateUpdateColorRequest(req *organizationpb.UpdateColorRequest) error {
	if req == nil {
		return validator.New().Var(nil, "required")
	}

	return v.validate.Struct(req)
}

func (v *ColorValidator) ValidateColorEntity(color *entity.Color) error {
	if color == nil {
		return validator.New().Var(nil, "required")
	}

	return v.validate.Struct(color)
} 