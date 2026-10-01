package controller

import (
	"context"
	"net/http"

	_utils "common/utils"
	"userv4/internal/httpapi"
	"userv4/modules/seo/dto"
	"userv4/modules/seo/services"
)

// NewsController handles planning news public/client/admin routes.
type NewsController struct {
	svc services.NewsService
}

func NewNewsController(svc services.NewsService) *NewsController {
	return &NewsController{svc: svc}
}

func (c *NewsController) RegisterHTTP(mux *http.ServeMux) {
	// Public news
	httpapi.GET(mux, "/v2/crm/public/seo/news", c.listNews)
	httpapi.GET(mux, "/v2/crm/public/seo/news/{id}", c.getNews)
	httpapi.POST(mux, "/v2/crm/public/seo/click", c.recordClick)
	// Public categories
	httpapi.GET(mux, "/v2/crm/public/content/categories", c.listCategories)
	// Public content search
	httpapi.GET(mux, "/v2/crm/public/content/search", c.searchContent)
	// Client saved planning-news (authenticated)
	httpapi.GET(mux, "/v2/crm/client/content/planning-news/saved", c.listSaved)
	httpapi.POST(mux, "/v2/crm/client/content/planning-news/{id}/saved", c.saveNews)
	mux.HandleFunc("DELETE /v2/crm/client/content/planning-news/{id}/saved", httpapi.Bind(c.unsaveNews))
	// Admin categories CRUD (CRM path parity)
	httpapi.GET(mux, "/v2/crm/admin/content/categories", c.adminListCategories)
	httpapi.GET(mux, "/v2/crm/admin/seo/categories", c.adminListCategories)
	httpapi.GET(mux, "/v2/crm/admin/seo/categories/{id}", c.adminGetCategory)
	httpapi.POST(mux, "/v2/crm/admin/content/categories", c.adminCreateCategory)
	httpapi.POST(mux, "/v2/crm/admin/seo/categories", c.adminCreateCategory)
	httpapi.PUT(mux, "/v2/crm/admin/content/categories/{id}", c.adminUpdateCategory)
	httpapi.PUT(mux, "/v2/crm/admin/seo/categories/{id}", c.adminUpdateCategory)
	mux.HandleFunc("DELETE /v2/crm/admin/content/categories/{id}", httpapi.Bind(c.adminDeleteCategory))
	mux.HandleFunc("DELETE /v2/crm/admin/seo/categories/{id}", httpapi.Bind(c.adminDeleteCategory))
}

// ── Public ────────────────────────────────────────────────────────────────────

func (c *NewsController) listNews(ctx context.Context, req dto.PlanningNewsListFilter) (*httpapi.Result, error) {
	result, err := c.svc.ListPlanningNews(ctx, req)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return listResult([]dto.PlanningNewsItem{}, 0)
	}
	return listResult(result.Data, result.Total)
}

func (c *NewsController) getNews(ctx context.Context, req dto.GetNewsDetailRequest) (*httpapi.Result, error) {
	result, err := c.svc.GetPlanningNews(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	return dataResult(result)
}

func (c *NewsController) recordClick(ctx context.Context, req dto.RecordClickRequest) (*httpapi.Result, error) {
	if req.UserID == 0 {
		req.UserID = _utils.GetProfileIdWithContext(ctx)
	}
	result, err := c.svc.RecordClick(ctx, req)
	if err != nil {
		return nil, err
	}
	return dataResult(result)
}

func (c *NewsController) listCategories(ctx context.Context, _ httpapi.Empty) (*httpapi.Result, error) {
	cats, err := c.svc.ListCategories(ctx)
	if err != nil {
		return nil, err
	}
	return listResult(cats, int64(len(cats)))
}

func (c *NewsController) searchContent(ctx context.Context, req dto.ContentSearchRequest) (*httpapi.Result, error) {
	result, err := c.svc.SearchContent(ctx, req)
	if err != nil {
		return nil, err
	}
	return dataResult(result)
}

// ── Client (authenticated) ────────────────────────────────────────────────────

func (c *NewsController) listSaved(ctx context.Context, req dto.SavedNewsFilter) (*httpapi.Result, error) {
	userID := _utils.GetProfileIdWithContext(ctx)
	result, err := c.svc.ListSavedPlanningNews(ctx, userID, req)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return listResult([]dto.PlanningNewsItem{}, 0)
	}
	return listResult(result.Data, result.Total)
}

func (c *NewsController) saveNews(ctx context.Context, req dto.NewsIDPath) (*httpapi.Result, error) {
	userID := _utils.GetProfileIdWithContext(ctx)
	item, err := c.svc.SavePlanningNews(ctx, userID, req.ID)
	if err != nil {
		return nil, err
	}
	return dataResult(item)
}

func (c *NewsController) unsaveNews(ctx context.Context, req dto.NewsIDPath) (*httpapi.Result, error) {
	userID := _utils.GetProfileIdWithContext(ctx)
	if err := c.svc.UnsavePlanningNews(ctx, userID, req.ID); err != nil {
		return nil, err
	}
	return dataResult(map[string]string{"message": "removed"})
}

// ── Admin ─────────────────────────────────────────────────────────────────────

func (c *NewsController) adminListCategories(ctx context.Context, _ httpapi.Empty) (*httpapi.Result, error) {
	return c.listCategories(ctx, httpapi.Empty{})
}

func (c *NewsController) adminGetCategory(ctx context.Context, req httpapi.PathUint64) (*httpapi.Result, error) {
	cat, err := c.svc.GetCategory(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	return dataResult(cat)
}

func (c *NewsController) adminCreateCategory(ctx context.Context, req dto.CreateCategoryRequest) (*httpapi.Result, error) {
	cat, err := c.svc.AdminCreateCategory(ctx, req)
	if err != nil {
		return nil, err
	}
	return dataResult(cat)
}

func (c *NewsController) adminUpdateCategory(ctx context.Context, req dto.UpdateCategoryRequest) (*httpapi.Result, error) {
	cat, err := c.svc.AdminUpdateCategory(ctx, req.ID, req)
	if err != nil {
		return nil, err
	}
	return dataResult(cat)
}

func (c *NewsController) adminDeleteCategory(ctx context.Context, req httpapi.PathUint64) (*httpapi.Result, error) {
	if err := c.svc.AdminDeleteCategory(ctx, req.ID); err != nil {
		return nil, err
	}
	return dataResult(map[string]string{"message": "deleted"})
}
