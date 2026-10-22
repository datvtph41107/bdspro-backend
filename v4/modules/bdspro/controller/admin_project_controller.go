package controller

import (
	"context"
	"net/http"

	"common/httpapi"
	"mdn4/modules/bdspro/dto"
	"mdn4/modules/bdspro/services"
	bdspropb "pb/types/bdspro"
	sharepb "pb/types/shared"
)

// AdminProjectController owns /v2/bdspro/admin/project/* (migrated from bdspro-service).
type AdminProjectController struct {
	admin services.AdminProjectService
}

func NewAdminProjectController(admin services.AdminProjectService) *AdminProjectController {
	return &AdminProjectController{admin: admin}
}

func (c *AdminProjectController) RegisterHTTP(mux *http.ServeMux) {
	httpapi.HandleProto(mux, http.MethodGet, "/v2/bdspro/admin/project/list", c.listProto)
	httpapi.HandleProto(mux, http.MethodPost, "/v2/bdspro/admin/project", c.createProto)
	httpapi.HandleProto(mux, http.MethodPut, "/v2/bdspro/admin/project/{id}", c.updateProto)
	httpapi.HandleProto(mux, http.MethodDelete, "/v2/bdspro/admin/project/{id}", c.deleteProto)
}

func (c *AdminProjectController) listProto(ctx context.Context, req *bdspropb.SearchQueryRequest) (*bdspropb.ProjectListResponse, error) {
	return c.admin.List(ctx, dto.CatalogSearchRequest{
		Page: req.GetPage(),
		Size: req.GetSize(),
		Text: req.GetText(),
	})
}

func (c *AdminProjectController) createProto(ctx context.Context, req *bdspropb.Project) (*bdspropb.Project, error) {
	return c.admin.Create(ctx, req)
}

func (c *AdminProjectController) updateProto(ctx context.Context, req *bdspropb.Project) (*bdspropb.Project, error) {
	return c.admin.Update(ctx, req)
}

func (c *AdminProjectController) deleteProto(ctx context.Context, req *sharepb.IdRequest) (*sharepb.SubmitResponse, error) {
	if err := c.admin.Delete(ctx, req.GetId()); err != nil {
		return nil, err
	}
	return &sharepb.SubmitResponse{Id: req.GetId(), Message: "success"}, nil
}
