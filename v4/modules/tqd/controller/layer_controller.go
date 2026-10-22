package controller

import (
	"context"
	"net/http"

	_dto "common/domain/dto"
	_errors "common/errors"
	"common/httpapi"
	"mdn4/modules/tqd/dto"
	"mdn4/modules/tqd/enums"
	"mdn4/modules/tqd/mapper"
	"mdn4/modules/tqd/repository"
	"mdn4/modules/tqd/services"
	tqdpb "pb/types/tqd"
)

// LayerController owns /v2/tqd/client/layers* routes.
type LayerController struct {
	layers services.LayerUsecase
}

func NewLayerController(layers services.LayerUsecase) *LayerController {
	return &LayerController{layers: layers}
}

func (c *LayerController) RegisterHTTP(mux *http.ServeMux) {
	httpapi.GET(mux, "/v2/tqd/client/layers", c.listLayers)
	httpapi.GET(mux, "/v2/tqd/client/layers/{id}", c.getLayer)
	httpapi.POST(mux, "/v2/tqd/client/layers/{layerId}/build-tile", c.buildTile)
}

func (c *LayerController) listLayers(ctx context.Context, req dto.ClientListLayersRequest) (*httpapi.Result, error) {
	filter := &repository.ClientLayerFilter{
		Pagable: _dto.Pagable{
			Page: req.Page,
			Size: req.Size,
		},
		Search:            req.Search,
		DocID:             req.DocID,
		PlanningProjectID: req.PlanningProjectID,
	}
	if req.Type != nil {
		t := enums.LayerType(*req.Type)
		filter.Type = &t
	}
	rows, total, err := c.layers.ListClient(ctx, filter)
	if err != nil {
		return nil, err
	}
	items := make([]*tqdpb.LayerResponse, 0, len(rows))
	for i := range rows {
		protoLayer, err := mapper.ToProtoLayer(&rows[i])
		if err != nil {
			continue
		}
		items = append(items, protoLayer)
	}
	return protoListResult(items, total, req.Page, req.Size)
}

func (c *LayerController) getLayer(ctx context.Context, req dto.ClientGetLayerRequest) (*httpapi.Result, error) {
	layer, err := c.layers.GetClientByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	if layer == nil {
		return nil, _errors.ReturnError(404, "layer not found")
	}
	protoLayer, err := mapper.ToProtoLayer(layer)
	if err != nil {
		return nil, err
	}
	return protoResult(protoLayer)
}

func (c *LayerController) buildTile(ctx context.Context, req dto.ClientBuildLayerTileRequest) (*httpapi.Result, error) {
	if req.LayerID == 0 {
		return nil, _errors.ReturnError(400, "layerId is required")
	}
	minZoom := req.MinZoom
	maxZoom := req.MaxZoom
	if maxZoom == 0 {
		maxZoom = 14
	}
	if maxZoom > 16 {
		maxZoom = 16
	}
	if minZoom > maxZoom {
		minZoom = maxZoom
	}
	outputDir := "./files/tiles"
	if req.OutputDir != nil && *req.OutputDir != "" {
		outputDir = *req.OutputDir
	}
	if err := c.layers.BuildPMTiles(ctx, req.LayerID, minZoom, maxZoom, outputDir, nil); err != nil {
		return nil, err
	}
	return jsonOK(map[string]any{
		"status":  "completed",
		"layerId": req.LayerID,
	})
}
