package validator

import (
	"bdspro/internal/dto"
	_routes "common/routes"
)

type ProductValidator struct {
}

func NewProductValidator() *ProductValidator {
	return &ProductValidator{}
}

func (v *ProductValidator) ValidateProductChild(dto *dto.SaveProductChildRequest) error {
	return nil
}

func (v *ProductValidator) ValidateDevideChildRequest(dto *dto.DevideChildRequest) error {
	if dto.ParentID == nil || *dto.ParentID == 0 || dto.Childs == nil || len(dto.Childs) == 0 {
		return &_routes.Except{
			Code:    400,
			Message: "parentId and childs are required",
		}
	}
	return nil
}

func (v *ProductValidator) ValidateSaveProductChild(dto *dto.ProductSaveRequest) error {
	if dto.ParentID == nil ||
		*dto.ParentID == 0 ||
		dto.Name == "" ||
		dto.Area == 0 ||
		dto.PriceData == nil ||
		(dto.PriceData.SalePrice == nil && dto.PriceData.RentPrice == nil) {
		return &_routes.Except{
			Code:    400,
			Message: "parentId, name, area, priceData are required",
		}
	}

	return nil
}
