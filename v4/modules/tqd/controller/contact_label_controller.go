package controller

import (
	"context"
	"net/http"

	_dto "common/domain/dto"
	_errors "common/errors"
	"userv4/internal/httpapi"
	"userv4/modules/tqd/dto"
	"userv4/modules/tqd/services"
)

// ContactLabelController owns /v2/tqd/contact/labels*.
type ContactLabelController struct {
	uc services.ContactLabelUsecase
}

func NewContactLabelController(uc services.ContactLabelUsecase) *ContactLabelController {
	return &ContactLabelController{uc: uc}
}

func (c *ContactLabelController) RegisterHTTP(mux *http.ServeMux) {
	httpapi.POST(mux, "/v2/tqd/contact/labels", c.create)
	httpapi.GET(mux, "/v2/tqd/contact/labels", c.list)
	httpapi.GET(mux, "/v2/tqd/contact/labels/{id}", c.get)
	httpapi.PUT(mux, "/v2/tqd/contact/labels/{id}", c.update)
	httpapi.DELETE(mux, "/v2/tqd/contact/labels/{id}", c.delete)
}

type contactLabelBody struct {
	ID          uint64 `path:"id" json:"id"`
	Name        string `json:"name"`
	Code        string `json:"code"`
	Description string `json:"description"`
	Color       string `json:"color"`
	Icon        string `json:"icon"`
	Type        int    `json:"type"`
	IsActive    bool   `json:"isActive"`
	SortOrder   int    `json:"sortOrder"`
}

type listContactLabelReq struct {
	Page   uint32 `form:"page"`
	Size   uint32 `form:"size"`
	Search string `form:"search"`
}

func (c *ContactLabelController) create(ctx context.Context, req contactLabelBody) (*httpapi.Result, error) {
	out, err := c.uc.Create(ctx, &dto.CreateContactLabelRequestDTO{
		Name: req.Name, Code: req.Code, Description: req.Description,
		Color: req.Color, Icon: req.Icon, Type: req.Type, IsActive: req.IsActive, SortOrder: req.SortOrder,
	})
	if err != nil {
		return nil, err
	}
	return jsonOK(out)
}

func (c *ContactLabelController) get(ctx context.Context, req dto.IDPath) (*httpapi.Result, error) {
	if req.ID == 0 {
		return nil, _errors.ReturnError(400, "id is required")
	}
	out, err := c.uc.GetByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	return jsonOK(out)
}

func (c *ContactLabelController) update(ctx context.Context, req contactLabelBody) (*httpapi.Result, error) {
	if req.ID == 0 {
		return nil, _errors.ReturnError(400, "id is required")
	}
	out, err := c.uc.Update(ctx, req.ID, &dto.UpdateContactLabelRequestDTO{
		Name: req.Name, Code: req.Code, Description: req.Description,
		Color: req.Color, Icon: req.Icon, Type: req.Type, IsActive: req.IsActive, SortOrder: req.SortOrder,
	})
	if err != nil {
		return nil, err
	}
	return jsonOK(out)
}

func (c *ContactLabelController) delete(ctx context.Context, req dto.IDPath) (*httpapi.Result, error) {
	if req.ID == 0 {
		return nil, _errors.ReturnError(400, "id is required")
	}
	if err := c.uc.Delete(ctx, req.ID); err != nil {
		return nil, err
	}
	return jsonOK(map[string]any{"success": true, "message": "deleted"})
}

func (c *ContactLabelController) list(ctx context.Context, req listContactLabelReq) (*httpapi.Result, error) {
	pagable := _dto.Pagable{Page: req.Page, Size: req.Size}
	pagable.Normalize()
	out, err := c.uc.List(ctx, &dto.ListContactLabelsRequestDTO{
		Pagable: pagable, Search: req.Search,
	})
	if err != nil {
		return nil, err
	}
	if out == nil {
		return jsonListResult([]dto.ContactLabelDTO{}, 0, pagable.GetPage(), pagable.GetSize())
	}
	return jsonListResult(out.Data, out.Total, pagable.GetPage(), pagable.GetSize())
}
