package controller

import (
	"context"
	"net/http"

	"common/httpapi"
	"mdn4/modules/hub/dto"
	"mdn4/modules/hub/services"
)

// LocationController owns v1 location HTTP routes under /v2/hub.
type LocationController struct {
	location services.LocationService
}

func NewLocationController(location services.LocationService) *LocationController {
	return &LocationController{location: location}
}

func (c *LocationController) RegisterHTTP(mux *http.ServeMux) {
	httpapi.GET(mux, "/v2/hub/provinces", c.getProvinces)
	httpapi.GET(mux, "/v2/hub/districts", c.getDistricts)
	httpapi.GET(mux, "/v2/hub/wards", c.getWards)
	httpapi.GET(mux, "/v2/hub/location/search", c.searchLocation)
}

func (c *LocationController) getProvinces(ctx context.Context, req dto.GetProvincesRequest) (*httpapi.Result, error) {
	result, err := c.location.GetProvinces(ctx, req)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return emptyProtoListResult()
	}
	return protoListWithTotalResult(result.GetData(), result.GetTotal())
}

func (c *LocationController) getDistricts(ctx context.Context, req dto.GetDistrictsRequest) (*httpapi.Result, error) {
	result, err := c.location.GetDistricts(ctx, req)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return emptyProtoListResult()
	}
	return protoListWithTotalResult(result.GetData(), result.GetTotal())
}

func (c *LocationController) getWards(ctx context.Context, req dto.GetWardsRequest) (*httpapi.Result, error) {
	result, err := c.location.GetWards(ctx, req)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return emptyProtoListResult()
	}
	return protoListWithTotalResult(result.GetData(), result.GetTotal())
}

func (c *LocationController) searchLocation(ctx context.Context, req dto.SearchLocationRequest) (*httpapi.Result, error) {
	result, err := c.location.SearchLocation(ctx, req)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return emptyProtoListResult()
	}
	return protoListWithTotalResult(result.GetData(), result.GetTotal())
}
