package router

import (
	"payment/internal/wallet/usecases"

	"github.com/gin-gonic/gin"
)

type TransactionRouter struct {
	uc *usecases.AcountUsecase
}

func NewTransactionRouter(uc *usecases.AcountUsecase) *TransactionRouter {
	return &TransactionRouter{
		uc: uc,
	}
}

func (r *TransactionRouter) RegisterRouter(router *gin.RouterGroup) {
	// group := router.Group("/account")
	// {
	// 	// todo:
	// 	group.POST("", func(c *gin.Context) {

	// 	})
	// 	group.GET("/dashboard/post/:id", func(c *gin.Context) {

	// 	})
	// 	group.GET("/dashboard/post/bar", func(c *gin.Context) {

	// 	})
	// }
}
