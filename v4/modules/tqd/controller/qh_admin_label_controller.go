package controller

import (
	"context"
	"net/http"
	"strings"

	_dto "common/domain/dto"
	_errors "common/errors"
	"common/httpapi"
	_utils "common/utils"
	"mdn4/modules/tqd/dto"
	qh "mdn4/modules/tqd/entity/qh"
	"mdn4/modules/tqd/enums"
	"mdn4/modules/tqd/services"
	tqdpb "pb/types/tqd"
)

// QHAdminLabelController owns /v2/tqd/qh/admin/labels*.
type QHAdminLabelController struct {
	labels services.QHLabelUsecase
}

func NewQHAdminLabelController(labels services.QHLabelUsecase) *QHAdminLabelController {
	return &QHAdminLabelController{labels: labels}
}

func (c *QHAdminLabelController) RegisterHTTP(mux *http.ServeMux) {
	httpapi.POST(mux, "/v2/tqd/qh/admin/labels", c.create)
	httpapi.POST(mux, "/v2/tqd/qh/admin/labels/batch", c.batchCreate)
	httpapi.POST(mux, "/v2/tqd/qh/admin/labels/batch-delete", c.batchDelete)
	httpapi.POST(mux, "/v2/tqd/qh/admin/labels/merge", c.merge)
	httpapi.GET(mux, "/v2/tqd/qh/admin/labels/list", c.list)
	httpapi.GET(mux, "/v2/tqd/qh/admin/labels/{id}", c.get)
	httpapi.PUT(mux, "/v2/tqd/qh/admin/labels/{id}", c.update)
	httpapi.DELETE(mux, "/v2/tqd/qh/admin/labels/{id}", c.delete)
	httpapi.GET(mux, "/v2/tqd/qh/admin/layers/{layerId}/labels", c.listByLayer)
}

type createLabelAdminReq struct {
	LayerID         uint64  `json:"layerId"`
	Name            string  `json:"name"`
	DisplayName     string  `json:"displayName"`
	Description     string  `json:"description"`
	Color           string  `json:"color"`
	FillOpacity     float32 `json:"fillOpacity"`
	StrokeColor     string  `json:"strokeColor"`
	StrokeWidth     int32   `json:"strokeWidth"`
	StrokeDashArray string  `json:"strokeDashArray"`
	DisplayOrder    int32   `json:"displayOrder"`
	IsVisible       bool    `json:"isVisible"`
	MinZoom         int32   `json:"minZoom"`
	MaxZoom         int32   `json:"maxZoom"`
	Status          *uint32 `json:"status"`
	StandardAt      *string `json:"standardAt"`
}

type updateLabelAdminReq struct {
	ID              uint64   `path:"id" json:"-"`
	Name            *string  `json:"name"`
	DisplayName     *string  `json:"displayName"`
	Description     *string  `json:"description"`
	Color           *string  `json:"color"`
	FillOpacity     *float32 `json:"fillOpacity"`
	StrokeColor     *string  `json:"strokeColor"`
	StrokeWidth     *int32   `json:"strokeWidth"`
	StrokeDashArray *string  `json:"strokeDashArray"`
	DisplayOrder    *int32   `json:"displayOrder"`
	IsVisible       *bool    `json:"isVisible"`
	MinZoom         *int32   `json:"minZoom"`
	MaxZoom         *int32   `json:"maxZoom"`
	Status          *uint32  `json:"status"`
	StandardAt      *string  `json:"standardAt"`
	LandUseID       *uint64  `json:"landUseId"`
	LegendID        *uint64  `json:"legendId"`
	LayerID         *uint64  `json:"layerId"`
}

type batchCreateLabelsReq struct {
	LayerID uint64                `json:"layerId"`
	Items   []createLabelAdminReq `json:"items"`
}

type batchDeleteLabelsReq struct {
	IDs []uint64 `json:"ids"`
}

type mergeLabelsReq struct {
	LayerID        uint64              `json:"layerId"`
	SourceLabelIDs []uint64            `json:"sourceLabelIds"`
	NewLabel       createLabelAdminReq `json:"newLabel"`
}

