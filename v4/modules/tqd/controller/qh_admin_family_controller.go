package controller

import (
	"context"
	"net/http"
	"strings"

	_dto "common/domain/dto"
	_errors "common/errors"
	tqdpb "pb/types/tqd"
	"userv4/internal/httpapi"
	"userv4/modules/tqd/dto"
	qh "userv4/modules/tqd/entity/qh"
	"userv4/modules/tqd/mapper"
	"userv4/modules/tqd/services"
)

// QHAdminFamilyController owns /v2/tqd/qh/admin/layer-families*.
type QHAdminFamilyController struct {
	uc services.QHLayerFamilyUsecase
}

func NewQHAdminFamilyController(uc services.QHLayerFamilyUsecase) *QHAdminFamilyController {
	return &QHAdminFamilyController{uc: uc}
}

func (c *QHAdminFamilyController) RegisterHTTP(mux *http.ServeMux) {
	httpapi.POST(mux, "/v2/tqd/qh/admin/layer-families", c.create)
	httpapi.GET(mux, "/v2/tqd/qh/admin/layer-families/list", c.list)
	httpapi.GET(mux, "/v2/tqd/qh/admin/layer-families/{id}", c.get)
	httpapi.PUT(mux, "/v2/tqd/qh/admin/layer-families/{id}", c.update)
	httpapi.DELETE(mux, "/v2/tqd/qh/admin/layer-families/{id}", c.delete)
}

type createFamilyReq struct {
	Name       string `json:"name"`
	SortNumber int32  `json:"sortNumber"`
}

type updateFamilyReq struct {
	ID         uint64  `path:"id" json:"-"`
	Name       *string `json:"name"`
	SortNumber *int32  `json:"sortNumber"`
}

type listFamilyReq struct {
	Page   uint32  `form:"page"`
	Size   uint32  `form:"size"`
	Search *string `form:"search"`
}

func (c *QHAdminFamilyController) create(ctx context.Context, req createFamilyReq) (*httpapi.Result, error) {
	out, err := c.uc.Create(ctx, &qh.QHLayerFamily{Name: strings.TrimSpace(req.Name), SortNumber: req.SortNumber})
	if err != nil {
		return nil, err
	}
	return protoResult(mapper.ToProtoQHLayerFamily(out))
}

func (c *QHAdminFamilyController) get(ctx context.Context, req dto.IDPath) (*httpapi.Result, error) {
	if req.ID == 0 {
		return nil, _errors.ReturnError(400, "id is required")
	}
	out, err := c.uc.GetByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	return protoResult(mapper.ToProtoQHLayerFamily(out))
}

func (c *QHAdminFamilyController) update(ctx context.Context, req updateFamilyReq) (*httpapi.Result, error) {
	if req.ID == 0 {
		return nil, _errors.ReturnError(400, "id is required")
	}
	out, err := c.uc.Update(ctx, req.ID, &services.QHLayerFamilyUpdateInput{Name: req.Name, SortNumber: req.SortNumber})
	if err != nil {
		return nil, err
	}
	return protoResult(mapper.ToProtoQHLayerFamily(out))
}

func (c *QHAdminFamilyController) delete(ctx context.Context, req dto.IDPath) (*httpapi.Result, error) {
	if req.ID == 0 {
		return nil, _errors.ReturnError(400, "id is required")
	}
	if err := c.uc.Delete(ctx, req.ID); err != nil {
		return nil, err
	}
	return jsonOK(map[string]string{"message": "deleted"})
}

func (c *QHAdminFamilyController) list(ctx context.Context, req listFamilyReq) (*httpapi.Result, error) {
	pagable := &_dto.Pagable{Page: req.Page, Size: req.Size}
	search := ""
	if req.Search != nil {
		search = *req.Search
	}
	rows, total, err := c.uc.List(ctx, pagable, search)
	if err != nil {
		return nil, err
	}
	items := make([]*tqdpb.QHLayerFamilyResponse, len(rows))
	for i := range rows {
		items[i] = mapper.ToProtoQHLayerFamily(&rows[i])
	}
	return protoListResult(items, total, pagable.GetPage(), pagable.GetSize())
}
