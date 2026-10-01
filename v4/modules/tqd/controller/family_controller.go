package controller

import (
	"context"
	"net/http"

	_dto "common/domain/dto"
	_errors "common/errors"
	tqdpb "pb/types/tqd"
	"userv4/internal/httpapi"
	"userv4/modules/tqd/dto"
	"userv4/modules/tqd/mapper"
	"userv4/modules/tqd/services"
)

// FamilyController owns /v2/tqd/client/families* and layer-families list routes.
type FamilyController struct {
	families services.QHLayerFamilyUsecase
}

func NewFamilyController(families services.QHLayerFamilyUsecase) *FamilyController {
	return &FamilyController{families: families}
}

func (c *FamilyController) RegisterHTTP(mux *http.ServeMux) {
	httpapi.GET(mux, "/v2/tqd/client/qh/layer-families/list", c.listLayerFamilies)
	httpapi.GET(mux, "/v2/tqd/client/families", c.listFamilies)
	httpapi.POST(mux, "/v2/tqd/client/families/{familyId}/build-tile", c.buildTile)
}

func (c *FamilyController) listLayerFamilies(ctx context.Context, req dto.ClientListFamiliesRequest) (*httpapi.Result, error) {
	return c.listFamilies(ctx, req)
}

func (c *FamilyController) listFamilies(ctx context.Context, req dto.ClientListFamiliesRequest) (*httpapi.Result, error) {
	pagable := &_dto.Pagable{Page: req.Page, Size: req.Size}
	search := ""
	if req.Search != nil {
		search = *req.Search
	}
	rows, total, err := c.families.ListClient(ctx, pagable, search)
	if err != nil {
		return nil, err
	}
	items := make([]*tqdpb.QHLayerFamilyResponse, 0, len(rows))
	for i := range rows {
		family := &rows[i]
		minZoom, maxZoom, zErr := c.families.GetZoomRangeByFamilyID(ctx, family.ID)
		if zErr != nil {
			minZoom, maxZoom = 6, 16
		}
		items = append(items, mapper.ToProtoQHLayerFamilyClient(family, minZoom, maxZoom))
	}
	return protoListResult(items, total, pagable.GetPage(), pagable.GetSize())
}

func (c *FamilyController) buildTile(ctx context.Context, req dto.ClientBuildFamilyTileRequest) (*httpapi.Result, error) {
	if req.FamilyID == 0 {
		return nil, _errors.ReturnError(400, "familyId is required")
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
	if err := c.families.BuildFamilyPMTiles(ctx, req.FamilyID, minZoom, maxZoom, outputDir, nil); err != nil {
		return nil, err
	}
	return jsonOK(map[string]any{
		"status":   "completed",
		"familyId": req.FamilyID,
	})
}
