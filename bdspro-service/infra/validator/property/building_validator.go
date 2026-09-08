package property_validator

import (
	"common/pkg/validate"
	bdspropb "pb/types/bdspro"
)

func validateBuildingPayload(p *bdspropb.UpdateBuildingPayload) *validate.Error {
	e := validate.New()
	if p == nil {
		return e
	}
	validate.StringMaxLen(e, "building_info.note", p.Note, 1000)
	validate.Float64NonNegative(e, "building_info.area_actual", p.AreaActual)
	validate.Float64NonNegative(e, "building_info.area_floor", p.AreaFloor)
	validate.Float64NonNegative(e, "building_info.area_construction", p.AreaConstruction)
	return e
}
