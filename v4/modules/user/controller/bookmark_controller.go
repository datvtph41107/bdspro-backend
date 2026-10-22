package controller

import (
	"context"
	"net/http"

	"common/httpapi"
	"mdn4/modules/user/dto"
	"mdn4/modules/user/services"
	userpb "pb/types/user"
)

// BookmarkController owns bookmark-user HTTP routes under /v2/user.
type BookmarkController struct {
	svc services.BookmarkService
}

func NewBookmarkController(svc services.BookmarkService) *BookmarkController {
	return &BookmarkController{svc: svc}
}

func (c *BookmarkController) RegisterHTTP(mux *http.ServeMux) {
	httpapi.GET(mux, "/v2/user/bookmark-users", c.getBookmarkUsers)
	httpapi.PUT(mux, "/v2/user/bookmark-users", c.updateBookmarkUsers)
}

func (c *BookmarkController) getBookmarkUsers(ctx context.Context, req dto.BookmarkUsersListRequest) (*httpapi.Result, error) {
	result, err := c.svc.GetBookmarkUsers(ctx, req.Page, req.Size)
	if err != nil {
		return nil, err
	}
	return protoMessageResult(result, false)
}

func (c *BookmarkController) updateBookmarkUsers(ctx context.Context, req userpb.UpdateBookmarkUsersRequest) (*httpapi.Result, error) {
	if err := c.svc.UpdateBookmarkUsers(ctx, &req); err != nil {
		return nil, err
	}
	return httpapi.OK(nil), nil
}
