package controller

import (
	"context"
	"net/http"

	_dto "common/domain/dto"
	_errors "common/errors"
	"userv4/internal/httpapi"
	"userv4/modules/tqd/dto"
	domain "userv4/modules/tqd/entity"
	"userv4/modules/tqd/services"
)

// DirectoryController owns directory sources/suppliers and service categories.
type DirectoryController struct {
	sources    *services.DirectorySourceUsecase
	suppliers  *services.DirectorySupplierUsecase
	categories *services.DirectoryCategoryUsecase
}

func NewDirectoryController(
	sources *services.DirectorySourceUsecase,
	suppliers *services.DirectorySupplierUsecase,
	categories *services.DirectoryCategoryUsecase,
) *DirectoryController {
	return &DirectoryController{sources: sources, suppliers: suppliers, categories: categories}
}

func (c *DirectoryController) RegisterHTTP(mux *http.ServeMux) {
	httpapi.POST(mux, "/v2/tqd/directory/sources", c.createSource)
	httpapi.GET(mux, "/v2/tqd/directory/sources", c.listSources)
	httpapi.GET(mux, "/v2/tqd/directory/sources/{id}", c.getSource)
	httpapi.PUT(mux, "/v2/tqd/directory/sources/{id}", c.updateSource)
	httpapi.DELETE(mux, "/v2/tqd/directory/sources/{id}", c.deleteSource)

	httpapi.POST(mux, "/v2/tqd/directory/suppliers", c.createSupplier)
	httpapi.GET(mux, "/v2/tqd/directory/suppliers", c.listSuppliers)
	httpapi.GET(mux, "/v2/tqd/directory/suppliers/{id}", c.getSupplier)
	httpapi.PUT(mux, "/v2/tqd/directory/suppliers/{id}", c.updateSupplier)
	httpapi.DELETE(mux, "/v2/tqd/directory/suppliers/{id}", c.deleteSupplier)

	// Canonical directory category routes (tqd-service parity).
	httpapi.POST(mux, "/v2/tqd/directory/categories", c.createCategory)
	httpapi.GET(mux, "/v2/tqd/directory/categories", c.listCategories)
	httpapi.GET(mux, "/v2/tqd/directory/categories/{id}", c.getCategory)
	httpapi.PUT(mux, "/v2/tqd/directory/categories/{id}", c.updateCategory)
	httpapi.DELETE(mux, "/v2/tqd/directory/categories/{id}", c.deleteCategory)

	// Legacy aliases kept for older clients/swagger that used /service/categories.
	httpapi.POST(mux, "/v2/tqd/service/categories", c.createCategory)
	httpapi.GET(mux, "/v2/tqd/service/categories", c.listCategories)
	httpapi.GET(mux, "/v2/tqd/service/categories/active", c.activeCategories)
	httpapi.GET(mux, "/v2/tqd/service/categories/level/{level}", c.categoriesByLevel)
	httpapi.GET(mux, "/v2/tqd/service/categories/{id}", c.getCategory)
	httpapi.PUT(mux, "/v2/tqd/service/categories/{id}", c.updateCategory)
	httpapi.DELETE(mux, "/v2/tqd/service/categories/{id}", c.deleteCategory)
}

func (c *DirectoryController) createSource(ctx context.Context, req dto.DirectorySourceDTO) (*httpapi.Result, error) {
	out, err := c.sources.Create(ctx, &req)
	if err != nil {
		return nil, err
	}
	return jsonOK(out)
}

func (c *DirectoryController) getSource(ctx context.Context, req dto.IDPath) (*httpapi.Result, error) {
	out, err := c.sources.GetByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	return jsonOK(out)
}

func (c *DirectoryController) updateSource(ctx context.Context, req dto.DirectorySourceDTO) (*httpapi.Result, error) {
	if req.ID == 0 {
		return nil, _errors.ReturnError(400, "id is required")
	}
	out, err := c.sources.Update(ctx, req.ToDomain())
	if err != nil {
		return nil, err
	}
	return jsonOK(out)
}

func (c *DirectoryController) deleteSource(ctx context.Context, req dto.IDPath) (*httpapi.Result, error) {
	if err := c.sources.Delete(ctx, req.ID); err != nil {
		return nil, err
	}
	return jsonOK(map[string]any{"success": true})
}

type listDirPage struct {
	Page     uint32 `form:"page"`
	Size     uint32 `form:"size"`
	Search   string `form:"search"`
	IsActive *bool  `form:"isActive"`
}

func (c *DirectoryController) listSources(ctx context.Context, req listDirPage) (*httpapi.Result, error) {
	pagable := _dto.Pagable{Page: req.Page, Size: req.Size}
	pagable.Normalize()
	rows, total, err := c.sources.List(ctx, &dto.ListDirectorySourcesRequestDTO{Pagable: pagable, Search: req.Search})
	if err != nil {
		return nil, err
	}
	return jsonListResult(rows, total, pagable.GetPage(), pagable.GetSize())
}

func (c *DirectoryController) createSupplier(ctx context.Context, req dto.CreateDirectorySupplierRequestDTO) (*httpapi.Result, error) {
	out, err := c.suppliers.Create(ctx, &req)
	if err != nil {
		return nil, err
	}
	return jsonOK(out)
}

