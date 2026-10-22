package controller

import (
	"context"
	"net/http"

	"common/httpapi"
	"mdn4/modules/bdspro/dto"
	"mdn4/modules/bdspro/services"
)

// CatalogController owns public list routes under /v2/bdspro/v2/list/* (project-build/items is on ProjectBuildController).
type CatalogController struct {
	catalog services.CatalogService
}

func NewCatalogController(catalog services.CatalogService) *CatalogController {
	return &CatalogController{catalog: catalog}
}

func (c *CatalogController) RegisterHTTP(mux *http.ServeMux) {
	httpapi.GET(mux, "/v2/bdspro/v2/list/project", c.listProjects)
	httpapi.GET(mux, "/v2/bdspro/v2/list/project-build", c.listProjectBuilds)
}

func (c *CatalogController) listProjects(ctx context.Context, req dto.CatalogSearchRequest) (*httpapi.Result, error) {
	resp, err := c.catalog.ListProjects(ctx, req)
	if err != nil {
		return nil, err
	}
	return listItemResponseResult(resp, req.Page, req.Size)
}

func (c *CatalogController) listProjectBuilds(ctx context.Context, req dto.CatalogSearchRequest) (*httpapi.Result, error) {
	return c.listProjectBuildItems(ctx, req)
}

func (c *CatalogController) listProjectBuildItems(ctx context.Context, req dto.CatalogSearchRequest) (*httpapi.Result, error) {
	resp, err := c.catalog.ListProjectBuildItems(ctx, req)
	if err != nil {
		return nil, err
	}
	return listItemResponseResult(resp, req.Page, req.Size)
}
