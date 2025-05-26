package iu_product

import (
	"bdspro/v1/models"

	"github.com/gin-gonic/gin"
)

type IUserProductService interface {
	RequiredOwnerProduct(c *gin.Context, id uint64) (*models.Product, error)
}
