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

// LabelController owns /v2/tqd/client label routes.
type LabelController struct {
	labels services.QHLabelUsecase
}

func NewLabelController(labels services.QHLabelUsecase) *LabelController {
	return &LabelController{labels: labels}
}

func (c *LabelController) RegisterHTTP(mux *http.ServeMux) {
	httpapi.GET(mux, "/v2/tqd/client/layers/{layerId}/labels", c.listByLayer)
	httpapi.GET(mux, "/v2/tqd/client/labels/{id}", c.getByID)
}

func (c *LabelController) listByLayer(ctx context.Context, req dto.ClientListLabelsByLayerRequest) (*httpapi.Result, error) {
	if req.LayerID == 0 {
		return nil, _errors.ReturnError(400, "layerId is required")
	}
	pagable := &_dto.Pagable{Page: req.Page, Size: req.Size}
	rows, total, err := c.labels.ListClientByLayerID(ctx, req.LayerID, pagable)
	if err != nil {
		return nil, err
	}
	items := make([]*tqdpb.QHLabelResponse, len(rows))
	for i := range rows {
		items[i] = toQHLabelResponse(&rows[i])
	}
	return protoListResult(items, total, pagable.GetPage(), pagable.GetSize())
}

func (c *LabelController) getByID(ctx context.Context, req dto.ClientGetLabelRequest) (*httpapi.Result, error) {
	if req.ID == 0 {
		return nil, _errors.ReturnError(400, "id is required")
	}
	label, err := c.labels.GetClientByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	if label == nil {
		return nil, _errors.ReturnError(404, "label not found")
	}
	return protoResult(toQHLabelResponse(label))
}

func toQHLabelResponse(l *qh.QHLabel) *tqdpb.QHLabelResponse {
	if l == nil {
		return nil
	}
	resp := &tqdpb.QHLabelResponse{
		Id:              l.ID,
		LayerId:         l.LayerID,
		Name:            l.Name,
		DisplayName:     l.DisplayName,
		Description:     l.Description,
		Color:           l.Color,
		FillOpacity:     l.FillOpacity,
		StrokeColor:     l.StrokeColor,
		StrokeWidth:     int32(l.StrokeWidth),
		StrokeDashArray: l.StrokeDashArray,
		DisplayOrder:    int32(l.DisplayOrder),
		IsVisible:       l.IsVisible,
		MinZoom:         int32(l.MinZoom),
		MaxZoom:         int32(l.MaxZoom),
		RegionCount:     l.RegionCount,
		Status:          uint32(l.Status),
		CreatedAt:       _utils.FormatTimeToString(&l.CreatedAt),
		UpdatedAt:       _utils.FormatTimeToString(&l.UpdatedAt),
	}
	resp.LandUseId = l.LandUseID
	resp.LegendId = l.LegendID
	if l.StandardAt != nil {
		resp.StandardAt = _utils.FormatTimeToString(l.StandardAt)
	}
	return resp
}
