package controller

import (
	"context"
	"net/http"

	"userv4/internal/httpapi"
	"userv4/modules/bdspro/dto"
	"userv4/modules/bdspro/services"
)

// PropertyController owns core property HTTP routes under /v2/bdspro/v2/property.
type PropertyController struct {
	property services.PropertyService
}

func NewPropertyController(property services.PropertyService) *PropertyController {
	return &PropertyController{property: property}
}

func (c *PropertyController) RegisterHTTP(mux *http.ServeMux) {
	httpapi.GET(mux, "/v2/bdspro/v2/property/me", c.listMe)
	httpapi.GET(mux, "/v2/bdspro/v2/property/{id}/detail", c.getDetail)
	httpapi.GET(mux, "/v2/bdspro/v2/property/{id}/timestamp", c.getTimestamps)
}

func (c *PropertyController) listMe(ctx context.Context, req dto.PropertyListRequest) (*httpapi.Result, error) {
	result, err := c.property.ListMe(ctx, req)
	if err != nil {
		return nil, err
	}
	return propertyListResult(result)
}

func (c *PropertyController) getDetail(ctx context.Context, req dto.PropertyDetailRequest) (*httpapi.Result, error) {
	result, err := c.property.GetDetail(ctx, req.ToListRequest())
	if err != nil {
		return nil, err
	}
	return protoResult(result)
}

func (c *PropertyController) getTimestamps(ctx context.Context, req httpapi.PathUint64) (*httpapi.Result, error) {
	result, err := c.property.GetTimestamps(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return httpapi.OK(map[string]int64{}), nil
	}
	return httpapi.OK(result.Timestamps), nil
}
