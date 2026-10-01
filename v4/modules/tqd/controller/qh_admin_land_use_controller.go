package controller

import (
	"context"
	"net/http"

	_dto "common/domain/dto"
	_errors "common/errors"
	tqdpb "pb/types/tqd"
	"userv4/internal/httpapi"
	"userv4/modules/tqd/dto"
	qh "userv4/modules/tqd/entity/qh"
	"userv4/modules/tqd/services"
)

// QHAdminLandUseController owns /v2/tqd/qh/admin/land-use*.
type QHAdminLandUseController struct {
	uc services.QHLayerLandUseGroupUsecase
}

func NewQHAdminLandUseController(uc services.QHLayerLandUseGroupUsecase) *QHAdminLandUseController {
	return &QHAdminLandUseController{uc: uc}
}

func (c *QHAdminLandUseController) RegisterHTTP(mux *http.ServeMux) {
	httpapi.POST(mux, "/v2/tqd/qh/admin/land-use", c.create)
	httpapi.GET(mux, "/v2/tqd/qh/admin/land-use/list", c.list)
	httpapi.GET(mux, "/v2/tqd/qh/admin/land-use/{id}", c.get)
	httpapi.PUT(mux, "/v2/tqd/qh/admin/land-use/{id}", c.update)
	httpapi.DELETE(mux, "/v2/tqd/qh/admin/land-use/{id}", c.delete)
}

type createLandUseAdminReq struct {
	Note         string `json:"note"`
	DisplayOrder int32  `json:"displayOrder"`
	IsVisible    bool   `json:"isVisible"`
	Color        string `json:"color"`
	Code         string `json:"code"`
	Name         string `json:"name"`
	LandUseID    uint64 `json:"landUseId"`
	LayerID      uint64 `json:"layerId"`
}

type updateLandUseAdminReq struct {
	ID    uint64 `path:"id" json:"-"`
	Note  string `json:"note"`
	Color string `json:"color"`
	Code  string `json:"code"`
	Name  string `json:"name"`
}

type listLandUseAdminReq struct {
	Page    uint32  `form:"page"`
	Size    uint32  `form:"size"`
	LayerID *uint64 `form:"layerId"`
}

func (c *QHAdminLandUseController) create(ctx context.Context, req createLandUseAdminReq) (*httpapi.Result, error) {
	row := &qh.QHLandUse{
		Note: req.Note, DisplayOrder: int(req.DisplayOrder), IsVisible: req.IsVisible,
		Color: req.Color, Code: req.Code, Name: req.Name,
	}
	if req.LandUseID > 0 {
		row.ID = req.LandUseID
	}
	if req.LayerID > 0 {
		row.Layers = []*qh.QHLayer{{ID: req.LayerID}}
	}
	out, err := c.uc.Create(ctx, row)
	if err != nil {
		return nil, err
	}
	return protoResult(toQHLandUseResponse(out))
}

func (c *QHAdminLandUseController) get(ctx context.Context, req dto.IDPath) (*httpapi.Result, error) {
	if req.ID == 0 {
		return nil, _errors.ReturnError(400, "id is required")
	}
	out, err := c.uc.GetByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	return protoResult(toQHLandUseResponse(out))
}

func (c *QHAdminLandUseController) update(ctx context.Context, req updateLandUseAdminReq) (*httpapi.Result, error) {
	if req.ID == 0 {
		return nil, _errors.ReturnError(400, "id is required")
	}
	out, err := c.uc.Update(ctx, req.ID, &qh.QHLandUse{Note: req.Note, Color: req.Color, Code: req.Code, Name: req.Name})
	if err != nil {
		return nil, err
	}
	return protoResult(toQHLandUseResponse(out))
}

func (c *QHAdminLandUseController) delete(ctx context.Context, req dto.IDPath) (*httpapi.Result, error) {
	if req.ID == 0 {
		return nil, _errors.ReturnError(400, "id is required")
	}
	if err := c.uc.Delete(ctx, req.ID); err != nil {
		return nil, err
	}
	return jsonOK(map[string]string{"message": "deleted"})
}

func (c *QHAdminLandUseController) list(ctx context.Context, req listLandUseAdminReq) (*httpapi.Result, error) {
	pagable := &_dto.Pagable{Page: req.Page, Size: req.Size}
	rows, total, err := c.uc.List(ctx, pagable, req.LayerID, nil)
	if err != nil {
		return nil, err
	}
	items := make([]*tqdpb.QHLandUseResponse, len(rows))
	for i := range rows {
		items[i] = toQHLandUseResponse(&rows[i])
	}
	return protoListResult(items, total, pagable.GetPage(), pagable.GetSize())
}
