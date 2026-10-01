package controller

import (
	"context"
	"net/http"

	_errors "common/errors"
	"userv4/internal/httpapi"
	"userv4/modules/crm/services"
)

type BlockController struct{ uc *services.BlockUsecase }

func NewBlockController(uc *services.BlockUsecase) *BlockController { return &BlockController{uc: uc} }

func (c *BlockController) RegisterHTTP(mux *http.ServeMux) {
	httpapi.GET(mux, "/v2/crm/block/list", c.list)
	httpapi.POST(mux, "/v2/crm/block/{profileId}", c.block)
	httpapi.DELETE(mux, "/v2/crm/block/{id}", c.unblock)
}

func (c *BlockController) require() error {
	if c == nil || c.uc == nil {
		return _errors.ReturnError(501, "block service not configured")
	}
	return nil
}

type blockProfilePath struct {
	ProfileID uint64 `path:"profileId"`
}

func (c *BlockController) list(ctx context.Context, _ httpapi.Empty) (*httpapi.Result, error) {
	if err := c.require(); err != nil {
		return nil, err
	}
	items, err := c.uc.ListBlockedUsers(ctx)
	if err != nil {
		return nil, err
	}
	return httpapi.OK(items), nil
}
func (c *BlockController) block(ctx context.Context, req blockProfilePath) (*httpapi.Result, error) {
	if err := c.require(); err != nil {
		return nil, err
	}
	out, err := c.uc.BlockUser(ctx, req.ProfileID)
	if err != nil {
		return nil, err
	}
	return httpapi.OK(out), nil
}
func (c *BlockController) unblock(ctx context.Context, req httpapi.PathUint64) (*httpapi.Result, error) {
	if err := c.require(); err != nil {
		return nil, err
	}
	id, err := c.uc.UnblockUser(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	return submitResult(id, "unblocked")
}
