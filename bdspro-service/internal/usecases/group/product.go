package group_usecase

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"

	"github.com/gin-gonic/gin"
)

type GroupProductUsecase struct {
}

func NewGroupProductUsecase() *GroupProductUsecase {
	return &GroupProductUsecase{}
}

func (u *GroupProductUsecase) GetProduct(c *gin.Context, groupId int64, dto dto.ProductSearchRequest) ([]domain.Product, int64, error) {
	// products, err := u.productRepo.GetProduct(groupId, page, limit, productId)
	// if err != nil {
	// 	return nil, err
	// }
	// return products, nil
	return nil, 0, nil
}
