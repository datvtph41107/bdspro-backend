package controller

import (
	"context"
	"net/http"

	_dto "common/domain/dto"
	_errors "common/errors"
	"common/httpapi"
	"mdn4/modules/tqd/dto"
	qh "mdn4/modules/tqd/entity/qh"
	"mdn4/modules/tqd/services"
	tqdpb "pb/types/tqd"
)

// QHAdminLegendController owns /v2/tqd/qh/admin/layer-legends* and /v2/tqd/admin/legends*.
type QHAdminLegendController struct {
	uc services.QHLayerLegendUsecase
}

func NewQHAdminLegendController(uc services.QHLayerLegendUsecase) *QHAdminLegendController {
	return &QHAdminLegendController{uc: uc}
}

func (c *QHAdminLegendController) RegisterHTTP(mux *http.ServeMux) {
	// qh/admin layer-legends
	httpapi.POST(mux, "/v2/tqd/qh/admin/layer-legends", c.create)
	httpapi.POST(mux, "/v2/tqd/qh/admin/layer-legends/batch", c.batchCreate)
	httpapi.GET(mux, "/v2/tqd/qh/admin/layer-legends/list", c.list)
	httpapi.GET(mux, "/v2/tqd/qh/admin/layer-legends/{id}", c.get)
	httpapi.PUT(mux, "/v2/tqd/qh/admin/layer-legends/{id}", c.update)
	httpapi.DELETE(mux, "/v2/tqd/qh/admin/layer-legends/{id}", c.delete)

	// admin/legends (same usecase)
	httpapi.POST(mux, "/v2/tqd/admin/legends", c.create)
	httpapi.GET(mux, "/v2/tqd/admin/legends", c.list)
	httpapi.GET(mux, "/v2/tqd/admin/legends/{id}", c.get)
	httpapi.PUT(mux, "/v2/tqd/admin/legends/{id}", c.update)
	httpapi.DELETE(mux, "/v2/tqd/admin/legends/{id}", c.delete)
	httpapi.GET(mux, "/v2/tqd/admin/layers/{layerId}/legends", c.listByLayer)
}

type createLegendReq struct {
	LayerID         uint64  `json:"layerId"`
	LabelID         uint64  `json:"labelId"`
	Note            string  `json:"note"`
	DisplayOrder    int32   `json:"displayOrder"`
	IsVisible       bool    `json:"isVisible"`
	LegendType      string  `json:"legendType"`
	GeometryType    string  `json:"geometryType"`
	StrokeDashArray string  `json:"strokeDashArray"`
	IconURL         string  `json:"iconUrl"`
	ImageURL        string  `json:"imageUrl"`
	Description     string  `json:"description"`
	Color           string  `json:"color"`
	LandUseID       *uint64 `json:"landUseId"`
}

type batchCreateLegendsReq struct {
	LayerID uint64            `json:"layerId"`
	Items   []createLegendReq `json:"items"`
}

type updateLegendReq struct {
	ID              uint64  `path:"id" json:"-"`
	LayerID         *uint64 `json:"layerId"`
	LabelID         *uint64 `json:"labelId"`
	Note            *string `json:"note"`
	DisplayOrder    *int32  `json:"displayOrder"`
	IsVisible       *bool   `json:"isVisible"`
	LegendType      *string `json:"legendType"`
	GeometryType    *string `json:"geometryType"`
	StrokeDashArray *string `json:"strokeDashArray"`
	IconURL         *string `json:"iconUrl"`
	ImageURL        *string `json:"imageUrl"`
	Description     *string `json:"description"`
	Color           *string `json:"color"`
	LandUseID       *uint64 `json:"landUseId"`
}

type listLegendsAdminReq struct {
	Page           uint32  `form:"page"`
	Size           uint32  `form:"size"`
	LayerID        *uint64 `form:"layerId"`
	LabelID        *uint64 `form:"labelId"`
	LegendType     *string `form:"legendType"`
	LandUseGroupID *uint64 `form:"landUseGroupId"`
}

type listLegendsByLayerAdminReq struct {
	LayerID uint64 `path:"layerId"`
	Page    uint32 `form:"page"`
	Size    uint32 `form:"size"`
}

func (c *QHAdminLegendController) create(ctx context.Context, req createLegendReq) (*httpapi.Result, error) {
	row := &qh.QHLayerLegend{
		LayerID: req.LayerID, LabelID: req.LabelID, Note: req.Note,
		DisplayOrder: int(req.DisplayOrder), IsVisible: req.IsVisible,
		LegendType: req.LegendType, GeometryType: req.GeometryType,
		StrokeDashArray: req.StrokeDashArray, IconURL: req.IconURL, ImageURL: req.ImageURL,
		Description: req.Description, Color: req.Color, LandUseID: req.LandUseID,
	}
	out, err := c.uc.Create(ctx, row)
	if err != nil {
		return nil, err
	}
	return protoResult(toQHLayerLegendResponse(out))
}

