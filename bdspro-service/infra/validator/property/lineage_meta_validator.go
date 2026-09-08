package property_validator

import (
	"common/pkg/validate"
	bdspropb "pb/types/bdspro"
)

func validateLineageMetaPayload(p *bdspropb.UpdateLineageMetaPayload) *validate.Error {
	e := validate.New()
	if p == nil {
		return e
	}
	validate.StringMaxLen(e, "lineage_meta.national_id", p.NationalId, 100)
	validate.DateStringOpt(e, "lineage_meta.verified_national_at", p.VerifiedNationalAt)
	return e
}
