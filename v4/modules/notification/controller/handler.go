package controller

import (
	"context"
	"net/http"

	"common/httpapi"
	_utils "common/utils"
	"mdn4/modules/notification/dto"
	"mdn4/modules/notification/services"
)

// Handler is the modules/notification HTTP entry for the user inbox.
type Handler struct {
	notification services.NotificationService
	sync         *_utils.SyncUtil
}

func NewHandler(notification services.NotificationService, sync *_utils.SyncUtil) *Handler {
	return &Handler{notification: notification, sync: sync}
}

// RegisterHTTP mounts phase-1 notification routes with gateway path parity.
func (h *Handler) RegisterHTTP(mux *http.ServeMux) {
	httpapi.GET(mux, "/v2/notification/count", h.count)
	httpapi.GET(mux, "/v2/notification/list", h.list)
	httpapi.PUT(mux, "/v2/notification/read/{id}", h.read)
	httpapi.PUT(mux, "/v2/notification/read-many", h.readMany)
	httpapi.PUT(mux, "/v2/notification/read-all", h.readAll)
	httpapi.DELETE(mux, "/v2/notification/remove/{id}", h.remove)
	httpapi.GET(mux, "/v2/admin-history/search", h.searchAdminHistory)
}

func (h *Handler) count(ctx context.Context, req dto.CountRequest) (*httpapi.Result, error) {
	result, err := h.notification.Count(ctx, req)
	if err != nil {
		return nil, err
	}

	return countResult(result)
}

func (h *Handler) list(ctx context.Context, req dto.ListRequest) (*httpapi.Result, error) {
	result, err := h.notification.List(ctx, req)
	if err != nil {
		return nil, err
	}
	return listResult(result)
}

func (h *Handler) read(ctx context.Context, req httpapi.PathUint64) (*httpapi.Result, error) {
	result, err := h.notification.Read(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	return protoDataResult(result)
}

func (h *Handler) readMany(ctx context.Context, req dto.ReadManyRequest) (*httpapi.Result, error) {
	result, err := h.notification.ReadMany(ctx, req.IDs)
	if err != nil {
		return nil, err
	}
	return protoDataResult(result)
}

func (h *Handler) readAll(ctx context.Context, _ httpapi.Empty) (*httpapi.Result, error) {
	result, err := h.notification.ReadAll(ctx)
	if err != nil {
		return nil, err
	}
	return protoDataResult(result)
}

func (h *Handler) remove(ctx context.Context, req httpapi.PathUint64) (*httpapi.Result, error) {
	result, err := h.notification.Remove(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	return protoDataResult(result)
}

func (h *Handler) searchAdminHistory(ctx context.Context, req dto.AdminHistorySearchRequest) (*httpapi.Result, error) {
	items, total, err := h.notification.SearchAdminHistory(ctx, req)
	if err != nil {
		return nil, err
	}
	data, err := protoMessages(items)
	if err != nil {
		return nil, err
	}
	if data == nil {
		data = []any{}
	}
	return httpapi.OK(data, httpapi.WithPagination(&httpapi.Pagination{
		Page:  req.Page,
		Size:  req.Size,
		Total: total,
	})), nil
}
