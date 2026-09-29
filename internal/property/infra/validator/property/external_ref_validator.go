package property_validator

import (
	"common/pkg/validate"
	"fmt"
	bdspropb "pb/types/bdspro"
)

func validateExternalRefPayload(p *bdspropb.UpdateExternalRefPayload) *validate.Error {
	e := validate.New()
	if p == nil {
		return e
	}
	for i, item := range p.Items {
		f := func(field string) string { return fmt.Sprintf("external_refs[%d].%s", i, field) }
		// INSERT (id=0): external_ref_id bắt buộc
		if item.Id == 0 {
			v := item.ExternalRefId
			validate.RequiredID(e, f("external_ref_id"), &v)
		}
	}
	return e
}
