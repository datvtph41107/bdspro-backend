package crud2

import "github.com/gin-gonic/gin"

// BaseRouterInterface định nghĩa các phương thức chung cho router
// type BaseRouterInterface interface {
// 	Create(c *gin.Context)
// 	Update(c *gin.Context)
// 	Delete(c *gin.Context)
// 	GetByID(c *gin.Context)
// 	GetAll(c *gin.Context)
// }

type IBaseService[T any] interface {
	Create(c *gin.Context) (*T, error)
	Update(c *gin.Context) (*T, error)
	Delete(c *gin.Context) (bool, error)
	GetByID(c *gin.Context) (*T, error)
	GetAll(c *gin.Context) ([]T, error)
}

type IBaseRepo[T any] interface {
	Create(c *gin.Context, entity *T) error
	Update(c *gin.Context, id uint64, entity *T) error
	Delete(id uint64) error
	GetAll() ([]T, error)
	GetByID(c *gin.Context, id uint64) (*T, error)
}
