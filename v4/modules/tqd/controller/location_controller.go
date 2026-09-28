package controller

import (
	"context"
	"net/http"

	_errors "common/errors"
	tqdpb "pb/types/tqd"
	"userv4/internal/httpapi"
	"userv4/modules/tqd/dto"
	"userv4/modules/tqd/services"
)

// LocationController owns TQD client location HTTP routes under /v2/tqd/client.
type LocationController struct {
	location services.LocationService
}

func NewLocationController(location services.LocationService) *LocationController {
	return &LocationController{location: location}
}

func (c *LocationController) RegisterHTTP(mux *http.ServeMux) {
	httpapi.GET(mux, "/v2/tqd/client/provinces", c.listProvinces)
	httpapi.GET(mux, "/v2/tqd/client/wards", c.listWards)
}

func (c *LocationController) listProvinces(ctx context.Context, req dto.ListProvincesRequest) (*httpapi.Result, error) {
	result, err := c.location.ListProvinces(ctx, req)
	if err != nil {
		return nil, err
	}
	return listProvincesResult(result)
}

func (c *LocationController) listWards(ctx context.Context, req dto.ListWardsRequest) (*httpapi.Result, error) {
	if req.Page != nil && *req.Page < 1 {
		return nil, _errors.ReturnError(400, "page must be >= 1")
	}
	result, err := c.location.ListWards(ctx, req)
	if err != nil {
		return nil, err
	}
	return listWardsResult(result)
}

func listProvincesResult(resp *tqdpb.ListProvincesResponse) (*httpapi.Result, error) {
	if resp == nil {
		return httpapi.OK([]any{}), nil
	}
	data, err := protoMessages(resp.GetProvinces())
	if err != nil {
		return nil, err
	}
	var opts []httpapi.OKOption
	if resp.GetTimestamp() > 0 {
		opts = append(opts, httpapi.WithTimestamp(resp.GetTimestamp()))
	}
	return httpapi.OK(data, opts...), nil
}

func listWardsResult(resp *tqdpb.ClientListWardsResponse) (*httpapi.Result, error) {
	if resp == nil {
		return httpapi.OK([]any{}), nil
	}
	data, err := protoMessages(resp.GetData())
	if err != nil {
		return nil, err
	}
	var opts []httpapi.OKOption
	if p := resp.GetPagination(); p != nil {
		opts = append(opts, httpapi.WithPagination(&httpapi.Pagination{
			Page:  p.GetPage(),
			Size:  p.GetSize(),
			Total: p.GetTotal(),
		}))
	}
	if resp.GetTimestamp() > 0 {
		opts = append(opts, httpapi.WithTimestamp(resp.GetTimestamp()))
	}
	return httpapi.OK(data, opts...), nil
}
