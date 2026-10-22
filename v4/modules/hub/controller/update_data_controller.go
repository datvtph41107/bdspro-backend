package controller

import (
	"context"
	"net/http"

	_enums "common/domain/enum"
	"common/httpapi"
	"mdn4/modules/hub/dto"
	"mdn4/modules/hub/services"
)

// UpdateDataController handles /v2/hub/update/{resource}/... sync endpoints.
type UpdateDataController struct {
	svc services.UpdateDataService
}

func NewUpdateDataController(svc services.UpdateDataService) *UpdateDataController {
	return &UpdateDataController{svc: svc}
}

func (c *UpdateDataController) RegisterHTTP(mux *http.ServeMux) {
	httpapi.GET(mux, "/v2/hub/update/{resource}/sync/status/{id}", c.checkVersionSyncByID)
	httpapi.GET(mux, "/v2/hub/update/{resource}/sync/status", c.checkVersionSync)
	httpapi.GET(mux, "/v2/hub/update/{resource}/refresh/ids", c.refreshSyncIDs)
	httpapi.POST(mux, "/v2/hub/update/{resource}/flush/ids", c.flushSyncIDs)
}

func (c *UpdateDataController) checkVersionSync(ctx context.Context, req dto.CheckVersionSyncQuery) (*httpapi.Result, error) {
	ownerID := req.ID
	if ownerID == 0 {
		if pid, ok := ctx.Value(_enums.ProfileIDKey).(uint64); ok {
			ownerID = pid
		}
	}
	resp, err := c.svc.CheckVersionSync(ctx, ownerID, req.Resource, req.LastSync, req.Limit)
	if err != nil {
		return nil, err
	}
	return httpapi.OK(resp), nil
}

func (c *UpdateDataController) checkVersionSyncByID(ctx context.Context, req dto.CheckVersionSyncByIDQuery) (*httpapi.Result, error) {
	resp, err := c.svc.CheckVersionSync(ctx, req.ID, req.Resource, req.LastSync, 1)
	if err != nil {
		return nil, err
	}
	return httpapi.OK(resp), nil
}

func (c *UpdateDataController) refreshSyncIDs(ctx context.Context, req dto.RefreshSyncIdsQuery) (*httpapi.Result, error) {
	ownerID := uint64(0)
	if pid, ok := ctx.Value(_enums.ProfileIDKey).(uint64); ok {
		ownerID = pid
	}
	resp, err := c.svc.RefreshSyncIDs(ctx, ownerID, req.Resource, req.LastSync, req.Limit)
	if err != nil {
		return nil, err
	}
	return httpapi.OK(resp), nil
}

func (c *UpdateDataController) flushSyncIDs(ctx context.Context, req dto.FlushSyncIdsQuery) (*httpapi.Result, error) {
	ownerID := uint64(0)
	if pid, ok := ctx.Value(_enums.ProfileIDKey).(uint64); ok {
		ownerID = pid
	}
	if err := c.svc.FlushSyncIDs(ctx, ownerID, req.Resource, req.Limit); err != nil {
		return nil, err
	}
	return httpapi.OK(map[string]string{"message": "ok"}), nil
}
