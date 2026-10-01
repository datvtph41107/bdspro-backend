package controller

import (
	"context"
	"net/http"

	_dto "common/domain/dto"
	_errors "common/errors"
	_utils "common/utils"
	tqdpb "pb/types/tqd"
	"userv4/internal/httpapi"
	"userv4/modules/tqd/dto"
	qh "userv4/modules/tqd/entity/qh"
	"userv4/modules/tqd/services"
)

// LandUseController owns /v2/tqd/client/layers/{layerId}/land-use routes.
type LandUseController struct {
	landUse services.QHLayerLandUseGroupUsecase
}

func NewLandUseController(landUse services.QHLayerLandUseGroupUsecase) *LandUseController {
	return &LandUseController{landUse: landUse}
}

func (c *LandUseController) RegisterHTTP(mux *http.ServeMux) {
	httpapi.GET(mux, "/v2/tqd/client/layers/{layerId}/land-use", c.listByLayer)
	httpapi.POST(mux, "/v2/tqd/client/layers/{layerId}/land-use", c.create)
}

func (c *LandUseController) listByLayer(ctx context.Context, req dto.ClientListLandUseRequest) (*httpapi.Result, error) {
	if req.LayerID == 0 {
		return nil, _errors.ReturnError(400, "layerId is required")
	}
	pagable := &_dto.Pagable{Page: req.Page, Size: req.Size}
	rows, total, err := c.landUse.ListClientByLayer(ctx, req.LayerID, pagable)
	if err != nil {
		return nil, err
	}
	items := make([]*tqdpb.QHLandUseResponse, len(rows))
	for i := range rows {
		items[i] = toQHLandUseResponse(&rows[i])
	}
	return protoListResult(items, total, pagable.GetPage(), pagable.GetSize())
}

func (c *LandUseController) create(ctx context.Context, req dto.ClientCreateLandUseRequest) (*httpapi.Result, error) {
	if req.LayerID == 0 {
		return nil, _errors.ReturnError(400, "layerId is required")
	}
	if req.LandUseID == 0 {
		return nil, _errors.ReturnError(400, "landUseId is required")
	}
	out, err := c.landUse.CreateLayerLandUse(ctx, req.LayerID, req.LandUseID)
	if err != nil {
		return nil, err
	}
	return protoResult(toQHLandUseResponse(out))
}

func toQHLandUseResponse(e *qh.QHLandUse) *tqdpb.QHLandUseResponse {
	if e == nil {
		return nil
	}
	return &tqdpb.QHLandUseResponse{
		Id:        e.ID,
		Note:      e.Note,
		Color:     e.Color,
		Code:      e.Code,
		Name:      e.Name,
		CreatedAt: _utils.FormatTimeToString(e.CreatedAt),
		UpdatedAt: _utils.FormatTimeToString(e.UpdatedAt),
	}
}