func (c *DirectoryController) getSupplier(ctx context.Context, req dto.IDPath) (*httpapi.Result, error) {
	out, err := c.suppliers.GetByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	return jsonOK(out)
}

func (c *DirectoryController) updateSupplier(ctx context.Context, req struct {
	ID uint64 `path:"id" json:"-"`
	dto.UpdateDirectorySupplierRequestDTO
}) (*httpapi.Result, error) {
	out, err := c.suppliers.Update(ctx, req.ID, &req.UpdateDirectorySupplierRequestDTO)
	if err != nil {
		return nil, err
	}
	return jsonOK(out)
}

func (c *DirectoryController) deleteSupplier(ctx context.Context, req dto.IDPath) (*httpapi.Result, error) {
	if err := c.suppliers.Delete(ctx, req.ID); err != nil {
		return nil, err
	}
	return jsonOK(map[string]any{"success": true})
}

func (c *DirectoryController) listSuppliers(ctx context.Context, req listDirPage) (*httpapi.Result, error) {
	pagable := _dto.Pagable{Page: req.Page, Size: req.Size}
	pagable.Normalize()
	out, err := c.suppliers.List(ctx, &dto.DirectorySupplierFilterDTO{Pagable: pagable, Search: req.Search})
	if err != nil {
		return nil, err
	}
	if out == nil {
		return jsonListResult([]dto.DirectorySupplierDTO{}, 0, pagable.GetPage(), pagable.GetSize())
	}
	return jsonListResult(out.Data, out.Total, pagable.GetPage(), pagable.GetSize())
}

func (c *DirectoryController) createCategory(ctx context.Context, req dto.CreateDirectoryCategoryRequestDTO) (*httpapi.Result, error) {
	cat := &domain.DirectoryCategory{
		Name: req.Name, Description: req.Description, Code: req.Code, Icon: req.Icon,
		Color: req.Color, IsActive: req.IsActive, SortOrder: req.SortOrder, Level: req.Level,
	}
	if err := c.categories.CreateDirectoryCategory(ctx, cat); err != nil {
		return nil, err
	}
	return jsonOK(cat)
}

func (c *DirectoryController) getCategory(ctx context.Context, req dto.IDPath) (*httpapi.Result, error) {
	out, err := c.categories.GetDirectoryCategoryByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	return jsonOK(out)
}

func (c *DirectoryController) updateCategory(ctx context.Context, req struct {
	ID          uint64 `path:"id" json:"-"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Code        string `json:"code"`
	Icon        string `json:"icon"`
	Color       string `json:"color"`
	IsActive    bool   `json:"isActive"`
	SortOrder   uint32 `json:"sortOrder"`
	Level       uint32 `json:"level"`
}) (*httpapi.Result, error) {
	out, err := c.categories.UpdateDirectoryCategory(ctx, req.ID, &domain.DirectoryCategory{
		Name: req.Name, Description: req.Description, Code: req.Code, Icon: req.Icon,
		Color: req.Color, IsActive: req.IsActive, SortOrder: req.SortOrder, Level: req.Level,
	})
	if err != nil {
		return nil, err
	}
	return jsonOK(out)
}

func (c *DirectoryController) deleteCategory(ctx context.Context, req dto.IDPath) (*httpapi.Result, error) {
	if err := c.categories.DeleteDirectoryCategory(ctx, req.ID); err != nil {
		return nil, err
	}
	return jsonOK(map[string]any{"success": true})
}

func (c *DirectoryController) listCategories(ctx context.Context, req listDirPage) (*httpapi.Result, error) {
	pagable := _dto.Pagable{Page: req.Page, Size: req.Size}
	pagable.Normalize()
	rows, total, err := c.categories.ListDirectoryCategories(ctx, &dto.ListDirectoryCategoriesRequest{
		Pagable:  pagable,
		Search:   req.Search,
		IsActive: req.IsActive,
	})
	if err != nil {
		return nil, err
	}
	return jsonListResult(rows, total, pagable.GetPage(), pagable.GetSize())
}

func (c *DirectoryController) activeCategories(ctx context.Context, _ dto.EmptyJSON) (*httpapi.Result, error) {
	active := true
	rows, total, err := c.categories.ListDirectoryCategories(ctx, &dto.ListDirectoryCategoriesRequest{
		Pagable: _dto.Pagable{Page: 1, Size: 1000}, IsActive: &active,
	})
	if err != nil {
		return nil, err
	}
	return jsonListResult(rows, total, 1, 1000)
}

func (c *DirectoryController) categoriesByLevel(ctx context.Context, req dto.LevelPath) (*httpapi.Result, error) {
	rows, _, err := c.categories.ListDirectoryCategories(ctx, &dto.ListDirectoryCategoriesRequest{
		Pagable: _dto.Pagable{Page: 1, Size: 1000},
	})
	if err != nil {
		return nil, err
	}
	filtered := make([]domain.DirectoryCategory, 0)
	for _, row := range rows {
		if row.Level == uint32(req.Level) {
			filtered = append(filtered, row)
		}
	}
	return jsonListResult(filtered, int64(len(filtered)), 1, uint32(len(filtered)))
}
