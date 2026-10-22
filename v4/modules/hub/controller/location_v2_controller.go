package controller

import (
	"context"
	"net/http"

	"common/httpapi"
	"mdn4/modules/hub/dto"
	"mdn4/modules/hub/services"
)

// LocationV2Controller owns v2 location HTTP routes under /v2/hub/v2.
type LocationV2Controller struct {
	location services.LocationV2Service
}

func NewLocationV2Controller(location services.LocationV2Service) *LocationV2Controller {
	return &LocationV2Controller{location: location}
}

func (c *LocationV2Controller) RegisterHTTP(mux *http.ServeMux) {
	httpapi.GET(mux, "/v2/hub/v2/provinces", c.getProvinces)
	httpapi.GET(mux, "/v2/hub/v2/province/{code}", c.getProvinceByCode)
	httpapi.GET(mux, "/v2/hub/v2/wards", c.getWards)
	httpapi.GET(mux, "/v2/hub/v2/ward/{code}", c.getWardByCode)
	httpapi.GET(mux, "/v2/hub/v2/location/search", c.searchLocation)
}

func (c *LocationV2Controller) getProvinces(ctx context.Context, req dto.GetProvincesV2Request) (*httpapi.Result, error) {
	result, err := c.location.GetProvincesV2(ctx, req)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return emptyProtoListResult()
	}
	return protoListWithTotalResult(result.GetData(), result.GetTotal())
}

func (c *LocationV2Controller) getProvinceByCode(ctx context.Context, req dto.GetProvinceV2ByCodeRequest) (*httpapi.Result, error) {
	result, err := c.location.GetProvinceV2ByCode(ctx, req.Code, req)
	if err != nil {
		return nil, err
	}
	return wrappedDataResult(result)
}

func (c *LocationV2Controller) getWards(ctx context.Context, req dto.GetWardsV2Request) (*httpapi.Result, error) {
	result, err := c.location.GetWardsV2(ctx, req)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return emptyProtoListResult()
	}
	return protoListWithTotalResult(result.GetData(), result.GetTotal())
}

func (c *LocationV2Controller) getWardByCode(ctx context.Context, req dto.GetWardV2ByCodeRequest) (*httpapi.Result, error) {
	result, err := c.location.GetWardV2ByCode(ctx, req.Code)
	if err != nil {
		return nil, err
	}
	return wrappedDataResult(result)
}

func (c *LocationV2Controller) searchLocation(ctx context.Context, req dto.SearchLocationV2Request) (*httpapi.Result, error) {
	result, err := c.location.SearchLocationV2(ctx, req)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return emptyProtoListResult()
	}
	return protoListWithTotalResult(result.GetData(), result.GetTotal())
}
