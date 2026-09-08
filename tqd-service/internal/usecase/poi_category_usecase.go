package usecase

import (
	"context"
	"fmt"

	_utils "common/utils"
	"tqd/infra/mapper"
	"tqd/internal/dto"
	"tqd/internal/interface/repo"
)

// PoiCategoryUsecase defines business logic for POI category
type PoiCategoryUsecase interface {
	// CRUD operations
	Create(ctx context.Context, req *dto.CreatePoiCategoryRequest) (*dto.PoiCategoryResponse, error)
	GetByID(ctx context.Context, id uint64) (*dto.PoiCategoryResponse, error)
	Update(ctx context.Context, id uint64, req *dto.UpdatePoiCategoryRequest) (*dto.PoiCategoryResponse, error)
	Delete(ctx context.Context, id uint64) error

	// List operations
	List(ctx context.Context, filter *dto.PoiCategoryFilter) (*dto.ListPoiCategoriesResponse, error)

	// Tree operations
	GetTree(ctx context.Context) ([]dto.PoiCategoryTreeResponse, error)

	// POI count operations
	IncrementPOICount(ctx context.Context, id uint64) error
	DecrementPOICount(ctx context.Context, id uint64) error
}

// poiCategoryUsecase implements PoiCategoryUsecase
type poiCategoryUsecase struct {
	repo   repo.IPoiCategoryRepo
	mapper *mapper.PoiCategoryMapper
}

// NewPoiCategoryUsecase creates new POI category usecase
func NewPoiCategoryUsecase(
	repo repo.IPoiCategoryRepo,
	mapper *mapper.PoiCategoryMapper,
) PoiCategoryUsecase {
	return &poiCategoryUsecase{
		repo:   repo,
		mapper: mapper,
	}
}

// Create implements PoiCategoryUsecase.Create
func (u *poiCategoryUsecase) Create(ctx context.Context, req *dto.CreatePoiCategoryRequest) (*dto.PoiCategoryResponse, error) {
	// Get user ID
	userID := _utils.GetProfileIdWithContext(ctx)
	if userID == 0 {
		return nil, fmt.Errorf("unauthorized")
	}

	// Validate required fields
	if req.Code == "" {
		return nil, fmt.Errorf("code is required")
	}
	if req.Name == "" {
		return nil, fmt.Errorf("name is required")
	}

	// Check if parent exists
	if req.ParentID != nil && *req.ParentID > 0 {
		parent, err := u.repo.GetByID(ctx, *req.ParentID)
		if err != nil {
			return nil, fmt.Errorf("failed to check parent category: %v", err)
		}
		if parent == nil {
			return nil, fmt.Errorf("parent category not found")
		}
	}

	// Convert to domain
	category, err := u.mapper.ToDomainFromCreate(req, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to map request: %v", err)
	}

	// Create in repository
	if err := u.repo.Create(ctx, category); err != nil {
		return nil, fmt.Errorf("failed to create category: %v", err)
	}

	// Get created category
	created, err := u.repo.GetByID(ctx, category.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get created category: %v", err)
	}

	return u.mapper.ToResponse(created), nil
}

// GetByID implements PoiCategoryUsecase.GetByID
func (u *poiCategoryUsecase) GetByID(ctx context.Context, id uint64) (*dto.PoiCategoryResponse, error) {
	if id == 0 {
		return nil, fmt.Errorf("invalid category ID")
	}

	category, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get category: %v", err)
	}
	if category == nil {
		return nil, fmt.Errorf("category not found")
	}

	return u.mapper.ToResponse(category), nil
}

// Update implements PoiCategoryUsecase.Update
func (u *poiCategoryUsecase) Update(ctx context.Context, id uint64, req *dto.UpdatePoiCategoryRequest) (*dto.PoiCategoryResponse, error) {
	// Get user ID
	userID := _utils.GetProfileIdWithContext(ctx)
	if userID == 0 {
		return nil, fmt.Errorf("unauthorized")
	}

	// Get existing category
	existing, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get category: %v", err)
	}
	if existing == nil {
		return nil, fmt.Errorf("category not found")
	}

	// Check if parent exists and not creating cycle
	if req.ParentID != nil && *req.ParentID > 0 {
		if *req.ParentID == id {
			return nil, fmt.Errorf("category cannot be its own parent")
		}

		parent, err := u.repo.GetByID(ctx, *req.ParentID)
		if err != nil {
			return nil, fmt.Errorf("failed to check parent category: %v", err)
		}
		if parent == nil {
			return nil, fmt.Errorf("parent category not found")
		}
	}

	// Convert to domain
	updated, err := u.mapper.ToDomainFromUpdate(req, existing, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to map request: %v", err)
	}

	// Update in repository
	if err := u.repo.Update(ctx, id, updated); err != nil {
		return nil, fmt.Errorf("failed to update category: %v", err)
	}

	// Get updated category
	result, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated category: %v", err)
	}

	return u.mapper.ToResponse(result), nil
}

// Delete implements PoiCategoryUsecase.Delete
func (u *poiCategoryUsecase) Delete(ctx context.Context, id uint64) error {
	if id == 0 {
		return fmt.Errorf("invalid category ID")
	}

	// Check if category exists
	existing, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get category: %v", err)
	}
	if existing == nil {
		return fmt.Errorf("category not found")
	}

	// Check if category has children
	children, err := u.repo.FindByParent(ctx, &id)
	if err != nil {
		return fmt.Errorf("failed to check children: %v", err)
	}
	if len(children) > 0 {
		return fmt.Errorf("cannot delete category with children")
	}

	// Delete from repository
	if err := u.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete category: %v", err)
	}

	return nil
}

// List implements PoiCategoryUsecase.List
func (u *poiCategoryUsecase) List(ctx context.Context, filter *dto.PoiCategoryFilter) (*dto.ListPoiCategoriesResponse, error) {
	// Validate pagination
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Size < 1 {
		filter.Size = 10
	}
	if filter.Size > 100 {
		filter.Size = 100
	}

	// Get from repository
	categories, total, err := u.repo.ListWithFilter(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to list categories: %v", err)
	}

	return &dto.ListPoiCategoriesResponse{
		Data:  u.mapper.ToResponseList(categories),
		Total: total,
		Page:  filter.Page,
		Size:  filter.Size,
	}, nil
}

// GetTree implements PoiCategoryUsecase.GetTree
func (u *poiCategoryUsecase) GetTree(ctx context.Context) ([]dto.PoiCategoryTreeResponse, error) {
	categories, err := u.repo.GetTree(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get category tree: %v", err)
	}

	return u.mapper.BuildTree(categories), nil
}

// IncrementPOICount implements PoiCategoryUsecase.IncrementPOICount
func (u *poiCategoryUsecase) IncrementPOICount(ctx context.Context, id uint64) error {
	if id == 0 {
		return fmt.Errorf("invalid category ID")
	}

	if err := u.repo.IncrementPOICount(ctx, id); err != nil {
		return fmt.Errorf("failed to increment POI count: %v", err)
	}

	return nil
}

// DecrementPOICount implements PoiCategoryUsecase.DecrementPOICount
func (u *poiCategoryUsecase) DecrementPOICount(ctx context.Context, id uint64) error {
	if id == 0 {
		return fmt.Errorf("invalid category ID")
	}

	if err := u.repo.DecrementPOICount(ctx, id); err != nil {
		return fmt.Errorf("failed to decrement POI count: %v", err)
	}

	return nil
}
