package router

import (
	"payment/internal/wallet/usecases"

	"github.com/gin-gonic/gin"
)

type AccountRouter struct {
	uc *usecases.AcountUsecase
}

func NewAccountRouter(uc *usecases.AcountUsecase) *AccountRouter {
	return &AccountRouter{
		uc: uc,
	}
}

func (r *AccountRouter) RegisterRouter(router *gin.RouterGroup) {
	// group := router.Group("/account")
	{
		// todo:
		// group.POST("/new", func(c *gin.Context) {
		// 	dto := &domain.Account{}
		// 	if err := _utils.ParseBodyWithValidator(c, dto); err != nil {
		// 		_routes.RouteResult(c, nil, err)
		// 		return
		// 	}
		// 	result, err := r.uc.Register(dto)
		// 	_routes.RouteResult(c, result, err)
		// })
		// group.GET("/dashboard/post/:id", func(c *gin.Context) {

		// })
		// group.GET("/dashboard/post/bar", func(c *gin.Context) {

		// })
	}
}
