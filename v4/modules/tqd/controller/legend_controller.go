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

// LegendController owns /v2/tqd/client/layers/{layerId}/legends.
type LegendController struct {
	legends services.QHLayerLegendUsecase
}

func NewLegendController(legends services.QHLayerLegendUsecase) *LegendController {
	return &LegendController{legends: legends}
}

func (c *LegendController) RegisterHTTP(mux *http.ServeMux) {
	httpapi.GET(mux, "/v2/tqd/client/layers/{layerId}/legends", c.listByLayer)
}

func (c *LegendController) listByLayer(ctx context.Context, req dto.ClientListLegendsRequest) (*httpapi.Result, error) {
	if req.LayerID == 0 {
		return nil, _errors.ReturnError(400, "layerId is required")
	}
	pagable := &_dto.Pagable{Page: req.Page, Size: req.Size}
	rows, total, err := c.legends.ListClientByLayer(ctx, req.LayerID, pagable)
	if err != nil {
		return nil, err
	}
	items := make([]*tqdpb.QHLayerLegendResponse, len(rows))
	for i := range rows {
		items[i] = toQHLayerLegendResponse(&rows[i])
	}
	return protoListResult(items, total, pagable.GetPage(), pagable.GetSize())
}

func toQHLayerLegendResponse(e *qh.QHLayerLegend) *tqdpb.QHLayerLegendResponse {
	if e == nil {
		return nil
	}
	resp := &tqdpb.QHLayerLegendResponse{
		Id:              e.ID,
		LayerId:         e.LayerID,
		LabelId:         e.LabelID,
		Note:            e.Note,
		DisplayOrder:    int32(e.DisplayOrder),
		IsVisible:       e.IsVisible,
		LegendType:      e.GetLegendType(),
		GeometryType:    e.GetGeometryType(),
		StrokeDashArray: e.GetEffectiveStrokeDash(),
		IconUrl:         e.IconURL,
		ImageUrl:        e.ImageURL,
		Description:     e.Description,
		Color:           e.GetDisplayColor(),
		LandUseGroupId:  e.LandUseID,
		CreatedAt:       _utils.FormatTimeToString(e.CreatedAt),
		UpdatedAt:       _utils.FormatTimeToString(e.UpdatedAt),
	}
	if e.Label != nil {
		resp.LabelStyle = &tqdpb.LegendLabelStyleInfo{
			Id:              e.Label.ID,
			Name:            e.Label.Name,
			DisplayName:     e.Label.DisplayName,
			Color:           e.Label.Color,
			FillOpacity:     e.Label.FillOpacity,
			StrokeColor:     e.Label.StrokeColor,
			StrokeWidth:     int32(e.Label.StrokeWidth),
			StrokeDashArray: e.Label.StrokeDashArray,
		}
	}
	if e.LandUse != nil {
		resp.LandUseGroupInfo = &tqdpb.LandUseGroupInfo{
			Id:          e.LandUse.ID,
			Code:        e.LandUse.Code,
			Name:        e.LandUse.Name,
			Description: e.LandUse.Description,
			Color:       e.LandUse.Color,
		}
	}
	return resp
}
