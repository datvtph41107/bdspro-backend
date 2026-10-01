package controller

import (
	"context"
	"net/http"
	"strings"

	_dto "common/domain/dto"
	_errors "common/errors"
	"userv4/internal/httpapi"
	"userv4/modules/crm/data"
	"userv4/modules/crm/enums"
	"userv4/modules/crm/services"
)

type FriendController struct{ uc *services.FriendUsecase }

func NewFriendController(uc *services.FriendUsecase) *FriendController { return &FriendController{uc: uc} }

func (c *FriendController) RegisterHTTP(mux *http.ServeMux) {
	httpapi.GET(mux, "/v2/crm/friend/list", c.list)
	httpapi.POST(mux, "/v2/crm/friend/request/{receiverId}", c.request)
	httpapi.DELETE(mux, "/v2/crm/friend/cancel/{id}", c.cancel)
	httpapi.PUT(mux, "/v2/crm/friend/accept/{id}", c.accept)
	httpapi.DELETE(mux, "/v2/crm/friend/reject/{id}", c.reject)
	httpapi.PUT(mux, "/v2/crm/friend/group/{requestId}/{groupId}", c.setGroup)
	httpapi.GET(mux, "/v2/crm/friend/request/sent", c.sent)
	httpapi.GET(mux, "/v2/crm/friend/request/received", c.received)
	httpapi.GET(mux, "/v2/crm/friend/info", c.info)
	httpapi.GET(mux, "/v2/crm/friend/with-lead", c.withLead)
	httpapi.GET(mux, "/v2/crm/friend/with-member", c.withMember)
}

func (c *FriendController) require() error {
	if c == nil || c.uc == nil {
		return _errors.ReturnError(501, "friend service not configured")
	}
	return nil
}

type friendListReq struct {
	Page uint32 `form:"page"`
	Size uint32 `form:"size"`
	Text string `form:"text"`
}
type friendReceiverPath struct {
	ReceiverID uint64 `path:"receiverId"`
}
type friendGroupPath struct {
	RequestID uint64 `path:"requestId"`
	GroupID   uint64 `path:"groupId"`
}

func (c *FriendController) list(ctx context.Context, req friendListReq) (*httpapi.Result, error) {
	if err := c.require(); err != nil {
		return nil, err
	}
	items, total, err := c.uc.FriendList(ctx, data.FriendRequest{
		Text: strings.TrimSpace(strings.ToLower(req.Text)),
		Pagable: _dto.Pagable{Page: req.Page, Size: req.Size},
	})
	if err != nil {
		return nil, err
	}
	return httpapi.OK(items, httpapi.WithPagination(&httpapi.Pagination{Total: total})), nil
}
func (c *FriendController) request(ctx context.Context, req friendReceiverPath) (*httpapi.Result, error) {
	if err := c.require(); err != nil {
		return nil, err
	}
	entity, err := c.uc.Request(ctx, req.ReceiverID)
	if err != nil {
		return nil, err
	}
	id := uint64(0)
	if entity != nil {
		id = entity.ID
	}
	return submitResult(id, "ok")
}
func (c *FriendController) cancel(ctx context.Context, req httpapi.PathUint64) (*httpapi.Result, error) {
	if err := c.require(); err != nil {
		return nil, err
	}
	if err := c.uc.CancelRequest(ctx, req.ID); err != nil {
		return nil, err
	}
	return submitResult(req.ID, "cancelled")
}
func (c *FriendController) accept(ctx context.Context, req httpapi.PathUint64) (*httpapi.Result, error) {
	if err := c.require(); err != nil {
		return nil, err
	}
	_, err := c.uc.ChangeStatus(ctx, req.ID, enums.FriendStatusAccepted)
	if err != nil {
		return nil, err
	}
	return submitResult(req.ID, "accepted")
}
func (c *FriendController) reject(ctx context.Context, req httpapi.PathUint64) (*httpapi.Result, error) {
	if err := c.require(); err != nil {
		return nil, err
	}
	_, err := c.uc.ChangeStatus(ctx, req.ID, enums.FriendStatusReject)
	if err != nil {
		return nil, err
	}
	return submitResult(req.ID, "rejected")
}
func (c *FriendController) setGroup(ctx context.Context, req friendGroupPath) (*httpapi.Result, error) {
	if err := c.require(); err != nil {
		return nil, err
	}
	_, err := c.uc.ChangeGroup(ctx, req.RequestID, req.GroupID)
	if err != nil {
		return nil, err
	}
	return submitResult(req.RequestID, "updated")
}
func (c *FriendController) sent(ctx context.Context, req friendListReq) (*httpapi.Result, error) {
	if err := c.require(); err != nil {
		return nil, err
	}
	text := req.Text
	items, err := c.uc.RequestSent(ctx, data.FriendDTO{Text: &text})
	if err != nil {
		return nil, err
	}
	return httpapi.OK(items), nil
}
func (c *FriendController) received(ctx context.Context, req friendListReq) (*httpapi.Result, error) {
	if err := c.require(); err != nil {
		return nil, err
	}
	text := req.Text
	items, err := c.uc.RequestReceived(ctx, data.FriendDTO{Text: &text})
	if err != nil {
		return nil, err
	}
	return httpapi.OK(items), nil
}
func (c *FriendController) info(ctx context.Context, _ httpapi.Empty) (*httpapi.Result, error) {
	if err := c.require(); err != nil {
		return nil, err
	}
	info, err := c.uc.FriendInfo(ctx)
	if err != nil {
		return nil, err
	}
	return httpapi.OK(info), nil
}
func (c *FriendController) withLead(ctx context.Context, req friendListReq) (*httpapi.Result, error) {
	return c.list(ctx, req)
}
func (c *FriendController) withMember(ctx context.Context, req friendListReq) (*httpapi.Result, error) {
	return c.list(ctx, req)
}
