package property_validator

import (
	"common/pkg/validate"
	bdspropb "pb/types/bdspro"
)

func validateInfoPayload(p *bdspropb.UpdateInfoPayload) *validate.Error {
	e := validate.New()
	if p == nil {
		return e
	}
	validate.StringMaxLen(e, "info.title", p.Title, 255)
	validate.StringMaxLen(e, "info.note", p.Note, 1000)
	validate.StringMaxLen(e, "info.unit_code", p.UnitCode, 50)
	validate.StringMaxLen(e, "info.identifier", p.Identifier, 100)
	validate.StringMaxLen(e, "info.level", p.Level, 100)
	return e
}
