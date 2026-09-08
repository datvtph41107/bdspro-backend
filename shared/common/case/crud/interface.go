package crud

// import "github.com/gin-gonic/gin"

// BaseRouterInterface định nghĩa các phương thức chung cho router
// type BaseRouterInterface interface {
// 	Create(c *gin.Context)
// 	Update(c *gin.Context)
// 	Delete(c *gin.Context)
// 	GetByID(c *gin.Context)
// 	GetAll(c *gin.Context)
// }

// type BaseServiceInterface[T any] interface {
// 	Create(c *gin.Context) (*T, error)
// 	Update(c *gin.Context) (*T, error)
// 	Delete(c *gin.Context) (bool, error)
// 	GetByID(c *gin.Context) (*T, error)
// 	GetAll(c *gin.Context) ([]T, error)
// }

// type BaseRepoInterface[T any] interface {
// 	Create(c *gin.Context, entity *T) error
// 	Update(c *gin.Context, id uint64, entity *T) error
// 	Delete(id uint64) error
// 	GetAll() ([]T, error)
// 	GetByID(id uint64) (*T, error)
// }

// type CrudRepo[T any] interface {
// 	Create(c *context.Context, entity *T) error
// 	Update(c *context.Context, id uint64, entity *T) error
// 	Delete(id uint64) error
// 	GetAll() ([]T, error)
// 	GetByID(id uint64) (*T, error)
// }
