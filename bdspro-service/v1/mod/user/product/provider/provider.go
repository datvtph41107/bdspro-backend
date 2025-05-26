package u_product_provider

// import (
// 	u_product "bdspro/v1/mod/user/product"
// 	"bdspro/v1/models"
// 	_jwt "common/jwt"
// 	_routes "common/routes"

// 	"github.com/gin-gonic/gin"
// )

// type ProductProvider struct {
// 	ProductRepo *u_product.UserProductRepo
// }

// func NewProductProvider(ProductRepo *u_product.UserProductRepo) *ProductProvider {
// 	return &ProductProvider{
// 		ProductRepo: ProductRepo,
// 	}
// }

// func (s *ProductProvider) RequiredOwnerProduct(c *gin.Context, id uint64) (*models.Product, error) {
// 	product, err := s.ProductRepo.GetByID(c, id)
// 	if err != nil {
// 		return nil, err
// 	}

// 	profileId := _jwt.GetProfileId(c)

// 	if *product.OwnerUserID != profileId {
// 		return nil, &_routes.Except{
// 			Code:    401,
// 			Message: "Bạn không có quyền truy cập",
// 		}
// 	}
// 	return product, nil
// }
