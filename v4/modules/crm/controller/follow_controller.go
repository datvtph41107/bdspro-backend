package controller

import (
	"context"
	"net/http"

	_errors "common/errors"
	"userv4/internal/httpapi"
	"userv4/modules/crm/services"
)

type FollowController struct{ uc *services.FollowUsecase }

func NewFollowController(uc *services.FollowUsecase) *FollowController { return &FollowController{uc: uc} }

func (c *FollowController) RegisterHTTP(mux *http.ServeMux) {
	httpapi.GET(mux, "/v2/crm/follow/followers", c.followers)
	httpapi.GET(mux, "/v2/crm/follow/following", c.following)
	httpapi.POST(mux, "/v2/crm/follow/{id}", c.follow)
	httpapi.DELETE(mux, "/v2/crm/follow/{id}", c.unfollow)
}

func (c *FollowController) require() error {
	if c == nil || c.uc == nil {
		return _errors.ReturnError(501, "follow service not configured")
	}
	return nil
}

func (c *FollowController) followers(ctx context.Context, _ httpapi.Empty) (*httpapi.Result, error) {
	if err := c.require(); err != nil {
		return nil, err
	}
	items, err := c.uc.FollowerUser(ctx)
	if err != nil {
		return nil, err
	}
	return httpapi.OK(items), nil
}
func (c *FollowController) following(ctx context.Context, _ httpapi.Empty) (*httpapi.Result, error) {
	if err := c.require(); err != nil {
		return nil, err
	}
	items, err := c.uc.FollowingUser(ctx)
	if err != nil {
		return nil, err
	}
	return httpapi.OK(items), nil
}
func (c *FollowController) follow(ctx context.Context, req httpapi.PathUint64) (*httpapi.Result, error) {
	if err := c.require(); err != nil {
		return nil, err
	}
	out, err := c.uc.FollowUser(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	return httpapi.OK(out), nil
}
func (c *FollowController) unfollow(ctx context.Context, req httpapi.PathUint64) (*httpapi.Result, error) {
	if err := c.require(); err != nil {
		return nil, err
	}
	if err := c.uc.UnfollowUser(ctx, req.ID); err != nil {
		return nil, err
	}
	return submitResult(req.ID, "unfollowed")
}
