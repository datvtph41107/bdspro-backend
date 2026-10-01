package controller

import (
	"context"
	"net/http"

	"userv4/internal/httpapi"
	"userv4/modules/hub/dto"
	"userv4/modules/hub/services"
)

// SystemConfigController handles system config HTTP routes.
type SystemConfigController struct {
	svc services.SystemConfigService
}

func NewSystemConfigController(svc services.SystemConfigService) *SystemConfigController {
	return &SystemConfigController{svc: svc}
}

func (c *SystemConfigController) RegisterHTTP(mux *http.ServeMux) {
	// Public endpoints
	httpapi.GET(mux, "/v2/hub/system-config/by-group/{groupKey}", c.getByGroup)
	httpapi.GET(mux, "/v2/hub/system-config/key/{key}", c.getByKey)
	httpapi.GET(mux, "/v2/hub/settings/{groupKey}", c.getUserSettings)

	// Admin endpoints
	httpapi.POST(mux, "/v2/hub/system-config/bulk-upsert", c.bulkUpsert)
	httpapi.POST(mux, "/v2/hub/system-config/reset-default/{groupKey}", c.resetDefault)
}

func (c *SystemConfigController) getByGroup(ctx context.Context, req dto.GetSystemConfigByGroupQuery) (*httpapi.Result, error) {
	items, err := c.svc.GetByGroup(ctx, req.GroupKey)
	if err != nil {
		return nil, err
	}
	return httpapi.OK(items), nil
}

func (c *SystemConfigController) getByKey(ctx context.Context, req dto.GetSystemConfigByKeyQuery) (*httpapi.Result, error) {
	item, err := c.svc.GetByKey(ctx, req.Key)
	if err != nil {
		return nil, err
	}
	return httpapi.OK(item), nil
}

func (c *SystemConfigController) getUserSettings(ctx context.Context, req dto.GetUserSettingsByGroupQuery) (*httpapi.Result, error) {
	settings, err := c.svc.GetUserSettingsByGroup(ctx, req.GroupKey)
	if err != nil {
		return nil, err
	}
	return httpapi.OK(settings), nil
}

func (c *SystemConfigController) bulkUpsert(ctx context.Context, req dto.BulkUpsertSystemConfigRequest) (*httpapi.Result, error) {
	results, created, updated, err := c.svc.BulkUpsert(ctx, req)
	if err != nil {
		return nil, err
	}
	return httpapi.OK(map[string]interface{}{
		"data":         results,
		"createdCount": created,
		"updatedCount": updated,
	}), nil
}

func (c *SystemConfigController) resetDefault(ctx context.Context, req dto.ResetDefaultByGroupQuery) (*httpapi.Result, error) {
	items, count, err := c.svc.ResetDefaultByGroup(ctx, req.GroupKey)
	if err != nil {
		return nil, err
	}
	return httpapi.OK(map[string]interface{}{
		"data":       items,
		"resetCount": count,
	}), nil
}
