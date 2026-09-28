package controller

import (
	"context"
	"net/http"

	"userv4/internal/httpapi"
	"userv4/modules/bdspro/services"
)

// SyncController owns GET /v2/bdspro/ts.
type SyncController struct {
	sync services.SyncService
}

func NewSyncController(sync services.SyncService) *SyncController {
	return &SyncController{sync: sync}
}

func (c *SyncController) RegisterHTTP(mux *http.ServeMux) {
	httpapi.GET(mux, "/v2/bdspro/ts", c.getAppTimestamps)
}

func (c *SyncController) getAppTimestamps(ctx context.Context, _ httpapi.Empty) (*httpapi.Result, error) {
	result, err := c.sync.GetAppTimestamps(ctx)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return httpapi.OK(map[string]int64{}), nil
	}
	return httpapi.OK(result.Timestamps), nil
}
