package admin_usecases

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"bdspro/internal/repo"
	admin_repo "bdspro/internal/repo/admin"
	shared_usecase "bdspro/internal/usecases/shared"
	"common/case/crud"
	"context"
)

type AdminProductUsecase struct {
	crud.BaseUsecase[domain.Product, admin_repo.AdminProductRepo]
	ProductUsecase *shared_usecase.ProductUsecase
	ProductRepo    repo.ProductRepo
}

func NewAdminProductUsecase(
	repo admin_repo.AdminProductRepo,
	productUsecase *shared_usecase.ProductUsecase,
	productRepo repo.ProductRepo,
) *AdminProductUsecase {
	return &AdminProductUsecase{
		BaseUsecase: crud.BaseUsecase[domain.Product, admin_repo.AdminProductRepo]{
			Repo: repo,
		},
		ProductUsecase: productUsecase,
		ProductRepo:    productRepo,
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

// GetList lấy danh sách sản phẩm với phân trang - search toàn bộ sản phẩm trong hệ thống
func (uc *AdminProductUsecase) GetList(ctx context.Context, searchRequest *dto.ProductSearchRequest) ([]domain.Product, int64, error) {
	// Gọi GetProductsByAdmin để lấy sản phẩm với PostIds và giới hạn 2 bản ghi gần nhất
	results, total, err := uc.ProductRepo.GetProductsByAdmin(ctx, searchRequest)
	if err != nil {
		return nil, 0, err
	}

	// Map product status cho từng sản phẩm
	for i := 0; i < len(results); i++ {
		uc.ProductUsecase.MapProductStatus(ctx, &results[i])
	}

	return results, total, nil
}

// GetDetail lấy chi tiết sản phẩm theo ID
func (uc *AdminProductUsecase) GetDetail(ctx context.Context, id uint64) (*domain.Product, error) {
	return uc.ProductUsecase.Detail(ctx, id)
}

// UpdateNote cập nhật ghi chú sản phẩm
func (uc *AdminProductUsecase) UpdateNote(ctx context.Context, id uint64, note string) error {
	return uc.ProductUsecase.UpdateNote(ctx, id, note)
}