type listLabelsAdminReq struct {
	Page            uint32  `form:"page"`
	Size            uint32  `form:"size"`
	LayerID         *uint64 `form:"layerId"`
	IncludeInactive bool    `form:"includeInactive"`
}

type listLabelsByLayerAdminReq struct {
	LayerID         uint64 `path:"layerId"`
	Page            uint32 `form:"page"`
	Size            uint32 `form:"size"`
	IncludeInactive bool   `form:"includeInactive"`
}

func (c *QHAdminLabelController) create(ctx context.Context, req createLabelAdminReq) (*httpapi.Result, error) {
	label := labelFromCreate(req)
	out, err := c.labels.Create(ctx, label)
	if err != nil {
		return nil, err
	}
	return protoResult(toQHLabelResponse(out))
}

func (c *QHAdminLabelController) batchCreate(ctx context.Context, req batchCreateLabelsReq) (*httpapi.Result, error) {
	if req.LayerID == 0 {
		return nil, _errors.ReturnError(400, "layerId is required")
	}
	items := make([]*qh.QHLabel, 0, len(req.Items))
	for _, it := range req.Items {
		it.LayerID = req.LayerID
		items = append(items, labelFromCreate(it))
	}
	created, skipped, err := c.labels.BatchCreate(ctx, req.LayerID, items)
	if err != nil {
		return nil, err
	}
	data := make([]*tqdpb.QHLabelResponse, len(created))
	for i := range created {
		data[i] = toQHLabelResponse(created[i])
	}
	encoded, err := protoMessages(data)
	if err != nil {
		return nil, err
	}
	return jsonOK(map[string]any{"created": encoded, "skippedNames": skipped})
}

func (c *QHAdminLabelController) batchDelete(ctx context.Context, req batchDeleteLabelsReq) (*httpapi.Result, error) {
	if err := c.labels.BatchDelete(ctx, req.IDs); err != nil {
		return nil, err
	}
	return jsonOK(map[string]string{"message": "deleted"})
}

func (c *QHAdminLabelController) merge(ctx context.Context, req mergeLabelsReq) (*httpapi.Result, error) {
	if req.LayerID == 0 {
		return nil, _errors.ReturnError(400, "layerId is required")
	}
	req.NewLabel.LayerID = req.LayerID
	out, err := c.labels.Merge(ctx, req.LayerID, req.SourceLabelIDs, labelFromCreate(req.NewLabel))
	if err != nil {
		return nil, err
	}
	return protoResult(toQHLabelResponse(out))
}

func (c *QHAdminLabelController) get(ctx context.Context, req dto.IDPath) (*httpapi.Result, error) {
	if req.ID == 0 {
		return nil, _errors.ReturnError(400, "id is required")
	}
	out, err := c.labels.GetByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	return protoResult(toQHLabelResponse(out))
}

func (c *QHAdminLabelController) update(ctx context.Context, req updateLabelAdminReq) (*httpapi.Result, error) {
	if req.ID == 0 {
		return nil, _errors.ReturnError(400, "id is required")
	}
	existing, err := c.labels.GetByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, _errors.ReturnError(404, "label not found")
	}
	applyLabelUpdate(existing, req)
	out, err := c.labels.Update(ctx, req.ID, existing)
	if err != nil {
		return nil, err
	}
	if req.LayerID != nil && *req.LayerID > 0 {
		_ = c.labels.UpdateLayerLinkMetadata(ctx, req.ID, *req.LayerID, req.LandUseID, req.LegendID)
	}
	return protoResult(toQHLabelResponse(out))
}

func (c *QHAdminLabelController) delete(ctx context.Context, req dto.IDPath) (*httpapi.Result, error) {
	if req.ID == 0 {
		return nil, _errors.ReturnError(400, "id is required")
	}
	if err := c.labels.Delete(ctx, req.ID); err != nil {
		return nil, err
	}
	return jsonOK(map[string]string{"message": "deleted"})
}

