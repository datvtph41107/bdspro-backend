package controller

import (
	"context"
	"net/http"

	"common/httpapi"
	"mdn4/modules/hub/dto"
	"mdn4/modules/hub/services"
)

// EventQueueController handles event queue HTTP routes.
type EventQueueController struct {
	svc services.EventQueueService
}

func NewEventQueueController(svc services.EventQueueService) *EventQueueController {
	return &EventQueueController{svc: svc}
}

func (c *EventQueueController) RegisterHTTP(mux *http.ServeMux) {
	httpapi.GET(mux, "/v2/hub/event-queue/list", c.list)
	httpapi.GET(mux, "/v2/hub/event-queue/detail/{id}", c.detail)
	mux.HandleFunc("POST /v2/hub/event-queue", httpapi.BindJSON(c.create))
	mux.HandleFunc("PUT /v2/hub/event-queue/edit/{id}", httpapi.BindJSON(c.update))
	mux.HandleFunc("DELETE /v2/hub/event-queue/{id}", httpapi.BindJSON(c.delete))
	mux.HandleFunc("PUT /v2/hub/event-queue/status/{id}", httpapi.BindJSON(c.updateStatus))
}

func (c *EventQueueController) list(ctx context.Context, req dto.EventQueueListQuery) (*httpapi.Result, error) {
	page := int(req.Page)
	size := int(req.Size)
	if size <= 0 {
		size = 20
	}
	items, total, err := c.svc.List(ctx, req.Text, page, size)
	if err != nil {
		return nil, err
	}
	return httpapi.OK(items, httpapi.WithPagination(&httpapi.Pagination{Total: total})), nil
}

func (c *EventQueueController) detail(ctx context.Context, req dto.EventQueueIDPath) (*httpapi.Result, error) {
	item, err := c.svc.GetByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	return httpapi.OK(item), nil
}

func (c *EventQueueController) create(ctx context.Context, req dto.EventQueueSaveRequest) (*httpapi.Result, error) {
	item, err := c.svc.Create(ctx, req)
	if err != nil {
		return nil, err
	}
	return httpapi.OK(item), nil
}

func (c *EventQueueController) update(ctx context.Context, req dto.EventQueueSaveRequest) (*httpapi.Result, error) {
	id := uint64(0)
	if req.ID != nil {
		id = *req.ID
	}
	item, err := c.svc.Update(ctx, id, req)
	if err != nil {
		return nil, err
	}
	return httpapi.OK(item), nil
}

func (c *EventQueueController) delete(ctx context.Context, req dto.EventQueueIDPath) (*httpapi.Result, error) {
	if err := c.svc.Delete(ctx, req.ID); err != nil {
		return nil, err
	}
	return httpapi.OK(map[string]string{"message": "deleted"}), nil
}

func (c *EventQueueController) updateStatus(ctx context.Context, req dto.EventQueueStatusRequest) (*httpapi.Result, error) {
	item, err := c.svc.UpdateStatus(ctx, req)
	if err != nil {
		return nil, err
	}
	return httpapi.OK(item), nil
}
