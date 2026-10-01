package controller

import (
	"context"
	"net/http"
	"strings"

	_dto "common/domain/dto"
	_errors "common/errors"
	_utils "common/utils"
	tqdpb "pb/types/tqd"
	"userv4/internal/httpapi"
	"userv4/modules/tqd/dto"
	qh "userv4/modules/tqd/entity/qh"
	"userv4/modules/tqd/services"
)

// QHAdminAuthorityController owns /v2/tqd/qh/admin/authority-issuring*.
type QHAdminAuthorityController struct {
	uc services.QHAuthorityIssuringUsecase
}

func NewQHAdminAuthorityController(uc services.QHAuthorityIssuringUsecase) *QHAdminAuthorityController {
	return &QHAdminAuthorityController{uc: uc}
}

func (c *QHAdminAuthorityController) RegisterHTTP(mux *http.ServeMux) {
	httpapi.POST(mux, "/v2/tqd/qh/admin/authority-issuring", c.create)
	httpapi.GET(mux, "/v2/tqd/qh/admin/authority-issuring/list", c.list)
	httpapi.GET(mux, "/v2/tqd/qh/admin/authority-issuring/{id}", c.get)
	httpapi.PUT(mux, "/v2/tqd/qh/admin/authority-issuring/{id}", c.update)
	httpapi.DELETE(mux, "/v2/tqd/qh/admin/authority-issuring/{id}", c.delete)
}

type createAuthorityReq struct {
	Name        string `json:"name"`
	Code        string `json:"code"`
	Description string `json:"description"`
}

type updateAuthorityReq struct {
	ID          uint64  `path:"id" json:"-"`
	Name        *string `json:"name"`
	Code        *string `json:"code"`
	Description *string `json:"description"`
}

func (c *QHAdminAuthorityController) create(ctx context.Context, req createAuthorityReq) (*httpapi.Result, error) {
	out, err := c.uc.Create(ctx, &qh.QHAuthorityIssuring{
		Name:        strings.TrimSpace(req.Name),
		Code:        strings.TrimSpace(req.Code),
		Description: req.Description,
	})
	if err != nil {
		return nil, err
	}
	return protoResult(toAuthorityPB(out))
}

func (c *QHAdminAuthorityController) get(ctx context.Context, req dto.IDPath) (*httpapi.Result, error) {
	if req.ID == 0 {
		return nil, _errors.ReturnError(400, "id is required")
	}
	out, err := c.uc.GetByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	return protoResult(toAuthorityPB(out))
}

func (c *QHAdminAuthorityController) update(ctx context.Context, req updateAuthorityReq) (*httpapi.Result, error) {
	if req.ID == 0 {
		return nil, _errors.ReturnError(400, "id is required")
	}
	in := &services.QHAuthorityIssuringUpdateInput{Name: req.Name, Code: req.Code, Description: req.Description}
	if in.Name == nil && in.Code == nil && in.Description == nil {
		return nil, _errors.ReturnError(400, "at least one field to update is required")
	}
	out, err := c.uc.Update(ctx, req.ID, in)
	if err != nil {
		return nil, err
	}
	return protoResult(toAuthorityPB(out))
}

func (c *QHAdminAuthorityController) delete(ctx context.Context, req dto.IDPath) (*httpapi.Result, error) {
	if req.ID == 0 {
		return nil, _errors.ReturnError(400, "id is required")
	}
	if err := c.uc.Delete(ctx, req.ID); err != nil {
		return nil, err
	}
	return jsonOK(map[string]string{"message": "deleted"})
}

func (c *QHAdminAuthorityController) list(ctx context.Context, req dto.PageQuery) (*httpapi.Result, error) {
	pagable := &_dto.Pagable{Page: req.Page, Size: req.Size}
	rows, total, err := c.uc.List(ctx, pagable)
	if err != nil {
		return nil, err
	}
	items := make([]*tqdpb.QHAuthorityIssuringResponse, len(rows))
	for i := range rows {
		items[i] = toAuthorityPB(&rows[i])
	}
	return protoListResult(items, total, pagable.GetPage(), pagable.GetSize())
}

func toAuthorityPB(e *qh.QHAuthorityIssuring) *tqdpb.QHAuthorityIssuringResponse {
	if e == nil {
		return nil
	}
	return &tqdpb.QHAuthorityIssuringResponse{
		Id: e.ID, Name: e.Name, Code: e.Code, Description: e.Description,
		CreatedAt: _utils.FormatTimeToString(e.CreatedAt),
		UpdatedAt: _utils.FormatTimeToString(e.UpdatedAt),
	}
}