func (c *QHAdminLabelController) list(ctx context.Context, req listLabelsAdminReq) (*httpapi.Result, error) {
	pagable := &_dto.Pagable{Page: req.Page, Size: req.Size}
	rows, total, err := c.labels.ListAdmin(ctx, req.LayerID, pagable, req.IncludeInactive)
	if err != nil {
		return nil, err
	}
	items := make([]*tqdpb.QHLabelResponse, len(rows))
	for i := range rows {
		items[i] = toQHLabelResponse(&rows[i])
	}
	return protoListResult(items, total, pagable.GetPage(), pagable.GetSize())
}

func (c *QHAdminLabelController) listByLayer(ctx context.Context, req listLabelsByLayerAdminReq) (*httpapi.Result, error) {
	if req.LayerID == 0 {
		return nil, _errors.ReturnError(400, "layerId is required")
	}
	pagable := &_dto.Pagable{Page: req.Page, Size: req.Size}
	rows, total, err := c.labels.ListByLayerID(ctx, req.LayerID, pagable, req.IncludeInactive)
	if err != nil {
		return nil, err
	}
	items := make([]*tqdpb.QHLabelResponse, len(rows))
	for i := range rows {
		items[i] = toQHLabelResponse(&rows[i])
	}
	return protoListResult(items, total, pagable.GetPage(), pagable.GetSize())
}

func labelFromCreate(req createLabelAdminReq) *qh.QHLabel {
	label := &qh.QHLabel{
		LayerID: req.LayerID, Name: req.Name, DisplayName: req.DisplayName, Description: req.Description,
		Color: req.Color, FillOpacity: float64(req.FillOpacity), StrokeColor: req.StrokeColor,
		StrokeWidth: int(req.StrokeWidth), StrokeDashArray: req.StrokeDashArray,
		DisplayOrder: int(req.DisplayOrder), IsVisible: req.IsVisible,
		MinZoom: int(req.MinZoom), MaxZoom: int(req.MaxZoom),
	}
	if req.Status != nil {
		label.Status = enums.LabelStatus(*req.Status)
	} else {
		label.Status = enums.LabelStatusActive
	}
	if req.StandardAt != nil {
		s := strings.TrimSpace(*req.StandardAt)
		if s != "" {
			label.StandardAt = _utils.ParseStringToTime(s)
		}
	}
	return label
}

func applyLabelUpdate(existing *qh.QHLabel, req updateLabelAdminReq) {
	if req.Name != nil {
		existing.Name = *req.Name
	}
	if req.DisplayName != nil {
		existing.DisplayName = *req.DisplayName
	}
	if req.Description != nil {
		existing.Description = *req.Description
	}
	if req.Color != nil {
		existing.Color = *req.Color
	}
	if req.FillOpacity != nil {
		existing.FillOpacity = float64(*req.FillOpacity)
	}
	if req.StrokeColor != nil {
		existing.StrokeColor = *req.StrokeColor
	}
	if req.StrokeWidth != nil {
		existing.StrokeWidth = int(*req.StrokeWidth)
	}
	if req.StrokeDashArray != nil {
		existing.StrokeDashArray = *req.StrokeDashArray
	}
	if req.DisplayOrder != nil {
		existing.DisplayOrder = int(*req.DisplayOrder)
	}
	if req.IsVisible != nil {
		existing.IsVisible = *req.IsVisible
	}
	if req.MinZoom != nil {
		existing.MinZoom = int(*req.MinZoom)
	}
	if req.MaxZoom != nil {
		existing.MaxZoom = int(*req.MaxZoom)
	}
	if req.Status != nil {
		existing.Status = enums.LabelStatus(*req.Status)
	}
	if req.StandardAt != nil {
		s := strings.TrimSpace(*req.StandardAt)
		if s == "" {
			existing.StandardAt = nil
		} else {
			existing.StandardAt = _utils.ParseStringToTime(s)
		}
	}
	if req.LandUseID != nil {
		existing.LandUseID = req.LandUseID
	}
	if req.LegendID != nil {
		existing.LegendID = req.LegendID
	}
}
