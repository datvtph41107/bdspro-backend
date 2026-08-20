package validator

import (
	_errors "common/errors"
	"tqd/internal/dto"
)

type ParcelValidator struct{}

func NewParcelValidator() *ParcelValidator {
	return &ParcelValidator{}
}

func (v *ParcelValidator) ValidateAdminUpdate(req *dto.AdminUpdateParcelRequest) error {
	if req == nil || req.ParcelID == 0 {
		return _errors.ReturnError(400, "parcelId is required")
	}
	if req.Lat < -90 || req.Lat > 90 {
		return _errors.ReturnError(400, "lat không hợp lệ")
	}
	if req.Lon < -180 || req.Lon > 180 {
		return _errors.ReturnError(400, "lon không hợp lệ")
	}
	if req.AreaSqm < 0 {
		return _errors.ReturnError(400, "areaSqm không hợp lệ")
	}
	return nil
}
