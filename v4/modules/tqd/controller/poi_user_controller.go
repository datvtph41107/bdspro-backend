package controller

import (
	"context"
	"net/http"

	_dto "common/domain/dto"
	"common/httpapi"
	"mdn4/modules/tqd/dto"
	"mdn4/modules/tqd/services"
)

// PoiUserController owns /v2/tqd/pois* (non-admin) and /v2/tqd/poi/categories*.
type PoiUserController struct {
	pois       services.PoiUsecase
	categories services.PoiCategoryUsecase
}

func NewPoiUserController(pois services.PoiUsecase, categories services.PoiCategoryUsecase) *PoiUserController {
	return &PoiUserController{pois: pois, categories: categories}
}

func (c *PoiUserController) RegisterHTTP(mux *http.ServeMux) {
	httpapi.POST(mux, "/v2/tqd/pois", c.createPoi)
	httpapi.GET(mux, "/v2/tqd/pois/nearby", c.nearby)
	httpapi.GET(mux, "/v2/tqd/pois/code/{code}", c.getByCode)
	httpapi.GET(mux, "/v2/tqd/pois/category/{categoryId}", c.byCategory)
	httpapi.GET(mux, "/v2/tqd/pois/{id}", c.getPoi)
	httpapi.PUT(mux, "/v2/tqd/pois/{id}/rating", c.updateRating)

	httpapi.GET(mux, "/v2/tqd/poi/categories", c.listCategories)
	httpapi.GET(mux, "/v2/tqd/poi/categories/icons", c.listIcons)
	httpapi.GET(mux, "/v2/tqd/poi/categories/tree", c.tree)
	httpapi.GET(mux, "/v2/tqd/poi/categories/{id}", c.getCategory)
	httpapi.PUT(mux, "/v2/tqd/poi/categories/{id}", c.updateCategory)
	httpapi.DELETE(mux, "/v2/tqd/poi/categories/{id}", c.deleteCategory)

	// admin create category
	httpapi.POST(mux, "/v2/tqd/admin/poi/categories", c.createCategory)
}

func (c *PoiUserController) createPoi(ctx context.Context, req dto.CreatePoiRequest) (*httpapi.Result, error) {
	out, err := c.pois.Create(ctx, &req)
	if err != nil {
		return nil, err
	}
	return jsonOK(out)
}

func (c *PoiUserController) getPoi(ctx context.Context, req dto.IDPath) (*httpapi.Result, error) {
	out, err := c.pois.GetByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	return jsonOK(out)
}

func (c *PoiUserController) getByCode(ctx context.Context, req dto.CodePath) (*httpapi.Result, error) {
	out, err := c.pois.GetByCode(ctx, req.Code)
	if err != nil {
		return nil, err
	}
	return jsonOK(out)
}

func (c *PoiUserController) byCategory(ctx context.Context, req dto.CategoryIDPath) (*httpapi.Result, error) {
	out, err := c.pois.ListByCategory(ctx, req.CategoryID, 100)
	if err != nil {
		return nil, err
	}
	return jsonListResult(out, int64(len(out)), 1, uint32(len(out)))
}

type nearbyPoiHTTP struct {
	Lat    float64 `form:"lat"`
	Lng    float64 `form:"lng"`
	Radius float64 `form:"radius"`
	Limit  int     `form:"limit"`
}

func (c *PoiUserController) nearby(ctx context.Context, req nearbyPoiHTTP) (*httpapi.Result, error) {
	out, err := c.pois.ListNearby(ctx, &dto.NearbyPoiRequest{
		Latitude: req.Lat, Longitude: req.Lng, Radius: req.Radius,
	})
	if err != nil {
		return nil, err
	}
	if out == nil {
		return jsonListResult([]dto.PoiResponse{}, 0, 1, 0)
	}
	return jsonListResult(out.Data, out.Total, out.Page, out.Size)
}

func (c *PoiUserController) updateRating(ctx context.Context, req dto.UpdateRatingRequest) (*httpapi.Result, error) {
	out, err := c.pois.UpdateRating(ctx, req.ID, req.Rating, req.ReviewCount)
	if err != nil {
		return nil, err
	}
	return jsonOK(out)
}

func (c *PoiUserController) createCategory(ctx context.Context, req dto.CreatePoiCategoryRequest) (*httpapi.Result, error) {
	out, err := c.categories.Create(ctx, &req)
	if err != nil {
		return nil, err
	}
	return jsonOK(out)
}

func (c *PoiUserController) getCategory(ctx context.Context, req dto.IDPath) (*httpapi.Result, error) {
	out, err := c.categories.GetByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	return jsonOK(out)
}

func (c *PoiUserController) updateCategory(ctx context.Context, req struct {
	ID uint64 `path:"id" json:"-"`
	dto.UpdatePoiCategoryRequest
}) (*httpapi.Result, error) {
	out, err := c.categories.Update(ctx, req.ID, &req.UpdatePoiCategoryRequest)
	if err != nil {
		return nil, err
	}
	return jsonOK(out)
}

func (c *PoiUserController) deleteCategory(ctx context.Context, req dto.IDPath) (*httpapi.Result, error) {
	if err := c.categories.Delete(ctx, req.ID); err != nil {
		return nil, err
	}
	return jsonOK(map[string]any{"success": true})
}

type listCatReq struct {
	Page   uint32 `form:"page"`
	Size   uint32 `form:"size"`
	Search string `form:"search"`
}

func (c *PoiUserController) listCategories(ctx context.Context, req listCatReq) (*httpapi.Result, error) {
	pagable := _dto.Pagable{Page: req.Page, Size: req.Size}
	pagable.Normalize()
	out, err := c.categories.List(ctx, &dto.PoiCategoryFilter{Pagable: pagable, Search: req.Search})
	if err != nil {
		return nil, err
	}
	if out == nil {
		return jsonListResult([]dto.PoiCategoryResponse{}, 0, pagable.GetPage(), pagable.GetSize())
	}
	return jsonListResult(out.Data, out.Total, pagable.GetPage(), pagable.GetSize())
}

func (c *PoiUserController) tree(ctx context.Context, _ dto.EmptyJSON) (*httpapi.Result, error) {
	out, err := c.categories.GetTree(ctx)
	if err != nil {
		return nil, err
	}
	return jsonListResult(out, int64(len(out)), 1, uint32(len(out)))
}

func (c *PoiUserController) listIcons(ctx context.Context, _ dto.EmptyJSON) (*httpapi.Result, error) {
	out, err := c.categories.ListIcons(ctx)
	if err != nil {
		return nil, err
	}
	if out == nil {
		return jsonListResult([]dto.PoiCategoryIconResponse{}, 0, 1, 0)
	}
	return jsonListResult(out.Data, int64(len(out.Data)), 1, uint32(len(out.Data)))
}
