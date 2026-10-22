package controller

import (
	"context"
	"net/http"

	base_enum "base/enum"
	_errors "common/errors"
	"common/httpapi"
	_utils "common/utils"
	"mdn4/modules/crm/dto"
	domain "mdn4/modules/crm/entity"
	"mdn4/modules/crm/services"
)

type RuleController struct{ uc *services.RuleUsecase }

func NewRuleController(uc *services.RuleUsecase) *RuleController { return &RuleController{uc: uc} }

func (c *RuleController) RegisterHTTP(mux *http.ServeMux) {
	httpapi.GET(mux, "/v2/crm/rule/list", c.list)
	httpapi.GET(mux, "/v2/crm/rule/detail/{id}", c.detail)
	httpapi.POST(mux, "/v2/crm/rule", c.create)
	httpapi.PUT(mux, "/v2/crm/rule/{id}", c.update)
	httpapi.DELETE(mux, "/v2/crm/rule/{id}", c.delete)
	httpapi.PUT(mux, "/v2/crm/rule/{id}/active", c.active)
}

func (c *RuleController) require() error {
	if c == nil || c.uc == nil {
		return _errors.ReturnError(501, "rule service not configured")
	}
	return nil
}

type ruleListReq struct {
	Page    uint32 `form:"page"`
	Size    uint32 `form:"size"`
	OwnerID uint64 `form:"ownerId"`
	OwnerOf int32  `form:"ownerOf"`
}
type ruleBody struct {
	ID     uint64 `path:"id" json:"id"`
	Active *bool  `json:"active"`
	domain.RuleEntity
}
type ruleActiveReq struct {
	ID     uint64 `path:"id"`
	Active bool   `json:"active" form:"active"`
}

func (c *RuleController) list(ctx context.Context, req ruleListReq) (*httpapi.Result, error) {
	if err := c.require(); err != nil {
		return nil, err
	}
	ownerID := req.OwnerID
	if ownerID == 0 {
		ownerID = _utils.GetProfileIdWithContext(ctx)
	}
	ownerOf := base_enum.EOwnerOf(req.OwnerOf)
	if ownerOf == 0 {
		ownerOf = base_enum.EOwnerOfMember
	}
	items, total, err := c.uc.Search(ctx, ownerID, ownerOf, dto.RuleSearchDTO{})
	if err != nil {
		return nil, err
	}
	return httpapi.OK(items, httpapi.WithPagination(&httpapi.Pagination{Total: total})), nil
}
func (c *RuleController) detail(ctx context.Context, req httpapi.PathUint64) (*httpapi.Result, error) {
	if err := c.require(); err != nil {
		return nil, err
	}
	out, err := c.uc.GetByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	return httpapi.OK(out), nil
}
func (c *RuleController) create(ctx context.Context, req ruleBody) (*httpapi.Result, error) {
	if err := c.require(); err != nil {
		return nil, err
	}
	out, err := c.uc.Create(ctx, &req.RuleEntity)
	if err != nil {
		return nil, err
	}
	return httpapi.OK(out), nil
}
func (c *RuleController) update(ctx context.Context, req ruleBody) (*httpapi.Result, error) {
	if err := c.require(); err != nil {
		return nil, err
	}
	out, err := c.uc.Update(ctx, req.ID, &req.RuleEntity)
	if err != nil {
		return nil, err
	}
	return httpapi.OK(out), nil
}
func (c *RuleController) delete(ctx context.Context, req httpapi.PathUint64) (*httpapi.Result, error) {
	if err := c.require(); err != nil {
		return nil, err
	}
	if err := c.uc.Delete(ctx, req.ID); err != nil {
		return nil, err
	}
	return submitResult(req.ID, "deleted")
}
func (c *RuleController) active(ctx context.Context, req ruleActiveReq) (*httpapi.Result, error) {
	if err := c.require(); err != nil {
		return nil, err
	}
	if err := c.uc.Active(ctx, req.ID, req.Active); err != nil {
		return nil, err
	}
	return submitResult(req.ID, "updated")
}
