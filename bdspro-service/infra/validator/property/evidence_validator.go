package property_validator

import (
	"common/pkg/validate"
	bdspropb "pb/types/bdspro"
)

func validateEvidencePayload(p *bdspropb.UpdateEvidencePayload) *validate.Error {
	e := validate.New()
	if p == nil {
		return e
	}
	validate.StringMaxLen(e, "evidence.title", p.Title, 255)
	validate.StringMaxLen(e, "evidence.description", p.Description, 2000)
	return e
}
