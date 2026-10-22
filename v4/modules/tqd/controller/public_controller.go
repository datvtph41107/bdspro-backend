package controller

import (
	"context"
	"net/http"

	"common/httpapi"
	"mdn4/modules/tqd/dto"
	"mdn4/modules/tqd/services"
	sharepb "pb/types/shared"
)

// PublicController owns public TQD HTTP routes under /v2/tqd/public.
type PublicController struct {
	public services.PublicService
}

func NewPublicController(public services.PublicService) *PublicController {
	return &PublicController{public: public}
}

func (c *PublicController) RegisterHTTP(mux *http.ServeMux) {
	httpapi.GET(mux, "/v2/tqd/public/enums", c.getEnums)
}

func (c *PublicController) getEnums(ctx context.Context, req dto.GetEnumsRequest) (*httpapi.Result, error) {
	result, err := c.public.GetEnums(ctx, req.Timestamp)
	if err != nil {
		return nil, err
	}
	return getEnumsResult(result)
}

func getEnumsResult(resp *sharepb.GetEnumsResponse) (*httpapi.Result, error) {
	if resp == nil {
		return httpapi.OK([]any{}), nil
	}
	data, err := protoMessages(resp.GetEnums())
	if err != nil {
		return nil, err
	}
	var opts []httpapi.OKOption
	if resp.GetTimestamp() > 0 {
		opts = append(opts, httpapi.WithTimestamp(resp.GetTimestamp()))
	}
	if !resp.GetLastest() {
		return httpapi.OK([]any{}, opts...), nil
	}
	return httpapi.OK(data, opts...), nil
}
