package crud

import "github.com/gin-gonic/gin"

// BaseRouterInterface định nghĩa các phương thức chung cho router
type BaseRouterInterface interface {
	Create(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	GetByID(c *gin.Context)
	GetAll(c *gin.Context)
}
