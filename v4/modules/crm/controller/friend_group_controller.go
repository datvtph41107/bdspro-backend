package controller

import (
	"context"
	"net/http"

	_errors "common/errors"
	"userv4/internal/httpapi"
	domain "userv4/modules/crm/entity"
	"userv4/modules/crm/services"
)

type FriendGroupController struct{ uc *services.FriendGroupUsecase }

func NewFriendGroupController(uc *services.FriendGroupUsecase) *FriendGroupController {
	return &FriendGroupController{uc: uc}
}

func (c *FriendGroupController) RegisterHTTP(mux *http.ServeMux) {
	httpapi.POST(mux, "/v2/crm/friend-group", c.create)
	httpapi.PUT(mux, "/v2/crm/friend-group/{id}", c.update)
	httpapi.DELETE(mux, "/v2/crm/friend-group", c.delete)
	httpapi.GET(mux, "/v2/crm/friend-group/list", c.list)
}

func (c *FriendGroupController) require() error {
	if c == nil || c.uc == nil {
		return _errors.ReturnError(501, "friend-group service not configured")
	}
	return nil
}

type friendGroupBody struct {
	ID   uint64 `path:"id" json:"id"`
	Name string `json:"name"`
}
type friendGroupDelete struct {
	ID uint64 `form:"id" json:"id"`
}

func (c *FriendGroupController) create(ctx context.Context, req friendGroupBody) (*httpapi.Result, error) {
	if err := c.require(); err != nil {
		return nil, err
	}
	out, err := c.uc.CreateGroup(ctx, domain.GroupEntity{Name: req.Name})
	if err != nil {
		return nil, err
	}
	return httpapi.OK(out), nil
}
func (c *FriendGroupController) update(ctx context.Context, req friendGroupBody) (*httpapi.Result, error) {
	if err := c.require(); err != nil {
		return nil, err
	}
	out, err := c.uc.UpdateGroup(ctx, req.ID, &domain.GroupEntity{Name: req.Name})
	if err != nil {
		return nil, err
	}
	return httpapi.OK(out), nil
}
func (c *FriendGroupController) delete(ctx context.Context, req friendGroupDelete) (*httpapi.Result, error) {
	if err := c.require(); err != nil {
		return nil, err
	}
	if err := c.uc.DeleteGroup(ctx, req.ID); err != nil {
		return nil, err
	}
	return submitResult(req.ID, "deleted")
}
func (c *FriendGroupController) list(ctx context.Context, _ httpapi.Empty) (*httpapi.Result, error) {
	if err := c.require(); err != nil {
		return nil, err
	}
	items, err := c.uc.ListGroups(ctx)
	if err != nil {
		return nil, err
	}
	return httpapi.OK(items), nil
}
