package controller

import (
	"context"
	"net/http"

	_errors "common/errors"
	tqdpb "pb/types/tqd"
	"userv4/internal/httpapi"
	"userv4/modules/tqd/dto"
	"userv4/modules/tqd/mapper"
	"userv4/modules/tqd/services"
)

// RegionController owns /v2/tqd/client/regions* routes.
type RegionController struct {
	regions services.RegionUsecase
	mapper  *mapper.RegionMapper
}

func NewRegionController(regions services.RegionUsecase) *RegionController {
	return &RegionController{
		regions: regions,
		mapper:  mapper.NewRegionMapper(),
	}
}

func (c *RegionController) RegisterHTTP(mux *http.ServeMux) {
	httpapi.GET(mux, "/v2/tqd/client/regions", c.listRegions)
	httpapi.GET(mux, "/v2/tqd/client/regions/by-point", c.findByPoint)
	httpapi.GET(mux, "/v2/tqd/client/regions/by-bbox", c.findByBBox)
}

func (c *RegionController) listRegions(ctx context.Context, req dto.ClientListRegionsRequest) (*httpapi.Result, error) {
	var bBox *services.BBox
	if req.MinLng != nil && req.MinLat != nil && req.MaxLng != nil && req.MaxLat != nil {
		bBox = &services.BBox{
			MinLng: *req.MinLng,
			MinLat: *req.MinLat,
			MaxLng: *req.MaxLng,
			MaxLat: *req.MaxLat,
		}
	}
	var layerID uint64
	if req.LayerID != nil {
		layerID = *req.LayerID
	}
	rows, total, err := c.regions.ListClient(ctx, layerID, bBox, req.Page, req.Size)
	if err != nil {
		return nil, err
	}
	items := make([]*tqdpb.RegionResponse, 0, len(rows))
	for i := range rows {
		items = append(items, c.mapper.ToProtoClientResponse(&rows[i]))
	}
	return protoListResult(items, total, req.Page, req.Size)
}

func (c *RegionController) findByPoint(ctx context.Context, req dto.ClientFindRegionByPointRequest) (*httpapi.Result, error) {
	if req.Lat < -90 || req.Lat > 90 || req.Lng < -180 || req.Lng > 180 {
		return nil, _errors.ReturnError(400, "invalid coordinates")
	}
	region, err := c.regions.FindByPoint(ctx, req.Lat, req.Lng, req.LayerID)
	if err != nil {
		return nil, err
	}
	if region == nil {
		return nil, _errors.ReturnError(404, "no region found")
	}
	return protoResult(c.mapper.ToProtoClientResponse(region))
}

func (c *RegionController) findByBBox(ctx context.Context, req dto.ClientFindRegionsByBBoxRequest) (*httpapi.Result, error) {
	rows, total, err := c.regions.FindByBBox(ctx, req.MinLng, req.MinLat, req.MaxLng, req.MaxLat, req.LayerID, req.LabelID, req.Page, req.Size)
	if err != nil {
		return nil, err
	}
	items := make([]*tqdpb.RegionResponse, 0, len(rows))
	for i := range rows {
		items = append(items, c.mapper.ToProtoClientResponse(&rows[i]))
	}
	return protoListResult(items, total, uint32(req.Page), uint32(req.Size))
}
