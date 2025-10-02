package usecase

import (
	_crud "common/domain/crud"
	_err "common/domain/err"
	"context"
	"tqd/internal/domain"
	"tqd/internal/dto"
	"tqd/internal/interface/repo"
)

// PoiCategoryUsecase defines the interface for POI category usecase
type PoiCategoryUsecase interface {
	_crud.IBaseUsecase[domain.PoiCategory]
	GetTree(ctx context.Context) ([]dto.PoiCategoryTreeResponse, *_err.ErrorDTO)
}

// poiCategoryUsecase implements PoiCategoryUsecase interface
type poiCategoryUsecase struct {
	_crud.BaseUsecase[domain.PoiCategory, repo.IPoiCategoryRepo]
}

// NewPoiCategoryUsecase creates a new PoiCategoryUsecase
func NewPoiCategoryUsecase(
	poiCategoryRepo repo.IPoiCategoryRepo,
) PoiCategoryUsecase {
	return &poiCategoryUsecase{
		BaseUsecase: _crud.BaseUsecase[domain.PoiCategory, repo.IPoiCategoryRepo]{
			Repo: poiCategoryRepo,
		},
	}
}

// GetTree lấy tất cả category dưới dạng cấu trúc cha-con
func (u *poiCategoryUsecase) GetTree(ctx context.Context) ([]dto.PoiCategoryTreeResponse, *_err.ErrorDTO) {
	categories, err := u.Repo.GetTree(ctx)
	if err != nil {
		return nil, &_err.ErrorDTO{
			Code:    500,
			Message: "Không thể lấy danh sách danh mục POI dạng cây",
		}
	}

	return u.mapToTreeResponse(categories), nil
}

// mapToTreeResponse chuyển đổi domain entity sang tree DTO
func (u *poiCategoryUsecase) mapToTreeResponse(categories []domain.PoiCategory) []dto.PoiCategoryTreeResponse {
	result := make([]dto.PoiCategoryTreeResponse, 0, len(categories))

	for _, category := range categories {
		treeItem := dto.PoiCategoryTreeResponse{
			ID:          category.ID,
			Code:        category.Code,
			Name:        category.Name,
			ParentID:    category.ParentID,
			Icon:        category.Icon,
			IsActive:    category.IsActive,
			SortOrder:   category.SortOrder,
			Description: category.Description,
		}

		// Map children nếu có
		if len(category.Children) > 0 {
			treeItem.Children = u.mapToTreeResponse(category.Children)
		}

		result = append(result, treeItem)
	}

	return result
}
