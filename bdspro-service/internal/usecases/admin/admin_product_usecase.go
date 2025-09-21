package admin_usecases

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	admin_repo "bdspro/internal/repo/admin"
	shared_usecase "bdspro/internal/usecases/shared"
	"common/case/crud"
	"context"
)

type AdminProductUsecase struct {
	crud.BaseUsecase[domain.Product, admin_repo.AdminProductRepo]
	ProductUsecase *shared_usecase.ProductUsecase // @bind: internal/usecases/shared
}

func NewAdminProductUsecase(
	repo admin_repo.AdminProductRepo,
	productUsecase *shared_usecase.ProductUsecase,
) *AdminProductUsecase {
	return &AdminProductUsecase{
		BaseUsecase: crud.BaseUsecase[domain.Product, admin_repo.AdminProductRepo]{
			Repo: repo,
		},
		ProductUsecase: productUsecase,
	}
}

// Create tạo sản phẩm mới với ownerId từ request body
func (uc *AdminProductUsecase) Create(ctx context.Context, product *dto.ProductSaveRequest) (*domain.Product, error) {
	// Sử dụng ProductUsecase hiện có để tạo sản phẩm
	// OwnerId đã được set từ request body
	result, err := uc.ProductUsecase.CreateProduct(ctx, product)
	if err != nil {
		return nil, err
	}
	return result.Product, nil
}

// GetList lấy danh sách sản phẩm với phân trang
func (uc *AdminProductUsecase) GetList(ctx context.Context, searchRequest *dto.ProductSearchRequest) ([]domain.Product, int64, error) {
	return uc.ProductUsecase.SearchProductByOwner(ctx, searchRequest)
}

// GetDetail lấy chi tiết sản phẩm theo ID
func (uc *AdminProductUsecase) GetDetail(ctx context.Context, id uint64) (*domain.Product, error) {
	return uc.ProductUsecase.Detail(ctx, id)
}
