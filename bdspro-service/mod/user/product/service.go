package u_product

import (
	"bdspro/data"
	"bdspro/enums"
	product_accecss "bdspro/mod/global/product_access"
	u_asset "bdspro/mod/user/asset"
	u_post "bdspro/mod/user/post"
	"bdspro/models"
	_db "common/db"
	_jwt "common/jwt"
	_routes "common/routes"
	_utils "common/utils"
	"context"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
)

type ProductService struct {
	productRepo       *UserProductRepo
	postRepo          *u_post.PostRepo
	assetRepo         *u_asset.AssetRepo
	productAccessRepo *product_accecss.ProductAccessRepo
}

func NewProductService(productRepo *UserProductRepo,
	postRepo *u_post.PostRepo,
	assetRepo *u_asset.AssetRepo,
	productAccessRepo *product_accecss.ProductAccessRepo,
) *ProductService {
	return &ProductService{
		productRepo:       productRepo,
		postRepo:          postRepo,
		assetRepo:         assetRepo,
		productAccessRepo: productAccessRepo,
	}
}

// func (s *ProductService) UpdateProductStatusToSelling(productId uint64) error {
// 	product, err := s.productRepo.GetByID(productId)
// 	if err != nil {
// 		return err
// 	}

// 	// Nếu chưa thuộc sở hữu thì không cập nhật tài sản
// 	if product.OwnerUserID == nil && product.OwnerOrgID == nil {
// 		return errors.New("product chưa thuộc sở hữu, không thể cập nhật")
// 	}

// 	// Chuyển trạng thái sản phẩm sang "Đang bán"
// 	product.SaleStatus = enums.ProductSelling

// 	// Kiểm tra visibility và tạo tin đăng
// 	// if product.SaleVisibility == enums.ProductS {
// 	// 	post := &models.Post{
// 	// 		ProductId: product.ID,
// 	// 		Status:    enums.PostPublic, // Giả sử 1 là trạng thái "Đang hiển thị"
// 	// 	}
// 	// 	if err := s.postRepo.CreatePost(post); err != nil {
// 	// 		return err
// 	// 	}
// 	// }

// 	// Cập nhật sản phẩm
// 	return s.productRepo.UpdateProduct(product)
// }

// Chuyển trạng thái sản phẩm về "Chưa bán"
// func (s *ProductService) UpdateProductStatusToNotSelling(productId uint64) error {
// 	product, err := s.productRepo.GetByID(productId)
// 	if err != nil {
// 		return err
// 	}

// 	if product.OwnerOrgID != nil {
// 		// Nếu công ty sở hữu SP, cập nhật trạng thái
// 		product.SaleStatus = 1 // Chưa bán
// 		if err := s.productRepo.UpdateProduct(product); err != nil {
// 			return err
// 		}
// 	}

// 	// Hủy hoặc ẩn tin đăng
// 	post, err := s.postRepo.GetPostByProductID(product.ID)
// 	if err == nil && post != nil {
// 		post.Status = 3 // Giả sử 3 là "Tạm ẩn"
// 		if err := s.postRepo.UpdatePost(post); err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }

func (s *ProductService) UpdateProductStatus(ctx context.Context,
	productID uint64,
	newStatus enums.ProductStatus,
	buyerID *uint64,
	buyerOrgID *uint64,
) error {
	product, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		return err
	}

	switch newStatus {
	case enums.ProductSelling:
		if product.SaleStatus == enums.ProductNotSold {
			// nếu thuộc sở hữu của công ty hoặc cá nhân thì đổi trạng thái tài sản sang đang bán
			if product.OwnerUserID != nil || product.OwnerOrgID != nil {
				if err := s.assetRepo.UpdateStatus(ctx, product.AssetId, enums.AssetStatusSelling); err != nil {
					return err
				}
			}

			if product.SaleVisibility == enums.ProductVisiblePublic {
				if err := s.postRepo.CreateOrUpdatePostVisibility(ctx, productID, enums.PostVisiblePublic); err != nil {
					return err
				}
			} else if product.SaleVisibility == enums.ProductVisibleInternal {
				if err := s.postRepo.CreateOrUpdatePostVisibility(ctx, productID, enums.PostVisibleInternal); err != nil {
					return err
				}
			}
			// nếu sản phẩm là riêng tư thì không tạo tin
		}

	case enums.ProductNotSold:
		if product.SaleStatus == enums.ProductSelling {
			// sản phẩm thuộc công ty -> đổi sang chưa bán
			if product.OwnerOrgID != nil {
				if err := s.assetRepo.UpdateStatus(ctx, product.AssetId, enums.AssetStatusNotSold); err != nil {
					return err
				}
			}
			// công ty không sở hữu tài sản không thay đổi

			post, err := s.postRepo.GetPostByProductID(product.ID)
			if err != nil {
				return err
			}

			// tin đang công khai => tạm ẩn
			if post.Visibility == enums.PostVisiblePublic {
				if err := s.postRepo.UpdatePostStatus(ctx, productID, enums.PostHide); err != nil {
					return err
				}
			}

			// tin nội bộ => dừng
			if post.Visibility == enums.PostVisiblePublic {
				if err := s.postRepo.UpdatePostStatus(ctx, productID, enums.PostStop); err != nil {
					return err
				}
			}
		}

	case enums.ProductSold:
		if buyerID != nil || buyerOrgID != nil {
			// Bán cho công ty hoặc người dùng
			product.OwnerUserID = buyerID
			product.OwnerOrgID = buyerOrgID
			if err := s.assetRepo.UpdateStatus(ctx, product.AssetId, enums.AssetStatusOwning); err != nil {
				return err
			}
		} else {
			// todo
			// Bán ra ngoài
			product.OwnerUserID = nil
			product.OwnerOrgID = nil
		}

		if err := s.postRepo.UpdatePostStatus(ctx, productID, enums.PostStatusCompleted); err != nil {
			return err
		}
		if err := s.postRepo.HideOrDeletePost(ctx, productID); err != nil {
			return err
		}
	}

	return s.productRepo.UpdateStatus(ctx, productID, newStatus)
}

func (s *ProductService) CreateProduct(c *gin.Context, product *data.NewProductRequest) (*models.Product, error) {
	profileId := _jwt.GetProfileId(c)
	entity := &models.Product{
		Name: product.Name,
	}
	_db.SaveWithAudit(c, entity)

	access := &models.ProductAccess{
		UserId:    profileId,
		ProductId: entity.ID,
	}
	err := s.productAccessRepo.Create(access)
	if err != nil {
		return nil, err
	}

	return entity, nil
}

func (s *ProductService) CheckAccessProduct(profileId, productId uint64) error {
	existed, err := s.productAccessRepo.ExistsByUserAndProduct(profileId, productId)
	if err != nil {
		return err
	}
	if !existed {
		return &_routes.Except{
			Code:    401,
			Message: "Bạn không có quyền truy cập",
		}
	}

	return nil
}

func (s *ProductService) SearchProductWithOwner(c *gin.Context, dto *data.ProductSearchRequest) ([]models.Product, error) {
	profileId := _jwt.GetProfileId(c)
	results, err := s.productRepo.SearchProductWithOwner(profileId, *dto, _utils.GetPage(c))
	return results, err
}

func (s *ProductService) Update(c *gin.Context, id uint64, dto *data.ProductSaveRequest) (*models.Product, error) {
	profileId := _jwt.GetProfileId(c)
	err := s.CheckAccessProduct(profileId, id)
	if err != nil {
		return nil, err
	}
	product, err := s.productRepo.GetByID(c, id)
	if err != nil {
		return nil, err
	}

	copier.Copy(&product, &dto)

	err = s.productRepo.UpdateProduct(c, product)
	return product, err
}
