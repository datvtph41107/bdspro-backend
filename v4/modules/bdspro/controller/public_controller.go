package controller

import (
	"context"
	"net/http"

	"common/httpapi"
	"mdn4/modules/bdspro/dto"
	"mdn4/modules/bdspro/services"
)

// PublicController owns public bdspro HTTP routes under /v2/bdspro/public.
type PublicController struct {
	public services.PublicService
}

func NewPublicController(public services.PublicService) *PublicController {
	return &PublicController{public: public}
}

func (c *PublicController) RegisterHTTP(mux *http.ServeMux) {
	httpapi.GET(mux, "/v2/bdspro/public/enums", c.getEnums)
	httpapi.GET(mux, "/v2/bdspro/public/enums/{text}", c.getEnumByName)
}

func (c *PublicController) getEnums(ctx context.Context, req dto.GetEnumsRequest) (*httpapi.Result, error) {
	result, err := c.public.GetEnums(ctx, req.Timestamp)
	if err != nil {
		return nil, err
	}
	return getEnumsResult(result)
}

func (c *PublicController) getEnumByName(ctx context.Context, req dto.GetEnumByNameRequest) (*httpapi.Result, error) {
	if req.Text == services.DealSettingEnumKey {
		return httpapi.OK(map[string]any{"data": services.DealSettingEnumPayload()}), nil
	}
	result, err := c.public.GetEnumByName(ctx, req.Text)
	if err != nil {
		return nil, err
	}
	return protoResult(result)
}
