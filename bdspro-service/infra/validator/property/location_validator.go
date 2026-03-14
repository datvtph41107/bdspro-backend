package property_validator

import (
	"common/pkg/validate"
	bdspropb "pb/types/bdspro"
)

func validateLocationPayload(p *bdspropb.UpdateLocationPayload) *validate.Error {
	e := validate.New()
	if p == nil {
		return e
	}
	validate.StringMaxLen(e, "location.address_detail", p.AddressDetail, 100)
	validate.URLOpt(e, "location.map_url", p.MapUrl)
	validate.Latitude(e, "location.latitude", p.Latitude)
	validate.Longitude(e, "location.longitude", p.Longitude)
	return e
}
