package property_validator

import (
	"common/pkg/validate"
	bdspropb "pb/types/bdspro"
)

func validateLandInfoPayload(p *bdspropb.UpdateLandInfoPayload) *validate.Error {
	e := validate.New()
	if p == nil {
		return e
	}
	validate.StringMaxLen(e, "land_info.note", p.Note, 1000)
	validate.StringMaxLen(e, "land_info.land_note", p.LandNote, 1000)
	// if p.Area != nil {
	validate.Float64NonNegative(e, "land_info.area.area_total", p.AreaTotal)
	validate.Float64NonNegative(e, "land_info.area.area_land", p.AreaLand)
	validate.Float64NonNegative(e, "land_info.area.area_plant", p.AreaPlant)
	// }
	// if p.Dimension != nil {
	validate.Float64NonNegative(e, "land_info.dimension.front_width", p.FrontWidth)
	validate.Float64NonNegative(e, "land_info.dimension.depth", p.Depth)
	validate.Float64NonNegative(e, "land_info.dimension.street_width", p.StreetWidth)
	// }
	// if p.Expiry != nil {
	validate.DateStringOpt(e, "land_info.expiry.expired_land", p.ExpiredLand)
	validate.DateStringOpt(e, "land_info.expiry.expired_plant", p.ExpiredPlant)
	// Cross-field: nếu cả 2 có giá trị → land phải trước plant
	validate.DateOrder(e,
		"land_info.expiry.expired_land",
		"land_info.expiry.expired_plant",
		p.ExpiredLand,
		p.ExpiredPlant,
	)
	// }
	return e
}
