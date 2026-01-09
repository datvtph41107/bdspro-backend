package usecases

import (
	"bdspro/internal/dto"
	"bdspro/internal/repo"
	_routes "common/routes"
	_utils "common/utils"
	"context"
	"fmt"
)

type ProductAssetUsecase struct {
	ProductAssetRepo repo.ProductAssetRepo
	ProductRepo      repo.ProductRepo
	AssetRepo        repo.AssetRepo
}

func NewProductAssetUsecase(
	productAssetRepo repo.ProductAssetRepo,
	productRepo repo.ProductRepo,
	assetRepo repo.AssetRepo,
) *ProductAssetUsecase {
	return &ProductAssetUsecase{
		ProductAssetRepo: productAssetRepo,
		ProductRepo:      productRepo,
		AssetRepo:        assetRepo,
	}
}

// LinkProductAsset liên kết nhiều product với asset
func (uc *ProductAssetUsecase) LinkProductAsset(ctx context.Context, req *dto.LinkProductAssetRequest) error {
	// Validate danh sách productIds không rỗng
	if len(req.ProductIDs) == 0 {
		return &_routes.Except{
			Code:    400,
			Message: "Danh sách sản phẩm không được để trống",
		}
	}

	// Validate asset tồn tại
	asset, err := uc.AssetRepo.GetByID(req.AssetID)
	if err != nil {
		return &_routes.Except{
			Code:    404,
			Message: "Tài sản không tồn tại",
		}
	}
	if asset == nil {
		return &_routes.Except{
			Code:    404,
			Message: "Tài sản không tồn tại",
		}
	}

	// Validate từng product và kiểm tra liên kết đã tồn tại
	var productsToLink []uint64
	for _, productID := range req.ProductIDs {
		// Validate product tồn tại
		product, err := uc.ProductRepo.GetByIDContext(ctx, productID)
		if err != nil {
			return &_routes.Except{
				Code:    404,
				Message: fmt.Sprintf("Sản phẩm ID %d không tồn tại", productID),
			}
		}
		if product == nil {
			return &_routes.Except{
				Code:    404,
				Message: fmt.Sprintf("Sản phẩm ID %d không tồn tại", productID),
			}
		}

		// Kiểm tra liên kết đã tồn tại chưa
		exists, err := uc.ProductAssetRepo.CheckExists(ctx, productID, req.AssetID)
		if err != nil {
			return &_routes.Except{
				Code:    500,
				Message: fmt.Sprintf("Lỗi kiểm tra liên kết cho sản phẩm ID %d", productID),
			}
		}
		if !exists {
			// Chỉ thêm vào danh sách nếu chưa tồn tại
			productsToLink = append(productsToLink, productID)
		}
	}

	// Tạo liên kết cho các product chưa được liên kết
	if len(productsToLink) == 0 {
		return &_routes.Except{
			Code:    400,
			Message: "Tất cả các sản phẩm đã được liên kết với tài sản này",
		}
	}

	// Tạo liên kết cho từng product
	for _, productID := range productsToLink {
		err = uc.ProductAssetRepo.Link(ctx, productID, req.AssetID)
		if err != nil {
			return &_routes.Except{
				Code:    500,
				Message: fmt.Sprintf("Lỗi tạo liên kết cho sản phẩm ID %d", productID),
			}
		}
	}

	// Cập nhật AssetCount cho các products được liên kết
	uc.updateAssetCountForProducts(ctx, productsToLink)

	return nil
}

// UnlinkProductAsset hủy liên kết product với asset
func (uc *ProductAssetUsecase) UnlinkProductAsset(ctx context.Context, req *dto.UnlinkProductAssetRequest) error {
	// Kiểm tra liên kết tồn tại
	exists, err := uc.ProductAssetRepo.CheckExists(ctx, req.ProductID, req.AssetID)
	if err != nil {
		return &_routes.Except{
			Code:    500,
			Message: "Lỗi kiểm tra liên kết",
		}
	}
	if !exists {
		return &_routes.Except{
			Code:    404,
			Message: "Liên kết không tồn tại",
		}
	}

	// Hủy liên kết
	err = uc.ProductAssetRepo.Unlink(ctx, req.ProductID, req.AssetID)
	if err != nil {
		return &_routes.Except{
			Code:    500,
			Message: "Lỗi hủy liên kết",
		}
	}

	// Cập nhật AssetCount cho product sau khi hủy liên kết
	uc.updateAssetCountForProducts(ctx, []uint64{req.ProductID})

	return nil
}

// GetAssetsByProductID lấy danh sách asset theo product ID
func (uc *ProductAssetUsecase) GetAssetsByProductID(ctx context.Context, productID uint64) ([]uint64, error) {
	assetIDs, err := uc.ProductAssetRepo.GetAssetIDsByProductID(ctx, productID)
	if err != nil {
		return nil, &_routes.Except{
			Code:    500,
			Message: "Lỗi lấy danh sách tài sản",
		}
	}
	return assetIDs, nil
}

// GetProductsByAssetID lấy danh sách product theo asset ID
func (uc *ProductAssetUsecase) GetProductsByAssetID(ctx context.Context, assetID uint64) ([]uint64, error) {
	productIDs, err := uc.ProductAssetRepo.GetProductIDsByAssetID(ctx, assetID)
	if err != nil {
		return nil, &_routes.Except{
			Code:    500,
			Message: "Lỗi lấy danh sách sản phẩm",
		}
	}
	return productIDs, nil
}

// updateAssetCountForProducts cập nhật AssetCount cho danh sách products
func (uc *ProductAssetUsecase) updateAssetCountForProducts(ctx context.Context, productIDs []uint64) {
	go func() {
		cloneCtx := _utils.CloneContext(ctx)
		for _, productID := range productIDs {
			count, err := uc.ProductAssetRepo.CountAssetsByProductID(cloneCtx, productID)
			if err == nil {
				_ = uc.ProductRepo.UpdateAssetCount(cloneCtx, productID, uint32(count))
			}
		}
	}()
}
