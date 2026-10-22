package controller

import (
	"context"
	"net/http"

	"common/httpapi"
	bdspropb "pb/types/bdspro"
	"mdn4/modules/bdspro/dto"
	"mdn4/modules/bdspro/services"
)

// PropertyController owns core property HTTP routes under /v2/bdspro/v2/property.
type PropertyController struct {
	property     services.PropertyService
	propertyType services.PropertyTypeService
}

func NewPropertyController(property services.PropertyService, propertyType services.PropertyTypeService) *PropertyController {
	return &PropertyController{property: property, propertyType: propertyType}
}

func (c *PropertyController) RegisterHTTP(mux *http.ServeMux) {
	httpapi.GET(mux, "/v2/bdspro/v2/property/me", c.listMe)
	httpapi.GET(mux, "/v2/bdspro/v2/property/detail/{id}", c.getDetailByPath)
	httpapi.GET(mux, "/v4/bdspro/property/detail/{id}", c.getDetailByPath)
	// Literal "timestamp" segment avoids ServeMux conflict with .../property/detail/{id}.
	httpapi.GET(mux, "/v2/bdspro/v2/property/timestamp/{id}", c.getTimestamps)
	httpapi.GET(mux, "/v4/bdspro/property/timestamp/{id}", c.getTimestamps)
	httpapi.POST(mux, "/v2/bdspro/property/new", c.createProperty)
	httpapi.GET(mux, "/v2/bdspro/v2/list/property-type", c.listPropertyTypes)
}

func (c *PropertyController) listPropertyTypes(ctx context.Context, req dto.PropertyTypeListRequest) (*httpapi.Result, error) {
	out, err := c.propertyType.List(ctx, req)
	if err != nil {
		return nil, err
	}
	return httpapi.OKResponse(&httpapi.ResultResponse{
		Data:      out.List.Data,
		Total:     out.List.TotalElements,
		Page:      int32(req.Page),
		Size:      int32(req.Size),
		Timestamp: out.Timestamp,
	})
}

func (c *PropertyController) createProperty(ctx context.Context, req bdspropb.CreatePropertyRequest) (*httpapi.Result, error) {
	result, err := c.property.CreateProperty(ctx, &req)
	if err != nil {
		return nil, err
	}
	return protoResult(result)
}

func (c *PropertyController) listMe(ctx context.Context, req dto.PropertyListRequest) (*httpapi.Result, error) {
	result, err := c.property.ListMe(ctx, req)
	if err != nil {
		return nil, err
	}
	return propertyListResult(result)
}

func (c *PropertyController) getDetailByPath(ctx context.Context, req dto.SyncListByIDRequest) (*httpapi.Result, error) {
	return c.propertyDetailResult(ctx, req.ToSyncListRequest())
}

func (c *PropertyController) propertyDetailResult(ctx context.Context, req dto.SyncListRequest) (*httpapi.Result, error) {
	result, ts, err := c.property.GetDetail(ctx, dto.PropertyListRequest{
		ID:        req.ID,
		Page:      req.Page,
		Size:      req.Size,
		Timestamp: req.Timestamp,
	})
	if err != nil {
		return nil, err
	}
	return protoResult(result, timestampOpt(ts)...)
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
