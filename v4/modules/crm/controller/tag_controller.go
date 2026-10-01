package controller

import (
	"context"
	"net/http"

	_errors "common/errors"
	"userv4/internal/httpapi"
	"userv4/modules/crm/dto"
	"userv4/modules/crm/services"
)

type TagController struct{ uc *services.TagUsecase }

func NewTagController(uc *services.TagUsecase) *TagController { return &TagController{uc: uc} }

func (c *TagController) RegisterHTTP(mux *http.ServeMux) {
	httpapi.POST(mux, "/v2/crm/tag/admin", c.adminCreate)
	httpapi.PUT(mux, "/v2/crm/tag/admin/{id}", c.adminUpdate)
	httpapi.DELETE(mux, "/v2/crm/tag/admin/{id}", c.adminDelete)
	httpapi.GET(mux, "/v2/crm/tag/list", c.list)
	httpapi.GET(mux, "/v2/crm/tag/admin/list", c.adminList)
	httpapi.GET(mux, "/v2/crm/tag/contact/{id}", c.byContact)
}

func (c *TagController) require() error {
	if c == nil || c.uc == nil {
		return _errors.ReturnError(501, "tag service not configured")
	}
	return nil
}

type tagAdminBody struct {
	ID        uint64 `path:"id" json:"id"`
	Name      string `json:"name"`
	IsActive  bool   `json:"isActive"`
	IsDefault bool   `json:"isDefault"`
}
type tagListReq struct {
	Page uint32 `form:"page"`
	Size uint32 `form:"size"`
	Text string `form:"text"`
}

func (c *TagController) adminCreate(ctx context.Context, req tagAdminBody) (*httpapi.Result, error) {
	if err := c.require(); err != nil { return nil, err }
	out, err := c.uc.AdminCreateTag(ctx, &dto.AdminTagRequest{Name: req.Name, IsActive: req.IsActive, IsDefault: req.IsDefault})
	if err != nil { return nil, err }
	return httpapi.OK(out), nil
}
func (c *TagController) adminUpdate(ctx context.Context, req tagAdminBody) (*httpapi.Result, error) {
	if err := c.require(); err != nil { return nil, err }
	out, err := c.uc.AdminUpdateTag(ctx, &dto.AdminTagRequest{ID: req.ID, Name: req.Name, IsActive: req.IsActive, IsDefault: req.IsDefault})
	if err != nil { return nil, err }
	return httpapi.OK(out), nil
}
func (c *TagController) adminDelete(ctx context.Context, req httpapi.PathUint64) (*httpapi.Result, error) {
	if err := c.require(); err != nil { return nil, err }
	if err := c.uc.AdminDeleteTag(ctx, req.ID); err != nil { return nil, err }
	return submitResult(req.ID, "deleted")
}
func (c *TagController) list(ctx context.Context, req tagListReq) (*httpapi.Result, error) {
	if err := c.require(); err != nil { return nil, err }
	items, total, err := c.uc.ListTagsForContact(ctx, &dto.TagListRequest{Page: req.Page, Size: req.Size, Text: req.Text})
	if err != nil { return nil, err }
	return httpapi.OK(items, httpapi.WithPagination(&httpapi.Pagination{Total: int64(total)})), nil
}
func (c *TagController) adminList(ctx context.Context, req tagListReq) (*httpapi.Result, error) {
	if err := c.require(); err != nil { return nil, err }
	items, total, err := c.uc.AdminListTags(ctx, &dto.AdminTagListRequest{Page: req.Page, Size: req.Size, Text: req.Text})
	if err != nil { return nil, err }
	return httpapi.OK(items, httpapi.WithPagination(&httpapi.Pagination{Total: int64(total)})), nil
}
func (c *TagController) byContact(ctx context.Context, req httpapi.PathUint64) (*httpapi.Result, error) {
	if err := c.require(); err != nil { return nil, err }
	items, err := c.uc.ListTagsByContactID(ctx, req.ID)
	if err != nil { return nil, err }
	return httpapi.OK(items), nil
}
