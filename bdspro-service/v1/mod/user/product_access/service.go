package u_product_access

import (
	"bdspro/v1/enums"
	product_accecss "bdspro/v1/mod/global/product_access"
	u_product "bdspro/v1/mod/user/product"
	"bdspro/v1/models"
	"common/case/crud2"
	_db "common/db"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"gorm.io/gorm/clause"
)

type ProductAccessService struct {
	crud2.BaseService[models.ProductAccess, *product_accecss.ProductAccessRepo]
	productService *u_product.ProductService
}

func NewUserProductAccessService(repo *product_accecss.ProductAccessRepo, productService *u_product.ProductService) *ProductAccessService {
	return &ProductAccessService{
		BaseService:    crud2.BaseService[models.ProductAccess, *product_accecss.ProductAccessRepo]{Repo: repo},
		productService: productService,
	}
}

func (s *ProductAccessService) CreateBatch(c *gin.Context, id uint64, dtos []AccessRequest) (*[]models.ProductAccess, error) {
	var entities []models.ProductAccess
	product, err := s.productService.RequiredOwner(c, id)
	if err != nil {
		return nil, err
	}

	copier.Copy(&entities, &dtos)
	for i := 0; i < len(entities); i++ {
		entities[i].ProductID = id
	}

	if product.SaleVisibility == enums.ProductVisiblePrivate {
		product.SaleVisibility = enums.ProductVisibleInternal
		_db.SaveWithAudit(c, product)
	}
	// err := _db.SaveWithAudit(c, &entities)
	err = _db.DB.WithContext(c).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "product_id"}, {Name: "target_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"fields", "commission", "target_type"}),
		}).
		Save(&entities).Error

	return &entities, err
}
