package property_validator

import (
	"common/pkg/validate"
	"fmt"
	bdspropb "pb/types/bdspro"
)

func validateMediaPayload(p *bdspropb.UpdateMediaPayload) *validate.Error {
	e := validate.New()
	if p == nil {
		return e
	}
	for i, item := range p.Items {
		f := func(field string) string { return fmt.Sprintf("media[%d].%s", i, field) }
		// INSERT (id=0): media_url và media_type bắt buộc
		if item.Id == 0 {
			url := item.MediaUrl
			validate.Required(e, f("media_url"), &url)
		}
		mt := item.MediaType
		validate.MediaType(e, f("media_type"), &mt)
	}
	return e
}