func (c *QHAdminLegendController) batchCreate(ctx context.Context, req batchCreateLegendsReq) (*httpapi.Result, error) {
	if req.LayerID == 0 {
		return nil, _errors.ReturnError(400, "layerId is required")
	}
	items := make([]qh.QHLayerLegend, 0, len(req.Items))
	for _, it := range req.Items {
		items = append(items, qh.QHLayerLegend{
			LayerID: req.LayerID, LabelID: it.LabelID, Note: it.Note,
			DisplayOrder: int(it.DisplayOrder), IsVisible: it.IsVisible,
			LegendType: it.LegendType, GeometryType: it.GeometryType,
			StrokeDashArray: it.StrokeDashArray, IconURL: it.IconURL, ImageURL: it.ImageURL,
			Description: it.Description, Color: it.Color, LandUseID: it.LandUseID,
		})
	}
	created, skipped, err := c.uc.BatchCreate(ctx, req.LayerID, items)
	if err != nil {
		return nil, err
	}
	data := make([]*tqdpb.QHLayerLegendResponse, len(created))
	for i := range created {
		data[i] = toQHLayerLegendResponse(&created[i])
	}
	encoded, err := protoMessages(data)
	if err != nil {
		return nil, err
	}
	return jsonOK(map[string]any{"created": encoded, "skippedLabelIds": skipped})
}

func (c *QHAdminLegendController) get(ctx context.Context, req dto.IDPath) (*httpapi.Result, error) {
	if req.ID == 0 {
		return nil, _errors.ReturnError(400, "id is required")
	}
	out, err := c.uc.GetByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	return protoResult(toQHLayerLegendResponse(out))
}

func (c *QHAdminLegendController) update(ctx context.Context, req updateLegendReq) (*httpapi.Result, error) {
	if req.ID == 0 {
		return nil, _errors.ReturnError(400, "id is required")
	}
	in := &services.QHLayerLegendUpdateInput{
		LayerID: req.LayerID, LabelID: req.LabelID, Note: req.Note,
		IsVisible: req.IsVisible, LegendType: req.LegendType, GeometryType: req.GeometryType,
		StrokeDashArray: req.StrokeDashArray, IconURL: req.IconURL, ImageURL: req.ImageURL,
		Description: req.Description, Color: req.Color, LandUseID: req.LandUseID,
	}
	if req.DisplayOrder != nil {
		v := int(*req.DisplayOrder)
		in.DisplayOrder = &v
	}
	out, err := c.uc.Update(ctx, req.ID, in)
	if err != nil {
		return nil, err
	}
	return protoResult(toQHLayerLegendResponse(out))
}

func (c *QHAdminLegendController) delete(ctx context.Context, req dto.IDPath) (*httpapi.Result, error) {
	if req.ID == 0 {
		return nil, _errors.ReturnError(400, "id is required")
	}
	if err := c.uc.Delete(ctx, req.ID); err != nil {
		return nil, err
	}
	return jsonOK(map[string]string{"message": "deleted"})
}

func (c *QHAdminLegendController) list(ctx context.Context, req listLegendsAdminReq) (*httpapi.Result, error) {
	pagable := &_dto.Pagable{Page: req.Page, Size: req.Size}
	rows, total, err := c.uc.List(ctx, pagable, req.LayerID, req.LabelID, req.LegendType, req.LandUseGroupID)
	if err != nil {
		return nil, err
	}
	items := make([]*tqdpb.QHLayerLegendResponse, len(rows))
	for i := range rows {
		items[i] = toQHLayerLegendResponse(&rows[i])
	}
	return protoListResult(items, total, pagable.GetPage(), pagable.GetSize())
}

func (c *QHAdminLegendController) listByLayer(ctx context.Context, req listLegendsByLayerAdminReq) (*httpapi.Result, error) {
	if req.LayerID == 0 {
		return nil, _errors.ReturnError(400, "layerId is required")
	}
	pagable := &_dto.Pagable{Page: req.Page, Size: req.Size}
	rows, total, err := c.uc.ListClientByLayer(ctx, req.LayerID, pagable)
	if err != nil {
		return nil, err
	}
	items := make([]*tqdpb.QHLayerLegendResponse, len(rows))
	for i := range rows {
		items[i] = toQHLayerLegendResponse(&rows[i])
	}
	return protoListResult(items, total, pagable.GetPage(), pagable.GetSize())
}
