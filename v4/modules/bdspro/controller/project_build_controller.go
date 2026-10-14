package controller

import (
	"context"
	"net/http"

	bdspropb "pb/types/bdspro"
	sharepb "pb/types/shared"
	"userv4/internal/httpapi"
	"userv4/modules/bdspro/dto"
	"userv4/modules/bdspro/services"
)

// ProjectBuildController owns /v2/bdspro/v2/project-build/* (proto parity with bdspro-service).
type ProjectBuildController struct {
	build services.ProjectBuildService
}

func NewProjectBuildController(build services.ProjectBuildService) *ProjectBuildController {
	return &ProjectBuildController{build: build}
}

func (c *ProjectBuildController) RegisterHTTP(mux *http.ServeMux) {
	httpapi.GET(mux, "/v2/bdspro/v2/project-build/me", c.listMine)
	httpapi.HandleProto(mux, http.MethodGet, "/v2/bdspro/v2/project-build/items", c.listItemsProto)
	httpapi.HandleProto(mux, http.MethodGet, "/v2/bdspro/v2/project-build/detail/{id}", c.getByIDProto)
	httpapi.HandleProto(mux, http.MethodPost, "/v2/bdspro/v2/project-build", c.createProto)
	httpapi.HandleProto(mux, http.MethodPut, "/v2/bdspro/v2/project-build/{id}", c.updateProto)
	httpapi.HandleProto(mux, http.MethodDelete, "/v2/bdspro/v2/project-build/{id}", c.deleteProto)
}

func (c *ProjectBuildController) listMine(ctx context.Context, req dto.CatalogSearchRequest) (*httpapi.Result, error) {
	resp, err := c.build.ListMine(ctx, req)
	if err != nil {
		return nil, err
	}
	return listItemResponseResult(resp, req.Page, req.Size)
}

func (c *ProjectBuildController) listItemsProto(ctx context.Context, req *bdspropb.SearchQueryRequest) (*bdspropb.ListItemResponse, error) {
	return c.build.ListItems(ctx, dto.CatalogSearchRequest{
		Page: req.GetPage(),
		Size: req.GetSize(),
		Text: req.GetText(),
	})
}

func (c *ProjectBuildController) getByIDProto(ctx context.Context, req *sharepb.IdRequest) (*bdspropb.ProjectBuild, error) {
	return c.build.GetByID(ctx, req.GetId())
}

func (c *ProjectBuildController) createProto(ctx context.Context, req *bdspropb.ProjectBuild) (*bdspropb.ProjectBuild, error) {
	return c.build.Create(ctx, req)
}

func (c *ProjectBuildController) updateProto(ctx context.Context, req *bdspropb.ProjectBuild) (*bdspropb.ProjectBuild, error) {
	return c.build.Update(ctx, req)
}

func (c *ProjectBuildController) deleteProto(ctx context.Context, req *sharepb.IdRequest) (*sharepb.Empty, error) {
	if err := c.build.Delete(ctx, req.GetId()); err != nil {
		return nil, err
	}
	return &sharepb.Empty{}, nil
}
