package controller

import (
	"context"
	"net/http"

	"common/httpapi"
	"mdn4/modules/hub/dto"
	"mdn4/modules/hub/services"
	hubpb "pb/types/hub"
)

// VersionController owns app-bundle version HTTP routes under /v2/hub/app-bundle.
type VersionController struct {
	version services.VersionService
}

func NewVersionController(version services.VersionService) *VersionController {
	return &VersionController{version: version}
}

func (c *VersionController) RegisterHTTP(mux *http.ServeMux) {
	httpapi.GET(mux, "/v2/hub/app-bundle/update-check/{appName}", c.checkAppUpdate)
	mux.HandleFunc("POST /v2/hub/app-bundle/version", withHubAPIKey(httpapi.BindJSON(c.createVersionFromDev)))
	// Admin CRUD
	httpapi.GET(mux, "/v2/hub/admin/versions/{id}", c.adminGet)
	httpapi.GET(mux, "/v2/hub/admin/versions", c.adminList)
	mux.HandleFunc("POST /v2/hub/admin/versions", httpapi.BindJSON(c.adminCreate))
	mux.HandleFunc("PUT /v2/hub/admin/versions/{id}", httpapi.BindJSON(c.adminUpdate))
	mux.HandleFunc("DELETE /v2/hub/admin/versions/{id}", httpapi.BindJSON(c.adminDelete))
}

func (c *VersionController) checkAppUpdate(ctx context.Context, req dto.CheckAppUpdateRequest) (*httpapi.Result, error) {
	result, err := c.version.CheckAppUpdate(ctx, req.AppName, req)
	if err != nil {
		return nil, err
	}
	return protoDataResult(result)
}

func (c *VersionController) createVersionFromDev(ctx context.Context, req hubpb.CreateVersionRequest) (*httpapi.Result, error) {
	result, err := c.version.CreateVersionFromDev(ctx, &req)
	if err != nil {
		return nil, err
	}
	return protoDataResult(result)
}

func (c *VersionController) adminList(ctx context.Context, req dto.VersionListQuery) (*httpapi.Result, error) {
	items, total, err := c.version.AdminList(ctx, req.AppName, req.Platform, int(req.Page), int(req.Size))
	if err != nil {
		return nil, err
	}
	data, err := protoMessages(items)
	if err != nil {
		return nil, err
	}
	return httpapi.OK(data, httpapi.WithPagination(&httpapi.Pagination{Total: total})), nil
}

func (c *VersionController) adminGet(ctx context.Context, req httpapi.PathUint64) (*httpapi.Result, error) {
	item, err := c.version.AdminGet(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	return protoDataResult(item)
}

func (c *VersionController) adminCreate(ctx context.Context, req dto.CreateVersionAdminRequest) (*httpapi.Result, error) {
	item, err := c.version.AdminCreate(ctx, req)
	if err != nil {
		return nil, err
	}
	return protoDataResult(item)
}

func (c *VersionController) adminUpdate(ctx context.Context, req dto.UpdateVersionRequest) (*httpapi.Result, error) {
	createReq := dto.CreateVersionAdminRequest{
		AppName: req.AppName, Platform: req.Platform, VersionName: req.VersionName,
		BuildNumber: req.BuildNumber, ForceUpdate: req.ForceUpdate, Active: req.Active,
		DownloadURL: req.DownloadURL, ReleaseNotes: req.ReleaseNotes,
		BundleName: req.BundleName, BundleSize: req.BundleSize,
	}
	item, err := c.version.AdminUpdate(ctx, req.ID, createReq)
	if err != nil {
		return nil, err
	}
	return protoDataResult(item)
}

func (c *VersionController) adminDelete(ctx context.Context, req dto.DeleteVersionPath) (*httpapi.Result, error) {
	if err := c.version.AdminDelete(ctx, req.ID); err != nil {
		return nil, err
	}
	return httpapi.OK(map[string]string{"message": "deleted"}), nil
}

func withHubAPIKey(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("api-key")
		if apiKey == "" {
			apiKey = r.Header.Get("x-api-key")
		}
		next(w, r.WithContext(services.WithAPIKey(r.Context(), apiKey)))
	}
}
